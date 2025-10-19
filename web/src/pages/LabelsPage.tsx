import { useEffect, useState } from "react";
import { hexToArgb } from "@/components/ui/utils/color-utils";
import { useLabelMutations } from "@/hooks/labels/useLabelMutations";
import { useFamiliesMutations } from "@/hooks/families/useFamiliesMutations";
import { PageHeader } from "@/components/ui/page-header";
import { useAllLabelsQuery } from "@/hooks/labels/useAllLabelsQuery";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";
import { LabelItem } from "@/components/labels/LabelItem";
import { EditableLabelItem } from "@/components/labels/EditableLabelItem";
import { AddLabelDialog } from "@/components/labels/AddLabelDialog";
import { DeleteLabelDialog } from "@/components/labels/DeleteLabelDialog";
import Label from "@/models/label";
import { OwnerType } from "@/models/ownerType";
import { Loader2 } from "lucide-react";
import { useLabelsQuotaQuery } from "@/hooks/labels/useLabelsQuotaQuery.ts";
import { QuotaButton } from "@/components/quotas/QuotaButton";
import { FeatureId } from "@/models/billing.ts";

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

    // State for search and editing
    const [searchText, setSearchText] = useState("");
    const [editingId, setEditingId] = useState<string | null>(null);
    const [editingName, setEditingName] = useState("");
    const [editingColor, setEditingColor] = useState("");
    const [labelToDelete, setLabelToDelete] = useState<Label | null>(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

    // Fetch the single family and adapt to an array for UI components expecting arrays
    const { data: familyData } = useFamilyQuery();
    const families = familyData ? [familyData] : [];

    // Fetch all labels (system, personal, and family)
    const {
        data,
        isLoading,
        error,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
    } = useAllLabelsQuery({
        ownerTypes: [OwnerType.System, OwnerType.Personal, OwnerType.Family],
        limit: 10,
        search: searchText,
    });

    // Keep requesting next pages until every label is fetched
    useEffect(() => {
        if (hasNextPage && !isFetchingNextPage) {
            fetchNextPage().catch(() => void 0);
        }
    }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

    // Combine all labels from all pages
    const allLabels = data?.pages.flatMap((page) => page.labels) || [];

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

        // Open the confirmation dialog
        if (label) {
            setLabelToDelete(label);
            setIsDeleteDialogOpen(true);
        }
    };

    const handleConfirmDelete = () => {
        if (labelToDelete) {
            deleteLabelMutation.mutate(labelToDelete.id, {
                onSettled: () => {
                    setIsDeleteDialogOpen(false);
                    setLabelToDelete(null);
                }
            });
        }
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


    // Create a map of family IDs to family objects for easy lookup
    const familyMap = new Map(families.map(family => [family.id, family]));

    if (isLoading) {
        return (
            <div className="container mx-auto py-6">
                <PageHeader
                    title="Labels"
                    searchText={searchText}
                    onSearchChange={setSearchText}
                    searchPlaceholder="Search labels..."
                    quotaButton={
                        <QuotaButton
                            useQuotaQuery={useLabelsQuotaQuery}
                            featureIds={[FeatureId.CustomLabelsCount]}
                            featureLabels={{
                                [FeatureId.CustomLabelsCount]: "Labels Used"
                            }}
                        />
                    }
                    actionButton={
                        <AddLabelDialog
                            onAddLabel={handleAddLabel}
                            isAdding={createLabelMutation.isPending || createFamilyLabelMutation.isPending}
                            families={families}
                        />
                    }
                />
                <div className="flex flex-col items-center justify-center h-64">
                    <Loader2 className="h-8 w-8 animate-spin text-primary mb-4" />
                    <p className="text-muted-foreground">Loading labels...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="container mx-auto py-6">
                <PageHeader
                    title="Labels"
                    searchText={searchText}
                    onSearchChange={setSearchText}
                    searchPlaceholder="Search labels..."
                    quotaButton={
                        <QuotaButton
                            useQuotaQuery={useLabelsQuotaQuery}
                            featureIds={[FeatureId.CustomLabelsCount]}
                            featureLabels={{
                                [FeatureId.CustomLabelsCount]: "Labels Used"
                            }}
                        />
                    }
                    actionButton={
                        <AddLabelDialog
                            onAddLabel={handleAddLabel}
                            isAdding={createLabelMutation.isPending || createFamilyLabelMutation.isPending}
                            families={families}
                        />
                    }
                />
                <div className="p-6 border rounded-md bg-destructive/10 mt-8">
                    <p className="text-destructive text-center">Error loading labels</p>
                </div>
            </div>
        );
    }

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Labels"
                searchText={searchText}
                onSearchChange={setSearchText}
                searchPlaceholder="Search labels..."
                quotaButton={
                    <QuotaButton
                        useQuotaQuery={useLabelsQuotaQuery}
                        featureIds={[FeatureId.CustomLabelsCount]}
                        featureLabels={{
                            [FeatureId.CustomLabelsCount]: "Labels Used"
                        }}
                    />
                }
                actionButton={
                    <AddLabelDialog
                        onAddLabel={handleAddLabel}
                        isAdding={createLabelMutation.isPending || createFamilyLabelMutation.isPending}
                        families={families}
                    />
                }
            />
            <div className="mt-8">
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
                    {allLabels.length === 0 ? (
                        <p className="text-muted-foreground col-span-full text-center py-4">
                            {searchText ? "No labels match your search" : "No labels found"}
                        </p>
                    ) : (
                        allLabels.map((label) => {
                            // Get family name for family labels
                            let ownerName;
                            if (label.owner.type === OwnerType.Family && label.owner.familyId) {
                                const family = familyMap.get(label.owner.familyId);
                                ownerName = family?.name || "Unknown Family";
                            }

                            return (
                                <div
                                    key={label.id}
                                    className={`${editingId === label.id ? 'bg-muted' : 'hover:bg-muted/50'}`}
                                >
                                    {editingId === label.id ? (
                                        <EditableLabelItem
                                            label={label}
                                            initialName={editingName}
                                            initialColor={editingColor}
                                            onSave={handleSaveEdit}
                                            onCancel={handleCancelEdit}
                                            isSaving={updateLabelMutation.isPending}
                                        />
                                    ) : (
                                        <LabelItem
                                            label={label}
                                            ownerName={ownerName}
                                            onEdit={canModifyLabel(label) ? handleStartEditing : undefined}
                                            onDelete={canModifyLabel(label) ? handleDeleteLabel : undefined}
                                            isDeleting={deleteLabelMutation.isPending && deleteLabelMutation.variables === label.id}
                                            isReadOnly={!canModifyLabel(label)}
                                        />
                                    )}
                                </div>
                            );
                        })
                    )}
                </div>
            </div>

            <DeleteLabelDialog
                label={labelToDelete}
                open={isDeleteDialogOpen}
                onOpenChange={setIsDeleteDialogOpen}
                onConfirm={handleConfirmDelete}
                isDeleting={deleteLabelMutation.isPending}
            />
        </div>
    );
};

export default LabelsPage;

