import { Button } from "@/components/ui/button";

interface ErrorStateProps {
  message?: string;
  onRetry?: () => void;
}

export const ErrorState = ({
  message = "Failed to load providers. Please try again later.",
  onRetry = () => window.location.reload()
}: ErrorStateProps) => {
  return (
    <div className="text-center py-10">
      <p className="text-red-500">{message}</p>
      <Button variant="outline" className="mt-4" onClick={onRetry}>
        Retry
      </Button>
    </div>
  );
};