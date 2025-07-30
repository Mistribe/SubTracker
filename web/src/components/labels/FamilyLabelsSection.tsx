import Label from "@/models/label";
import { LabelItem } from "./LabelItem";
import { EditableLabelItem } from "./EditableLabelItem";
import { AddLabelForm } from "./AddLabelForm";
import { OwnerType } from "@/models/ownerType";
import { Loader2 } from "lucide-react";
import Family from "@/models/family";
import { useLabelsQuery } from "@/hooks/labels/useLabelsQuery";
import { useFamiliesQuery } from "@/hooks/families/useFamiliesQuery";
import { useState } from "react";

interface FamilyLabelsSectionProps {
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
  // Fetch all families
  const {
    data: familiesResponse,
    isLoading: isLoadingFamilies,
    error: familiesError
  } = useFamiliesQuery();
  
  const families = familiesResponse?.families || [];
  const hasFamilies = families.length > 0;
  
  // State to track which family's add form is open
  const [activeFamilyId, setActiveFamilyId] = useState<string | null>(null);
  
  if (isLoadingFamilies) {
    return (
      <div>
        <h2 className="text-xl font-semibold mb-4">Family Labels</h2>
        <div className="flex flex-col items-center justify-center h-32">
          <Loader2 className="h-6 w-6 animate-spin text-primary mb-2"/>
          <p className="text-muted-foreground text-sm">Loading families...</p>
        </div>
      </div>
    );
  }

  if (familiesError) {
    return (
      <div>
        <h2 className="text-xl font-semibold mb-4">Family Labels</h2>
        <div className="p-4 border rounded-md bg-destructive/10">
          <p className="text-destructive text-sm">Error loading families</p>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">Family Labels</h2>
      
      {!hasFamilies ? (
        <div className="p-4 border rounded-md bg-muted/50">
          <p className="text-center text-muted-foreground mb-2">
            You don't have a family yet.
          </p>
          <p className="text-center text-muted-foreground">
            You need to create a family first or join one to manage family labels.
          </p>
        </div>
      ) : (
        <div className="space-y-8">
          {families.map((family) => (
            <FamilyLabelsGroup
              key={family.id}
              family={family}
              editingId={editingId}
              editingName={editingName}
              editingColor={editingColor}
              onStartEditing={onStartEditing}
              onSaveEdit={onSaveEdit}
              onCancelEdit={onCancelEdit}
              onDeleteLabel={onDeleteLabel}
              onAddLabel={onAddLabel}
              isUpdating={isUpdating}
              isDeletingId={isDeletingId}
              isAdding={isAdding}
              isActive={activeFamilyId === family.id}
              onToggleActive={() => setActiveFamilyId(activeFamilyId === family.id ? null : family.id)}
            />
          ))}
        </div>
      )}
    </div>
  );
};

interface FamilyLabelsGroupProps {
  family: Family;
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
  isActive: boolean;
  onToggleActive: () => void;
}

const FamilyLabelsGroup = ({
  family,
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
  isActive,
  onToggleActive
}: FamilyLabelsGroupProps) => {
  // Fetch labels for this specific family
  const {
    data: labelsResponse,
    isLoading,
    error
  } = useLabelsQuery({
    ownerTypes: [OwnerType.Family],
    familyId: family.id
  });
  
  const labels = labelsResponse?.labels || [];
  const hasLabels = labels.length > 0;
  
  const handleAddLabel = (name: string, color: string) => {
    onAddLabel(name, color, OwnerType.Family, family.id);
  };
  
  return (
    <div className="border rounded-lg p-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-medium">{family.name}</h3>
        <button 
          onClick={onToggleActive}
          className="text-sm text-primary hover:underline"
        >
          {isActive ? "Hide" : "Manage labels"}
        </button>
      </div>
      
      {isActive && (
        <>
          {isLoading ? (
            <div className="flex items-center justify-center h-24">
              <Loader2 className="h-5 w-5 animate-spin text-primary mr-2"/>
              <p className="text-sm text-muted-foreground">Loading labels...</p>
            </div>
          ) : error ? (
            <div className="p-3 border rounded-md bg-destructive/10">
              <p className="text-destructive text-sm">Error loading labels</p>
            </div>
          ) : (
            <>
              <div className="mb-4">
                <AddLabelForm 
                  onAddLabel={handleAddLabel}
                  isAdding={isAdding}
                  title={`Add label to ${family.name}`}
                />
              </div>
              
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {!hasLabels ? (
                  <p className="text-muted-foreground col-span-full text-center py-4">
                    No labels for this family yet
                  </p>
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
            </>
          )}
        </>
      )}
    </div>
  );
};