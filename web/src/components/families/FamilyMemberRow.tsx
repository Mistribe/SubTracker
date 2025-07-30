import { useState } from "react";
import { TableCell, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { CheckIcon, Loader2, XIcon } from "lucide-react";
import FamilyMember from "@/models/familyMember.ts";
import { useApiClient } from "@/hooks/use-api-client.ts";
import { useQueryClient } from "@tanstack/react-query";
import type { UpdateFamilyMemberModel } from "@/api/models";

interface FamilyMemberRowProps {
  member: FamilyMember;
  familyId: string;
  isOwner: boolean;
}

export const FamilyMemberRow = ({ member, familyId, isOwner }: FamilyMemberRowProps) => {
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();
  
  const [editingMemberId, setEditingMemberId] = useState<string | null>(null);
  const [editedMemberName, setEditedMemberName] = useState<string>("");
  const [editedMemberIsKid, setEditedMemberIsKid] = useState<boolean>(false);
  const [isUpdatingMember, setIsUpdatingMember] = useState<boolean>(false);
  const [isDeletingMember, setIsDeletingMember] = useState<boolean>(false);

  // Function to start editing a family member
  const startEditingMember = (member: FamilyMember) => {
    setEditingMemberId(member.id);
    setEditedMemberName(member.name);
    setEditedMemberIsKid(member.isKid);
  };

  // Function to cancel editing a family member
  const cancelEditingMember = () => {
    setEditingMemberId(null);
    setEditedMemberName("");
    setEditedMemberIsKid(false);
  };

  // Function to save family member changes
  const saveMemberChanges = async (familyId: string, member: FamilyMember) => {
    if (!apiClient) return;

    try {
      setIsUpdatingMember(true);

      // Only update if values have changed
      if (editedMemberName === member.name && editedMemberIsKid === member.isKid) {
        cancelEditingMember();
        return;
      }

      const patchModel: Partial<UpdateFamilyMemberModel> = {
        name: editedMemberName,
        isKid: editedMemberIsKid,
      };

      await apiClient.families.byFamilyId(familyId).members.byId(member.id).put(patchModel);

      // Invalidate and refetch the families query
      await queryClient.invalidateQueries({ queryKey: ['families'] });

      // Reset editing state
      cancelEditingMember();
    } catch (error) {
      console.error("Failed to update family member:", error);
    } finally {
      setIsUpdatingMember(false);
    }
  };

  // Function to remove a family member
  const removeMember = async (familyId: string, member: FamilyMember, isOwner: boolean) => {
    if (!apiClient) return;

    // Prevent removing yourself if you are the owner of the family
    if (member.isYou && isOwner) {
      console.error("Cannot remove yourself as the owner of the family");
      return;
    }

    try {
      setIsDeletingMember(true);

      await apiClient.families.byFamilyId(familyId).members.byId(member.id).delete();

      // Invalidate and refetch the families query
      await queryClient.invalidateQueries({ queryKey: ['families'] });
    } catch (error) {
      console.error("Failed to remove family member:", error);
    } finally {
      setIsDeletingMember(false);
    }
  };

  return (
    <TableRow key={member.id}>
      <TableCell className="font-medium">
        {editingMemberId === member.id ? (
          <Input
            value={editedMemberName}
            onChange={(e) => setEditedMemberName(e.target.value)}
            className="w-full"
            placeholder="Member name"
          />
        ) : (
          member.isYou ? (
            <span>{member.name} <i>(You)</i></span>
          ) : (
            <span>{member.name}</span>
          )
        )}
      </TableCell>
      <TableCell>
        {editingMemberId === member.id ? (
          <div className="flex items-center space-x-2">
            <Checkbox
              checked={editedMemberIsKid}
              onCheckedChange={(checked) => setEditedMemberIsKid(checked === true)}
              id={`kid-checkbox-${member.id}`}
            />
            <label htmlFor={`kid-checkbox-${member.id}`}>Kid</label>
          </div>
        ) : (
          member.isKid ? 'Kid' : 'Adult'
        )}
      </TableCell>
      <TableCell>
        {editingMemberId === member.id ? (
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => saveMemberChanges(familyId, member)}
              disabled={isUpdatingMember}
            >
              {isUpdatingMember ? (
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
              ) : (
                <CheckIcon className="h-4 w-4 mr-2" />
              )}
              Save
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={cancelEditingMember}
              disabled={isUpdatingMember}
            >
              <XIcon className="h-4 w-4 mr-2" />
              Cancel
            </Button>
          </div>
        ) : (
          <div className="flex gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => startEditingMember(member)}
            >
              Edit
            </Button>
            {/* Don't show remove button for yourself if you're the owner */}
            {!(member.isYou && isOwner) && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeMember(familyId, member, isOwner)}
                disabled={isDeletingMember}
                className="text-red-500 hover:text-red-700 hover:bg-red-50"
              >
                {isDeletingMember ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                ) : (
                  <XIcon className="h-4 w-4 mr-2" />
                )}
                Remove
              </Button>
            )}
          </div>
        )}
      </TableCell>
    </TableRow>
  );
};