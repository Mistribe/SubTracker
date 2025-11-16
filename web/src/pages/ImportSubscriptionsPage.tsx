import { useState, useCallback, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Download } from 'lucide-react';
import { FileUploadZone } from '@/components/import/FileUploadZone';
import { ImportPreviewTable } from '@/components/import/ImportPreviewTable';
import { ImportHelp } from '@/components/import/ImportHelp';
import { TemplateDownloadSection } from '@/components/import/TemplateDownloadSection';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/ui/page-header';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Skeleton } from '@/components/ui/skeleton';
import { fileParser, FileParseError, FileSizeError } from '@/services/fileParser';
import { SubscriptionFieldMapper } from '@/services/importMapper';
import { useImportManager } from '@/hooks/import/useImportManager';
import { useSubscriptionsMutations } from '@/hooks/subscriptions/useSubscriptionsMutations';
import { useProvidersByIds } from '@/hooks/providers/useProvidersByIds';
import { toast } from 'sonner';
import type { ParsedImportRecord, ImportColumnDef } from '@/types/import';
import type { DtoCreateSubscriptionRequest } from '@/api';

const ACCEPTED_FORMATS = ['.csv', '.json', '.yaml', '.yml'];

const subscriptionMapper = new SubscriptionFieldMapper();

// Small cell component to display provider information for a given providerId
function ProviderInfoCell({ providerId }: { providerId?: string | null }) {
  if (!providerId) {
    return <span className="text-sm text-muted-foreground">-</span>;
  }

  const { providerMap, isLoading } = useProvidersByIds([providerId]);
  const provider = providerMap.get(providerId);

  if (isLoading && !provider) {
    return (
      <div className="flex items-center gap-2">
        <Skeleton className="size-6 rounded-full" />
        <Skeleton className="h-4 w-28" />
      </div>
    );
  }

  if (!provider) {
    // Fallback to truncated ID when provider isn't found
    return (
      <span className="text-sm font-mono text-muted-foreground">
        {String(providerId).length > 11
          ? String(providerId).substring(0, 8) + '...'
          : String(providerId)}
      </span>
    );
  }

  const initial = provider.name?.charAt(0)?.toUpperCase() || '?';

  return (
    <div className="flex items-center gap-2 min-w-0">
      <Avatar className="size-6">
        <AvatarImage src={provider.iconUrl ?? undefined} alt={provider.name} />
        <AvatarFallback className="text-xs">{initial}</AvatarFallback>
      </Avatar>
      <span className="text-sm truncate max-w-[220px]" title={provider.name}>
        {provider.name}
      </span>
    </div>
  );
}

// Define columns for the preview table
const subscriptionColumns: ImportColumnDef<DtoCreateSubscriptionRequest>[] = [
  {
    key: 'friendlyName',
    label: 'Name',
    render: (value) => (
      <span className="text-sm">
        {value ? String(value) : '-'}
      </span>
    ),
  },
  {
    key: 'providerId',
    label: 'Provider',
    render: (value) => <ProviderInfoCell providerId={value ? String(value) : undefined} />,
  },
  {
    key: 'recurrency',
    label: 'Recurrency',
    render: (value) => (
      <span className="text-sm capitalize">
        {value ? String(value) : '-'}
      </span>
    ),
  },
  {
    key: 'startDate',
    label: 'Start Date',
    render: (value) => {
      if (!value) return '-';
      const date = value instanceof Date ? value : new Date(value as string);
      return (
        <span className="text-sm">
          {date.toLocaleDateString()}
        </span>
      );
    },
  },
  {
    key: 'price',
    label: 'Price',
    render: (value: any) => {
      if (!value || typeof value !== 'object') return '-';
      return (
        <span className="text-sm">
          {value.currency} {value.value}
        </span>
      );
    },
  },
  {
    key: 'owner',
    label: 'Owner',
    render: (value: any) => {
      if (!value) return 'Personal';
      return value === 'family' ? 'Family' : 'Personal';
    },
  },
];

