import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import Label from "@/models/label";
import { OwnerType } from "@/models/ownerType";

export const useLabelMutations = () => {
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();

  // Create label mutation
  const createLabelMutation = useMutation({
    mutationFn: async (labelData: { name: string, color: string, ownerType?: OwnerType }) => {
      const ownerType = labelData.ownerType ?? OwnerType.Personal;
      const payload = {
        name: labelData.name,
        color: labelData.color,
        owner: ownerType,
      };
      
      return apiClient?.labels.labelsPost({ dtoCreateLabelRequest: payload });
    },
    onSuccess: async () => {
      // Invalidate and refetch
      await queryClient.invalidateQueries({ queryKey: ['labels'] });
    }
  });

  // Update label mutation
  const updateLabelMutation = useMutation({
    mutationFn: async ({ id, name, color }: { id: string, name: string, color: string }) => {
      return apiClient?.labels.labelsLabelIdPut({
        labelId: id,
        dtoUpdateLabelRequest: {
          name,
          color,
        },
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['labels'] });
    }
  });

  // Delete label mutation
  const deleteLabelMutation = useMutation({
    mutationFn: async (id: string) => {
      return apiClient?.labels.labelsLabelIdDelete({ labelId: id });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['labels'] });
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