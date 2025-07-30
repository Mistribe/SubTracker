import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import Label from "@/models/label";
import { OwnerType } from "@/models/ownerType";

// Define types for the payload
interface LabelOwner {
  type: OwnerType;
  family_id?: string;
}

interface CreateLabelPayload {
  name: string;
  color: string;
  owner?: LabelOwner;
}

export const useLabelMutations = () => {
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();

  // Create label mutation
  const createLabelMutation = useMutation({
    mutationFn: async (labelData: { name: string, color: string, ownerType?: OwnerType, familyId?: string }) => {
      const payload: CreateLabelPayload = {
        name: labelData.name,
        color: labelData.color
      };
      
      // Add owner information if specified
      if (labelData.ownerType) {
        payload.owner = {
          type: labelData.ownerType
        };
        
        // Add family ID if owner type is family and family ID is provided
        if (labelData.ownerType === OwnerType.Family && labelData.familyId) {
          payload.owner.family_id = labelData.familyId;
        }
      }
      
      return apiClient?.labels.post(payload);
    },
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ['labels'] });
    }
  });

  // Update label mutation
  const updateLabelMutation = useMutation({
    mutationFn: async ({ id, name, color }: { id: string, name: string, color: string }) => {
      return apiClient?.labels.byId(id).put({
        name: name,
        color: color,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['labels'] });
    }
  });

  // Delete label mutation
  const deleteLabelMutation = useMutation({
    mutationFn: async (id: string) => {
      return apiClient?.labels.byId(id).delete();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['labels'] });
    }
  });

  // Helper function to check if a label can be edited or deleted
  const canModifyLabel = (label: Label): boolean => {
    return !label.owner.isSystem;
  };

  return {
    createLabelMutation,
    updateLabelMutation,
    deleteLabelMutation,
    canModifyLabel
  };
};