export default function ImportSubscriptionsPage() {
  const navigate = useNavigate();
  const { createSubscriptionMutation } = useSubscriptionsMutations();

  const [isParsingFile, setIsParsingFile] = useState(false);
  const [parseProgress, setParseProgress] = useState(0);
  const [parseError, setParseError] = useState<string | undefined>();
  const [parsedRecords, setParsedRecords] = useState<ParsedImportRecord<DtoCreateSubscriptionRequest>[]>([]);
  const [selectedRecords, setSelectedRecords] = useState<Set<number>>(new Set());
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

  // Validate provider existence: gather unique provider IDs from parsed records
  const providerIds = useMemo(
    () => Array.from(new Set(parsedRecords
      .map(r => (r.data as any)?.providerId as string | undefined)
      .filter((id): id is string => !!id && id.length > 0)
    )),
    [parsedRecords]
  );

  // Fetch providers for existence checks (cached via React Query)
  const { providerMap: providerLookupMap, isLoading: isLoadingProviders } = useProvidersByIds(providerIds);

  // Derive records adding validation error if provider does not exist
  const computedRecords = useMemo(() => {
    if (parsedRecords.length === 0) return parsedRecords;

    return parsedRecords.map((record) => {
      const providerId = (record.data as any)?.providerId as string | undefined;
      let validationErrors = record.validationErrors;
      let isValid = record.isValid;

      // Only validate existence when we have providerId and loading is finished
      if (providerId && !isLoadingProviders) {
        const exists = providerLookupMap.get(providerId);
        if (!exists) {
          // Append a validation error without mutating original record
          validationErrors = [
            ...record.validationErrors,
            { field: 'providerId', message: 'Provider not found', severity: 'error' },
          ];
          isValid = false;
        }
      }

      if (validationErrors === record.validationErrors && isValid === record.isValid) {
        // No change
        return record;
      }

      return {
        ...record,
        validationErrors,
        isValid,
      } as ParsedImportRecord<DtoCreateSubscriptionRequest>;
    });
  }, [parsedRecords, providerLookupMap, isLoadingProviders]);

  // Create a wrapper mutation that accepts DtoCreateSubscriptionRequest
  const wrappedMutation = {
    ...createSubscriptionMutation,
    mutateAsync: async (data: DtoCreateSubscriptionRequest) => {
      return createSubscriptionMutation.mutateAsync({
        friendlyName: data.friendlyName,
        providerId: data.providerId!,
        recurrency: data.recurrency as any,
        customRecurrency: data.customRecurrency,
        startDate: data.startDate!,
        endDate: data.endDate,
        ownerType: data.owner,
        payer: data.payer ? {
          type: data.payer.type as any,
          memberId: data.payer.memberId,
        } : undefined,
        familyUsers: data.familyUsers,
        price: data.price ? {
          amount: data.price.value,
          currency: data.price.currency,
        } : undefined,
        freeTrial: data.freeTrial ? {
          startDate: data.freeTrial.startDate,
          endDate: data.freeTrial.endDate,
        } : undefined,
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
  } = useImportManager({
    records: computedRecords,
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
      const mapped: ParsedImportRecord<DtoCreateSubscriptionRequest>[] = rawRecords.map((rawRecord, index) => {
        const data = subscriptionMapper.mapFields(rawRecord);
        const validation = subscriptionMapper.validate(data);

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
    if (isLoadingProviders) {
      toast.info('Validating providers…', {
        description: 'Please wait until provider validation finishes before importing.',
      });
      return;
    }
    const validSelectedIndices = Array.from(selectedRecords).filter(
      index => computedRecords[index]?.isValid
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
      toast.success(`Successfully imported ${successCount} subscriptions`);
      setHasUnsavedChanges(false);
    }
  }, [selectedRecords, computedRecords, importRecords, progress, isLoadingProviders]);

  // Handle import all valid records
  const handleImportAll = useCallback(async () => {
    if (isLoadingProviders) {
      toast.info('Validating providers…', {
        description: 'Please wait until provider validation finishes before importing.',
      });
      return;
    }
    const validIndices = computedRecords
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
      toast.success(`Successfully imported ${successCount} subscriptions`);
      setHasUnsavedChanges(false);
    }
  }, [computedRecords, importRecords, progress, isLoadingProviders]);

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

    navigate('/subscriptions');
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
        title="Import Subscriptions"
        description="Import subscriptions from CSV, JSON, or YAML files"
        actionButton={
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleBack}
              disabled={isImporting}
            >
              <ArrowLeft className="size-4 mr-2" />
              Back to Subscriptions
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
      <TemplateDownloadSection entityType="subscriptions" />

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
          <ImportHelp entityType="subscriptions" />
        </div>
      )}

      {/* Preview table */}
      {parsedRecords.length > 0 && (
        <div className="space-y-4">
          <ImportPreviewTable
            records={computedRecords}
            columns={subscriptionColumns}
            selectedRecords={selectedRecords}
            onSelectionChange={setSelectedRecords}
            onImportSelected={handleImportSelected}
            onImportAll={handleImportAll}
            importStatus={importStatus}
            isImporting={isImporting}
            progress={progress}
            onRetryRecord={retryRecord}
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
                Done - Return to Subscriptions
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
