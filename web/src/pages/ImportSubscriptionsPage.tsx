import {useCallback, useEffect, useMemo, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {ArrowLeft} from 'lucide-react';
import {FileUploadZone} from '@/components/import/FileUploadZone';
import {ImportPreviewTable} from '@/components/import/ImportPreviewTable';
import {ImportHelp} from '@/components/import/ImportHelp';
import {TemplateDownloadSection} from '@/components/import/TemplateDownloadSection';
import {Button} from '@/components/ui/button';
import {PageHeader} from '@/components/ui/page-header';
import {Avatar, AvatarFallback, AvatarImage} from '@/components/ui/avatar';
import {Skeleton} from '@/components/ui/skeleton';
import {FileParseError, fileParser, FileSizeError} from '@/services/fileParser';
import {SubscriptionFieldMapper} from '@/services/importMapper';
import {useImportManager} from '@/hooks/import/useImportManager';
import {useSubscriptionsMutations} from '@/hooks/subscriptions/useSubscriptionsMutations';
import {useAllProvidersQuery} from '@/hooks/providers/useAllProvidersQuery';
import {toast} from 'sonner';
import type {ImportColumnDef, ParsedImportRecord} from '@/types/import';
import type {DtoCreateSubscriptionRequest} from '@/api';
import {useProvidersByIds} from "@/hooks/providers/useProvidersByIds.ts";
import {useProviderQuery} from "@/hooks/providers/useProviderQuery.ts";

const ACCEPTED_FORMATS = ['.csv', '.json', '.yaml', '.yml'];

const subscriptionMapper = new SubscriptionFieldMapper();

// Display provider avatar and name based on providerKey
function ProviderInfoCell({providerKey}: { providerKey?: string | null }) {
    if (!providerKey) {
        return <span className="text-sm text-muted-foreground">-</span>;
    }

    const {data, isLoading} = useProviderQuery(providerKey)
    const provider = data;

    if (isLoading && !provider) {
        return (
            <div className="flex items-center gap-2">
                <Skeleton className="size-6 rounded-full" />
                <Skeleton className="h-4 w-28" />
            </div>
        );
    }

    if (!provider) {
        return (
            <span className="text-sm font-mono text-muted-foreground" title={providerKey}>
                {providerKey}
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
        key: 'providerKey',
        label: 'Provider',
        render: (value) => <ProviderInfoCell providerKey={value ? String(value) : undefined} />,
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
        key: 'endDate',
        label: 'End Date',
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
    const {createSubscriptionMutation} = useSubscriptionsMutations();

    const [isParsingFile, setIsParsingFile] = useState(false);
    const [parseProgress, setParseProgress] = useState(0);
    const [parseError, setParseError] = useState<string | undefined>();
    const [parsedRecords, setParsedRecords] = useState<ParsedImportRecord<DtoCreateSubscriptionRequest>[]>([]);
    const [selectedRecords, setSelectedRecords] = useState<Set<number>>(new Set());
    const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

    // Derive records (provider existence checks are handled server-side; we only validate structure here)
    const computedRecords = useMemo(() => parsedRecords, [parsedRecords]);

    // Create a wrapper mutation that accepts DtoCreateSubscriptionRequest
    const wrappedMutation = {
        ...createSubscriptionMutation,
        mutateAsync: async (data: DtoCreateSubscriptionRequest) => {
            return createSubscriptionMutation.mutateAsync({
                subscriptionId: data.id,
                friendlyName: data.friendlyName,
                providerKey: (data as any).providerKey,
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
    resetImport,
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
                toast.error('File too large', {description: error.message});
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
                toast.error('Failed to parse file', {description: errorMsg});
            }
        } finally {
            setIsParsingFile(false);
            setParseProgress(0);
        }
    }, []);

    // Handle import selected records
    const handleImportSelected = useCallback(async () => {
        const validSelectedIndices = Array.from(selectedRecords).filter(
            index => computedRecords[index]?.isValid
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
                : `${successCount} subscriptions imported`;
            toast.success(`Import completed: ${baseMsg}`);
            setHasUnsavedChanges(false);
        }
    }, [selectedRecords, computedRecords, importRecords, progress]);

    // Handle import all valid records
    const handleImportAll = useCallback(async () => {
        const validIndices = computedRecords
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
                : `${successCount} subscriptions imported`;
            toast.success(`Import completed: ${baseMsg}`);
            setHasUnsavedChanges(false);
        }
    }, [computedRecords, importRecords, progress]);

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
                            <ArrowLeft className="size-4 mr-2"/>
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
            <TemplateDownloadSection entityType="subscriptions"/>

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
                    <ImportHelp entityType="subscriptions"/>
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
                        onRemoveSelected={(indices) => {
                          if (isImporting) return;
                          // Remove selected records based on provided indices
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
                                Done - Return to Subscriptions
                            </Button>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}
