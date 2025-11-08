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
    <Card className={cn('', className)}>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          {isFullSuccess && (
            <>
              <CheckCircle2 className="size-5 text-green-600 dark:text-green-400" />
              Import Completed Successfully
            </>
          )}
          {isPartialSuccess && (
            <>
              <AlertCircle className="size-5 text-yellow-600 dark:text-yellow-400" />
              Import Completed with Errors
            </>
          )}
          {isFullFailure && (
            <>
              <XCircle className="size-5 text-red-600 dark:text-red-400" />
              Import Failed
            </>
          )}
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-6">
        {/* Summary stats */}
        <div className="grid grid-cols-3 gap-4">
          <div className="flex flex-col items-center rounded-lg border p-4">
            <span className="text-3xl font-bold">{progress.total}</span>
            <span className="text-sm text-muted-foreground">Total Records</span>
          </div>
          <div className="flex flex-col items-center rounded-lg border border-green-200 bg-green-50 p-4 dark:border-green-900 dark:bg-green-950/20">
            <span className="text-3xl font-bold text-green-600 dark:text-green-400">
              {successCount}
            </span>
            <span className="text-sm text-muted-foreground">Successful</span>
          </div>
          <div className="flex flex-col items-center rounded-lg border border-red-200 bg-red-50 p-4 dark:border-red-900 dark:bg-red-950/20">
            <span className="text-3xl font-bold text-red-600 dark:text-red-400">
              {progress.failed}
            </span>
            <span className="text-sm text-muted-foreground">Failed</span>
          </div>
        </div>

        {/* Success message */}
        {isFullSuccess && (
          <Alert>
            <CheckCircle2 className="size-4" />
            <AlertTitle>All records imported successfully</AlertTitle>
            <AlertDescription>
              {successCount} {getEntityLabel().toLowerCase()} {successCount === 1 ? 'has' : 'have'} been
              added to your account.
            </AlertDescription>
          </Alert>
        )}

        {/* Partial success message */}
        {isPartialSuccess && (
          <Alert>
            <AlertCircle className="size-4" />
            <AlertTitle>Import completed with some errors</AlertTitle>
            <AlertDescription>
              {successCount} {getEntityLabel().toLowerCase()} {successCount === 1 ? 'was' : 'were'}{' '}
              imported successfully, but {progress.failed}{' '}
              {progress.failed === 1 ? 'record' : 'records'} failed. Review the errors below and try
              again.
            </AlertDescription>
          </Alert>
        )}

        {/* Full failure message */}
        {isFullFailure && (
          <Alert variant="destructive">
            <XCircle className="size-4" />
            <AlertTitle>Import failed</AlertTitle>
            <AlertDescription>
              All {progress.total} {progress.total === 1 ? 'record' : 'records'} failed to import.
              Review the errors below and try again.
            </AlertDescription>
          </Alert>
        )}

        {/* Failed records list */}
        {failedRecords.length > 0 && (
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold">Failed Records</h3>
              {onRetryFailed && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onRetryFailed}
                  className="gap-2"
                >
                  <RotateCcw className="size-3" />
                  Retry All Failed
                </Button>
              )}
            </div>

            <ScrollArea className="h-[300px] rounded-lg border">
              <div className="space-y-2 p-4">
                {failedRecords.map(({ index, record, error }, idx) => (
                  <div key={index}>
                    <div className="flex items-start justify-between gap-3 rounded-lg bg-red-50 p-3 dark:bg-red-950/20">
                      <div className="flex-1 space-y-1">
                        <div className="flex items-center gap-2">
                          <Badge variant="outline" className="text-xs">
                            #{index + 1}
                          </Badge>
                          <span className="font-medium text-sm">
                            {getRecordName(record)}
                          </span>
                        </div>
                        <p className="text-xs text-red-600 dark:text-red-400">
                          {error}
                        </p>
                      </div>
                      {onRetryRecord && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => onRetryRecord(index)}
                          className="shrink-0 gap-1"
                        >
                          <RotateCcw className="size-3" />
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
        >
          <ArrowLeft className="size-4" />
          Return to {getEntityLabel()}
        </Button>
        {isFullSuccess && (
          <Button onClick={onReturnToList}>
            View {getEntityLabel()}
          </Button>
        )}
      </CardFooter>
    </Card>
  );
}
