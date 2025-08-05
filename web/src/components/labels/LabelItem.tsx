import { Button } from "@/components/ui/button";
import { argbToRgba } from "@/components/ui/utils/color-utils";
import { Loader2, Pencil, X } from "lucide-react";
import Label from "@/models/label";
import { OwnerType } from "@/models/ownerType";
import { Badge } from "@/components/ui/badge";

interface LabelItemProps {
  label: Label;
  onEdit?: (label: Label) => void;
  onDelete?: (id: string, label?: Label) => void;
  isDeleting?: boolean;
  isReadOnly?: boolean;
  ownerName?: string; // For family labels, pass the family name
}

export const LabelItem = ({
  label,
  onEdit,
  onDelete,
  isDeleting = false,
  isReadOnly = false,
  ownerName
}: LabelItemProps) => {
  // Determine the badge text based on owner type
  const getBadgeText = () => {
    switch (label.owner.type) {
      case OwnerType.System:
        return 'System';
      case OwnerType.Personal:
        return 'Personal';
      case OwnerType.Family:
        return ownerName || 'Family';
      default:
        return label.owner.type;
    }
  };

  // Determine badge variant based on owner type
  const getBadgeVariant = (): 'default' | 'secondary' | 'destructive' | 'outline' => {
    switch (label.owner.type) {
      case OwnerType.System:
        return 'secondary';
      case OwnerType.Personal:
        return 'default';
      case OwnerType.Family:
        return 'outline';
      default:
        return 'secondary';
    }
  };

  return (
    <div
      className="flex flex-col p-2 border rounded-md hover:bg-muted/50"
    >
      <div className="flex items-center mb-1">
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
      
      <Badge variant={getBadgeVariant()} className="self-start text-xs">
        {getBadgeText()}
      </Badge>
    </div>
  );
};