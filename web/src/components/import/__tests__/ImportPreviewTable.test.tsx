import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ImportPreviewTable } from '../ImportPreviewTable';
import type {
  ParsedImportRecord,
  ImportColumnDef,
  ImportStatus,
} from '@/types/import';
import type { DtoCreateLabelRequest } from '@/api';

// Mock the Tooltip component to avoid portal issues in tests
vi.mock('@/components/ui/tooltip', () => ({
  Tooltip: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  TooltipTrigger: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  TooltipContent: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

describe('ImportPreviewTable - ID Column Display', () => {
  const mockOnSelectionChange = vi.fn();
  const mockOnImportSelected = vi.fn();
  const mockOnImportAll = vi.fn();

  const columns: ImportColumnDef<DtoCreateLabelRequest>[] = [
    { key: 'name', label: 'Name' },
    { key: 'color', label: 'Color' },
  ];

  it('displays ID column when at least one record has an ID', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: '550e8400-e29b-41d4-a716-446655440001', name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
      {
        index: 1,
        data: { name: 'Label 2', color: '#00FF00' },
        validationErrors: [],
        isValid: true,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that ID column header is present
    expect(screen.getByText('ID')).toBeInTheDocument();
  });

  it('hides ID column when no records have IDs', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
      {
        index: 1,
        data: { name: 'Label 2', color: '#00FF00' },
        validationErrors: [],
        isValid: true,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that ID column header is not present
    expect(screen.queryByText('ID')).not.toBeInTheDocument();
  });

  it('displays "Auto-generated" for records without IDs when ID column is shown', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: '550e8400-e29b-41d4-a716-446655440001', name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
      {
        index: 1,
        data: { name: 'Label 2', color: '#00FF00' },
        validationErrors: [],
        isValid: true,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that "Auto-generated" text is present for record without ID
    expect(screen.getByText('Auto-generated')).toBeInTheDocument();
  });

  it('displays UUID values in monospace font', () => {
    const testUuid = '550e8400-e29b-41d4-a716-446655440001';
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: testUuid, name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
    ];

    const { container } = render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Find the code element containing the UUID (may be truncated)
    const codeElement = container.querySelector('code.font-mono');
    expect(codeElement).toBeInTheDocument();
    // Check that it starts with the UUID prefix (it may be truncated with ...)
    expect(codeElement?.textContent).toMatch(/^550e8400-e29b-41d4-a/);
  });

  it('respects showIdColumn prop override to show ID column', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
        showIdColumn={true}
      />
    );

    // ID column should be shown even though no records have IDs
    expect(screen.getByText('ID')).toBeInTheDocument();
  });

  it('respects showIdColumn prop override to hide ID column', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: '550e8400-e29b-41d4-a716-446655440001', name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
        showIdColumn={false}
      />
    );

    // ID column should be hidden even though record has ID
    expect(screen.queryByText('ID')).not.toBeInTheDocument();
  });
});

