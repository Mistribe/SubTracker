import { Button } from "@/components/ui/button";
import { argbToRgba } from "@/components/ui/utils/color-utils";
import { Loader2, Pencil, X } from "lucide-react";
import Label from "@/models/label";

interface LabelItemProps {
  label: Label;
  onEdit?: (label: Label) => void;
  onDelete?: (id: string, label?: Label) => void;
  isDeleting?: boolean;
  isReadOnly?: boolean;
}

export const LabelItem = ({
  label,
  onEdit,
  onDelete,
  isDeleting = false,
  isReadOnly = false
}: LabelItemProps) => {
  return (
    <div
      className="flex items-center p-2 border rounded-md hover:bg-muted/50"
    >
      <div
        className="w-3 h-3 rounded-full mr-2 flex-shrink-0"
        style={{ backgroundColor: argbToRgba(label.color) }}
      />
      <span className="text-sm font-medium truncate flex-grow">{label.name}</span>
      
      {!isReadOnly && (
        <div className="flex items-center gap-1 ml-1">
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6"
            onClick={() => onEdit?.(label)}
            title="Edit label"
          >
            <Pencil className="h-3 w-3" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6 text-destructive hover:text-destructive"
            onClick={() => onDelete?.(label.id, label)}
            disabled={isDeleting}
            title="Delete label"
          >
            {isDeleting ? (
              <Loader2 className="h-3 w-3 animate-spin" />
            ) : (
              <X className="h-3 w-3" />
            )}
          </Button>
        </div>
      )}
    </div>
  );
};