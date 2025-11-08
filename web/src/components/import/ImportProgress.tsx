import * as React from 'react';
import { Loader2, CheckCircle2, XCircle, Clock } from 'lucide-react';
import { Progress } from '@/components/ui/progress';
import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';
import type { ImportProgress as ImportProgressType } from '@/types/import';

export interface ImportProgressProps {
  progress: ImportProgressType;
  startTime?: number;
  className?: string;
}

export function ImportProgress({ progress, startTime, className }: ImportProgressProps) {
  const [estimatedTimeRemaining, setEstimatedTimeRemaining] = React.useState<number | null>(null);

  React.useEffect(() => {
    if (!startTime || !progress.inProgress || progress.completed === 0) {
      setEstimatedTimeRemaining(null);
      return;
    }

    const elapsedTime = Date.now() - startTime;
    const averageTimePerRecord = elapsedTime / progress.completed;
    const remainingRecords = progress.total - progress.completed;
    const estimated = Math.ceil((averageTimePerRecord * remainingRecords) / 1000); // in seconds

    setEstimatedTimeRemaining(estimated);
  }, [startTime, progress.completed, progress.total, progress.inProgress]);

  const percentage = progress.total > 0 ? (progress.completed / progress.total) * 100 : 0;
  const successCount = progress.completed - progress.failed;

  const formatTime = (seconds: number): string => {
    if (seconds < 60) {
      return `${seconds}s`;
    }
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}m ${remainingSeconds}s`;
  };

  return (
    <Card className={cn('border-primary/20', className)}>
      <CardContent className="space-y-4">
        {/* Progress bar */}
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="font-medium">
              {progress.inProgress ? 'Importing...' : 'Import Complete'}
            </span>
            <span className="text-muted-foreground">
              {Math.round(percentage)}%
            </span>
          </div>
          <Progress value={percentage} className="h-3" />
        </div>

        {/* Stats grid */}
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          {/* Total */}
          <div className="flex items-center gap-2">
            <div className="flex size-8 items-center justify-center rounded-full bg-muted">
              <Loader2
                className={cn(
                  'size-4',
                  progress.inProgress && 'animate-spin text-primary'
                )}
              />
            </div>
            <div className="flex flex-col">
              <span className="text-2xl font-bold">{progress.total}</span>
              <span className="text-xs text-muted-foreground">Total</span>
            </div>
          </div>

          {/* Completed */}
          <div className="flex items-center gap-2">
            <div className="flex size-8 items-center justify-center rounded-full bg-muted">
              <CheckCircle2 className="size-4 text-muted-foreground" />
            </div>
            <div className="flex flex-col">
              <span className="text-2xl font-bold">{progress.completed}</span>
              <span className="text-xs text-muted-foreground">Completed</span>
            </div>
          </div>

          {/* Success */}
          <div className="flex items-center gap-2">
            <div className="flex size-8 items-center justify-center rounded-full bg-green-100 dark:bg-green-950">
              <CheckCircle2 className="size-4 text-green-600 dark:text-green-400" />
            </div>
            <div className="flex flex-col">
              <span className="text-2xl font-bold text-green-600 dark:text-green-400">
                {successCount}
              </span>
              <span className="text-xs text-muted-foreground">Success</span>
            </div>
          </div>

          {/* Failed */}
          <div className="flex items-center gap-2">
            <div className="flex size-8 items-center justify-center rounded-full bg-red-100 dark:bg-red-950">
              <XCircle className="size-4 text-red-600 dark:text-red-400" />
            </div>
            <div className="flex flex-col">
              <span className="text-2xl font-bold text-red-600 dark:text-red-400">
                {progress.failed}
              </span>
              <span className="text-xs text-muted-foreground">Failed</span>
            </div>
          </div>
        </div>

        {/* Estimated time remaining */}
        {progress.inProgress && estimatedTimeRemaining !== null && (
          <div className="flex items-center gap-2 rounded-lg bg-muted/50 px-3 py-2">
            <Clock className="size-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              Estimated time remaining: <span className="font-medium">{formatTime(estimatedTimeRemaining)}</span>
            </span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
