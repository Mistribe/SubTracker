import {useState} from "react";
import {hexToArgb} from "@/components/ui/utils/color-utils";
import {useLabelMutations} from "@/hooks/labels/useLabelMutations";
import {useFamiliesMutations} from "@/hooks/families/useFamiliesMutations";
import {LabelHeader} from "@/components/labels/LabelHeader";
import {SystemLabelsSection} from "@/components/labels/SystemLabelsSection";
import {PersonalLabelsSection} from "@/components/labels/PersonalLabelsSection";
import {FamilyLabelsSection} from "@/components/labels/FamilyLabelsSection";
import Label from "@/models/label";
import {OwnerType} from "@/models/ownerType";

const LabelsPage = () => {

    const {
        createLabelMutation,
        updateLabelMutation,
        deleteLabelMutation,
        canModifyLabel
    } = useLabelMutations();

    const {
        createFamilyLabelMutation
    } = useFamiliesMutations();

    // State for editing
    const [editingId, setEditingId] = useState<string | null>(null);
    const [editingName, setEditingName] = useState("");
    const [editingColor, setEditingColor] = useState("");

    const handleStartEditing = (label: Label) => {
        // Prevent editing system labels
        if (!canModifyLabel(label)) {
            return;
        }

        setEditingId(label.id);
        setEditingName(label.name);

        // Convert hex color to ARGB if needed
        let color = label.color;

        // If it's a standard hex color (#RRGGBB)
        if (label.color.startsWith('#') && label.color.length === 7) {
            color = hexToArgb(label.color);
        }
        // If it's not in ARGB format (#AARRGGBB), make sure it is
        else if (!label.color.startsWith('#') || label.color.length !== 9) {
            // Default to fully opaque black if invalid format
            color = hexToArgb("#000000");
        }

        setEditingColor(color);
    };

    const handleSaveEdit = (id: string, name: string, color: string) => {
        if (id && name.trim()) {
            updateLabelMutation.mutate({
                id,
                name,
                color
            }, {
                onSuccess: () => {
                    setEditingId(null);
                    setEditingName("");
                }
            });
        }
    };

    const handleCancelEdit = () => {
        setEditingId(null);
        setEditingName("");
    };

    const handleDeleteLabel = (id: string, label?: Label) => {
        // If label is provided, check if it can be modified
        if (label && !canModifyLabel(label)) {
            return;
        }

        deleteLabelMutation.mutate(id);
    };

    const handleAddLabel = (name: string, color: string, ownerType?: OwnerType, familyId?: string) => {
        if (!name.trim()) {
            return;
        }

        // Use different mutations based on the owner type
        if (ownerType === OwnerType.Family && familyId) {
            // Use the family-specific mutation for family labels
            createFamilyLabelMutation.mutate({
                familyId,
                name,
                color
            });
        } else {
            // Use the general label mutation for other types
            createLabelMutation.mutate({
                name,
                color,
                ownerType,
                familyId
            });
        }
    };

    // We'll handle loading and error states for each section separately

    return (
        <div>
            <LabelHeader/>

            <div className="mt-8">
                {/* System Labels Section */}
                <SystemLabelsSection />

                {/* Personal and Family Labels Section */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Personal Labels Column */}
                    <PersonalLabelsSection
                        editingId={editingId}
                        editingName={editingName}
                        editingColor={editingColor}
                        onStartEditing={handleStartEditing}
                        onSaveEdit={handleSaveEdit}
                        onCancelEdit={handleCancelEdit}
                        onDeleteLabel={handleDeleteLabel}
                        onAddLabel={handleAddLabel}
                        isUpdating={updateLabelMutation.isPending}
                        isDeletingId={deleteLabelMutation.isPending ? deleteLabelMutation.variables as string : null}
                        isAdding={createLabelMutation.isPending}
                    />

                    {/* Family Labels Column */}
                    <FamilyLabelsSection
                        editingId={editingId}
                        editingName={editingName}
                        editingColor={editingColor}
                        onStartEditing={handleStartEditing}
                        onSaveEdit={handleSaveEdit}
                        onCancelEdit={handleCancelEdit}
                        onDeleteLabel={handleDeleteLabel}
                        onAddLabel={handleAddLabel}
                        isUpdating={updateLabelMutation.isPending}
                        isDeletingId={deleteLabelMutation.isPending ? deleteLabelMutation.variables as string : null}
                        isAdding={createLabelMutation.isPending || createFamilyLabelMutation.isPending}
                    />
                </div>
            </div>
        </div>
    );
};

export default LabelsPage;