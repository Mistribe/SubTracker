import { Button } from "@/components/ui/button";

interface NoMoreProvidersProps {
  show: boolean;
}

export const NoMoreProviders = ({ show }: NoMoreProvidersProps) => {
  if (!show) return null;
  
  return (
    <p className="text-center text-muted-foreground mt-8 pb-8">
      No more providers to load
    </p>
  );
};

interface NoProvidersProps {
  show: boolean;
  onAddProvider: () => void;
}

export const NoProviders = ({ show, onAddProvider }: NoProvidersProps) => {
  if (!show) return null;
  
  return (
    <div className="text-center py-10">
      <p className="text-muted-foreground">No providers found</p>
      <Button variant="outline" className="mt-4" onClick={onAddProvider}>
        Add Your First Provider
      </Button>
    </div>
  );
};