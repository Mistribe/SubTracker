import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useImportManager } from '../useImportManager';
import type { ParsedImportRecord } from '../useImportManager';
import type { UseMutationResult } from '@tanstack/react-query';

describe('useImportManager - API Integration Tests', () => {
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

  const createMockRecords = (count: number): ParsedImportRecord<any>[] => {
    return Array.from({ length: count }, (_, i) => ({
      index: i,
      data: { name: `Record ${i}`, value: i * 10 },
      validationErrors: [],
      isValid: true,
    }));
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Successful API calls', () => {
    it('should make successful API calls for all records', async () => {
      const apiResponses: any[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        const response = { id: `id-${data.name}`, ...data, createdAt: new Date().toISOString() };
        apiResponses.push(response);
        return response;
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

      expect(mutateAsync).toHaveBeenCalledTimes(3);
      expect(apiResponses).toHaveLength(3);
      
      // Verify each API call received correct data
      expect(mutateAsync).toHaveBeenNthCalledWith(1, records[0].data);
      expect(mutateAsync).toHaveBeenNthCalledWith(2, records[1].data);
      expect(mutateAsync).toHaveBeenNthCalledWith(3, records[2].data);
    });

    it('should handle API responses with different data structures', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        // Simulate different API response structures
        return {
          success: true,
          data: {
            id: `uuid-${Date.now()}`,
            attributes: data,
            metadata: {
              version: 1,
              timestamp: new Date().toISOString(),
            },
          },
        };
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

    it('should handle empty API responses', async () => {
      const mutateAsync = vi.fn(async () => {
        return {}; // Empty response
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(1);
      });

      expect(result.current.importStatus.get(0)?.status).toBe('success');
    });
  });

  describe('API error handling', () => {
    it('should handle 400 Bad Request errors', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Bad Request');
        error.response = {
          status: 400,
          data: { message: 'Invalid field: name is required' },
        };
        throw error;
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
      expect(result.current.importStatus.get(0)?.error).toContain('Validation error');
      expect(result.current.importStatus.get(0)?.error).toContain('Invalid field');
    });

    it('should handle 401 Unauthorized errors', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Unauthorized');
        error.response = {
          status: 401,
          data: { message: 'Invalid authentication token' },
        };
        throw error;
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

      expect(result.current.importStatus.get(0)?.error).toContain('Authentication error');
    });

    it('should handle 404 Not Found errors', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Not Found');
        error.response = {
          status: 404,
          data: { message: 'Resource not found' },
        };
        throw error;
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

      expect(result.current.importStatus.get(0)?.error).toContain('Resource not found');
    });

    it('should handle 409 Conflict errors', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Conflict');
        error.response = {
          status: 409,
          data: { message: 'Record with this name already exists' },
        };
        throw error;
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

      expect(result.current.importStatus.get(0)?.error).toContain('Conflict');
      expect(result.current.importStatus.get(0)?.error).toContain('already exists');
    });

    it('should handle 429 Rate Limit errors with retry', async () => {
      let callCount = 0;
      const mutateAsync = vi.fn(async () => {
        callCount++;
        if (callCount <= 2) {
          const error: any = new Error('Too Many Requests');
          error.response = {
            status: 429,
            data: { message: 'Rate limit exceeded' },
          };
          throw error;
        }
        return { id: 'success' };
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

      expect(callCount).toBe(3); // Initial + 2 retries
      expect(result.current.importStatus.get(0)?.status).toBe('success');
    });

    it('should handle 500 Internal Server Error with retry', async () => {
      let callCount = 0;
      const mutateAsync = vi.fn(async () => {
        callCount++;
        if (callCount === 1) {
          const error: any = new Error('Internal Server Error');
          error.response = {
            status: 500,
            data: { message: 'Database connection failed' },
          };
          throw error;
        }
        return { id: 'success' };
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
        expect(result.current.progress.completed).toBe(1);
      });

      expect(callCount).toBe(2); // Initial + 1 retry
      expect(result.current.importStatus.get(0)?.status).toBe('success');
    });

    it('should handle network errors', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Network Error');
        error.request = {}; // Indicates network error (request made but no response)
        throw error;
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
          maxRetries: 1,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.failed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(2); // Initial + 1 retry
      expect(result.current.importStatus.get(0)?.error).toContain('Network error');
    });

    it('should handle timeout errors', async () => {
      const mutateAsync = vi.fn(async () => {
        const error: any = new Error('Request timeout');
        error.code = 'ECONNABORTED';
        error.request = {};
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
        expect(result.current.progress.failed).toBe(1);
      });

      expect(mutateAsync).toHaveBeenCalledTimes(3); // Initial + 2 retries
    });
  });

  describe('Rate limiting behavior', () => {
    it('should add delay between sequential API calls', async () => {
      const timestamps: number[] = [];
      const mutateAsync = vi.fn(async () => {
        timestamps.push(Date.now());
        return { id: 'test' };
      });

      const records = createMockRecords(4);
      const delayBetweenCalls = 100;
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls,
        })
      );

      await result.current.importRecords([0, 1, 2, 3]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(4);
      });

      // Verify delays between calls - allow generous tolerance for timing variations in CI
      for (let i = 1; i < timestamps.length; i++) {
        const actualDelay = timestamps[i] - timestamps[i - 1];
        expect(actualDelay).toBeGreaterThanOrEqual(delayBetweenCalls - 40); // Allow 40ms tolerance
      }
    });

    it('should respect configurable delay between calls', async () => {
      const timestamps: number[] = [];
      const mutateAsync = vi.fn(async () => {
        timestamps.push(Date.now());
        return { id: 'test' };
      });

      const records = createMockRecords(3);
      const customDelay = 200;
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: customDelay,
        })
      );

      await result.current.importRecords([0, 1, 2]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(3);
      });

      // Verify custom delay is respected - allow more tolerance for timing variations
      for (let i = 1; i < timestamps.length; i++) {
        const actualDelay = timestamps[i] - timestamps[i - 1];
        expect(actualDelay).toBeGreaterThanOrEqual(customDelay - 60);
      }
    });

    it('should implement exponential backoff on retries', async () => {
      const retryTimestamps: number[] = [];
      let attemptCount = 0;

      const mutateAsync = vi.fn(async () => {
        attemptCount++;
        retryTimestamps.push(Date.now());
        
        if (attemptCount <= 3) {
          const error: any = new Error('Temporary error');
          error.response = { status: 500 };
          throw error;
        }
        return { id: 'success' };
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
        expect(result.current.progress.completed).toBe(1);
      });

      // Verify exponential backoff (each retry should take longer)
      if (retryTimestamps.length >= 3) {
        const delay1 = retryTimestamps[1] - retryTimestamps[0];
        const delay2 = retryTimestamps[2] - retryTimestamps[1];
        const delay3 = retryTimestamps[3] - retryTimestamps[2];
        
        expect(delay2).toBeGreaterThan(delay1);
        expect(delay3).toBeGreaterThan(delay2);
      }
    });

    it('should cap exponential backoff at maximum delay', async () => {
      const retryTimestamps: number[] = [];
      let attemptCount = 0;

      const mutateAsync = vi.fn(async () => {
        attemptCount++;
        retryTimestamps.push(Date.now());
        
        if (attemptCount <= 5) {
          const error: any = new Error('Persistent error');
          error.response = { status: 500 };
          throw error;
        }
        return { id: 'success' };
      });

      const records = createMockRecords(1);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 100,
          maxRetries: 5,
        })
      );

      await result.current.importRecords([0]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(1);
      });

      // Verify that backoff doesn't exceed 10 seconds (10000ms)
      for (let i = 1; i < retryTimestamps.length; i++) {
        const delay = retryTimestamps[i] - retryTimestamps[i - 1];
        expect(delay).toBeLessThanOrEqual(10100); // 10s + 100ms tolerance
      }
    });
  });

  describe('Sequential request processing', () => {
    it('should process requests one at a time', async () => {
      let activeRequests = 0;
      let maxConcurrentRequests = 0;

      const mutateAsync = vi.fn(async () => {
        activeRequests++;
        maxConcurrentRequests = Math.max(maxConcurrentRequests, activeRequests);
        
        await new Promise(resolve => setTimeout(resolve, 30));
        
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

      // Verify only one request was active at a time
      expect(maxConcurrentRequests).toBe(1);
    });

    it('should maintain order of API calls', async () => {
      const callOrder: number[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        callOrder.push(index);
        return { id: `test-${index}` };
      });

      const records = createMockRecords(6);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      await result.current.importRecords([0, 1, 2, 3, 4, 5]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(6);
      });

      // Verify calls were made in order
      expect(callOrder).toEqual([0, 1, 2, 3, 4, 5]);
    });

    it('should not skip records after a failure', async () => {
      const processedIndices: number[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        processedIndices.push(index);
        
        if (index === 2) {
          throw new Error('Failed at index 2');
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

      await result.current.importRecords([0, 1, 2, 3, 4]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(4);
        expect(result.current.progress.failed).toBe(1);
      });

      // Verify all records were processed in order
      expect(processedIndices).toEqual([0, 1, 2, 3, 4]);
    });

    it('should handle selective record import in order', async () => {
      const callOrder: number[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        const index = parseInt(data.name.split(' ')[1]);
        callOrder.push(index);
        return { id: `test-${index}` };
      });

      const records = createMockRecords(10);
      const { result } = renderHook(() =>
        useImportManager({
          records,
          createMutation: createMockMutation(mutateAsync),
          delayBetweenCalls: 10,
        })
      );

      // Import only specific indices
      await result.current.importRecords([1, 3, 5, 7, 9]);

      await waitFor(() => {
        expect(result.current.progress.completed).toBe(5);
      });

      // Verify only selected records were called in order
      expect(callOrder).toEqual([1, 3, 5, 7, 9]);
    });
  });
});
