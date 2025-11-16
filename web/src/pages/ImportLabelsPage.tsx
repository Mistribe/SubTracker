import { useState, useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { FileUploadZone } from '@/components/import/FileUploadZone';
import { ImportPreviewTable } from '@/components/import/ImportPreviewTable';
import { ImportHelp } from '@/components/import/ImportHelp';
import { TemplateDownloadSection } from '@/components/import/TemplateDownloadSection';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/ui/page-header';
import { fileParser, FileParseError, FileSizeError } from '@/services/fileParser';
import { LabelFieldMapper } from '@/services/importMapper';
import { useImportManager } from '@/hooks/import/useImportManager';
import { useLabelMutations } from '@/hooks/labels/useLabelMutations';
import { toast } from 'sonner';
import type { ParsedImportRecord, ImportColumnDef } from '@/types/import';
import type { DtoCreateLabelRequest } from '@/api';

const ACCEPTED_FORMATS = ['.csv', '.json', '.yaml', '.yml'];

const labelMapper = new LabelFieldMapper();

// Define columns for the preview table
const labelColumns: ImportColumnDef<DtoCreateLabelRequest>[] = [
  {
    key: 'name',
    label: 'Name',
  },
  {
    key: 'color',
    label: 'Color',
    render: (value) => (
      <div className="flex items-center gap-2">
        {value && (
          <div
            className="size-4 rounded border"
            style={{ backgroundColor: value as string }}
          />
        )}
        <span>{value as string}</span>
      </div>
    ),
  },
  {
    key: 'owner',
    label: 'Owner',
    render: (value: any) => {
      if (!value) return 'Personal';
      return value.type === 'family' ? `Family (${value.familyId || 'N/A'})` : 'Personal';
    },
  },
];

export default function ImportLabelsPage() {
  const navigate = useNavigate();
  const { createLabelMutation } = useLabelMutations();

  const [isParsingFile, setIsParsingFile] = useState(false);
  const [parseProgress, setParseProgress] = useState(0);
  const [parseError, setParseError] = useState<string | undefined>();
  const [parsedRecords, setParsedRecords] = useState<ParsedImportRecord<DtoCreateLabelRequest>[]>([]);
  const [selectedRecords, setSelectedRecords] = useState<Set<number>>(new Set());
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

  // Create a wrapper mutation that accepts DtoCreateLabelRequest
  const wrappedMutation = {
    ...createLabelMutation,
    mutateAsync: async (data: DtoCreateLabelRequest) => {
      return createLabelMutation.mutateAsync({
        name: data.name!,
        color: data.color!,
        ownerType: data.owner,
      });
    },
  } as any;

  const {
    importRecords,
    retryRecord,
    importStatus,
    progress,
    isImporting,
    cancelImport,
    resetImport,
  } = useImportManager({
    records: parsedRecords,
    createMutation: wrappedMutation,
  });

  // Handle file selection and parsing
  const handleFileSelected = useCallback(async (file: File) => {
    setIsParsingFile(true);
    setParseProgress(0);
    setParseError(undefined);
    setParsedRecords([]);
    setSelectedRecords(new Set());
    setHasUnsavedChanges(false);

    try {
      // Parse the file with progress tracking
      const rawRecords = await fileParser.parse(file, (progress) => {
        setParseProgress(progress);
      });

      if (rawRecords.length === 0) {
        setParseError('The file contains no records');
        setIsParsingFile(false);
        return;
      }

      // Map and validate records
      const mapped: ParsedImportRecord<DtoCreateLabelRequest>[] = rawRecords.map((rawRecord, index) => {
        const data = labelMapper.mapFields(rawRecord);
        const validation = labelMapper.validate(data);

        return {
          index,
          data,
          validationErrors: validation.errors,
          isValid: validation.isValid,
        };
      });

      setParsedRecords(mapped);
      setHasUnsavedChanges(true);

      // Show success message
      const validCount = mapped.filter(r => r.isValid).length;
      const invalidCount = mapped.length - validCount;

      if (invalidCount > 0) {
        toast.warning(`Parsed ${mapped.length} records`, {
          description: `${validCount} valid, ${invalidCount} invalid`,
        });
      } else {
        toast.success(`Parsed ${mapped.length} valid records`);
      }
    } catch (error) {
      if (error instanceof FileSizeError) {
        setParseError(error.message);
        toast.error('File too large', { description: error.message });
      } else if (error instanceof FileParseError) {
        // Format error message with line numbers if available
        let errorMsg = error.message;
        if (error.errors.length > 0) {
          const errorDetails = error.errors
            .slice(0, 3) // Show first 3 errors
            .map(err => err.line ? `Line ${err.line}: ${err.message}` : err.message)
            .join('\n');
          errorMsg = `${error.message}\n\n${errorDetails}`;
          if (error.errors.length > 3) {
            errorMsg += `\n... and ${error.errors.length - 3} more errors`;
          }
        }
        setParseError(errorMsg);
        toast.error('Failed to parse file', { 
          description: error.errors.length > 0 
            ? `${error.errors.length} error(s) found. Check the upload zone for details.`
            : error.message 
        });
      } else {
        const errorMsg = error instanceof Error ? error.message : 'Unknown error';
        setParseError(errorMsg);
        toast.error('Failed to parse file', { description: errorMsg });
      }
    } finally {
      setIsParsingFile(false);
      setParseProgress(0);
    }
  }, []);

  // Handle import selected records
  const handleImportSelected = useCallback(async () => {
    const validSelectedIndices = Array.from(selectedRecords).filter(
      index => parsedRecords[index]?.isValid
    );

    if (validSelectedIndices.length === 0) {
      toast.error('No valid records selected');
      return;
    }

    await importRecords(validSelectedIndices);

    // Show completion message
    const skipped = progress.skipped ?? 0;
    const successCount = Math.max(0, progress.completed - progress.failed - skipped);
    if (progress.failed > 0) {
      toast.warning('Import completed with issues', {
        description: `${successCount} created, ${skipped} already existed, ${progress.failed} failed`,
      });
    } else {
      const baseMsg = skipped > 0
        ? `${successCount} created, ${skipped} already existed`
        : `${successCount} labels imported`;
      toast.success(`Import completed: ${baseMsg}`);
      setHasUnsavedChanges(false);
    }
  }, [selectedRecords, parsedRecords, importRecords, progress]);

  // Handle import all valid records
  const handleImportAll = useCallback(async () => {
    const validIndices = parsedRecords
      .map((record, index) => (record.isValid ? index : -1))
      .filter(index => index !== -1);

    if (validIndices.length === 0) {
      toast.error('No valid records to import');
      return;
    }

    await importRecords(validIndices);

    // Show completion message
    const skipped = progress.skipped ?? 0;
    const successCount = Math.max(0, progress.completed - progress.failed - skipped);
    if (progress.failed > 0) {
      toast.warning('Import completed with issues', {
        description: `${successCount} created, ${skipped} already existed, ${progress.failed} failed`,
      });
    } else {
      const baseMsg = skipped > 0
        ? `${successCount} created, ${skipped} already existed`
        : `${successCount} labels imported`;
      toast.success(`Import completed: ${baseMsg}`);
      setHasUnsavedChanges(false);
    }
  }, [parsedRecords, importRecords, progress]);

  // Handle navigation back
  const handleBack = useCallback(() => {
    if (hasUnsavedChanges && parsedRecords.length > 0 && !isImporting) {
      const unimportedCount = parsedRecords.filter(
        (_, index) => !importStatus.get(index) || importStatus.get(index)?.status === 'pending'
      ).length;

      if (unimportedCount > 0) {
        const confirmed = window.confirm(
          `You have ${unimportedCount} unimported records. Are you sure you want to leave?`
        );
        if (!confirmed) return;
      }
    }

    navigate('/labels');
  }, [hasUnsavedChanges, parsedRecords, isImporting, importStatus, navigate]);

  // Warn before leaving page with unsaved changes
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (hasUnsavedChanges && parsedRecords.length > 0 && !isImporting) {
        const unimportedCount = parsedRecords.filter(
          (_, index) => !importStatus.get(index) || importStatus.get(index)?.status === 'pending'
        ).length;

        if (unimportedCount > 0) {
          e.preventDefault();
          e.returnValue = '';
        }
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [hasUnsavedChanges, parsedRecords, isImporting, importStatus]);

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="Import Labels"
        description="Import labels from CSV, JSON, or YAML files"
        actionButton={
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleBack}
              disabled={isImporting}
            >
              <ArrowLeft className="size-4 mr-2" />
              Back to Labels
            </Button>
            {isImporting && (
              <Button
                variant="destructive"
                size="sm"
                onClick={cancelImport}
              >
                Cancel Import
              </Button>
            )}
          </div>
        }
      />

      {/* Template download section */}
      <TemplateDownloadSection entityType="labels" />

      {/* File upload zone */}
      {parsedRecords.length === 0 && (
        <div className="space-y-4">
          <FileUploadZone
            onFileSelected={handleFileSelected}
            acceptedFormats={ACCEPTED_FORMATS}
            isLoading={isParsingFile}
            error={parseError}
            parseProgress={parseProgress}
          />

          {/* Help documentation */}
          <ImportHelp entityType="labels" />
        </div>
      )}

      {/* Preview table */}
      {parsedRecords.length > 0 && (
        <div className="space-y-4">
          <ImportPreviewTable
            records={parsedRecords}
            columns={labelColumns}
            selectedRecords={selectedRecords}
            onSelectionChange={setSelectedRecords}
            onImportSelected={handleImportSelected}
            onImportAll={handleImportAll}
            importStatus={importStatus}
            isImporting={isImporting}
            progress={progress}
            onRetryRecord={retryRecord}
            onRemoveSelected={(indices) => {
              if (isImporting) return;
              const removeSet = new Set(indices);
              resetImport();
              setParsedRecords(prev => prev.filter((_, idx) => !removeSet.has(idx)));
              setSelectedRecords(new Set());
              setHasUnsavedChanges(true);
            }}
          />

          {/* Action buttons at bottom */}
          <div className="flex items-center justify-between">
            <Button
              variant="outline"
              onClick={() => {
                if (hasUnsavedChanges) {
                  const confirmed = window.confirm(
                    'Are you sure you want to select a new file? Current progress will be lost.'
                  );
                  if (!confirmed) return;
                }
                setParsedRecords([]);
                setSelectedRecords(new Set());
                setHasUnsavedChanges(false);
              }}
              disabled={isImporting}
            >
              Select New File
            </Button>

            {progress.completed > 0 && !progress.inProgress && (
              <Button onClick={handleBack}>
                Done - Return to Labels
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
