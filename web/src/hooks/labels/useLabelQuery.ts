import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import Label from "@/models/label";
import type {LabelModel} from "@/api/models/label";

export const useLabelQuery = (id: string | undefined | null) => {
  const { apiClient } = useApiClient();

  return useQuery({
    queryKey: ["label", id],
    enabled: !!apiClient && !!id,
    staleTime: 5 * 60 * 1000,
    queryFn: async () => {
      if (!apiClient || !id) {
        throw new Error("API client not initialized or invalid label id");
      }
      const result = await apiClient.labels.byId(id).get();
      if (result) {
        return Label.fromModel(result as LabelModel);
      }
      return undefined;
    },
  });
};
