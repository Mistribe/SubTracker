import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useImportManager } from '../useImportManager';
import type { ParsedImportRecord } from '../useImportManager';
import type { UseMutationResult } from '@tanstack/react-query';

describe('useImportManager', () => {
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
      data: { name: `Record ${i}` },
      validationErrors: allValid ? [] : [{ field: 'name', message: 'Invalid name' }],
      isValid: allValid,
    }));
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Sequential import logic', () => {
    it('should import records sequentially', async () => {
      const callOrder: number[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        callOrder.push(parseInt(data.name.split(' ')[1]));
        return { id: data.name };
      });

      const records = createMockRecords(3);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1, 2]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
      });

      expect(callOrder).toEqual([0, 1, 2]);
      expect(mutateAsync).toHaveBeenCalledTimes(3);
    });

    it('should add delay between API calls', async () => {
      const timestamps: number[] = [];
      const mutateAsync = vi.fn(async () => {
        timestamps.push(Date.now());
        return { id: 'test' };
      });

      const records = createMockRecords(3);
      const delayBetweenCalls = 50;
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls,
        })
      );

      await result.current.importRecords([0, 1, 2]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
      });

      // Check that there's at least the delay between calls
      for (let i = 1; i < timestamps.length; i++) {
        const timeDiff = timestamps[i] - timestamps[i - 1];
        expect(timeDiff).toBeGreaterThanOrEqual(delayBetweenCalls - 10); // Allow 10ms tolerance
      }
    });

    it('should not make concurrent requests', async () => {
      let activeRequests = 0;
      let maxConcurrent = 0;

      const mutateAsync = vi.fn(async () => {
        activeRequests++;
        maxConcurrent = Math.max(maxConcurrent, activeRequests);
        await new Promise(resolve => setTimeout(resolve, 20));
        activeRequests--;
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

      await result.current.importRecords([0, 1, 2, 3, 4]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(5);
      });

      expect(maxConcurrent).toBe(1);
    });
  });

  describe('Error handling and continuation', () => {
    it('should continue importing after a failed record', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        if (index === 1) {
          throw new Error('API Error');
        }
        return { id: data.name };
      });

      const records = createMockRecords(3);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 3,
        })
      );

      await result.current.importRecords([0, 1, 2]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(2);
        expect(result.current.progress.failed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(6); // 3 records: 1 + 4 (1 initial + 3 retries) + 1
      expect(result.current.importStatus.get(0)?.status).toBe('success');
      expect(result.current.importStatus.get(1)?.status).toBe('error');
      expect(result.current.importStatus.get(2)?.status).toBe('success');
    });

    it('should retry failed records with exponential backoff', async () => {
      let attemptCount = 0;
      const mutateAsync = vi.fn(async () => {
        attemptCount++;
        if (attemptCount <= 2) {
          throw new Error('Temporary error');
        }
        return { id: 'test' };
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 3,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(1);
      });

      expect(attemptCount).toBe(3);
      expect(result.current.importStatus.get(0)?.status).toBe('success');
    });

    it('should mark record as error after max retries', async () => {
      const mutateAsync = vi.fn(async () => {
        throw new Error('Persistent error');
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
        expect(result.current.progress.failed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(3); // Initial + 2 retries
      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toBe('Persistent error');
    });

    it('should skip invalid records', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        return { id: data.name };
      });

      const records = [
        { index: 0, data: { name: 'Valid' }, validationErrors: [], isValid: true },
        { index: 1, data: { name: 'Invalid' }, validationErrors: [{ field: 'name', message: 'Error' }], isValid: false },
        { index: 2, data: { name: 'Valid' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1, 2]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(2);
        expect(result.current.progress.failed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(2);
      expect(result.current.importStatus.get(1)?.status).toBe('error');
      expect(result.current.importStatus.get(1)?.error).toBe('Record has validation errors');
    });
  });

  describe('Cancellation mechanism', () => {
    it('should cancel import when cancelImport is called', async () => {
      const mutateAsync = vi.fn(async () => {
        await new Promise(resolve => setTimeout(resolve, 50));
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

      const importPromise = result.current.importRecords([0, 1, 2, 3, 4]);

      // Cancel after a short delay
      setTimeout(() => {
        result.current.cancelImport();
      }, 100);

      await importPromise;

      await waitFor(() => {
        expect(result.current.isImporting).toBe(false);
      });

      // Should have imported fewer than all records
      const totalProcessed = result.current.progress.completed + result.current.progress.failed;
      expect(totalProcessed).toBeLessThan(5);
    });

    it('should not start new import while one is in progress', async () => {
      const mutateAsync = vi.fn(async () => {
        await new Promise(resolve => setTimeout(resolve, 20));
        return { id: 'test' };
      });

      const records = createMockRecords(3);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      const firstImport = result.current.importRecords([0, 1]);
      
      // Wait a bit to ensure first import has started
      await new Promise(resolve => setTimeout(resolve, 10));
      
      const secondImport = result.current.importRecords([2]); // Should be ignored

      await Promise.all([firstImport, secondImport]);

      await waitFor(() => {
        expect(result.current.isImporting).toBe(false);
      });

      // Only the first import should have been processed
      expect(mutateAsync).toHaveBeenCalledTimes(2);
    });
  });

  describe('Progress tracking', () => {
    it('should track progress correctly', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        if (index === 1) {
          throw new Error('Error');
        }
        return { id: data.name };
      });

      const records = createMockRecords(4);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 0,
        })
      );

      await result.current.importRecords([0, 1, 2, 3]);

      await waitFor(() => {
        expect(result.current.progress.inProgress).toBe(false);
      });

      expect(result.current.progress.total).toBe(4);
      expect(result.current.progress.completed).toBe(3);
      expect(result.current.progress.failed).toBe(1);
    });

    it('should update status for each record', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        return { id: data.name };
      });

      const records = createMockRecords(2);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(2);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('success');
      expect(result.current.importStatus.get(1)?.status).toBe('success');
    });

    it('should set inProgress to true during import', async () => {
      const mutateAsync = vi.fn(async () => {
        await new Promise(resolve => setTimeout(resolve, 50));
        return { id: 'test' };
      });

      const records = createMockRecords(2);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      const importPromise = result.current.importRecords([0, 1]);

      // Check that inProgress is true during import
      await waitFor(() => {
        expect(result.current.progress.inProgress).toBe(true);
      });

      await importPromise;

      // Check that inProgress is false after import
      await waitFor(() => {
        expect(result.current.progress.inProgress).toBe(false);
      });
    });
  });

  describe('UUID conflict handling', () => {
    it('should handle 409 Conflict responses for duplicate UUIDs', async () => {
      const testUuid = '550e8400-e29b-41d4-a716-446655440001';
      const mutateAsync = vi.fn(async (data: any) => {
        if (data.id === testUuid) {
          const error: any = new Error('Conflict');
          error.response = {
            status: 409,
            data: { message: 'Entity with this UUID already exists' },
          };
          throw error;
        }
        return { id: data.id || 'generated-id' };
      });

      const records = [
        { index: 0, data: { name: 'Record 0', id: testUuid }, validationErrors: [], isValid: true },
        { index: 1, data: { name: 'Record 1' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 2,
        })
      );

      await result.current.importRecords([0, 1]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(1);
        expect(result.current.progress.failed).toBe(1);
      });

      // Should not retry 409 errors
      expect(mutateAsync).toHaveBeenCalledTimes(2); // One for each record, no retries
      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toContain('UUID conflict');
      expect(result.current.importStatus.get(0)?.error).toContain(testUuid);
      expect(result.current.importStatus.get(1)?.status).toBe('success');
    });

    it('should handle 400 Bad Request responses for invalid UUIDs from API', async () => {
      const invalidUuid = 'not-a-valid-uuid';
      const mutateAsync = vi.fn(async (data: any) => {
        if (data.id === invalidUuid) {
          const error: any = new Error('Bad Request');
          error.response = {
            status: 400,
            data: { message: 'Invalid UUID format provided' },
          };
          throw error;
        }
        return { id: data.id || 'generated-id' };
      });

      const records = [
        { index: 0, data: { name: 'Record 0', id: invalidUuid }, validationErrors: [], isValid: true },
        { index: 1, data: { name: 'Record 1' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 2,
        })
      );

      await result.current.importRecords([0, 1]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(1);
        expect(result.current.progress.failed).toBe(1);
      });

      // Should not retry 400 errors
      expect(mutateAsync).toHaveBeenCalledTimes(2); // One for each record, no retries
      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toContain('Invalid UUID format');
      expect(result.current.importStatus.get(0)?.error).toContain(invalidUuid);
      expect(result.current.importStatus.get(1)?.status).toBe('success');
    });

    it('should generate user-friendly error messages for UUID conflicts', async () => {
      const testUuid = '550e8400-e29b-41d4-a716-446655440002';
      const mutateAsync = vi.fn(async (data: any) => {
        const error: any = new Error('Conflict');
        error.response = {
          status: 409,
          data: { message: 'Duplicate UUID detected' },
        };
        throw error;
      });

      const records = [
        { index: 0, data: { name: 'Record 0', id: testUuid }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.failed).toBe(1);
      });

      const errorStatus = result.current.importStatus.get(0);
      expect(errorStatus?.status).toBe('error');
      expect(errorStatus?.error).toBe(`UUID conflict: Entity with ID ${testUuid} already exists`);
    });

    it('should continue importing remaining records after UUID conflict', async () => {
      const conflictUuid = '550e8400-e29b-41d4-a716-446655440003';
      const mutateAsync = vi.fn(async (data: any) => {
        if (data.id === conflictUuid) {
          const error: any = new Error('Conflict');
          error.response = {
            status: 409,
            data: { message: 'Entity already exists' },
          };
          throw error;
        }
        return { id: data.id || 'generated-id' };
      });

      const records = [
        { index: 0, data: { name: 'Record 0' }, validationErrors: [], isValid: true },
        { index: 1, data: { name: 'Record 1', id: conflictUuid }, validationErrors: [], isValid: true },
        { index: 2, data: { name: 'Record 2' }, validationErrors: [], isValid: true },
        { index: 3, data: { name: 'Record 3' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1, 2, 3]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
        expect(result.current.progress.failed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(4);
      expect(result.current.importStatus.get(0)?.status).toBe('success');
      expect(result.current.importStatus.get(1)?.status).toBe('error');
      expect(result.current.importStatus.get(2)?.status).toBe('success');
      expect(result.current.importStatus.get(3)?.status).toBe('success');
    });

    it('should allow retry functionality for failed UUID imports', async () => {
      const testUuid = '550e8400-e29b-41d4-a716-446655440004';
      let shouldFail = true;
      
      const mutateAsync = vi.fn(async (data: any) => {
        if (data.id === testUuid && shouldFail) {
          const error: any = new Error('Conflict');
          error.response = {
            status: 409,
            data: { message: 'Entity already exists' },
          };
          throw error;
        }
        return { id: data.id || 'generated-id' };
      });

      const records = [
        { index: 0, data: { name: 'Record 0', id: testUuid }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      // First import attempt - should fail
      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.failed).toBe(1);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('error');

      // User fixes the UUID conflict (simulated by changing shouldFail)
      shouldFail = false;

      // Retry the failed record
      await result.current.retryRecord(0);

      await waitFor(() => {
        expect(result.current.importStatus.get(0)?.status).toBe('success');
      });

      expect(mutateAsync).toHaveBeenCalledTimes(2); // Initial attempt + retry
    });

    it('should handle 409 conflict without UUID in message', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Conflict');
        error.response = {
          status: 409,
          data: { message: 'Resource conflict' },
        };
        throw error;
      });

      const records = [
        { index: 0, data: { name: 'Record 0' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.failed).toBe(1);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toBe('Conflict: Resource conflict');
    });

    it('should handle 400 validation error without UUID in message', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Bad Request');
        error.response = {
          status: 400,
          data: { message: 'Name is required' },
        };
        throw error;
      });

      const records = [
        { index: 0, data: { name: '' }, validationErrors: [], isValid: true },
      ];

      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.failed).toBe(1);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toBe('Validation error: Name is required');
    });
  });
});
