import { useState, useCallback, useRef } from 'react';
import { UseMutationResult } from '@tanstack/react-query';

export interface ImportStatus {
  status: 'pending' | 'importing' | 'success' | 'error';
  error?: string;
}

export interface ImportProgress {
  total: number;
  completed: number;
  failed: number;
  inProgress: boolean;
}

export interface ParsedImportRecord<T> {
  index: number;
  data: Partial<T>;
  validationErrors: Array<{ field: string; message: string }>;
  isValid: boolean;
}

interface UseImportManagerProps<T> {
  records: ParsedImportRecord<T>[];
  createMutation: UseMutationResult<any, any, T, any>;
  delayBetweenCalls?: number; // milliseconds
  maxRetries?: number;
}

interface UseImportManagerReturn {
  importRecords: (indices: number[]) => Promise<void>;
  importStatus: Map<number, ImportStatus>;
  progress: ImportProgress;
  isImporting: boolean;
  cancelImport: () => void;
}

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

const calculateBackoffDelay = (retryCount: number, baseDelay: number): number => {
  return Math.min(baseDelay * Math.pow(2, retryCount), 10000); // Max 10 seconds
};

export function useImportManager<T>({
  records,
  createMutation,
  delayBetweenCalls = 150,
  maxRetries = 3,
}: UseImportManagerProps<T>): UseImportManagerReturn {
  const [importStatus, setImportStatus] = useState<Map<number, ImportStatus>>(new Map());
  const [progress, setProgress] = useState<ImportProgress>({
    total: 0,
    completed: 0,
    failed: 0,
    inProgress: false,
  });
  const [isImporting, setIsImporting] = useState(false);
  const cancelledRef = useRef(false);

  const updateStatus = useCallback((index: number, status: ImportStatus) => {
    setImportStatus(prev => {
      const newMap = new Map(prev);
      newMap.set(index, status);
      return newMap;
    });
  }, []);

  const updateProgress = useCallback((completed: number, failed: number, total: number) => {
    setProgress({
      total,
      completed,
      failed,
      inProgress: completed + failed < total,
    });
  }, []);

  const importSingleRecord = async (
    index: number,
    record: ParsedImportRecord<T>,
    retryCount = 0
  ): Promise<boolean> => {
    if (cancelledRef.current) {
      return false;
    }

    // Skip invalid records
    if (!record.isValid) {
      updateStatus(index, {
        status: 'error',
        error: 'Record has validation errors',
      });
      return false;
    }

    updateStatus(index, { status: 'importing' });

    try {
      await createMutation.mutateAsync(record.data as T);
      updateStatus(index, { status: 'success' });
      return true;
    } catch (error: any) {
      const errorMessage = error?.message || 'Failed to import record';
      
      // Retry with exponential backoff if we haven't exceeded max retries
      if (retryCount < maxRetries) {
        const backoffDelay = calculateBackoffDelay(retryCount, delayBetweenCalls);
        await sleep(backoffDelay);
        return importSingleRecord(index, record, retryCount + 1);
      }

      updateStatus(index, {
        status: 'error',
        error: errorMessage,
      });
      return false;
    }
  };

  const importRecords = useCallback(
    async (indices: number[]) => {
      if (isImporting) {
        return;
      }

      // Reset cancellation flag
      cancelledRef.current = false;
      setIsImporting(true);

      // Initialize status for all records
      const initialStatusMap = new Map<number, ImportStatus>();
      indices.forEach(index => {
        initialStatusMap.set(index, { status: 'pending' });
      });
      setImportStatus(initialStatusMap);

      // Initialize progress
      const total = indices.length;
      let completed = 0;
      let failed = 0;

      updateProgress(completed, failed, total);

      // Process records sequentially
      for (const index of indices) {
        if (cancelledRef.current) {
          // Mark remaining records as pending
          break;
        }

        const record = records[index];
        if (!record) {
          failed++;
          updateProgress(completed, failed, total);
          continue;
        }

        const success = await importSingleRecord(index, record);

        if (success) {
          completed++;
        } else {
          failed++;
        }

        updateProgress(completed, failed, total);

        // Add delay between calls (except for the last one)
        if (index !== indices[indices.length - 1] && !cancelledRef.current) {
          await sleep(delayBetweenCalls);
        }
      }

      setIsImporting(false);
      updateProgress(completed, failed, total);
    },
    [isImporting, records, delayBetweenCalls, maxRetries, createMutation, updateProgress, updateStatus]
  );

  const cancelImport = useCallback(() => {
    cancelledRef.current = true;
  }, []);

  return {
    importRecords,
    importStatus,
    progress,
    isImporting,
    cancelImport,
  };
}
