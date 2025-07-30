import { useState } from "react";
import { CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { CheckIcon, Loader2, XIcon } from "lucide-react";
import Family from "@/models/family.ts";
import { useApiClient } from "@/hooks/use-api-client.ts";
import { useQueryClient } from "@tanstack/react-query";
import type { PatchFamilyModel } from "@/api/models";
import { AddFamilyMemberDialog } from "./AddFamilyMemberDialog.tsx";

interface FamilyHeaderProps {
  family: Family;
}

export const FamilyHeader = ({ family }: FamilyHeaderProps) => {
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();
  
  const [editingFamilyId, setEditingFamilyId] = useState<string | null>(null);
  const [editedName, setEditedName] = useState<string>("");
  const [isUpdating, setIsUpdating] = useState<boolean>(false);

  // Function to start editing a family
  const startEditing = (family: Family) => {
    setEditingFamilyId(family.id);
    setEditedName(family.name);
  };

  // Function to cancel editing
  const cancelEditing = () => {
    setEditingFamilyId(null);
    setEditedName("");
  };

  // Function to save changes
  const saveChanges = async (family: Family) => {
    if (!apiClient) return;

    try {
      setIsUpdating(true);

      // Only update if values have changed
      if (editedName === family.name) {
        cancelEditing();
        return;
      }

      const patchModel: Partial<PatchFamilyModel> = {
        id: family.id,
        name: editedName,
        updatedAt: new Date()
      };

      await apiClient.families.patch(patchModel);

      // Invalidate and refetch the families query
      await queryClient.invalidateQueries({ queryKey: ['families'] });

      // Reset editing state
      cancelEditing();
    } catch (error) {
      console.error("Failed to update family:", error);
    } finally {
      setIsUpdating(false);
    }
  };

  return (
    <div className="flex justify-between items-center">
      <div>
        {editingFamilyId === family.id ? (
          <div className="flex flex-col gap-2">
            <div className="flex items-center gap-2">
              <Input
                value={editedName}
                onChange={(e) => setEditedName(e.target.value)}
                className="w-full"
                placeholder="Family name"
              />
            </div>
          </div>
        ) : (
          <>
            <div className="flex items-center gap-2">
              <CardTitle>{family.name}</CardTitle>
            </div>
            <CardDescription>{family.members.length} members</CardDescription>
          </>
        )}
      </div>
      <div className="flex gap-2">
        {editingFamilyId === family.id ? (
          <>
            <Button
              variant="outline"
              size="sm"
              onClick={() => saveChanges(family)}
              disabled={isUpdating}
            >
              {isUpdating ? (
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
              ) : (
                <CheckIcon className="h-4 w-4 mr-2" />
              )}
              Save
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={cancelEditing}
              disabled={isUpdating}
            >
              <XIcon className="h-4 w-4 mr-2" />
              Cancel
            </Button>
          </>
        ) : (
          <>
            <Button
              variant="outline"
              size="sm"
              onClick={() => startEditing(family)}
            >
              Edit
            </Button>
            <AddFamilyMemberDialog familyId={family.id} />
          </>
        )}
      </div>
    </div>
  );
};