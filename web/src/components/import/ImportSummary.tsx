import * as React from 'react';
import { CheckCircle2, XCircle, AlertCircle, RotateCcw, ArrowLeft } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';
import type { ImportProgress, ImportStatus, ParsedImportRecord } from '@/types/import';

export interface ImportSummaryProps<T> {
  progress: ImportProgress;
  records: ParsedImportRecord<T>[];
  importStatus: Map<number, ImportStatus>;
  onRetryFailed?: () => void;
  onRetryRecord?: (index: number) => void;
  onReturnToList: () => void;
  entityType: 'labels' | 'providers' | 'subscriptions';
  className?: string;
}

export function ImportSummary<T>({
  progress,
  records,
  importStatus,
  onRetryFailed,
  onRetryRecord,
  onReturnToList,
  entityType,
  className,
}: ImportSummaryProps<T>) {
  const successCount = progress.completed - progress.failed;
  const failedRecords = React.useMemo(() => {
    return Array.from(importStatus.entries())
      .filter(([_, status]) => status.status === 'error')
      .map(([index, status]) => ({
        index,
        record: records[index],
        error: status.error || 'Unknown error',
      }));
  }, [importStatus, records]);

  const isFullSuccess = progress.completed === progress.total && progress.failed === 0;
  const isPartialSuccess = successCount > 0 && progress.failed > 0;
  const isFullFailure = progress.failed === progress.total;

  const getEntityLabel = () => {
    switch (entityType) {
      case 'labels':
        return 'Labels';
      case 'providers':
        return 'Providers';
      case 'subscriptions':
        return 'Subscriptions';
    }
  };

  const getRecordName = (record: ParsedImportRecord<T>): string => {
    const data = record.data as any;
    return data?.name || data?.friendlyName || `Record ${record.index + 1}`;
  };

  return (
    <Card className={cn('', className)} role="region" aria-label="Import summary">
      {/* Screen reader announcement */}
      <div className="sr-only" role="status" aria-live="polite" aria-atomic="true">
        {isFullSuccess && `Import completed successfully. ${successCount} records imported.`}
        {isPartialSuccess && `Import completed with errors. ${successCount} successful, ${progress.failed} failed.`}
        {isFullFailure && `Import failed. All ${progress.total} records failed to import.`}
      </div>

      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          {isFullSuccess && (
            <>
              <CheckCircle2 className="size-5 text-green-600 dark:text-green-400" aria-hidden="true" />
              Import Completed Successfully
            </>
          )}
          {isPartialSuccess && (
            <>
              <AlertCircle className="size-5 text-yellow-600 dark:text-yellow-400" aria-hidden="true" />
              Import Completed with Errors
            </>
          )}
          {isFullFailure && (
            <>
              <XCircle className="size-5 text-red-600 dark:text-red-400" aria-hidden="true" />
              Import Failed
            </>
          )}
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-6">
        {/* Summary stats */}
        <div className="grid grid-cols-3 gap-4" role="group" aria-label="Import summary statistics">
          <div className="flex flex-col items-center rounded-lg border p-4" role="group" aria-label={`Total records: ${progress.total}`}>
            <span className="text-3xl font-bold" aria-label={`${progress.total}`}>{progress.total}</span>
            <span className="text-sm text-muted-foreground">Total Records</span>
          </div>
          <div className="flex flex-col items-center rounded-lg border border-green-200 bg-green-50 p-4 dark:border-green-900 dark:bg-green-950/20" role="group" aria-label={`Successful imports: ${successCount}`}>
            <span className="text-3xl font-bold text-green-600 dark:text-green-400" aria-label={`${successCount}`}>
              {successCount}
            </span>
            <span className="text-sm text-muted-foreground">Successful</span>
          </div>
          <div className="flex flex-col items-center rounded-lg border border-red-200 bg-red-50 p-4 dark:border-red-900 dark:bg-red-950/20" role="group" aria-label={`Failed imports: ${progress.failed}`}>
            <span className="text-3xl font-bold text-red-600 dark:text-red-400" aria-label={`${progress.failed}`}>
              {progress.failed}
            </span>
            <span className="text-sm text-muted-foreground">Failed</span>
          </div>
        </div>

        {/* Success message */}
        {isFullSuccess && (
          <Alert className="border-2 border-green-300 dark:border-green-800 bg-green-50 dark:bg-green-950/30">
            <CheckCircle2 className="size-4 text-green-700 dark:text-green-300" aria-hidden="true" />
            <AlertTitle className="text-green-900 dark:text-green-100 font-bold">All records imported successfully</AlertTitle>
            <AlertDescription className="text-green-900 dark:text-green-100 font-medium">
              {successCount} {getEntityLabel().toLowerCase()} {successCount === 1 ? 'has' : 'have'} been
              added to your account.
            </AlertDescription>
          </Alert>
        )}

        {/* Partial success message */}
        {isPartialSuccess && (
          <Alert className="border-2 border-yellow-300 dark:border-yellow-800 bg-yellow-50 dark:bg-yellow-950/30">
            <AlertCircle className="size-4 text-yellow-700 dark:text-yellow-300" aria-hidden="true" />
            <AlertTitle className="text-yellow-900 dark:text-yellow-100 font-bold">Import completed with some errors</AlertTitle>
            <AlertDescription className="text-yellow-900 dark:text-yellow-100 font-medium">
              {successCount} {getEntityLabel().toLowerCase()} {successCount === 1 ? 'was' : 'were'}{' '}
              imported successfully, but {progress.failed}{' '}
              {progress.failed === 1 ? 'record' : 'records'} failed. Review the errors below and try
              again.
            </AlertDescription>
          </Alert>
        )}

        {/* Full failure message */}
        {isFullFailure && (
          <Alert variant="destructive" className="border-2 border-red-300 dark:border-red-800 bg-red-50 dark:bg-red-950/30">
            <XCircle className="size-4 text-red-700 dark:text-red-300" aria-hidden="true" />
            <AlertTitle className="text-red-900 dark:text-red-100 font-bold">Import failed</AlertTitle>
            <AlertDescription className="text-red-900 dark:text-red-100 font-medium">
              All {progress.total} {progress.total === 1 ? 'record' : 'records'} failed to import.
              Review the errors below and try again.
            </AlertDescription>
          </Alert>
        )}

        {/* Failed records list */}
        {failedRecords.length > 0 && (
          <div className="space-y-3" role="region" aria-label="Failed records">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold">Failed Records</h3>
              {onRetryFailed && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onRetryFailed}
                  className="gap-2"
                  aria-label={`Retry all ${failedRecords.length} failed records`}
                >
                  <RotateCcw className="size-3" aria-hidden="true" />
                  Retry All Failed
                </Button>
              )}
            </div>

            <ScrollArea className="h-[300px] rounded-lg border" aria-label="List of failed records">
              <div className="space-y-2 p-4" role="list">
                {failedRecords.map(({ index, record, error }, idx) => (
                  <div key={index} role="listitem">
                    <div className="flex items-start justify-between gap-3 rounded-lg bg-red-100 dark:bg-red-950/30 border-2 border-red-300 dark:border-red-800 p-3">
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs border-red-600 dark:border-red-400 bg-red-50 dark:bg-red-950/50 text-red-900 dark:text-red-100 font-bold" aria-label={`Record number ${index + 1}`}>
                            #{index + 1}
                          </Badge>
                          <span className="font-bold text-sm text-red-900 dark:text-red-100">
                            {getRecordName(record)}
                          </span>
                        </div>
                        <p className="text-xs text-red-900 dark:text-red-100 font-semibold" role="alert">
                          {error}
                        </p>
                      </div>
                      {onRetryRecord && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => onRetryRecord(index)}
                          className="shrink-0 gap-1 hover:bg-red-200 dark:hover:bg-red-900/50 focus-visible:ring-4 focus-visible:ring-red-500 focus-visible:ring-offset-2"
                          aria-label={`Retry importing ${getRecordName(record)}`}
                        >
                          <RotateCcw className="size-3" aria-hidden="true" />
                          Retry
                        </Button>
                      )}
                    </div>
                    {idx < failedRecords.length - 1 && <Separator className="my-2" />}
                  </div>
                ))}
              </div>
            </ScrollArea>
          </div>
        )}
      </CardContent>

      <CardFooter className="flex justify-between gap-3">
        <Button
          variant="outline"
          onClick={onReturnToList}
          className="gap-2"
          aria-label={`Return to ${getEntityLabel()} list`}
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Return to {getEntityLabel()}
        </Button>
        {isFullSuccess && (
          <Button 
            onClick={onReturnToList}
            aria-label={`View imported ${getEntityLabel()}`}
          >
            View {getEntityLabel()}
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}
