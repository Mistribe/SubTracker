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
  enableAdaptiveDelay?: boolean; // Adjust delay based on response times
  minDelay?: number; // Minimum delay between calls
  maxDelay?: number; // Maximum delay between calls
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
  // Exponential backoff with jitter to avoid thundering herd
  const exponentialDelay = baseDelay * Math.pow(2, retryCount);
  const jitter = Math.random() * 0.3 * exponentialDelay; // Add up to 30% jitter
  return Math.min(exponentialDelay + jitter, 10000); // Max 10 seconds
};

// Calculate adaptive delay based on recent response times
const calculateAdaptiveDelay = (
  recentResponseTimes: number[],
  baseDelay: number,
  minDelay: number,
  maxDelay: number
): number => {
  if (recentResponseTimes.length === 0) {
    return baseDelay;
  }

  // Calculate average response time from recent requests
  const avgResponseTime = recentResponseTimes.reduce((a, b) => a + b, 0) / recentResponseTimes.length;

  // If responses are fast (< 200ms), we can reduce delay
  // If responses are slow (> 1000ms), we should increase delay
  let adaptiveDelay = baseDelay;
  
  if (avgResponseTime < 200) {
    adaptiveDelay = Math.max(minDelay, baseDelay * 0.7);
  } else if (avgResponseTime > 1000) {
    adaptiveDelay = Math.min(maxDelay, baseDelay * 1.5);
  }

  return Math.round(adaptiveDelay);
};

export function useImportManager<T>({
  records,
  createMutation,
  delayBetweenCalls = 150,
  maxRetries = 3,
  enableAdaptiveDelay = true,
  minDelay = 50,
  maxDelay = 1000,
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
  const responseTimesRef = useRef<number[]>([]);
  const currentDelayRef = useRef(delayBetweenCalls);
  const rateLimitResetRef = useRef<number | null>(null);

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

    // Check if we need to wait for rate limit reset
    if (rateLimitResetRef.current && Date.now() < rateLimitResetRef.current) {
      const waitTime = rateLimitResetRef.current - Date.now();
      await sleep(waitTime);
      rateLimitResetRef.current = null;
    }

    updateStatus(index, { status: 'importing' });

    const startTime = Date.now();
    
    try {
      await createMutation.mutateAsync(record.data as T);
      
      // Track response time for adaptive delay
      const responseTime = Date.now() - startTime;
      if (enableAdaptiveDelay) {
        responseTimesRef.current.push(responseTime);
        // Keep only last 10 response times
        if (responseTimesRef.current.length > 10) {
          responseTimesRef.current.shift();
        }
        
        // Update current delay based on response times
        currentDelayRef.current = calculateAdaptiveDelay(
          responseTimesRef.current,
          delayBetweenCalls,
          minDelay,
          maxDelay
        );
      }
      
      updateStatus(index, { status: 'success' });
      return true;
    } catch (error: any) {
      // Determine error type and create appropriate message
      let errorMessage = 'Failed to import record';
      let isRateLimitError = false;
      
      if (error?.response) {
        // HTTP error response
        const status = error.response.status;
        const data = error.response.data;
        const headers = error.response.headers;
        
        if (status === 400) {
          // Check if this is a UUID-related validation error
          const message = data?.message || 'Invalid data';
          if (message.toLowerCase().includes('uuid') || message.toLowerCase().includes('id')) {
            // Extract UUID from record data if available (best-effort)
            const recordId = (record.data as any)?.id;
            errorMessage = recordId 
              ? `Invalid UUID format: ${recordId}. ${message}`
              : `Validation error: ${message}`;
          } else {
            errorMessage = `Validation error: ${message}`;
          }
        } else if (status === 401 || status === 403) {
          errorMessage = 'Authentication error: Please check your permissions';
        } else if (status === 404) {
          errorMessage = 'Resource not found';
        } else if (status === 409) {
          // Check if this is a UUID conflict error
          const message = data?.message || 'Record already exists';
          const recordId = (record.data as any)?.id;
          
          if (recordId && (message.toLowerCase().includes('uuid') || 
                          message.toLowerCase().includes('id') || 
                          message.toLowerCase().includes('already exists') ||
                          message.toLowerCase().includes('duplicate'))) {
            errorMessage = `UUID conflict: Entity with ID ${recordId} already exists`;
          } else {
            errorMessage = `Conflict: ${message}`;
          }
        } else if (status === 429) {
          isRateLimitError = true;
          errorMessage = 'Rate limit exceeded: Too many requests';
          
          // Check for Retry-After header
          const retryAfter = headers?.['retry-after'] || headers?.['x-ratelimit-reset'];
          if (retryAfter) {
            const resetTime = parseInt(retryAfter, 10);
            if (!isNaN(resetTime)) {
              // If it's a timestamp, use it directly; otherwise treat as seconds
              rateLimitResetRef.current = resetTime > 1000000000 
                ? resetTime * 1000 
                : Date.now() + (resetTime * 1000);
            }
          }
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
      // Specifically don't retry on 400 (validation/UUID errors) or 409 (UUID conflicts)
      const status = error?.response?.status;
      const shouldRetry = retryCount < maxRetries && (
        !error?.response || 
        error.response.status >= 500 || 
        isRateLimitError ||
        !error.response.status
      ) && status !== 400 && status !== 409;
      
      if (shouldRetry) {
        const backoffDelay = isRateLimitError && rateLimitResetRef.current
          ? Math.max(0, rateLimitResetRef.current - Date.now())
          : calculateBackoffDelay(retryCount, currentDelayRef.current);
        
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
        // Use adaptive delay if enabled
        if (index !== indices[indices.length - 1] && !cancelledRef.current) {
          const delay = enableAdaptiveDelay ? currentDelayRef.current : delayBetweenCalls;
          await sleep(delay);
        }
      }

      setIsImporting(false);
      updateProgress(completed, failed, total);
      
      // Reset adaptive delay tracking after import completes
      if (enableAdaptiveDelay) {
        responseTimesRef.current = [];
        currentDelayRef.current = delayBetweenCalls;
      }
    },
    [
      isImporting,
      records,
      delayBetweenCalls,
      maxRetries,
      enableAdaptiveDelay,
      minDelay,
      maxDelay,
      createMutation,
      updateProgress,
      updateStatus,
    ]
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
