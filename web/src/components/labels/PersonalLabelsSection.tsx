import Label from "@/models/label";
import { LabelItem } from "./LabelItem";
import { EditableLabelItem } from "./EditableLabelItem";
import { AddLabelForm } from "./AddLabelForm";
import { OwnerType } from "@/models/ownerType";

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
  isAdding
}: PersonalLabelsSectionProps) => {
  // Filter to only include personal labels
  const personalLabels = labels.filter(label => label.owner.isPersonal);

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
        {personalLabels.map((label) => (
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
    </div>
  );
};