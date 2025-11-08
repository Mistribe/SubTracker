import { useCallback, useState } from "react";
import { Upload, FileText, AlertCircle, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface FileUploadZoneProps {
  onFileSelected: (file: File) => void;
  acceptedFormats: string[]; // ['.csv', '.json', '.yaml', '.yml']
  isLoading?: boolean;
  error?: string;
}

export function FileUploadZone({
  onFileSelected,
  acceptedFormats,
  isLoading = false,
  error,
}: FileUploadZoneProps) {
  const [isDragOver, setIsDragOver] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  // Validate file format
  const validateFileFormat = (file: File): boolean => {
    const fileName = file.name.toLowerCase();
    return acceptedFormats.some((format) => fileName.endsWith(format));
  };

  // Handle file selection
  const handleFileSelect = useCallback(
    (file: File) => {
      if (!validateFileFormat(file)) {
        return;
      }
      setSelectedFile(file);
      onFileSelected(file);
    },
    [acceptedFormats, onFileSelected]
  );

  // Handle drag events
  const handleDragEnter = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(false);
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragOver(false);

      const files = e.dataTransfer.files;
      if (files && files.length > 0) {
        handleFileSelect(files[0]);
      }
    },
    [handleFileSelect]
  );

  // Handle file input change
  const handleFileInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = e.target.files;
      if (files && files.length > 0) {
        handleFileSelect(files[0]);
      }
    },
    [handleFileSelect]
  );

  // Handle click to open file picker
  const handleClick = () => {
    if (!isLoading) {
      document.getElementById("file-upload-input")?.click();
    }
  };

  // Handle keyboard interaction
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if ((e.key === "Enter" || e.key === " ") && !isLoading) {
      e.preventDefault();
      handleClick();
    }
  };

  // Format file size
  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + " " + sizes[i];
  };

  return (
    <div className="w-full">
      <div
        role="button"
        tabIndex={isLoading ? -1 : 0}
        aria-label="Upload file zone. Click or drag and drop a file to upload."
        aria-disabled={isLoading}
        className={cn(
          "relative flex flex-col items-center justify-center w-full min-h-[300px] px-6 py-10",
          "border-2 border-dashed rounded-lg transition-all duration-200",
          "focus:outline-none focus-visible:ring-4 focus-visible:ring-primary focus-visible:ring-offset-4 focus-visible:ring-offset-background",
          isDragOver && !isLoading
            ? "border-primary bg-primary/10 scale-[1.02]"
            : "border-border bg-background",
          !isLoading && !isDragOver && "hover:border-primary/60 hover:bg-accent/60 cursor-pointer",
          isLoading && "opacity-60 cursor-not-allowed",
          error && "border-destructive bg-destructive/10"
        )}
        onDragEnter={handleDragEnter}
        onDragLeave={handleDragLeave}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
      >
        <input
          id="file-upload-input"
          type="file"
          className="hidden"
          accept={acceptedFormats.join(",")}
          onChange={handleFileInputChange}
          disabled={isLoading}
          aria-label="File input"
        />

        {/* Icon */}
        <div className="mb-4" role="img" aria-label={
          isLoading ? "Loading spinner" :
          error ? "Error icon" :
          selectedFile ? "File selected icon" :
          "Upload icon"
        }>
          {isLoading ? (
            <Loader2 className="h-12 w-12 text-primary animate-spin" aria-hidden="true" />
          ) : error ? (
            <AlertCircle className="h-12 w-12 text-destructive" aria-hidden="true" />
          ) : selectedFile ? (
            <FileText className="h-12 w-12 text-primary" aria-hidden="true" />
          ) : (
            <Upload
              className={cn(
                "h-12 w-12 transition-colors",
                isDragOver ? "text-primary" : "text-muted-foreground"
              )}
              aria-hidden="true"
            />
          )}
        </div>

        {/* Text content */}
        <div className="text-center space-y-2">
          {isLoading ? (
            <>
              <p className="text-sm font-medium text-foreground">Processing file...</p>
              <p className="text-xs text-muted-foreground">Please wait while we parse your file</p>
            </>
          ) : error ? (
            <>
              <p className="text-sm font-semibold text-destructive">Failed to parse file</p>
              <p className="text-xs text-destructive font-medium max-w-md mx-auto whitespace-pre-wrap">{error}</p>
              <p className="text-xs text-foreground mt-3 font-medium">
                Click or drag another file to try again
              </p>
            </>
          ) : selectedFile ? (
            <>
              <p className="text-sm font-medium text-foreground">{selectedFile.name}</p>
              <p className="text-xs text-muted-foreground">{formatFileSize(selectedFile.size)}</p>
              <p className="text-xs text-muted-foreground mt-2">
                Click or drag another file to replace
              </p>
            </>
          ) : (
            <>
              <p className="text-sm font-medium text-foreground">
                {isDragOver ? "Drop file here" : "Drag and drop your file here"}
              </p>
              <p className="text-xs text-muted-foreground">or click to browse</p>
              <p className="text-xs text-muted-foreground mt-2">
                Supported formats: {acceptedFormats.join(", ")}
              </p>
            </>
          )}
        </div>

        {/* Drag overlay */}
        {isDragOver && !isLoading && (
          <div
            className="absolute inset-0 bg-primary/10 rounded-lg pointer-events-none"
            aria-hidden="true"
          />
        )}
      </div>
    </div>
  );
}
