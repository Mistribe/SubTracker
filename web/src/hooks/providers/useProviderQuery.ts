import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import Provider from "@/models/provider";

export const useProviderQuery = (id?: string | null) => {
  const { apiClient } = useApiClient();

  return useQuery<Provider | undefined>({
    queryKey: ["provider", id],
    enabled: !!apiClient && !!id && id.length > 0,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 30 * 60 * 1000, // 30 minutes
    refetchOnWindowFocus: false,
    queryFn: async () => {
      if (!apiClient) throw new Error("API client not initialized");
      if (!id) return undefined;
      const model = await apiClient.providers.providersProviderIdGet({ providerId: id });
      return model ? Provider.fromModel(model) : undefined;
    },
  });
};
