import { useInfiniteQuery } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import Provider from "@/models/provider";
import { OwnerType } from "@/models/ownerType.ts";
import type { DtoProviderModel as ProviderModel } from "@/api/models/DtoProviderModel";

interface ProvidersQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    limit?: number;
}

export const useProvidersQuery = (options: ProvidersQueryOptions = {}) => {
    const {
        ownerTypes,
        familyId,
        limit = 10
    } = options;

    const { apiClient } = useApiClient();

    return useInfiniteQuery({
        queryKey: ['providers', ownerTypes, familyId, limit],
        queryFn: async ({ pageParam = 0 }) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters = {
                limit: limit,
                offset: pageParam,
            };

            const result = await apiClient?.providers.providersGet(queryParameters);

            if (result && result.data) {
                return {
                    providers: result.data.map((model: ProviderModel) => {
                        return Provider.fromModel(model);
                    }),
                    nextOffset: result.data.length === limit ? pageParam + limit : undefined,
                    length: result.data.length,
                    total: result.total ?? 0
                }
            }
            return { providers: [], nextOffset: undefined, length: 0, total: 0 };
        },
        initialPageParam: 0,
        getNextPageParam: (lastPage) => lastPage.nextOffset,
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};