import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ColorPicker } from "@/components/ui/color-picker";
import { argbToRgba } from "@/components/ui/utils/color-utils";
import { Loader2, PlusIcon } from "lucide-react";
import { OwnerType } from "@/models/ownerType";

interface AddLabelFormProps {
  onAddLabel: (name: string, color: string, ownerType?: OwnerType, familyId?: string) => void;
  isAdding: boolean;
  ownerType?: OwnerType;
  familyId?: string;
  title?: string;
}

export const AddLabelForm = ({ 
  onAddLabel, 
  isAdding, 
  ownerType, 
  familyId,
  title = "Add Label" 
}: AddLabelFormProps) => {
  const [newLabel, setNewLabel] = useState("");
  const [newLabelColor, setNewLabelColor] = useState("#FF000000"); // Default ARGB color

  const handleAddLabel = () => {
    if (newLabel.trim()) {
      onAddLabel(newLabel, newLabelColor, ownerType, familyId);
      setNewLabel(""); // Clear the input after adding
    }
  };

  return (
    <div className="w-full">
      <h3 className="text-sm font-medium mb-2">{title}</h3>
      <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center">
        <div className="flex items-center gap-2 w-full sm:w-auto">
          <div
            className="w-4 h-4 rounded-full flex-shrink-0"
            style={{ backgroundColor: argbToRgba(newLabelColor) }}
          />
          <Input
            placeholder="Enter label name"
            value={newLabel}
            onChange={(e) => setNewLabel(e.target.value)}
            className="h-9 text-sm w-full sm:w-48 md:w-64"
          />
        </div>
        <div className="flex items-center gap-2">
          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" size="sm" className="h-9">
                <div
                  className="w-4 h-4 rounded-md mr-1"
                  style={{ backgroundColor: argbToRgba(newLabelColor) }}
                />
                Color
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="start">
              <ColorPicker
                color={newLabelColor}
                onChange={setNewLabelColor}
              />
            </PopoverContent>
          </Popover>
          <Button
            size="sm"
            className="h-9"
            onClick={handleAddLabel}
            disabled={isAdding}
          >
            {isAdding ? (
              <Loader2 className="h-4 w-4 mr-1 animate-spin" />
            ) : (
              <PlusIcon className="h-4 w-4 mr-1" />
            )}
            Add
          </Button>
        </div>
      </div>
    </div>
  );
};