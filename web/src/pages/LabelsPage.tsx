import {useState} from "react";
import {hexToArgb} from "@/components/ui/utils/color-utils";
import {useLabelsQuery} from "@/hooks/labels/useLabelsQuery";
import {useLabelMutations} from "@/hooks/labels/useLabelMutations";
import {useFamiliesMutations} from "@/hooks/families/useFamiliesMutations";
import {LabelHeader} from "@/components/labels/LabelHeader";
import {SystemLabelsSection} from "@/components/labels/SystemLabelsSection";
import {PersonalLabelsSection} from "@/components/labels/PersonalLabelsSection";
import {FamilyLabelsSection} from "@/components/labels/FamilyLabelsSection";
import Label from "@/models/label";
import {OwnerType} from "@/models/ownerType";

const LabelsPage = () => {
    const [page] = useState(1);
    const [pageSize] = useState(10);

    // Use custom hooks for data fetching and mutations with separate queries for each label type
    const {
        data: systemLabelsResponse,
        isLoading: isLoadingSystemLabels,
        error: systemLabelsError
    } = useLabelsQuery({
        ownerTypes: [OwnerType.System],
        page,
        pageSize
    });

    const {
        data: personalLabelsResponse,
        isLoading: isLoadingPersonalLabels,
        error: personalLabelsError
    } = useLabelsQuery({
        ownerTypes: [OwnerType.Personal],
        page,
        pageSize
    });

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

    const handleDeleteLabel = (id: string) => {
        // Find the label to check if it's a system label
        const labelToDelete =
            systemLabelsResponse?.labels?.find(label => label.id === id) ||
            personalLabelsResponse?.labels?.find(label => label.id === id);

        // Prevent deleting system labels
        if (labelToDelete && !canModifyLabel(labelToDelete)) {
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
                <SystemLabelsSection
                    labels={systemLabelsResponse?.labels || []}
                    isLoading={isLoadingSystemLabels}
                    error={systemLabelsError}
                />

                {/* Personal and Family Labels Section */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Personal Labels Column */}
                    <PersonalLabelsSection
                        labels={personalLabelsResponse?.labels || []}
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
                        isLoading={isLoadingPersonalLabels}
                        error={personalLabelsError}
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