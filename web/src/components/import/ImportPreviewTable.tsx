import * as React from 'react';
import { AlertCircle, CheckCircle2, Loader2, XCircle } from 'lucide-react';
import { useVirtualizer } from '@tanstack/react-virtual';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Checkbox } from '@/components/ui/checkbox';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { Progress } from '@/components/ui/progress';
import { cn } from '@/lib/utils';
import type {
  ParsedImportRecord,
  ImportColumnDef,
  ImportStatus,
  ImportProgress,
} from '@/types/import';

export interface ImportPreviewTableProps<T> {
  records: ParsedImportRecord<T>[];
  columns: ImportColumnDef<T>[];
  selectedRecords: Set<number>;
  onSelectionChange: (indices: Set<number>) => void;
  onImportSelected: () => void;
  onImportAll: () => void;
  importStatus: Map<number, ImportStatus>;
  isImporting: boolean;
  progress?: ImportProgress;
  enableVirtualization?: boolean;
  virtualizationThreshold?: number;
}

export function ImportPreviewTable<T>({
  records,
  columns,
  selectedRecords,
  onSelectionChange,
  onImportSelected,
  onImportAll,
  importStatus,
  isImporting,
  progress,
  enableVirtualization = true,
  virtualizationThreshold = 100,
}: ImportPreviewTableProps<T>) {
  const parentRef = React.useRef<HTMLDivElement>(null);
  const allRecordsSelected = records.length > 0 && selectedRecords.size === records.length;
  const someRecordsSelected = selectedRecords.size > 0 && !allRecordsSelected;
  const validRecords = records.filter((r) => r.isValid);
  const selectedValidRecords = Array.from(selectedRecords).filter(
    (index) => records[index]?.isValid
  );

  // Enable virtualization only if records exceed threshold
  const shouldVirtualize = enableVirtualization && records.length > virtualizationThreshold;

  const rowVirtualizer = useVirtualizer({
    count: records.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 53, // Estimated row height in pixels
    overscan: 10, // Number of items to render outside visible area
    enabled: shouldVirtualize,
  });

  const handleSelectAll = () => {
    if (allRecordsSelected) {
      onSelectionChange(new Set());
    } else {
      onSelectionChange(new Set(records.map((_, index) => index)));
    }
  };

  const handleSelectRecord = (index: number) => {
    const newSelection = new Set(selectedRecords);
    if (newSelection.has(index)) {
      newSelection.delete(index);
    } else {
      newSelection.add(index);
    }
    onSelectionChange(newSelection);
  };

  const getRowClassName = (record: ParsedImportRecord<T>, index: number) => {
    const status = importStatus.get(index);
    
    if (status?.status === 'success') {
      return 'bg-green-50 dark:bg-green-950/20 hover:bg-green-100 dark:hover:bg-green-950/30';
    }
    if (status?.status === 'error') {
      return 'bg-red-50 dark:bg-red-950/20 hover:bg-red-100 dark:hover:bg-red-950/30';
    }
    if (status?.status === 'importing') {
      return 'bg-yellow-50 dark:bg-yellow-950/20 hover:bg-yellow-100 dark:hover:bg-yellow-950/30';
    }
    if (!record.isValid) {
      return 'bg-red-50/50 dark:bg-red-950/10 hover:bg-red-100/50 dark:hover:bg-red-950/20';
    }
    
    return '';
  };

  const renderStatusCell = (record: ParsedImportRecord<T>, index: number) => {
    const status = importStatus.get(index);

    if (status?.status === 'importing') {
      return (
        <Badge variant="outline" className="gap-1">
          <Loader2 className="size-3 animate-spin" />
          Importing
        </Badge>
      );
    }

    if (status?.status === 'success') {
      return (
        <Badge variant="default" className="gap-1 bg-green-600 hover:bg-green-700">
          <CheckCircle2 className="size-3" />
          Success
        </Badge>
      );
    }

    if (status?.status === 'error') {
      return (
        <Tooltip>
          <TooltipTrigger asChild>
            <Badge variant="destructive" className="gap-1 cursor-help">
              <XCircle className="size-3" />
              Failed
            </Badge>
          </TooltipTrigger>
          <TooltipContent>
            <p className="max-w-xs">{status.error || 'Import failed'}</p>
          </TooltipContent>
        </Tooltip>
      );
    }

    if (!record.isValid) {
      return (
        <Tooltip>
          <TooltipTrigger asChild>
            <Badge variant="destructive" className="gap-1 cursor-help">
              <AlertCircle className="size-3" />
              Invalid
            </Badge>
          </TooltipTrigger>
          <TooltipContent>
            <div className="max-w-xs space-y-1">
              {record.validationErrors.map((error, idx) => (
                <p key={idx} className="text-xs">
                  <span className="font-semibold">{error.field}:</span> {error.message}
                </p>
              ))}
            </div>
          </TooltipContent>
        </Tooltip>
      );
    }

    return (
      <Badge variant="outline" className="gap-1">
        Pending
      </Badge>
    );
  };

  if (records.length === 0) {
    return (
      <div className="text-center py-12 text-muted-foreground">
        <p>No records to display</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Action buttons and progress */}
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Button
            onClick={onImportSelected}
            disabled={selectedValidRecords.length === 0 || isImporting}
            size="sm"
          >
            Import Selected ({selectedValidRecords.length})
          </Button>
          <Button
            onClick={onImportAll}
            disabled={validRecords.length === 0 || isImporting}
            variant="outline"
            size="sm"
          >
            Import All Valid ({validRecords.length})
          </Button>
        </div>

        {progress && progress.inProgress && (
          <div className="flex items-center gap-3 text-sm text-muted-foreground">
            <span>
              {progress.completed} / {progress.total} completed
              {progress.failed > 0 && ` (${progress.failed} failed)`}
            </span>
          </div>
        )}
      </div>

      {/* Progress bar */}
      {progress && progress.inProgress && (
        <div className="space-y-2">
          <Progress value={(progress.completed / progress.total) * 100} />
          <p className="text-xs text-muted-foreground text-center">
            Importing records... {Math.round((progress.completed / progress.total) * 100)}%
          </p>
        </div>
      )}

      {/* Summary stats */}
      <div className="flex items-center gap-4 text-sm text-muted-foreground">
        <span>Total: {records.length}</span>
        <span className="text-green-600 dark:text-green-400">
          Valid: {validRecords.length}
        </span>
        {records.length - validRecords.length > 0 && (
          <span className="text-red-600 dark:text-red-400">
            Invalid: {records.length - validRecords.length}
          </span>
        )}
        {progress && progress.completed > 0 && (
          <>
            <span className="text-green-600 dark:text-green-400">
              Imported: {progress.completed - progress.failed}
            </span>
            {progress.failed > 0 && (
              <span className="text-red-600 dark:text-red-400">
                Failed: {progress.failed}
              </span>
            )}
          </>
        )}
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        {shouldVirtualize ? (
          <div
            ref={parentRef}
            className="overflow-auto"
            style={{ height: '600px', width: '100%' }}
          >
            <Table>
              <TableHeader className="sticky top-0 z-10 bg-background">
                <TableRow>
                  <TableHead className="w-12">
                    <Checkbox
                      checked={allRecordsSelected}
                      onCheckedChange={handleSelectAll}
                      aria-label="Select all records"
                      ref={(el) => {
                        if (el) {
                          (el as any).indeterminate = someRecordsSelected;
                        }
                      }}
                    />
                  </TableHead>
                  <TableHead className="w-24">Status</TableHead>
                  {columns.map((column) => (
                    <TableHead key={String(column.key)}>{column.label}</TableHead>
                  ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                <tr style={{ height: `${rowVirtualizer.getTotalSize()}px` }}>
                  <td />
                </tr>
                {rowVirtualizer.getVirtualItems().map((virtualRow) => {
                  const index = virtualRow.index;
                  const record = records[index];
                  return (
                    <TableRow
                      key={index}
                      className={cn(getRowClassName(record, index))}
                      data-state={selectedRecords.has(index) ? 'selected' : undefined}
                      style={{
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        width: '100%',
                        height: `${virtualRow.size}px`,
                        transform: `translateY(${virtualRow.start}px)`,
                      }}
                    >
                      <TableCell>
                        <Checkbox
                          checked={selectedRecords.has(index)}
                          onCheckedChange={() => handleSelectRecord(index)}
                          aria-label={`Select record ${index + 1}`}
                          disabled={!record.isValid}
                        />
                      </TableCell>
                      <TableCell>{renderStatusCell(record, index)}</TableCell>
                      {columns.map((column) => {
                        const value = record.data[column.key as keyof T];
                        return (
                          <TableCell key={String(column.key)}>
                            {column.render ? column.render(value, record) : String(value ?? '')}
                          </TableCell>
                        );
                      })}
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12">
                  <Checkbox
                    checked={allRecordsSelected}
                    onCheckedChange={handleSelectAll}
                    aria-label="Select all records"
                    ref={(el) => {
                      if (el) {
                        (el as any).indeterminate = someRecordsSelected;
                      }
                    }}
                  />
                </TableHead>
                <TableHead className="w-24">Status</TableHead>
                {columns.map((column) => (
                  <TableHead key={String(column.key)}>{column.label}</TableHead>
                ))}
              </TableRow>
            </TableHeader>
            <TableBody>
              {records.map((record, index) => (
                <TableRow
                  key={index}
                  className={cn(getRowClassName(record, index))}
                  data-state={selectedRecords.has(index) ? 'selected' : undefined}
                >
                  <TableCell>
                    <Checkbox
                      checked={selectedRecords.has(index)}
                      onCheckedChange={() => handleSelectRecord(index)}
                      aria-label={`Select record ${index + 1}`}
                      disabled={!record.isValid}
                    />
                  </TableCell>
                  <TableCell>{renderStatusCell(record, index)}</TableCell>
                  {columns.map((column) => {
                    const value = record.data[column.key as keyof T];
                    return (
                      <TableCell key={String(column.key)}>
                        {column.render ? column.render(value, record) : String(value ?? '')}
                      </TableCell>
                    );
                  })}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
}
