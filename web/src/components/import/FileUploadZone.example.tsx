/**
 * Example usage of FileUploadZone component
 * 
 * This file demonstrates how to use the FileUploadZone component
 * in different scenarios.
 */

import { useState } from "react";
import { FileUploadZone } from "./FileUploadZone";

export function FileUploadZoneExample() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>();

  const handleFileSelected = async (file: File) => {
    setError(undefined);
    setIsLoading(true);

    try {
      // Simulate file processing
      await new Promise((resolve) => setTimeout(resolve, 2000));
      
      // Validate file format
      const validFormats = [".csv", ".json", ".yaml", ".yml"];
      const isValid = validFormats.some((format) => 
        file.name.toLowerCase().endsWith(format)
      );

      if (!isValid) {
        setError(`Invalid file format. Please upload a CSV, JSON, or YAML file.`);
        return;
      }

      console.log("File selected:", file.name);
      // Process the file here...
      
    } catch (err) {
      setError("Failed to process file. Please try again.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h2 className="text-2xl font-bold mb-4">Import Data</h2>
      <p className="text-muted-foreground mb-6">
        Upload a CSV, JSON, or YAML file to import your data.
      </p>
      
      <FileUploadZone
        onFileSelected={handleFileSelected}
        acceptedFormats={[".csv", ".json", ".yaml", ".yml"]}
        isLoading={isLoading}
        error={error}
      />
    </div>
  );
}
