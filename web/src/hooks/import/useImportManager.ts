import { useState, useCallback, useRef } from 'react';
import type { UseMutationResult } from '@tanstack/react-query';

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
  retryRecord: (index: number) => Promise<void>;
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
      // Determine error type and create appropriate message
      let errorMessage = 'Failed to import record';
      
      if (error?.response) {
        // HTTP error response
        const status = error.response.status;
        const data = error.response.data;
        
        if (status === 400) {
          errorMessage = `Validation error: ${data?.message || 'Invalid data'}`;
        } else if (status === 401 || status === 403) {
          errorMessage = 'Authentication error: Please check your permissions';
        } else if (status === 404) {
          errorMessage = 'Resource not found';
        } else if (status === 409) {
          errorMessage = `Conflict: ${data?.message || 'Record already exists'}`;
        } else if (status === 429) {
          errorMessage = 'Rate limit exceeded: Too many requests';
        } else if (status >= 500) {
          errorMessage = `Server error (${status}): Please try again later`;
        } else {
          errorMessage = data?.message || `HTTP error ${status}`;
        }
      } else if (error?.request) {
        // Network error (request made but no response)
        errorMessage = 'Network error: Unable to reach server';
      } else if (error?.message) {
        // Other errors
        errorMessage = error.message;
      }
      
      // Retry with exponential backoff if we haven't exceeded max retries
      // Don't retry on client errors (4xx) except rate limiting
      const shouldRetry = retryCount < maxRetries && (
        !error?.response || 
        error.response.status >= 500 || 
        error.response.status === 429 ||
        !error.response.status
      );
      
      if (shouldRetry) {
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

  const retryRecord = useCallback(
    async (index: number) => {
      if (isImporting) {
        return;
      }

      const record = records[index];
      if (!record) {
        return;
      }

      // Reset the status to pending before retrying
      updateStatus(index, { status: 'pending' });

      // Import the single record
      const success = await importSingleRecord(index, record);

      // Update progress counts
      const currentStatus = importStatus.get(index);
      if (currentStatus?.status === 'error') {
        // If it was previously failed, decrement failed count
        setProgress(prev => ({
          ...prev,
          failed: Math.max(0, prev.failed - 1),
          completed: success ? prev.completed + 1 : prev.completed,
        }));
      }
    },
    [isImporting, records, importStatus, updateStatus]
  );

  return {
    importRecords,
    retryRecord,
    importStatus,
    progress,
    isImporting,
    cancelImport,
  };
}
