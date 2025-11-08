/**
 * Example usage of ImportSummary component
 * 
 * This file demonstrates how to use the ImportSummary component
 * to display final import results with success/failure counts.
 */

import { ImportSummary } from './ImportSummary';
import type { ImportProgress, ImportStatus, ParsedImportRecord } from '@/types/import';
import type { DtoCreateLabelRequest } from '@/api';

// Mock data for examples
const mockLabelRecords: ParsedImportRecord<DtoCreateLabelRequest>[] = [
  {
    index: 0,
    data: { name: 'Entertainment', color: '#FF5733' },
    validationErrors: [],
    isValid: true,
  },
  {
    index: 1,
    data: { name: 'Utilities', color: '#33FF57' },
    validationErrors: [],
    isValid: true,
  },
  {
    index: 2,
    data: { name: 'Subscriptions', color: '#3357FF' },
    validationErrors: [],
    isValid: true,
  },
  {
    index: 3,
    data: { name: 'Work', color: '#FF33F5' },
    validationErrors: [],
    isValid: true,
  },
  {
    index: 4,
    data: { name: 'Personal', color: '#F5FF33' },
    validationErrors: [],
    isValid: true,
  },
];

// Example 1: All records imported successfully
export function ImportSummaryFullSuccessExample() {
  const progress: ImportProgress = {
    total: 5,
    completed: 5,
    failed: 0,
    inProgress: false,
  };

  const importStatus = new Map<number, ImportStatus>([
    [0, { status: 'success' }],
    [1, { status: 'success' }],
    [2, { status: 'success' }],
    [3, { status: 'success' }],
    [4, { status: 'success' }],
  ]);

  return (
    <div className="max-w-3xl mx-auto p-6">
      <ImportSummary
        progress={progress}
        records={mockLabelRecords}
        importStatus={importStatus}
        onReturnToList={() => console.log('Return to list')}
        entityType="labels"
      />
    </div>
  );
}

// Example 2: Partial success with some failures
export function ImportSummaryPartialSuccessExample() {
  const progress: ImportProgress = {
    total: 5,
    completed: 5,
    failed: 2,
    inProgress: false,
  };

  const importStatus = new Map<number, ImportStatus>([
    [0, { status: 'success' }],
    [1, { status: 'error', error: 'Label with this name already exists' }],
    [2, { status: 'success' }],
    [3, { status: 'error', error: 'Network error: Failed to connect to server' }],
    [4, { status: 'success' }],
  ]);

  return (
    <div className="max-w-3xl mx-auto p-6">
      <ImportSummary
        progress={progress}
        records={mockLabelRecords}
        importStatus={importStatus}
        onRetryFailed={() => console.log('Retry all failed')}
        onRetryRecord={(index) => console.log('Retry record', index)}
        onReturnToList={() => console.log('Return to list')}
        entityType="labels"
      />
    </div>
  );
}

// Example 3: All records failed
export function ImportSummaryFullFailureExample() {
  const progress: ImportProgress = {
    total: 5,
    completed: 5,
    failed: 5,
    inProgress: false,
  };

  const importStatus = new Map<number, ImportStatus>([
    [0, { status: 'error', error: 'Invalid color format' }],
    [1, { status: 'error', error: 'Label name is required' }],
    [2, { status: 'error', error: 'Quota limit exceeded' }],
    [3, { status: 'error', error: 'Network timeout' }],
    [4, { status: 'error', error: 'Server error: 500 Internal Server Error' }],
  ]);

  return (
    <div className="max-w-3xl mx-auto p-6">
      <ImportSummary
        progress={progress}
        records={mockLabelRecords}
        importStatus={importStatus}
        onRetryFailed={() => console.log('Retry all failed')}
        onRetryRecord={(index) => console.log('Retry record', index)}
        onReturnToList={() => console.log('Return to list')}
        entityType="labels"
      />
    </div>
  );
}

// Example 4: Providers import summary
export function ImportSummaryProvidersExample() {
  const mockProviderRecords: ParsedImportRecord<any>[] = [
    {
      index: 0,
      data: { name: 'Netflix', description: 'Streaming service' },
      validationErrors: [],
      isValid: true,
    },
    {
      index: 1,
      data: { name: 'Spotify', description: 'Music streaming' },
      validationErrors: [],
      isValid: true,
    },
    {
      index: 2,
      data: { name: 'AWS', description: 'Cloud services' },
      validationErrors: [],
      isValid: true,
    },
  ];

  const progress: ImportProgress = {
    total: 3,
    completed: 3,
    failed: 1,
    inProgress: false,
  };

  const importStatus = new Map<number, ImportStatus>([
    [0, { status: 'success' }],
    [1, { status: 'error', error: 'Provider already exists' }],
    [2, { status: 'success' }],
  ]);

  return (
    <div className="max-w-3xl mx-auto p-6">
      <ImportSummary
        progress={progress}
        records={mockProviderRecords}
        importStatus={importStatus}
        onRetryRecord={(index) => console.log('Retry record', index)}
        onReturnToList={() => console.log('Return to providers list')}
        entityType="providers"
      />
    </div>
  );
}
