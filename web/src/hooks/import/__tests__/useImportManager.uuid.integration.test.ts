import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useImportManager } from '../useImportManager';
import { fileParser } from '../../../services/fileParser';
import { LabelFieldMapper, ProviderFieldMapper, SubscriptionFieldMapper } from '../../../services/importMapper';
import type { ParsedImportRecord } from '../useImportManager';
import type { UseMutationResult } from '@tanstack/react-query';
import type { DtoCreateLabelRequest, DtoCreateProviderRequest, DtoCreateSubscriptionRequest } from '../../../api';

describe('useImportManager - UUID Import Integration Tests', () => {
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

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('9.1 Test file upload and parsing with IDs', () => {
    it('should parse CSV file with ID fields', async () => {
      const csvContent = `id,name,color,ownerType
550e8400-e29b-41d4-a716-446655440001,Entertainment,#FF5733,personal
550e8400-e29b-41d4-a716-446655440002,Utilities,#33FF57,personal`;

      const file = new File([csvContent], 'labels.csv', { type: 'text/csv' });
      const parsedData = await fileParser.parseCSV(file);

      expect(parsedData).toHaveLength(2);
      expect(parsedData[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(parsedData[0].name).toBe('Entertainment');
      expect(parsedData[1].id).toBe('550e8400-e29b-41d4-a716-446655440002');
      expect(parsedData[1].name).toBe('Utilities');
    });

    it('should parse JSON file with ID fields', async () => {
      const jsonContent = JSON.stringify([
        {
          id: '550e8400-e29b-41d4-a716-446655440001',
          name: 'Netflix',
          description: 'Streaming service',
          url: 'https://netflix.com',
        },
        {
          id: '550e8400-e29b-41d4-a716-446655440002',
          name: 'Spotify',
          description: 'Music streaming',
          url: 'https://spotify.com',
        },
      ]);

      const file = new File([jsonContent], 'providers.json', { type: 'application/json' });
      const parsedData = await fileParser.parseJSON(file);

      expect(parsedData).toHaveLength(2);
      expect(parsedData[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(parsedData[0].name).toBe('Netflix');
      expect(parsedData[1].id).toBe('550e8400-e29b-41d4-a716-446655440002');
      expect(parsedData[1].name).toBe('Spotify');
    });

    it('should parse YAML file with ID fields', async () => {
      const yamlContent = `- id: "550e8400-e29b-41d4-a716-446655440001"
  providerId: "provider-123"
  friendlyName: "Netflix Premium"
  startDate: "2024-01-01"
  recurrency: "monthly"
  customPriceAmount: 15.99
  customPriceCurrency: "USD"
  ownerType: "personal"
- id: "550e8400-e29b-41d4-a716-446655440002"
  providerId: "provider-456"
  friendlyName: "Spotify Family"
  startDate: "2024-02-15"
  recurrency: "monthly"
  customPriceAmount: 16.99
  customPriceCurrency: "USD"
  ownerType: "personal"`;

      const file = new File([yamlContent], 'subscriptions.yaml', { type: 'text/yaml' });
      const parsedData = await fileParser.parseYAML(file);

      expect(parsedData).toHaveLength(2);
      expect(parsedData[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(parsedData[0].friendlyName).toBe('Netflix Premium');
      expect(parsedData[1].id).toBe('550e8400-e29b-41d4-a716-446655440002');
      expect(parsedData[1].friendlyName).toBe('Spotify Family');
    });

    it('should parse files with mixed ID presence (some with, some without)', async () => {
      const csvContent = `id,name,color,ownerType
550e8400-e29b-41d4-a716-446655440001,Entertainment,#FF5733,personal
,Utilities,#33FF57,personal
550e8400-e29b-41d4-a716-446655440003,Productivity,#3357FF,personal
,Health,#FF33A1,personal`;

      const file = new File([csvContent], 'labels.csv', { type: 'text/csv' });
      const parsedData = await fileParser.parseCSV(file);

      expect(parsedData).toHaveLength(4);
      expect(parsedData[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(parsedData[1].id).toBe('');
      expect(parsedData[2].id).toBe('550e8400-e29b-41d4-a716-446655440003');
      expect(parsedData[3].id).toBe('');
    });

    it('should correctly extract and validate ID values through field mapper', async () => {
      const csvContent = `id,name,color,ownerType
550e8400-e29b-41d4-a716-446655440001,Entertainment,#FF5733,personal
invalid-uuid,Utilities,#33FF57,personal
,Productivity,#3357FF,personal`;

      const file = new File([csvContent], 'labels.csv', { type: 'text/csv' });
      const parsedData = await fileParser.parseCSV(file);

      const mapper = new LabelFieldMapper();
      
      // First record: valid UUID - should be mapped
      const data0 = mapper.mapFields(parsedData[0]);
      expect((data0 as any).id).toBe('550e8400-e29b-41d4-a716-446655440001');
      const validation0 = mapper.validate(data0);
      expect(validation0.isValid).toBe(true);

      // Second record: invalid UUID - should not be mapped (ID will be undefined)
      const data1 = mapper.mapFields(parsedData[1]);
      expect((data1 as any).id).toBeUndefined(); // Invalid UUIDs are not mapped
      const validation1 = mapper.validate(data1);
      expect(validation1.isValid).toBe(true); // Valid because ID is not present

      // Third record: no UUID - should be valid
      const data2 = mapper.mapFields(parsedData[2]);
      expect((data2 as any).id).toBeUndefined();
      const validation2 = mapper.validate(data2);
      expect(validation2.isValid).toBe(true);
    });
  });

  describe('9.2 Test import with valid UUIDs', () => {
    it('should import records with valid UUIDs', async () => {
      const importedRecords: any[] = [];
      const mutateAsync = vi.fn(async (data: DtoCreateLabelRequest) => {
        importedRecords.push(data);
        return { id: data.id || 'generated-id', ...data };
      });

      const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
        {
          index: 0,
          data: {
            id: '550e8400-e29b-41d4-a716-446655440001',
            name: 'Entertainment',
            color: '#FF5733',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 1,
          data: {
            id: '550e8400-e29b-41d4-a716-446655440002',
            name: 'Utilities',
            color: '#33FF57',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
      ];

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

      // Verify API calls included ID field
      expect(importedRecords).toHaveLength(2);
      expect(importedRecords[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(importedRecords[1].id).toBe('550e8400-e29b-41d4-a716-446655440002');
    });

    it('should import records without UUIDs', async () => {
      const importedRecords: any[] = [];
      const mutateAsync = vi.fn(async (data: DtoCreateProviderRequest) => {
        importedRecords.push(data);
        return { id: 'auto-generated-uuid', ...data };
      });

      const records: ParsedImportRecord<DtoCreateProviderRequest>[] = [
        {
          index: 0,
          data: {
            name: 'Netflix',
            description: 'Streaming service',
            url: 'https://netflix.com',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 1,
          data: {
            name: 'Spotify',
            description: 'Music streaming',
            url: 'https://spotify.com',
          },
          validationErrors: [],
          isValid: true,
        },
      ];

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

      // Verify API calls omitted ID field
      expect(importedRecords).toHaveLength(2);
      expect(importedRecords[0].id).toBeUndefined();
      expect(importedRecords[1].id).toBeUndefined();
    });

    it('should handle mixed import (some with UUIDs, some without)', async () => {
      const importedRecords: any[] = [];
      const mutateAsync = vi.fn(async (data: DtoCreateSubscriptionRequest) => {
        importedRecords.push(data);
        return { id: data.id || 'auto-generated-uuid', ...data };
      });

      const records: ParsedImportRecord<DtoCreateSubscriptionRequest>[] = [
        {
          index: 0,
          data: {
            id: '550e8400-e29b-41d4-a716-446655440001',
            providerId: 'provider-123',
            friendlyName: 'Netflix Premium',
            startDate: '2024-01-01',
            recurrency: 'monthly',
            customPriceAmount: 15.99,
            customPriceCurrency: 'USD',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 1,
          data: {
            providerId: 'provider-456',
            friendlyName: 'Spotify Family',
            startDate: '2024-02-15',
            recurrency: 'monthly',
            customPriceAmount: 16.99,
            customPriceCurrency: 'USD',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 2,
          data: {
            id: '550e8400-e29b-41d4-a716-446655440003',
            providerId: 'provider-789',
            friendlyName: 'GitHub Pro',
            startDate: '2024-03-01',
            recurrency: 'monthly',
            customPriceAmount: 4.00,
            customPriceCurrency: 'USD',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
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
        expect(result.current.progress.completed).toBe(3);
      });

      // Verify API calls include ID when present, omit when not
      expect(importedRecords).toHaveLength(3);
      expect(importedRecords[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(importedRecords[1].id).toBeUndefined();
      expect(importedRecords[2].id).toBe('550e8400-e29b-41d4-a716-446655440003');
    });
  });

  describe('9.3 Test UUID validation errors', () => {
    it('should detect invalid UUID format', async () => {
      const csvContent = `id,name,color,ownerType
not-a-uuid,Entertainment,#FF5733,personal
550e8400-INVALID-FORMAT,Utilities,#33FF57,personal
12345,Productivity,#3357FF,personal`;

      const file = new File([csvContent], 'labels.csv', { type: 'text/csv' });
      const parsedData = await fileParser.parseCSV(file);

      const mapper = new LabelFieldMapper();
      
      // All records have invalid UUIDs, so they should not be mapped
      const data0 = mapper.mapFields(parsedData[0]);
      expect((data0 as any).id).toBeUndefined(); // Invalid UUID not mapped
      
      const data1 = mapper.mapFields(parsedData[1]);
      expect((data1 as any).id).toBeUndefined(); // Invalid UUID not mapped
      
      const data2 = mapper.mapFields(parsedData[2]);
      expect((data2 as any).id).toBeUndefined(); // Invalid UUID not mapped
      
      // All should be valid since invalid IDs are simply not included
      expect(mapper.validate(data0).isValid).toBe(true);
      expect(mapper.validate(data1).isValid).toBe(true);
      expect(mapper.validate(data2).isValid).toBe(true);
    });

    it('should prevent import of records with invalid UUIDs', async () => {
      const mutateAsync = vi.fn(async (data: any) => {
        return { id: data.id || 'generated-id', ...data };
      });

      const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
        {
          index: 0,
          data: {
            id: 'invalid-uuid',
            name: 'Entertainment',
            color: '#FF5733',
            ownerType: 'personal',
          },
          validationErrors: [
            {
              field: 'id',
              message: 'Invalid UUID format: "invalid-uuid"',
              severity: 'error',
            },
          ],
          isValid: false,
        },
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

      // API should not be called for invalid records
      expect(mutateAsync).not.toHaveBeenCalled();
      expect(result.current.importStatus.get(0)?.status).toBe('error');
      expect(result.current.importStatus.get(0)?.error).toBe('Record has validation errors');
    });

    it('should allow import of valid records while blocking invalid ones', async () => {
      const importedRecords: any[] = [];
      const mutateAsync = vi.fn(async (data: any) => {
        importedRecords.push(data);
        return { id: data.id || 'generated-id', ...data };
      });

      const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
        {
          index: 0,
          data: {
            id: '550e8400-e29b-41d4-a716-446655440001',
            name: 'Entertainment',
            color: '#FF5733',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 1,
          data: {
            id: 'invalid-uuid',
            name: 'Utilities',
            color: '#33FF57',
            ownerType: 'personal',
          },
          validationErrors: [
            {
              field: 'id',
              message: 'Invalid UUID format',
              severity: 'error',
            },
          ],
          isValid: false,
        },
        {
          index: 2,
          data: {
            id: '550e8400-e29b-41d4-a716-446655440003',
            name: 'Productivity',
            color: '#3357FF',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
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

      // Only valid records should be imported
      expect(importedRecords).toHaveLength(2);
      expect(importedRecords[0].id).toBe('550e8400-e29b-41d4-a716-446655440001');
      expect(importedRecords[1].id).toBe('550e8400-e29b-41d4-a716-446655440003');

      // Invalid record should have error status
      expect(result.current.importStatus.get(1)?.status).toBe('error');
    });
  });

  describe('9.4 Test UUID conflict handling', () => {
    it('should handle duplicate UUID API error (409 Conflict)', async () => {
      const conflictUuid = '550e8400-e29b-41d4-a716-446655440001';
      const mutateAsync = vi.fn(async (data: any) => {
        if (data.id === conflictUuid) {
          const error: any = new Error('Conflict');
          error.response = {
            status: 409,
            data: { message: 'Entity with this UUID already exists' },
          };
          throw error;
        }
        return { id: data.id || 'generated-id', ...data };
      });

      const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
        {
          index: 0,
          data: {
            id: conflictUuid,
            name: 'Entertainment',
            color: '#FF5733',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
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
      expect(result.current.importStatus.get(0)?.error).toContain('UUID conflict');
      expect(result.current.importStatus.get(0)?.error).toContain(conflictUuid);
    });

    it('should display error message for UUID conflicts', async () => {
      const conflictUuid = '550e8400-e29b-41d4-a716-446655440002';
      const mutateAsync = vi.fn(async (data: any) => {
        const error: any = new Error('Conflict');
        error.response = {
          status: 409,
          data: { message: 'Duplicate UUID detected' },
        };
        throw error;
      });

      const records: ParsedImportRecord<DtoCreateProviderRequest>[] = [
        {
          index: 0,
          data: {
            id: conflictUuid,
            name: 'Netflix',
            description: 'Streaming service',
            url: 'https://netflix.com',
          },
          validationErrors: [],
          isValid: true,
        },
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
      expect(errorStatus?.error).toBe(`UUID conflict: Entity with ID ${conflictUuid} already exists`);
    });

    it('should continue import after UUID conflict', async () => {
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
        return { id: data.id || 'generated-id', ...data };
      });

      const records: ParsedImportRecord<DtoCreateSubscriptionRequest>[] = [
        {
          index: 0,
          data: {
            providerId: 'provider-123',
            friendlyName: 'Netflix Premium',
            startDate: '2024-01-01',
            recurrency: 'monthly',
            customPriceAmount: 15.99,
            customPriceCurrency: 'USD',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 1,
          data: {
            id: conflictUuid,
            providerId: 'provider-456',
            friendlyName: 'Spotify Family',
            startDate: '2024-02-15',
            recurrency: 'monthly',
            customPriceAmount: 16.99,
            customPriceCurrency: 'USD',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
        {
          index: 2,
          data: {
            providerId: 'provider-789',
            friendlyName: 'GitHub Pro',
            startDate: '2024-03-01',
            recurrency: 'monthly',
            customPriceAmount: 4.00,
            customPriceCurrency: 'USD',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
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

      // Verify import continued after conflict
      expect(result.current.importStatus.get(0)?.status).toBe('success');
      expect(result.current.importStatus.get(1)?.status).toBe('error');
      expect(result.current.importStatus.get(2)?.status).toBe('success');
    });

    it('should allow retry functionality after UUID conflict resolution', async () => {
      const conflictUuid = '550e8400-e29b-41d4-a716-446655440004';
      let shouldFail = true;

      const mutateAsync = vi.fn(async (data: any) => {
        if (data.id === conflictUuid && shouldFail) {
          const error: any = new Error('Conflict');
          error.response = {
            status: 409,
            data: { message: 'Entity already exists' },
          };
          throw error;
        }
        return { id: data.id || 'generated-id', ...data };
      });

      const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
        {
          index: 0,
          data: {
            id: conflictUuid,
            name: 'Entertainment',
            color: '#FF5733',
            ownerType: 'personal',
          },
          validationErrors: [],
          isValid: true,
        },
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
      expect(result.current.importStatus.get(0)?.error).toContain('UUID conflict');

      // Simulate user resolving the conflict
      shouldFail = false;

      // Retry the failed record
      await result.current.retryRecord(0);

      await waitFor(() => {
        expect(result.current.importStatus.get(0)?.status).toBe('success');
      });

      expect(mutateAsync).toHaveBeenCalledTimes(2); // Initial attempt + retry
    });
  });
});
