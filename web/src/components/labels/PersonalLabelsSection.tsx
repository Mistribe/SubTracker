import Label from "@/models/label";
import { LabelItem } from "./LabelItem";
import { EditableLabelItem } from "./EditableLabelItem";
import { AddLabelForm } from "./AddLabelForm";
import { OwnerType } from "@/models/ownerType";
import { Loader2 } from "lucide-react";

interface PersonalLabelsSectionProps {
  labels: Label[];
  editingId: string | null;
  editingName: string;
  editingColor: string;
  onStartEditing: (label: Label) => void;
  onSaveEdit: (id: string, name: string, color: string) => void;
  onCancelEdit: () => void;
  onDeleteLabel: (id: string) => void;
  onAddLabel: (name: string, color: string, ownerType?: OwnerType, familyId?: string) => void;
  isUpdating: boolean;
  isDeletingId: string | null;
  isAdding: boolean;
  isLoading?: boolean;
  error?: Error | null;
}

export const PersonalLabelsSection = ({
  labels,
  editingId,
  editingName,
  editingColor,
  onStartEditing,
  onSaveEdit,
  onCancelEdit,
  onDeleteLabel,
  onAddLabel,
  isUpdating,
  isDeletingId,
  isAdding,
  isLoading = false,
  error = null
}: PersonalLabelsSectionProps) => {
  // No need to filter labels anymore since we're getting only personal labels from the query

  if (isLoading) {
    return (
      <div>
        <h2 className="text-xl font-semibold mb-4">Personal Labels</h2>
        <div className="flex flex-col items-center justify-center h-32">
          <Loader2 className="h-6 w-6 animate-spin text-primary mb-2"/>
          <p className="text-muted-foreground text-sm">Loading personal labels...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div>
        <h2 className="text-xl font-semibold mb-4">Personal Labels</h2>
        <div className="p-4 border rounded-md bg-destructive/10">
          <p className="text-destructive text-sm">Error loading personal labels</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">Personal Labels</h2>
      
      <div className="mb-4">
        <AddLabelForm 
          onAddLabel={onAddLabel}
          isAdding={isAdding}
          ownerType={OwnerType.Personal}
          title="Add Personal Label"
        />
      </div>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {labels.length === 0 ? (
          <p className="text-muted-foreground col-span-full text-center py-4">No personal labels found</p>
        ) : (
          labels.map((label) => (
            <div
              key={label.id}
              className={`border rounded-md ${editingId === label.id ? 'bg-muted' : 'hover:bg-muted/50'}`}
            >
              {editingId === label.id ? (
                <EditableLabelItem
                  label={label}
                  initialName={editingName}
                  initialColor={editingColor}
                  onSave={onSaveEdit}
                  onCancel={onCancelEdit}
                  isSaving={isUpdating}
                />
              ) : (
                <LabelItem
                  label={label}
                  onEdit={onStartEditing}
                  onDelete={onDeleteLabel}
                  isDeleting={isDeletingId === label.id}
                />
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};