import { useState } from "react";
import { Loader2 } from "lucide-react";
import { hexToArgb } from "@/components/ui/utils/color-utils";
import { useLabelsQuery } from "@/hooks/labels/useLabelsQuery";
import { useLabelMutations } from "@/hooks/labels/useLabelMutations";
import { LabelHeader } from "@/components/labels/LabelHeader";
import { SystemLabelsSection } from "@/components/labels/SystemLabelsSection";
import { PersonalLabelsSection } from "@/components/labels/PersonalLabelsSection";
import { FamilyLabelsSection } from "@/components/labels/FamilyLabelsSection";
import Label from "@/models/label";
import { OwnerType } from "@/models/ownerType";

const LabelsPage = () => {
    const [page] = useState(1);
    const [pageSize] = useState(10);

    // Use custom hooks for data fetching and mutations
    const {
        data: queryResponse,
        isLoading,
        error
    } = useLabelsQuery(page, pageSize);

    const {
        createLabelMutation,
        updateLabelMutation,
        deleteLabelMutation,
        canModifyLabel
    } = useLabelMutations();

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
        const labelToDelete = queryResponse?.labels?.find(label => label.id === id);

        // Prevent deleting system labels
        if (labelToDelete && !canModifyLabel(labelToDelete)) {
            return;
        }

        deleteLabelMutation.mutate(id);
    };

    const handleAddLabel = (name: string, color: string, ownerType?: OwnerType, familyId?: string) => {
        if (name.trim()) {
            createLabelMutation.mutate({
                name,
                color,
                ownerType,
                familyId
            });
        }
    };

    // Show loading or error states
    if (isLoading) {
        return (
            <div className="flex flex-col items-center justify-center h-64">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-4"/>
                <p className="text-muted-foreground">Loading labels...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center h-64">
                <p className="text-destructive mb-2">Error loading labels</p>
                <p className="text-muted-foreground text-sm">{error instanceof Error ? error.message : 'Unknown error'}</p>
            </div>
        );
    }

    return (
        <div>
            <LabelHeader totalCount={queryResponse?.total || 0} />

            <div className="mt-8">
                {/* System Labels Section */}
                <SystemLabelsSection labels={queryResponse?.labels || []} />

                {/* Personal and Family Labels Section */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Personal Labels Column */}
                    <PersonalLabelsSection 
                        labels={queryResponse?.labels || []}
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
                        labels={queryResponse?.labels || []}
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
                </div>
            </div>
        </div>
    );
};

export default LabelsPage;