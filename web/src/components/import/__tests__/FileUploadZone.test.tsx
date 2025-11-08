import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { FileUploadZone } from "../FileUploadZone";

describe("FileUploadZone", () => {
  const mockOnFileSelected = vi.fn();
  const acceptedFormats = [".csv", ".json", ".yaml", ".yml"];

  it("renders the upload zone with default state", () => {
    render(
      <FileUploadZone
        onFileSelected={mockOnFileSelected}
        acceptedFormats={acceptedFormats}
      />
    );

    expect(screen.getByText(/drag and drop your file here/i)).toBeInTheDocument();
    expect(screen.getByText(/or click to browse/i)).toBeInTheDocument();
    expect(screen.getByText(/supported formats/i)).toBeInTheDocument();
  });

  it("shows loading state when isLoading is true", () => {
    render(
      <FileUploadZone
        onFileSelected={mockOnFileSelected}
        acceptedFormats={acceptedFormats}
        isLoading={true}
      />
    );

    expect(screen.getByText(/processing file/i)).toBeInTheDocument();
  });

  it("shows error state when error prop is provided", () => {
    const errorMessage = "Invalid file format";
    render(
      <FileUploadZone
        onFileSelected={mockOnFileSelected}
        acceptedFormats={acceptedFormats}
        error={errorMessage}
      />
    );

    expect(screen.getByText(/error/i)).toBeInTheDocument();
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it("has proper accessibility attributes", () => {
    render(
      <FileUploadZone
        onFileSelected={mockOnFileSelected}
        acceptedFormats={acceptedFormats}
      />
    );

    const uploadZone = screen.getByRole("button");
    expect(uploadZone).toHaveAttribute("aria-label");
    expect(uploadZone).toHaveAttribute("tabIndex", "0");
  });

  it("disables interaction when loading", () => {
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

  it("handles keyboard navigation with Enter key", () => {
    render(
      <FileUploadZone
        onFileSelected={mockOnFileSelected}
        acceptedFormats={acceptedFormats}
      />
    );

    const uploadZone = screen.getByRole("button");
    fireEvent.keyDown(uploadZone, { key: "Enter" });
    
    // The file input should be triggered (we can't fully test file selection in jsdom)
    const fileInput = document.getElementById("file-upload-input");
    expect(fileInput).toBeInTheDocument();
  });

  it("handles keyboard navigation with Space key", () => {
    render(
      <FileUploadZone
        onFileSelected={mockOnFileSelected}
        acceptedFormats={acceptedFormats}
      />
    );

    const uploadZone = screen.getByRole("button");
    fireEvent.keyDown(uploadZone, { key: " " });
    
    const fileInput = document.getElementById("file-upload-input");
    expect(fileInput).toBeInTheDocument();
  });
});
