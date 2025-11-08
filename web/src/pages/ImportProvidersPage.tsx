import { useState, useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Download } from 'lucide-react';
import { FileUploadZone } from '@/components/import/FileUploadZone';
import { ImportPreviewTable } from '@/components/import/ImportPreviewTable';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/ui/page-header';
import { fileParser, FileParseError } from '@/services/fileParser';
import { ProviderFieldMapper } from '@/services/importMapper';
import { useImportManager } from '@/hooks/import/useImportManager';
import { useProvidersMutations } from '@/hooks/providers/useProvidersMutations';
import { toast } from 'sonner';
import type { ParsedImportRecord, ImportColumnDef } from '@/types/import';
import type { DtoCreateProviderRequest } from '@/api';

const ACCEPTED_FORMATS = ['.csv', '.json', '.yaml', '.yml'];

const providerMapper = new ProviderFieldMapper();

// Define columns for the preview table
const providerColumns: ImportColumnDef<DtoCreateProviderRequest>[] = [
  {
    key: 'name',
    label: 'Name',
  },
  {
    key: 'description',
    label: 'Description',
    render: (value) => (
      <span className="text-sm text-muted-foreground">
        {value ? String(value) : '-'}
      </span>
    ),
  },
  {
    key: 'url',
    label: 'URL',
    render: (value) => {
      if (!value) return '-';
      return (
        <a
          href={value as string}
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm text-blue-600 hover:underline"
          onClick={(e) => e.stopPropagation()}
        >
          {value as string}
        </a>
      );
    },
  },
  {
    key: 'labels',
    label: 'Labels',
    render: (value) => {
      if (!value || !Array.isArray(value) || value.length === 0) return '-';
      return (
        <span className="text-sm">
          {value.join(', ')}
        </span>
      );
    },
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

export default function ImportProvidersPage() {
  const navigate = useNavigate();
  const { createProviderMutation } = useProvidersMutations();

  const [isParsingFile, setIsParsingFile] = useState(false);
  const [parseError, setParseError] = useState<string | undefined>();
  const [parsedRecords, setParsedRecords] = useState<ParsedImportRecord<DtoCreateProviderRequest>[]>([]);
  const [selectedRecords, setSelectedRecords] = useState<Set<number>>(new Set());
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

  // Create a wrapper mutation that accepts DtoCreateProviderRequest
  const wrappedMutation = {
    ...createProviderMutation,
    mutateAsync: async (data: DtoCreateProviderRequest) => {
      return createProviderMutation.mutateAsync({
        name: data.name!,
        description: data.description,
        url: data.url,
        iconUrl: data.iconUrl,
        pricingPageUrl: data.pricingPageUrl,
        labels: data.labels,
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
      const mapped: ParsedImportRecord<DtoCreateProviderRequest>[] = rawRecords.map((rawRecord, index) => {
        const data = providerMapper.mapFields(rawRecord);
        const validation = providerMapper.validate(data);

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
      toast.success(`Successfully imported ${successCount} providers`);
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
      toast.success(`Successfully imported ${successCount} providers`);
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

    navigate('/providers');
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
        title="Import Providers"
        description="Import providers from CSV, JSON, or YAML files"
        actionButton={
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleBack}
              disabled={isImporting}
            >
              <ArrowLeft className="size-4 mr-2" />
              Back to Providers
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
                <li><span className="font-medium">name</span> (required) - Provider name</li>
                <li><span className="font-medium">description</span> (optional) - Provider description</li>
                <li><span className="font-medium">url</span> (optional) - Provider website URL</li>
                <li><span className="font-medium">iconUrl</span> (optional) - Provider icon URL</li>
                <li><span className="font-medium">pricingPageUrl</span> (optional) - Pricing page URL</li>
                <li><span className="font-medium">labels</span> (optional) - Comma-separated label names</li>
                <li><span className="font-medium">ownerType</span> (optional) - Owner type: personal, family, or system</li>
                <li><span className="font-medium">ownerFamilyId</span> (optional) - Family ID if ownerType is family</li>
              </ul>
              <p className="mt-3">
                <Button
                  variant="link"
                  size="sm"
                  className="h-auto p-0"
                  onClick={() => {
                    const jsonContent = JSON.stringify([
                      {
                        name: 'Netflix',
                        description: 'Streaming service',
                        url: 'https://netflix.com',
                        labels: 'entertainment,streaming'
                      },
                      {
                        name: 'AWS',
                        description: 'Cloud services',
                        url: 'https://aws.amazon.com',
                        labels: 'cloud,infrastructure'
                      }
                    ], null, 2);
                    const blob = new Blob([jsonContent], { type: 'application/json' });
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = 'providers-template.json';
                    a.click();
                    URL.revokeObjectURL(url);
                  }}
                >
                  <Download className="size-3 mr-1" />
                  Download JSON template
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
            columns={providerColumns}
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
                Done - Return to Providers
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
