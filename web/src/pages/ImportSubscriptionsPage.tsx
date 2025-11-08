import { useState, useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Download } from 'lucide-react';
import { FileUploadZone } from '@/components/import/FileUploadZone';
import { ImportPreviewTable } from '@/components/import/ImportPreviewTable';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/ui/page-header';
import { fileParser, FileParseError } from '@/services/fileParser';
import { SubscriptionFieldMapper } from '@/services/importMapper';
import { useImportManager } from '@/hooks/import/useImportManager';
import { useSubscriptionsMutations } from '@/hooks/subscriptions/useSubscriptionsMutations';
import { toast } from 'sonner';
import type { ParsedImportRecord, ImportColumnDef } from '@/types/import';
import type { DtoCreateSubscriptionRequest } from '@/api';

const ACCEPTED_FORMATS = ['.csv', '.json', '.yaml', '.yml'];

const subscriptionMapper = new SubscriptionFieldMapper();

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
    label: 'Provider ID',
    render: (value) => (
      <span className="text-sm font-mono text-muted-foreground">
        {value ? String(value).substring(0, 8) + '...' : '-'}
      </span>
    ),
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
    key: 'customPrice',
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
      return value.type === 'family' ? `Family (${value.familyId || 'N/A'})` : 'Personal';
    },
  },
];

export default function ImportSubscriptionsPage() {
  const navigate = useNavigate();
  const { createSubscriptionMutation } = useSubscriptionsMutations();

  const [isParsingFile, setIsParsingFile] = useState(false);
  const [parseError, setParseError] = useState<string | undefined>();
  const [parsedRecords, setParsedRecords] = useState<ParsedImportRecord<DtoCreateSubscriptionRequest>[]>([]);
  const [selectedRecords, setSelectedRecords] = useState<Set<number>>(new Set());
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

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
        ownerType: data.owner?.type as any,
        familyId: data.owner?.familyId,
        payer: data.payer ? {
          type: data.payer.type as any,
          familyId: data.payer.familyId,
          memberId: data.payer.memberId,
        } : undefined,
        familyUsers: data.familyUsers,
        customPrice: data.customPrice ? {
          amount: data.customPrice.value,
          currency: data.customPrice.currency,
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
      toast.success(`Successfully imported ${successCount} subscriptions`);
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
      toast.success(`Successfully imported ${successCount} subscriptions`);
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
                <li><span className="font-medium">providerId</span> (required) - Provider UUID</li>
                <li><span className="font-medium">startDate</span> (required) - Start date (YYYY-MM-DD)</li>
                <li><span className="font-medium">recurrency</span> (required) - Billing cycle: daily, weekly, monthly, quarterly, yearly, custom</li>
                <li><span className="font-medium">friendlyName</span> (optional) - Subscription name</li>
                <li><span className="font-medium">endDate</span> (optional) - End date (YYYY-MM-DD)</li>
                <li><span className="font-medium">customRecurrency</span> (optional) - Custom recurrency in days (if recurrency is custom)</li>
                <li><span className="font-medium">customPriceAmount</span> (optional) - Price amount</li>
                <li><span className="font-medium">customPriceCurrency</span> (optional) - Currency code (e.g., USD, EUR)</li>
                <li><span className="font-medium">labels</span> (optional) - Comma-separated label names</li>
                <li><span className="font-medium">ownerType</span> (optional) - Owner type: personal, family, or system</li>
                <li><span className="font-medium">ownerFamilyId</span> (optional) - Family ID if ownerType is family</li>
                <li><span className="font-medium">payerType</span> (optional) - Payer type: family or family_member</li>
                <li><span className="font-medium">payerFamilyId</span> (optional) - Payer family ID</li>
                <li><span className="font-medium">payerMemberId</span> (optional) - Payer member ID (if payerType is family_member)</li>
                <li><span className="font-medium">freeTrialStartDate</span> (optional) - Free trial start date</li>
                <li><span className="font-medium">freeTrialEndDate</span> (optional) - Free trial end date</li>
                <li><span className="font-medium">familyUsers</span> (optional) - Comma-separated family user IDs</li>
              </ul>
              <p className="mt-3">
                <Button
                  variant="link"
                  size="sm"
                  className="h-auto p-0"
                  onClick={() => {
                    const yamlContent = `- providerId: "00000000-0000-0000-0000-000000000000"
  friendlyName: "Netflix Premium"
  startDate: "2024-01-01"
  recurrency: "monthly"
  customPriceAmount: 15.99
  customPriceCurrency: "USD"
  ownerType: "personal"
- providerId: "00000000-0000-0000-0000-000000000001"
  friendlyName: "AWS Cloud Services"
  startDate: "2024-02-01"
  recurrency: "monthly"
  customPriceAmount: 50.00
  customPriceCurrency: "USD"
  ownerType: "personal"`;
                    const blob = new Blob([yamlContent], { type: 'text/yaml' });
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = 'subscriptions-template.yaml';
                    a.click();
                    URL.revokeObjectURL(url);
                  }}
                >
                  <Download className="size-3 mr-1" />
                  Download YAML template
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
            columns={subscriptionColumns}
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
                Done - Return to Subscriptions
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
