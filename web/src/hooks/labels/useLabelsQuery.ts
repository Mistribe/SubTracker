import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import type { LabelModel } from "@/api/models";
import Label from "@/models/label";

export const useLabelsQuery = (page = 1, pageSize = 10) => {
  const { apiClient } = useApiClient();

  return useQuery({
    queryKey: ['labels'],
    queryFn: async () => {
      if (!apiClient) {
        throw new Error('API client not initialized');
      }
      const result = await apiClient?.labels.get({
        queryParameters: {
          withDefault: true,
          page: page,
          size: pageSize
        }
      });
      if (result && result.data) {
        return {
          labels: result.data.map((model: LabelModel) => {
            return Label.fromModel(model);
          }),
          length: result.data.length,
          total: result.total ?? 0
        }
      }
      return { labels: [], length: 0, total: 0 };
    },
    enabled: !!apiClient,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: true,
  });
};