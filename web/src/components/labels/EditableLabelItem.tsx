import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ColorPicker } from "@/components/ui/color-picker";
import { argbToRgba } from "@/components/ui/utils/color-utils";
import { Loader2 } from "lucide-react";
import { useState, useEffect } from "react";
import Label from "@/models/label";

interface EditableLabelItemProps {
  label: Label;
  onSave: (id: string, name: string, color: string) => void;
  onCancel: () => void;
  isSaving?: boolean;
  initialColor?: string;
  initialName?: string;
}

export const EditableLabelItem = ({
  label,
  onSave,
  onCancel,
  isSaving = false,
  initialColor,
  initialName
}: EditableLabelItemProps) => {
  const [name, setName] = useState(initialName || label.name);
  const [color, setColor] = useState(initialColor || label.color);

  // Update state if props change
  useEffect(() => {
    if (initialName) setName(initialName);
    if (initialColor) setColor(initialColor);
  }, [initialName, initialColor]);

  const handleSave = () => {
    if (name.trim()) {
      onSave(label.id, name, color);
    }
  };

  return (
    <div className="w-full space-y-2">
      <div className="flex items-center gap-2">
        <div
          className="w-4 h-4 rounded-full flex-shrink-0"
          style={{ backgroundColor: argbToRgba(color) }}
        />
        <Input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="h-7 text-sm"
          autoFocus
        />
      </div>
      <div className="flex items-center gap-1 justify-between">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" size="sm" className="h-7 px-2">
              <div
                className="w-4 h-4 rounded-md mr-1"
                style={{ backgroundColor: argbToRgba(color) }}
              />
              Color
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0" align="start">
            <ColorPicker
              color={color}
              onChange={setColor}
            />
          </PopoverContent>
        </Popover>
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="sm"
            className="h-7 px-2"
            onClick={handleSave}
            disabled={isSaving}
          >
            {isSaving ? (
              <Loader2 className="h-3 w-3 mr-1 animate-spin" />
            ) : null}
            Save
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="h-7 px-2"
            onClick={onCancel}
          >
            Cancel
          </Button>
        </div>
      </div>
    </div>
  );
};