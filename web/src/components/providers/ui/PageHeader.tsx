import { Button } from "@/components/ui/button";

interface PageHeaderProps {
  title: string;
  onAddProvider: () => void;
}

export const PageHeader = ({ title, onAddProvider }: PageHeaderProps) => {
  return (
    <div className="flex justify-between items-center mb-6">
      <h1 className="text-3xl font-bold">{title}</h1>
      <Button onClick={onAddProvider}>Add Provider</Button>
    </div>
  );
};