import * as React from 'react';
import { ImportPreviewTable } from './ImportPreviewTable';
import type {
  ParsedImportRecord,
  ImportColumnDef,
  ImportStatus,
  ImportProgress,
} from '@/types/import';
import type { DtoCreateLabelRequest } from '@/api';

/**
 * Example usage of ImportPreviewTable component
 * This demonstrates how to use the component with label imports
 */
export function ImportPreviewTableExample() {
  // Example data
  const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
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
      data: { name: 'Invalid Label', color: 'not-a-color' },
      validationErrors: [
        { field: 'color', message: 'Invalid hex color format', severity: 'error' },
      ],
      isValid: false,
    },
  ];

  // Column definitions
  const columns: ImportColumnDef<DtoCreateLabelRequest>[] = [
    {
      key: 'name',
      label: 'Name',
    },
    {
      key: 'color',
      label: 'Color',
      render: (value) => (
        <div className="flex items-center gap-2">
          <div
            className="size-4 rounded border"
            style={{ backgroundColor: value as string }}
          />
          <span>{value as string}</span>
        </div>
      ),
    },
  ];

  // State
  const [selectedRecords, setSelectedRecords] = React.useState<Set<number>>(new Set());
  const [importStatus, setImportStatus] = React.useState<Map<number, ImportStatus>>(
    new Map()
  );
  const [isImporting, setIsImporting] = React.useState(false);
  const [progress, setProgress] = React.useState<ImportProgress>({
    total: 0,
    completed: 0,
    failed: 0,
    inProgress: false,
  });

  const handleImportSelected = async () => {
    setIsImporting(true);
    setProgress({
      total: selectedRecords.size,
      completed: 0,
      failed: 0,
      inProgress: true,
    });

    // Simulate importing selected records
    for (const index of Array.from(selectedRecords)) {
      const record = records[index];
      if (!record.isValid) continue;

      // Update status to importing
      setImportStatus((prev) => new Map(prev).set(index, { status: 'importing' }));

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 500));

      // Simulate success/failure
      const success = Math.random() > 0.2;
      setImportStatus((prev) =>
        new Map(prev).set(index, {
          status: success ? 'success' : 'error',
          error: success ? undefined : 'Failed to create label',
        })
      );

      setProgress((prev) => ({
        ...prev,
        completed: prev.completed + 1,
        failed: success ? prev.failed : prev.failed + 1,
      }));
    }

    setProgress((prev) => ({ ...prev, inProgress: false }));
    setIsImporting(false);
  };

  const handleImportAll = async () => {
    const validIndices = records
      .map((r, i) => (r.isValid ? i : -1))
      .filter((i) => i !== -1);
    setSelectedRecords(new Set(validIndices));
    await handleImportSelected();
  };

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-2xl font-bold mb-6">Import Preview Table Example</h1>
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={selectedRecords}
        onSelectionChange={setSelectedRecords}
        onImportSelected={handleImportSelected}
        onImportAll={handleImportAll}
        importStatus={importStatus}
        isImporting={isImporting}
        progress={progress}
        enableVirtualization={true}
        virtualizationThreshold={100}
      />
    </div>
  );
}
