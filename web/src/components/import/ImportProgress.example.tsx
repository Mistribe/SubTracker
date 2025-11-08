/**
 * Example usage of ImportProgress component
 * 
 * This file demonstrates how to use the ImportProgress component
 * to display real-time progress during bulk import operations.
 */

import { ImportProgress } from './ImportProgress';
import type { ImportProgress as ImportProgressType } from '@/types/import';

// Example 1: Import in progress
export function ImportProgressInProgressExample() {
  const progress: ImportProgressType = {
    total: 100,
    completed: 45,
    failed: 3,
    inProgress: true,
  };

  const startTime = Date.now() - 30000; // Started 30 seconds ago

  return (
    <div className="max-w-2xl mx-auto p-6">
      <ImportProgress progress={progress} startTime={startTime} />
    </div>
  );
}

// Example 2: Import completed successfully
export function ImportProgressCompletedExample() {
  const progress: ImportProgressType = {
    total: 50,
    completed: 50,
    failed: 0,
    inProgress: false,
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <ImportProgress progress={progress} />
    </div>
  );
}

// Example 3: Import completed with some failures
export function ImportProgressWithFailuresExample() {
  const progress: ImportProgressType = {
    total: 75,
    completed: 75,
    failed: 12,
    inProgress: false,
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <ImportProgress progress={progress} />
    </div>
  );
}

// Example 4: Just started import
export function ImportProgressJustStartedExample() {
  const progress: ImportProgressType = {
    total: 200,
    completed: 5,
    failed: 0,
    inProgress: true,
  };

  const startTime = Date.now() - 2000; // Started 2 seconds ago

  return (
    <div className="max-w-2xl mx-auto p-6">
      <ImportProgress progress={progress} startTime={startTime} />
    </div>
  );
}
