import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { FileUploadZone } from "../FileUploadZone";

describe("FileUploadZone - Integration Tests", () => {
  const mockOnFileSelected = vi.fn();
  const acceptedFormats = [".csv", ".json", ".yaml", ".yml"];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Drag-and-drop functionality", () => {
    it("should handle drag enter and show visual feedback", () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      
      // Simulate drag enter
      fireEvent.dragEnter(uploadZone, {
        dataTransfer: {
          files: [],
        },
      });

      // Check for visual feedback (the component should have drag-over state)
      expect(uploadZone).toHaveClass("border-primary");
    });

    it("should handle drag leave and remove visual feedback", () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      
      // Simulate drag enter then drag leave
      fireEvent.dragEnter(uploadZone);
      fireEvent.dragLeave(uploadZone);

      // Visual feedback should be removed
      expect(uploadZone).not.toHaveClass("border-primary");
    });

    it("should accept valid CSV file via drag and drop", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["name,color\nTest,#FF0000"], "test.csv", {
        type: "text/csv",
      });

      // Simulate drag and drop
      fireEvent.dragEnter(uploadZone);
      fireEvent.dragOver(uploadZone);
      fireEvent.drop(uploadZone, {
        dataTransfer: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should accept valid JSON file via drag and drop", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(['[{"name":"Test"}]'], "test.json", {
        type: "application/json",
      });

      fireEvent.drop(uploadZone, {
        dataTransfer: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should accept valid YAML file via drag and drop", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["name: Test\ncolor: '#FF0000'"], "test.yaml", {
        type: "text/yaml",
      });

      fireEvent.drop(uploadZone, {
        dataTransfer: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should accept .yml extension via drag and drop", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["name: Test"], "test.yml", {
        type: "text/yaml",
      });

      fireEvent.drop(uploadZone, {
        dataTransfer: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should reject invalid file format via drag and drop", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["test content"], "test.txt", {
        type: "text/plain",
      });

      fireEvent.drop(uploadZone, {
        dataTransfer: {
          files: [file],
        },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).not.toHaveBeenCalled();
      });
    });

    it("should handle multiple files and only accept the first one", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file1 = new File(["data1"], "test1.csv", { type: "text/csv" });
      const file2 = new File(["data2"], "test2.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: {
          files: [file1, file2],
        },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledTimes(1);
        expect(mockOnFileSelected).toHaveBeenCalledWith(file1);
      });
    });
  });

  describe("File picker functionality", () => {
    it("should open file picker on click", async () => {
      const user = userEvent.setup();
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      await user.click(uploadZone);

      // File input should exist
      const fileInput = document.getElementById("file-upload-input");
      expect(fileInput).toBeInTheDocument();
    });

    it("should accept valid file through file picker", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const fileInput = document.getElementById("file-upload-input") as HTMLInputElement;
      const file = new File(["name,color\nTest,#FF0000"], "test.csv", {
        type: "text/csv",
      });

      // Simulate file selection
      Object.defineProperty(fileInput, "files", {
        value: [file],
        writable: false,
      });

      fireEvent.change(fileInput);

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should display selected file name and size", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const fileInput = document.getElementById("file-upload-input") as HTMLInputElement;
      const file = new File(["name,color\nTest,#FF0000"], "test-file.csv", {
        type: "text/csv",
      });

      Object.defineProperty(fileInput, "files", {
        value: [file],
        writable: false,
      });

      fireEvent.change(fileInput);

      await waitFor(() => {
        expect(screen.getByText("test-file.csv")).toBeInTheDocument();
        expect(screen.getByText(/bytes/i)).toBeInTheDocument();
      });
    });

    it("should not open file picker when loading", async () => {
      const user = userEvent.setup();
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
          isLoading={true}
        />
      );

      const uploadZone = screen.getByRole("button");
      await user.click(uploadZone);

      // File input should be disabled
      const fileInput = document.getElementById("file-upload-input") as HTMLInputElement;
      expect(fileInput).toBeDisabled();
    });
  });

  describe("Format detection and validation", () => {
    it("should validate CSV file extension", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["data"], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should validate JSON file extension", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["{}"], "test.json", { type: "application/json" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should validate YAML file extension", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["key: value"], "test.yaml", { type: "text/yaml" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should reject file with unsupported extension", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["data"], "test.xml", { type: "text/xml" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).not.toHaveBeenCalled();
      });
    });

    it("should be case-insensitive for file extensions", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["data"], "test.CSV", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should display supported formats in the UI", () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      expect(screen.getByText(/supported formats/i)).toBeInTheDocument();
      expect(screen.getByText(/\.csv, \.json, \.yaml, \.yml/i)).toBeInTheDocument();
    });
  });

  describe("File size limit handling", () => {
    it("should accept file within size limit", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      // Create a 1MB file
      const content = "a".repeat(1024 * 1024);
      const file = new File([content], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should display file size in human-readable format", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      // Create a file with known size
      const content = "a".repeat(1024); // 1KB
      const file = new File([content], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(screen.getByText(/1 KB/i)).toBeInTheDocument();
      });
    });

    it("should format bytes correctly", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["test"], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(screen.getByText(/bytes/i)).toBeInTheDocument();
      });
    });

    it("should format megabytes correctly", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      // Create a 2MB file
      const content = "a".repeat(2 * 1024 * 1024);
      const file = new File([content], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(screen.getByText(/2 MB/i)).toBeInTheDocument();
      });
    });
  });

  describe("Loading and error states", () => {
    it("should show loading state", () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
          isLoading={true}
        />
      );

      expect(screen.getByText(/processing file/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/loading spinner/i)).toBeInTheDocument();
    });

    it("should show error state with message", () => {
      const errorMessage = "Failed to parse CSV file";
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
          error={errorMessage}
        />
      );

      expect(screen.getByText(/failed to parse file/i)).toBeInTheDocument();
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });

    it("should allow file replacement after error", async () => {
      const { rerender } = render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
          error="Parse error"
        />
      );

      // Clear error and allow new file
      rerender(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["data"], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(mockOnFileSelected).toHaveBeenCalledWith(file);
      });
    });

    it("should prevent interaction during loading", () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
          isLoading={true}
        />
      );

      const uploadZone = screen.getByRole("button");
      expect(uploadZone).toHaveAttribute("aria-disabled", "true");
      expect(uploadZone).toHaveAttribute("tabIndex", "-1");
    });
  });

  describe("File replacement", () => {
    it("should allow replacing selected file with new file", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      
      // Upload first file
      const file1 = new File(["data1"], "test1.csv", { type: "text/csv" });
      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file1] },
      });

      await waitFor(() => {
        expect(screen.getByText("test1.csv")).toBeInTheDocument();
      });

      // Upload second file to replace
      const file2 = new File(["data2"], "test2.csv", { type: "text/csv" });
      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file2] },
      });

      await waitFor(() => {
        expect(screen.getByText("test2.csv")).toBeInTheDocument();
        expect(mockOnFileSelected).toHaveBeenCalledTimes(2);
        expect(mockOnFileSelected).toHaveBeenLastCalledWith(file2);
      });
    });

    it("should show replacement hint when file is selected", async () => {
      render(
        <FileUploadZone
          onFileSelected={mockOnFileSelected}
          acceptedFormats={acceptedFormats}
        />
      );

      const uploadZone = screen.getByRole("button");
      const file = new File(["data"], "test.csv", { type: "text/csv" });

      fireEvent.drop(uploadZone, {
        dataTransfer: { files: [file] },
      });

      await waitFor(() => {
        expect(screen.getByText(/click or drag another file to replace/i)).toBeInTheDocument();
      });
    });
  });
});
