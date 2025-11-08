import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useImportManager } from '../useImportManager';
import type { ParsedImportRecord } from '../useImportManager';
import type { UseMutationResult } from '@tanstack/react-query';

describe('useImportManager - Integration Tests', () => {
  const createMockMutation = (
    mutateAsyncFn: (data: any) => Promise<any>
  ): UseMutationResult<any, any, any, any> => {
    return {
      mutateAsync: mutateAsyncFn,
      mutate: vi.fn(),
      reset: vi.fn(),
      isIdle: true,
      isPending: false,
      isError: false,
      isSuccess: false,
      data: undefined,
      error: null,
      variables: undefined,
      status: 'idle',
      failureCount: 0,
      failureReason: null,
      isPaused: false,
      submittedAt: 0,
      context: undefined,
    } as UseMutationResult<any, any, any, any>;
  };

  const createMockRecords = (count: number, allValid = true): ParsedImportRecord<any>[] => {
    return Array.from({ length: count }, (_, i) => ({
      index: i,
      data: { name: `Record ${i}`, color: `#${i}${i}${i}${i}${i}${i}` },
      validationErrors: allValid ? [] : [{ field: 'name', message: 'Invalid name' }],
      isValid: allValid,
    }));
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Single record import', () => {
    it('should successfully import a single record', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        return { id: '123', ...data };
      });

      const records = createMockRecords(3);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      // Import only the second record
      await result.current.importRecords([1]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(1);
      expect(mutateAsync).toHaveBeenCalledWith(records[1].data);
      expect(result.current.importStatus.get(1)?.status).toBe('success');
      expect(result.current.progress.total).toBe(1);
      expect(result.current.progress.failed).toBe(0);
    });

    it('should handle single record import failure', async () => {
      const mutateAsync = vi.fn(async () => {
        throw new Error('API Error: Record already exists');
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 0,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.failed).toBe(1);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toContain('API Error');
    });

    it('should update status from pending to importing to success', async () => {
      const mutateAsync = vi.fn(async () => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return { id: 'test' };
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      // Start import
      const importPromise = result.current.importRecords([0]);

      // Wait for pending status
      await waitFor(() => {
        return result.current.importStatus.get(0)?.status === 'pending';
      });

      // Wait for importing status
      await waitFor(() => {
        return result.current.importStatus.get(0)?.status === 'importing';
      });

      await importPromise;

      // Verify final success status
      await waitFor(() => {
        expect(result.current.importStatus.get(0)?.status).toBe('success');
      });
    });
  });

  describe('Bulk import', () => {
    it('should successfully import all records in bulk', async () => {
      const importedRecords: any[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        importedRecords.push(data);
        return { id: `id-${data.name}`, ...data };
      });

      const records = createMockRecords(5);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1, 2, 3, 4]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(5);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(5);
      expect(importedRecords).toHaveLength(5);
      expect(result.current.progress.total).toBe(5);
      expect(result.current.progress.failed).toBe(0);
      
      // Verify all records have success status
      for (let i = 0; i < 5; i++) {
        expect(result.current.importStatus.get(i)?.status).toBe('success');
      }
    });

    it('should handle partial failures in bulk import', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        // Fail records at indices 1 and 3
        if (index === 1 || index === 3) {
          throw new Error(`Failed to import record ${index}`);
        }
        return { id: `id-${index}`, ...data };
      });

      const records = createMockRecords(5);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 0,
        })
      );

      await result.current.importRecords([0, 1, 2, 3, 4]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
        expect(result.current.progress.failed).toBe(2);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('success');
      expect(result.current.importStatus.get(1)?.status).toBe('error');
      expect(result.current.importStatus.get(2)?.status).toBe('success');
      expect(result.current.importStatus.get(3)?.status).toBe('error');
      expect(result.current.importStatus.get(4)?.status).toBe('success');
    });

    it('should track progress during bulk import', async () => {
      const progressSnapshots: Array<{ completed: number; failed: number; total: number }> = [];
      const mutateAsync = vi.fn(async (data: any) => {
        await new Promise(resolve => setTimeout(resolve, 20));
        return { id: 'test', ...data };
      });

      const records = createMockRecords(3);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      const importPromise = result.current.importRecords([0, 1, 2]);

      // Capture progress snapshots
      const captureProgress = async () => {
        for (let i = 0; i < 10; i++) {
          await new Promise(resolve => setTimeout(resolve, 15));
          progressSnapshots.push({ ...result.current.progress });
        }
      };

      await Promise.all([importPromise, captureProgress()]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
      });

      // Verify progress increased over time
      const completedCounts = progressSnapshots.map(s => s.completed);
      const uniqueCounts = [...new Set(completedCounts)];
      expect(uniqueCounts.length).toBeGreaterThan(1); // Progress should change
      expect(Math.max(...completedCounts)).toBe(3);
    });

    it('should handle bulk import with mixed valid and invalid records', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        return { id: 'test', ...data };
      });

      const records = [
        { index: 0, data: { name: 'Valid 1' }, validationErrors: [], isValid: true },
        { index: 1, data: { name: 'Invalid' }, validationErrors: [{ field: 'name', message: 'Required' }], isValid: false },
        { index: 2, data: { name: 'Valid 2' }, validationErrors: [], isValid: true },
        { index: 3, data: { name: 'Invalid 2' }, validationErrors: [{ field: 'color', message: 'Invalid format' }], isValid: false },
        { index: 4, data: { name: 'Valid 3' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1, 2, 3, 4]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
        expect(result.current.progress.failed).toBe(2);
      });

      // Only valid records should be sent to API
      expect(mutateAsync).toHaveBeenCalledTimes(3);
      expect(result.current.importStatus.get(1)?.error).toBe('Record has validation errors');
      expect(result.current.importStatus.get(3)?.error).toBe('Record has validation errors');
    });
  });

  describe('Import cancellation', () => {
    it('should cancel bulk import mid-process', async () => {
      let callCount = 0;
      const mutateAsync = vi.fn(async () => {
        callCount++;
        await new Promise(resolve => setTimeout(resolve, 30));
        return { id: `test-${callCount}` };
      });

      const records = createMockRecords(10);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      const importPromise = result.current.importRecords([0, 1, 2, 3, 4, 5, 6, 7, 8, 9]);

      // Cancel after a short delay (should cancel after 2-3 records)
      setTimeout(() => {
        result.current.cancelImport();
      }, 80);

      await importPromise;

      await waitFor(() => {
        expect(result.current.isImporting).toBe(false);
      });

      // Should have imported fewer than all records
      const totalProcessed = result.current.progress.completed + result.current.progress.failed;
      expect(totalProcessed).toBeLessThan(10);
      expect(totalProcessed).toBeGreaterThan(0);
    });

    it('should not process remaining records after cancellation', async () => {
      const processedIndices: number[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        processedIndices.push(index);
        await new Promise(resolve => setTimeout(resolve, 25));
        return { id: `test-${index}` };
      });

      const records = createMockRecords(8);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      const importPromise = result.current.importRecords([0, 1, 2, 3, 4, 5, 6, 7]);

      // Cancel after processing 2 records
      setTimeout(() => {
        result.current.cancelImport();
      }, 70);

      await importPromise;

      await waitFor(() => {
        expect(result.current.isImporting).toBe(false);
      });

      expect(processedIndices.length).toBeLessThan(8);
      // Verify indices are sequential (no skipping)
      for (let i = 1; i < processedIndices.length; i++) {
        expect(processedIndices[i]).toBe(processedIndices[i - 1] + 1);
      }
    });

    it('should allow new import after cancellation', async () => {
      const mutateAsync = vi.fn(async () => {
        await new Promise(resolve => setTimeout(resolve, 30));
        return { id: 'test' };
      });

      const records = createMockRecords(5);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      // First import - will be cancelled
      const firstImport = result.current.importRecords([0, 1, 2, 3, 4]);
      setTimeout(() => result.current.cancelImport(), 50);
      await firstImport;

      await waitFor(() => {
        expect(result.current.isImporting).toBe(false);
      });

      const firstImportCount = result.current.progress.completed + result.current.progress.failed;

      // Second import - should work normally
      await result.current.importRecords([0, 1]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(2);
      });

      expect(result.current.progress.total).toBe(2);
      expect(result.current.isImporting).toBe(false);
    });
  });

  describe('Error recovery and retry', () => {
    it('should retry individual failed record', async () => {
      let attemptCount = 0;
      const mutateAsync = vi.fn(async () => {
        attemptCount++;
        if (attemptCount === 1) {
          throw new Error('First attempt failed');
        }
        return { id: 'test' };
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 0, // No automatic retries
        })
      );

      // First import attempt - will fail
      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.importStatus.get(0)?.status).toBe('error');
      });

      expect(result.current.progress.failed).toBe(1);

      // Retry the failed record
      await result.current.retryRecord(0);

      await waitFor(() => {
        expect(result.current.importStatus.get(0)?.status).toBe('success');
      });

      expect(attemptCount).toBe(2);
      expect(result.current.progress.completed).toBe(1);
    });

    it('should handle automatic retry with exponential backoff', async () => {
      const timestamps: number[] = [];
      let attemptCount = 0;

      const mutateAsync = vi.fn(async () => {
        attemptCount++;
        timestamps.push(Date.now());
        
        if (attemptCount <= 2) {
          throw new Error('Temporary failure');
        }
        return { id: 'test' };
      });

      const records = createMockRecords(1);
      const baseDelay = 50;
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: baseDelay,
          maxRetries: 3,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.importStatus.get(0)?.status).toBe('success');
      });

      expect(attemptCount).toBe(3);
      
      // Verify exponential backoff (each retry should take longer)
      if (timestamps.length >= 3) {
        const delay1 = timestamps[1] - timestamps[0];
        const delay2 = timestamps[2] - timestamps[1];
        expect(delay2).toBeGreaterThan(delay1);
      }
    });

    it('should recover from network errors and continue', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        
        if (index === 1) {
          const error: any = new Error('Network error');
          error.request = {}; // Simulate network error
          throw error;
        }
        
        return { id: `test-${index}` };
      });

      const records = createMockRecords(4);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 1,
        })
      );

      await result.current.importRecords([0, 1, 2, 3]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
        expect(result.current.progress.failed).toBe(1);
      });

      expect(result.current.importStatus.get(1)?.error).toContain('Network error');
      expect(result.current.importStatus.get(0)?.status).toBe('success');
      expect(result.current.importStatus.get(2)?.status).toBe('success');
      expect(result.current.importStatus.get(3)?.status).toBe('success');
    });

    it('should handle different HTTP error codes appropriately', async () => {
      const errorCodes = [400, 401, 404, 409, 429, 500];
      const results: Array<{ code: number; error: string; retries: number }> = [];

      for (const code of errorCodes) {
        let callCount = 0;
        const mutateAsync = vi.fn(async () => {
          callCount++;
          const error: any = new Error(`HTTP ${code}`);
          error.response = {
            status: code,
            data: { message: `Error ${code}` },
          };
          throw error;
        });

        const records = createMockRecords(1);
        const { result } = renderHook(() =>
          useImportManager({
            records,
            createMutation: createMockMutation(mutateAsync),
            delayBetweenCalls: 10,
            maxRetries: 2,
          })
        );

        await result.current.importRecords([0]);

        await waitFor(() => {
          expect(result.current.importStatus.get(0)?.status).toBe('error');
        });

        results.push({
          code,
          error: result.current.importStatus.get(0)?.error || '',
          retries: callCount - 1,
        });
      }

      // 4xx errors (except 429) should not retry
      expect(results.find(r => r.code === 400)?.retries).toBe(0);
      expect(results.find(r => r.code === 401)?.retries).toBe(0);
      expect(results.find(r => r.code === 404)?.retries).toBe(0);
      expect(results.find(r => r.code === 409)?.retries).toBe(0);

      // 429 and 5xx should retry
      expect(results.find(r => r.code === 429)?.retries).toBe(2);
      expect(results.find(r => r.code === 500)?.retries).toBe(2);

      // Verify error messages are descriptive
      expect(results.find(r => r.code === 400)?.error).toContain('Validation error');
      expect(results.find(r => r.code === 401)?.error).toContain('Authentication error');
      expect(results.find(r => r.code === 429)?.error).toContain('Rate limit');
      expect(results.find(r => r.code === 500)?.error).toContain('Server error');
    });

    it('should allow selective retry of multiple failed records', async () => {
      let failedIndices = new Set([1, 3]);
      
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        if (failedIndices.has(index)) {
          failedIndices.delete(index); // Succeed on retry
          throw new Error(`Failed ${index}`);
        }
        return { id: `test-${index}` };
      });

      const records = createMockRecords(5);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 0,
        })
      );

      // Initial import
      await result.current.importRecords([0, 1, 2, 3, 4]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
        expect(result.current.progress.failed).toBe(2);
      });

      // Retry failed records
      await result.current.retryRecord(1);
      await result.current.retryRecord(3);

      await waitFor(() => {
        expect(result.current.importStatus.get(1)?.status).toBe('success');
        expect(result.current.importStatus.get(3)?.status).toBe('success');
      });

      expect(result.current.progress.completed).toBe(5);
    });
  });
});
