import {useMemo} from "react";
import {useQueries} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Provider from "@/models/provider";

export const useProvidersByIds = (ids?: (string | null | undefined)[]) => {
  const uniqueIds = useMemo(() => {
    return Array.from(new Set((ids ?? []).filter((id): id is string => !!id && id.length > 0)));
  }, [ids]);

  const {apiClient} = useApiClient();

  const results = useQueries({
    queries: uniqueIds.map((id) => ({
      queryKey: ["provider", id],
      enabled: !!apiClient && !!id,
      staleTime: 5 * 60 * 1000,
      gcTime: 30 * 60 * 1000,
      refetchOnWindowFocus: false,
      queryFn: async () => {
        if (!apiClient) throw new Error("API client not initialized");
        const model = await apiClient.providers.providersProviderIdGet({ providerId: id });
        return model ? Provider.fromModel(model) : undefined;
      },
    })),
  });

  const providerMap = useMemo(() => {
    const map = new Map<string, Provider>();
    results.forEach((r, idx) => {
      const id = uniqueIds[idx];
      if (id && r.data) {
        map.set(id, r.data);
      }
    });
    return map;
  }, [results, uniqueIds]);

  const isLoading = results.some((r) => r.isLoading || r.isFetching);

  return { providerMap, isLoading } as const;
};
