import Label from "@/models/label";
import { LabelItem } from "./LabelItem";
import { EditableLabelItem } from "./EditableLabelItem";
import { AddLabelForm } from "./AddLabelForm";
import { OwnerType } from "@/models/ownerType";

interface FamilyLabelsSectionProps {
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
}

export const FamilyLabelsSection = ({
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
  isAdding
}: FamilyLabelsSectionProps) => {
  // Filter to only include family labels
  const familyLabels = labels.filter(label => label.owner.isFamily);
  const hasFamilyLabels = familyLabels.length > 0;

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">Family Labels</h2>
      
      {!hasFamilyLabels ? (
        <div className="p-4 border rounded-md bg-muted/50">
          <p className="text-center text-muted-foreground mb-2">
            You don't have a family yet.
          </p>
          <p className="text-center text-muted-foreground">
            You need to create a family first or join one to manage family labels.
          </p>
        </div>
      ) : (
        <>
          <div className="mb-4">
            <AddLabelForm 
              onAddLabel={onAddLabel}
              isAdding={isAdding}
              ownerType={OwnerType.Family}
              title="Add Family Label"
            />
          </div>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {familyLabels.map((label) => (
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
            ))}
          </div>
        </>
      )}
    </div>
  );
};