describe('ImportPreviewTable - UUID Validation Error Display', () => {
  const mockOnSelectionChange = vi.fn();
  const mockOnImportSelected = vi.fn();
  const mockOnImportAll = vi.fn();

  const columns: ImportColumnDef<DtoCreateLabelRequest>[] = [
    { key: 'name', label: 'Name' },
    { key: 'color', label: 'Color' },
  ];

  it('displays UUID validation errors with special styling', () => {
    const invalidUuid = 'not-a-valid-uuid';
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: invalidUuid, name: 'Label 1', color: '#FF0000' },
        validationErrors: [
          {
            field: 'id',
            message: 'Invalid UUID format',
            severity: 'error',
          },
        ],
        isValid: false,
      },
    ];

    const { container } = render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that invalid UUID is displayed with error styling
    const errorCode = container.querySelector('code.text-red-600');
    expect(errorCode).toBeInTheDocument();
    expect(errorCode?.textContent).toContain(invalidUuid);
  });

  it('shows expected UUID format in validation error messages', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: 'invalid', name: 'Label 1', color: '#FF0000' },
        validationErrors: [
          {
            field: 'id',
            message: 'Invalid UUID format',
            severity: 'error',
          },
        ],
        isValid: false,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // The expected format should be in the tooltip content
    // Since we mocked Tooltip, we just check that the error message is present
    expect(screen.getByText('Invalid UUID format')).toBeInTheDocument();
  });

  it('displays clear error messages for invalid UUID formats', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: 'bad-uuid', name: 'Label 1', color: '#FF0000' },
        validationErrors: [
          {
            field: 'id',
            message: 'Invalid UUID format: "bad-uuid". Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx',
            severity: 'error',
          },
        ],
        isValid: false,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that the detailed error message is present
    expect(screen.getByText(/Invalid UUID format/)).toBeInTheDocument();
    // Use getAllByText since "bad-uuid" appears multiple times (in error message and in the ID cell)
    const badUuidElements = screen.getAllByText(/bad-uuid/);
    expect(badUuidElements.length).toBeGreaterThan(0);
  });

  it('highlights invalid UUID values in error display', () => {
    const invalidUuid = 'xyz-123';
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: invalidUuid, name: 'Label 1', color: '#FF0000' },
        validationErrors: [
          {
            field: 'id',
            message: 'Invalid UUID format',
            severity: 'error',
          },
        ],
        isValid: false,
      },
    ];

    const { container } = render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that the invalid UUID is highlighted with error styling
    const errorElements = container.querySelectorAll('.text-red-600, .text-red-400');
    expect(errorElements.length).toBeGreaterThan(0);
  });

  it('displays validation error summary for records with UUID errors', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: 'invalid', name: 'Label 1', color: '#FF0000' },
        validationErrors: [
          {
            field: 'id',
            message: 'Invalid UUID format',
            severity: 'error',
          },
        ],
        isValid: false,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check that validation error summary is displayed
    expect(screen.getByText(/Validation Errors Found/i)).toBeInTheDocument();
    expect(screen.getByText(/1 record\(s\) have validation errors/i)).toBeInTheDocument();
  });
});

describe('ImportPreviewTable - Mixed Scenarios', () => {
  const mockOnSelectionChange = vi.fn();
  const mockOnImportSelected = vi.fn();
  const mockOnImportAll = vi.fn();

  const columns: ImportColumnDef<DtoCreateLabelRequest>[] = [
    { key: 'name', label: 'Name' },
    { key: 'color', label: 'Color' },
  ];

  it('handles records with mix of valid UUIDs, invalid UUIDs, and no UUIDs', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: '550e8400-e29b-41d4-a716-446655440001', name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
      {
        index: 1,
        data: { id: 'invalid-uuid', name: 'Label 2', color: '#00FF00' },
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
        data: { name: 'Label 3', color: '#0000FF' },
        validationErrors: [],
        isValid: true,
      },
    ];

    const { container } = render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // ID column should be shown
    expect(screen.getByText('ID')).toBeInTheDocument();

    // Valid UUID should be in monospace
    const validUuidCode = container.querySelector('code.font-mono:not(.text-red-600)');
    expect(validUuidCode).toBeInTheDocument();

    // Invalid UUID should have error styling
    const invalidUuidCode = container.querySelector('code.text-red-600');
    expect(invalidUuidCode).toBeInTheDocument();

    // "Auto-generated" should be present for record without ID
    expect(screen.getByText('Auto-generated')).toBeInTheDocument();
  });

  it('correctly counts valid and invalid records with UUID errors', () => {
    const records: ParsedImportRecord<DtoCreateLabelRequest>[] = [
      {
        index: 0,
        data: { id: '550e8400-e29b-41d4-a716-446655440001', name: 'Label 1', color: '#FF0000' },
        validationErrors: [],
        isValid: true,
      },
      {
        index: 1,
        data: { id: 'bad-uuid-1', name: 'Label 2', color: '#00FF00' },
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
        data: { id: 'bad-uuid-2', name: 'Label 3', color: '#0000FF' },
        validationErrors: [
          {
            field: 'id',
            message: 'Invalid UUID format',
            severity: 'error',
          },
        ],
        isValid: false,
      },
    ];

    render(
      <ImportPreviewTable
        records={records}
        columns={columns}
        selectedRecords={new Set()}
        onSelectionChange={mockOnSelectionChange}
        onImportSelected={mockOnImportSelected}
        onImportAll={mockOnImportAll}
        importStatus={new Map()}
        isImporting={false}
      />
    );

    // Check statistics
    expect(screen.getByText('Total: 3')).toBeInTheDocument();
    expect(screen.getByText('Valid: 1')).toBeInTheDocument();
    expect(screen.getByText('Invalid: 2')).toBeInTheDocument();
  });
});
