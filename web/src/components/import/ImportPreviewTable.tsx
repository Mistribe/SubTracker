import * as React from 'react';
import { AlertCircle, CheckCircle2, ChevronDown, Loader2, XCircle } from 'lucide-react';
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
  onRetryRecord?: (index: number) => void;
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
  onRetryRecord,
}: ImportPreviewTableProps<T>) {
  const parentRef = React.useRef<HTMLDivElement>(null);
  const [expandedRows, setExpandedRows] = React.useState<Set<number>>(new Set());
  const [focusedRowIndex, setFocusedRowIndex] = React.useState<number>(-1);
  const allRecordsSelected = records.length > 0 && selectedRecords.size === records.length;
  const someRecordsSelected = selectedRecords.size > 0 && !allRecordsSelected;
  const validRecords = records.filter((r) => r.isValid);
  const selectedValidRecords = Array.from(selectedRecords).filter(
    (index) => records[index]?.isValid
  );

  const toggleRowExpansion = (index: number) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(index)) {
      newExpanded.delete(index);
    } else {
      newExpanded.add(index);
    }
    setExpandedRows(newExpanded);
  };

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

  // Keyboard navigation handler
  const handleTableKeyDown = (e: React.KeyboardEvent, rowIndex: number) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        if (rowIndex < records.length - 1) {
          setFocusedRowIndex(rowIndex + 1);
          // Focus the next row's checkbox
          setTimeout(() => {
            const nextCheckbox = document.querySelector(
              `[data-row-index="${rowIndex + 1}"] input[type="checkbox"]`
            ) as HTMLElement;
            nextCheckbox?.focus();
          }, 0);
        }
        break;
      case 'ArrowUp':
        e.preventDefault();
        if (rowIndex > 0) {
          setFocusedRowIndex(rowIndex - 1);
          // Focus the previous row's checkbox
          setTimeout(() => {
            const prevCheckbox = document.querySelector(
              `[data-row-index="${rowIndex - 1}"] input[type="checkbox"]`
            ) as HTMLElement;
            prevCheckbox?.focus();
          }, 0);
        }
        break;
      case ' ':
        // Space bar toggles selection
        if (records[rowIndex]?.isValid) {
          e.preventDefault();
          handleSelectRecord(rowIndex);
        }
        break;
      case 'Enter':
        // Enter expands error details if invalid
        if (!records[rowIndex]?.isValid && records[rowIndex]?.validationErrors.length > 0) {
          e.preventDefault();
          toggleRowExpansion(rowIndex);
        }
        break;
    }
  };

  const getRowClassName = (record: ParsedImportRecord<T>, index: number) => {
    const status = importStatus.get(index);
    
    if (status?.status === 'success') {
      return 'bg-green-100 dark:bg-green-950/30 hover:bg-green-200 dark:hover:bg-green-950/40 border-l-4 border-l-green-600 dark:border-l-green-400';
    }
    if (status?.status === 'error') {
      return 'bg-red-100 dark:bg-red-950/30 hover:bg-red-200 dark:hover:bg-red-950/40 border-l-4 border-l-red-600 dark:border-l-red-400';
    }
    if (status?.status === 'importing') {
      return 'bg-yellow-100 dark:bg-yellow-950/30 hover:bg-yellow-200 dark:hover:bg-yellow-950/40 border-l-4 border-l-yellow-600 dark:border-l-yellow-400';
    }
    if (!record.isValid) {
      return 'bg-red-50 dark:bg-red-950/20 hover:bg-red-100 dark:hover:bg-red-950/30 border-l-4 border-l-red-400 dark:border-l-red-500';
    }
    
    return 'hover:bg-accent/50';
  };

  const renderStatusCell = (record: ParsedImportRecord<T>, index: number) => {
    const status = importStatus.get(index);

    if (status?.status === 'importing') {
      return (
        <Badge variant="outline" className="gap-1 border-yellow-600 dark:border-yellow-400 bg-yellow-50 dark:bg-yellow-950/30 text-yellow-900 dark:text-yellow-100 font-semibold">
          <Loader2 className="size-3 animate-spin" aria-hidden="true" />
          Importing
        </Badge>
      );
    }

    if (status?.status === 'success') {
      return (
        <Badge variant="default" className="gap-1 bg-green-700 hover:bg-green-800 dark:bg-green-600 dark:hover:bg-green-700 font-semibold">
          <CheckCircle2 className="size-3" aria-hidden="true" />
          Success
        </Badge>
      );
    }

    if (status?.status === 'error') {
      return (
        <div className="flex items-center gap-2">
          <Tooltip>
            <TooltipTrigger asChild>
              <Badge variant="destructive" className="gap-1 cursor-help bg-red-700 hover:bg-red-800 dark:bg-red-600 dark:hover:bg-red-700 font-semibold focus-visible:ring-4 focus-visible:ring-red-500 focus-visible:ring-offset-2">
                <XCircle className="size-3" aria-hidden="true" />
                Failed
              </Badge>
            </TooltipTrigger>
            <TooltipContent className="bg-red-900 dark:bg-red-950 text-white border-red-700 dark:border-red-800">
              <p className="max-w-xs font-medium">{status.error || 'Import failed'}</p>
            </TooltipContent>
          </Tooltip>
          {onRetryRecord && !isImporting && (
            <Button
              size="sm"
              variant="ghost"
              className="h-6 px-2 text-xs focus-visible:ring-4 focus-visible:ring-primary focus-visible:ring-offset-2"
              onClick={(e) => {
                e.stopPropagation();
                onRetryRecord(index);
              }}
              aria-label={`Retry importing record ${index + 1}`}
            >
              Retry
            </Button>
          )}
        </div>
      );
    }

    if (!record.isValid) {
      return (
        <Tooltip>
          <TooltipTrigger asChild>
            <Badge variant="destructive" className="gap-1 cursor-help bg-red-700 hover:bg-red-800 dark:bg-red-600 dark:hover:bg-red-700 font-semibold focus-visible:ring-4 focus-visible:ring-red-500 focus-visible:ring-offset-2">
              <AlertCircle className="size-3" aria-hidden="true" />
              Invalid
            </Badge>
          </TooltipTrigger>
          <TooltipContent className="bg-red-900 dark:bg-red-950 text-white border-red-700 dark:border-red-800">
            <div className="max-w-xs space-y-1">
              {record.validationErrors.map((error, idx) => (
                <p key={idx} className="text-xs font-medium">
                  <span className="font-bold">{error.field}:</span> {error.message}
                </p>
              ))}
            </div>
          </TooltipContent>
        </Tooltip>
      );
    }

    return (
      <Badge variant="outline" className="gap-1 font-semibold">
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
      {/* Screen reader announcements */}
      <div className="sr-only" role="status" aria-live="polite" aria-atomic="true">
        {progress && progress.inProgress && (
          <span>
            Importing records: {progress.completed} of {progress.total} completed
            {progress.failed > 0 && `, ${progress.failed} failed`}
          </span>
        )}
        {progress && !progress.inProgress && progress.completed > 0 && (
          <span>
            Import complete: {progress.completed - progress.failed} successful, {progress.failed} failed
          </span>
        )}
      </div>

      {/* Action buttons and progress */}
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Button
            onClick={onImportSelected}
            disabled={selectedValidRecords.length === 0 || isImporting}
            size="sm"
            aria-label={`Import ${selectedValidRecords.length} selected record${selectedValidRecords.length !== 1 ? 's' : ''}`}
            className="focus-visible:ring-4 focus-visible:ring-primary focus-visible:ring-offset-2"
          >
            {isImporting && <Loader2 className="size-4 mr-2 animate-spin" aria-hidden="true" />}
            Import Selected ({selectedValidRecords.length})
          </Button>
          <Button
            onClick={onImportAll}
            disabled={validRecords.length === 0 || isImporting}
            variant="outline"
            size="sm"
            aria-label={`Import all ${validRecords.length} valid record${validRecords.length !== 1 ? 's' : ''}`}
            className="focus-visible:ring-4 focus-visible:ring-primary focus-visible:ring-offset-2"
          >
            {isImporting && <Loader2 className="size-4 mr-2 animate-spin" aria-hidden="true" />}
            Import All Valid ({validRecords.length})
          </Button>
          {progress && progress.failed > 0 && onRetryRecord && !isImporting && (
            <Button
              onClick={() => {
                // Retry all failed records
                const failedIndices = Array.from(importStatus.entries())
                  .filter(([_, status]) => status.status === 'error')
                  .map(([index]) => index);
                failedIndices.forEach(index => onRetryRecord(index));
              }}
              variant="outline"
              size="sm"
              className="text-orange-600 hover:text-orange-700 dark:text-orange-400 dark:hover:text-orange-300 focus-visible:ring-4 focus-visible:ring-orange-500 focus-visible:ring-offset-2"
              aria-label={`Retry all ${progress.failed} failed record${progress.failed !== 1 ? 's' : ''}`}
            >
              Retry All Failed ({progress.failed})
            </Button>
          )}
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
        <div className="space-y-2" role="region" aria-label="Import progress">
          <Progress 
            value={(progress.completed / progress.total) * 100}
            aria-label={`Import progress: ${Math.round((progress.completed / progress.total) * 100)}%`}
          />
          <p className="text-xs text-muted-foreground text-center" aria-live="polite">
            Importing records... {Math.round((progress.completed / progress.total) * 100)}%
          </p>
        </div>
      )}

      {/* Summary stats */}
      <div 
        className="flex items-center gap-4 text-sm text-muted-foreground"
        role="region"
        aria-label="Import statistics"
      >
        <span aria-label={`Total records: ${records.length}`}>Total: {records.length}</span>
        <span 
          className="text-green-600 dark:text-green-400"
          aria-label={`Valid records: ${validRecords.length}`}
        >
          Valid: {validRecords.length}
        </span>
        {records.length - validRecords.length > 0 && (
          <span 
            className="text-red-600 dark:text-red-400"
            aria-label={`Invalid records: ${records.length - validRecords.length}`}
          >
            Invalid: {records.length - validRecords.length}
          </span>
        )}
        {progress && progress.completed > 0 && (
          <>
            <span 
              className="text-green-600 dark:text-green-400"
              aria-label={`Successfully imported: ${progress.completed - progress.failed}`}
            >
              Imported: {progress.completed - progress.failed}
            </span>
            {progress.failed > 0 && (
              <span 
                className="text-red-600 dark:text-red-400"
                aria-label={`Failed imports: ${progress.failed}`}
              >
                Failed: {progress.failed}
              </span>
            )}
          </>
        )}
      </div>

      {/* Validation error summary */}
      {records.length - validRecords.length > 0 && (
        <div 
          className="rounded-lg border-2 border-red-300 dark:border-red-800 bg-red-50 dark:bg-red-950/30 p-4"
          role="alert"
          aria-live="polite"
        >
          <div className="flex items-start gap-3">
            <AlertCircle className="size-5 text-red-700 dark:text-red-300 flex-shrink-0 mt-0.5" aria-hidden="true" />
            <div className="flex-1 space-y-2">
              <h4 className="text-sm font-bold text-red-900 dark:text-red-100">
                Validation Errors Found
              </h4>
              <p className="text-sm font-medium text-red-900 dark:text-red-100">
                {records.length - validRecords.length} record(s) have validation errors and cannot be imported. 
                Hover over the "Invalid" badge in the Status column to see specific errors for each record.
              </p>
              <p className="text-xs font-medium text-red-800 dark:text-red-200">
                Common issues: missing required fields, invalid formats, or incorrect data types.
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Table */}
      <div className="border rounded-lg">
        {shouldVirtualize ? (
          <div
            ref={parentRef}
            className="overflow-auto"
            style={{ height: '600px', width: '100%' }}
            role="region"
            aria-label="Import preview table with virtualization"
          >
            <Table aria-label="Records to import">
              <TableHeader className="sticky top-0 z-10 bg-background">
                <TableRow>
                  <TableHead className="w-12">
                    <Checkbox
                      checked={allRecordsSelected}
                      onCheckedChange={handleSelectAll}
                      onKeyDown={(e) => {
                        if (e.key === ' ' || e.key === 'Enter') {
                          e.preventDefault();
                          handleSelectAll();
                        }
                      }}
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
                      data-row-index={index}
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
                          onKeyDown={(e) => handleTableKeyDown(e, index)}
                          aria-label={`Select record ${index + 1}${!record.isValid ? ' (invalid)' : ''}`}
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
          <Table aria-label="Records to import">
            <TableHeader>
              <TableRow>
                <TableHead className="w-12">
                  <Checkbox
                    checked={allRecordsSelected}
                    onCheckedChange={handleSelectAll}
                    onKeyDown={(e) => {
                      if (e.key === ' ' || e.key === 'Enter') {
                        e.preventDefault();
                        handleSelectAll();
                      }
                    }}
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
                <React.Fragment key={index}>
                  <TableRow
                    className={cn(getRowClassName(record, index))}
                    data-state={selectedRecords.has(index) ? 'selected' : undefined}
                    data-row-index={index}
                  >
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Checkbox
                          checked={selectedRecords.has(index)}
                          onCheckedChange={() => handleSelectRecord(index)}
                          onKeyDown={(e) => handleTableKeyDown(e, index)}
                          aria-label={`Select record ${index + 1}${!record.isValid ? ' (invalid)' : ''}`}
                          disabled={!record.isValid}
                        />
                        {!record.isValid && record.validationErrors.length > 0 && (
                          <button
                            onClick={() => toggleRowExpansion(index)}
                            onKeyDown={(e) => {
                              if (e.key === 'Enter' || e.key === ' ') {
                                e.preventDefault();
                                toggleRowExpansion(index);
                              }
                            }}
                            className="p-1 hover:bg-accent rounded transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                            aria-label={expandedRows.has(index) ? 'Collapse errors' : 'Expand errors'}
                            aria-expanded={expandedRows.has(index)}
                          >
                            <ChevronDown
                              className={cn(
                                'size-4 text-muted-foreground transition-transform',
                                expandedRows.has(index) && 'rotate-180'
                              )}
                              aria-hidden="true"
                            />
                          </button>
                        )}
                      </div>
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
                  {/* Expandable error details row */}
                  {!record.isValid && expandedRows.has(index) && (
                    <TableRow className={cn(getRowClassName(record, index))}>
                      <TableCell colSpan={columns.length + 2} className="p-0">
                        <div className="px-4 py-3 bg-red-100 dark:bg-red-950/30 border-t-2 border-red-300 dark:border-red-800">
                          <div className="space-y-2">
                            <h5 className="text-xs font-bold text-red-900 dark:text-red-100 flex items-center gap-2">
                              <AlertCircle className="size-3" aria-hidden="true" />
                              Validation Errors
                            </h5>
                            <ul className="space-y-1 text-xs text-red-900 dark:text-red-100" role="list">
                              {record.validationErrors.map((error, errorIdx) => (
                                <li key={errorIdx} className="flex items-start gap-2" role="listitem">
                                  <span className="font-bold min-w-[100px]">{error.field}:</span>
                                  <span className="font-medium">{error.message}</span>
                                </li>
                              ))}
                            </ul>
                          </div>
                        </div>
                      </TableCell>
                    </TableRow>
                  )}
                </React.Fragment>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
}
