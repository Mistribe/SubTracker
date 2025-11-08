import { useState, useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Download } from 'lucide-react';
import { FileUploadZone } from '@/components/import/FileUploadZone';
import { ImportPreviewTable } from '@/components/import/ImportPreviewTable';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/ui/page-header';
import { fileParser, FileParseError } from '@/services/fileParser';
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
        ownerType: data.owner?.type as any,
        familyId: data.owner?.familyId,
      });
    },
  } as any;

  const {
    importRecords,
    importStatus,
    progress,
    isImporting,
    cancelImport,
  } = useImportManager({
    records: parsedRecords,
    createMutation: wrappedMutation,
  });

  // Handle file selection and parsing
  const handleFileSelected = useCallback(async (file: File) => {
    setIsParsingFile(true);
    setParseError(undefined);
    setParsedRecords([]);
    setSelectedRecords(new Set());
    setHasUnsavedChanges(false);

    try {
      // Parse the file
      const rawRecords = await fileParser.parse(file);

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
      if (error instanceof FileParseError) {
        const errorMsg = error.errors.length > 0
          ? `${error.message}: ${error.errors[0].message}`
          : error.message;
        setParseError(errorMsg);
        toast.error('Failed to parse file', { description: errorMsg });
      } else {
        const errorMsg = error instanceof Error ? error.message : 'Unknown error';
        setParseError(errorMsg);
        toast.error('Failed to parse file', { description: errorMsg });
      }
    } finally {
      setIsParsingFile(false);
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
    const successCount = progress.completed - progress.failed;
    if (progress.failed > 0) {
      toast.warning('Import completed with errors', {
        description: `${successCount} succeeded, ${progress.failed} failed`,
      });
    } else {
      toast.success(`Successfully imported ${successCount} labels`);
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
    const successCount = progress.completed - progress.failed;
    if (progress.failed > 0) {
      toast.warning('Import completed with errors', {
        description: `${successCount} succeeded, ${progress.failed} failed`,
      });
    } else {
      toast.success(`Successfully imported ${successCount} labels`);
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

      {/* File upload zone */}
      {parsedRecords.length === 0 && (
        <div className="space-y-4">
          <FileUploadZone
            onFileSelected={handleFileSelected}
            acceptedFormats={ACCEPTED_FORMATS}
            isLoading={isParsingFile}
            error={parseError}
          />

          {/* Help text and template download */}
          <div className="rounded-lg border bg-muted/50 p-4 space-y-3">
            <h3 className="font-semibold text-sm">File Format Requirements</h3>
            <div className="text-sm text-muted-foreground space-y-2">
              <p>Your file must include the following fields:</p>
              <ul className="list-disc list-inside space-y-1 ml-2">
                <li><span className="font-medium">name</span> (required) - Label name</li>
                <li><span className="font-medium">color</span> (required) - Hex color code (e.g., #FF5733)</li>
                <li><span className="font-medium">ownerType</span> (optional) - Owner type: personal, family, or system</li>
                <li><span className="font-medium">ownerFamilyId</span> (optional) - Family ID if ownerType is family</li>
              </ul>
              <p className="mt-3">
                <Button
                  variant="link"
                  size="sm"
                  className="h-auto p-0"
                  onClick={() => {
                    const csvContent = 'name,color\nEntertainment,#FF5733\nUtilities,#33FF57\nSubscriptions,#3357FF';
                    const blob = new Blob([csvContent], { type: 'text/csv' });
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = 'labels-template.csv';
                    a.click();
                    URL.revokeObjectURL(url);
                  }}
                >
                  <Download className="size-3 mr-1" />
                  Download CSV template
                </Button>
              </p>
            </div>
          </div>
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
