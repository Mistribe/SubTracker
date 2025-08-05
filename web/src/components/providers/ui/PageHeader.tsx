import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import Label from "@/models/label";
import { X } from "lucide-react";

interface PageHeaderProps {
  title: string;
  onAddProvider: () => void;
  searchText?: string;
  onSearchChange?: (value: string) => void;
  selectedLabels?: Label[];
  availableLabels?: Label[];
  onLabelSelect?: (label: Label) => void;
  onLabelRemove?: (labelId: string) => void;
}

export const PageHeader = ({ 
  title, 
  onAddProvider, 
  searchText = "", 
  onSearchChange,
  selectedLabels = [],
  availableLabels = [],
  onLabelSelect,
  onLabelRemove
}: PageHeaderProps) => {
  return (
    <div className="flex flex-col gap-4 mb-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">{title}</h1>
        <Button onClick={onAddProvider}>Add Provider</Button>
      </div>
      
      {onSearchChange && (
        <div className="flex justify-center">
          <Input
            placeholder="Search providers, labels..."
            value={searchText}
            onChange={(e) => onSearchChange(e.target.value)}
            className="max-w-md"
          />
        </div>
      )}
      
      {onLabelSelect && availableLabels.length > 0 && (
        <div className="flex flex-col gap-2">
          <div className="flex justify-center">
            <div className="flex flex-wrap gap-2 max-w-2xl justify-center">
              {availableLabels
                .filter(label => !selectedLabels.some(selected => selected.id === label.id))
                .map(label => (
                  <Badge
                    key={label.id}
                    variant="outline"
                    className="cursor-pointer hover:bg-secondary/80"
                    style={{ backgroundColor: label.color || undefined }}
                    onClick={() => onLabelSelect(label)}
                  >
                    {label.name}
                  </Badge>
                ))
              }
            </div>
          </div>
          
          {selectedLabels.length > 0 && (
            <div className="flex justify-center">
              <div className="flex flex-wrap gap-2 max-w-2xl justify-center">
                {selectedLabels.map(label => (
                  <Badge
                    key={label.id}
                    variant="secondary"
                    className="cursor-pointer flex items-center gap-1"
                    style={{ backgroundColor: label.color || undefined }}
                  >
                    {label.name}
                    <X 
                      className="h-3 w-3 hover:text-destructive" 
                      onClick={() => onLabelRemove?.(label.id)}
                    />
                  </Badge>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};