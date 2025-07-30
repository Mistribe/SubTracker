import { Loader2 } from "lucide-react";

interface FamiliesLoadingErrorProps {
  isLoading: boolean;
  error: unknown;
}

export const FamiliesLoadingError = ({ isLoading, error }: FamiliesLoadingErrorProps) => {
  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin" />
        <span className="ml-2">Loading families...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong className="font-bold">Error!</strong>
        <span className="block sm:inline"> {error instanceof Error ? error.message : 'Failed to load families'}</span>
      </div>
    );
  }

  return null;
};