import {useInfiniteQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {ProviderModel} from "@/api/models";
import type {ProvidersRequestBuilderGetQueryParameters} from "@/api/providers";
import Provider from "@/models/provider";
import {OwnerType} from "@/models/ownerType";

interface AllProvidersQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    limit?: number; // page size, API maximum is 10
}

/**
 * Fetches **all** providers that match the given filters by requesting as
 * many pages as required. Internally relies on `useInfiniteQuery`.
 */
export const useAllProvidersQuery = (options: AllProvidersQueryOptions = {}) => {
    const {
        ownerTypes,
        familyId,
        limit = 10,
    } = options;

    const {apiClient} = useApiClient();

    return useInfiniteQuery({
        queryKey: ['providers', 'all', ownerTypes, familyId, limit],
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
        initialPageParam: 0,
        queryFn: async ({pageParam = 0}) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters: ProvidersRequestBuilderGetQueryParameters = {
                limit,
                offset: pageParam,
            };

            if (ownerTypes && ownerTypes.length > 0) {
                // Convert OwnerType enum values to strings for the API
                queryParameters.ownerType = ownerTypes.map(type => type.toString());
            }

            if (familyId) {
                queryParameters.familyId = familyId;
            }

            try {
                const result = await apiClient.providers.get({queryParameters});

                if (result && result.data) {
                    return {
                        providers: result.data.map((model: ProviderModel) => Provider.fromModel(model)),
                        length: result.data.length,
                        total: result.total ?? 0,
                        nextOffset: pageParam + result.data.length,
                    };
                }

                return {providers: [], length: 0, total: 0, nextOffset: pageParam};
            } catch (error) {
                console.error('Failed to fetch providers:', error);
                throw error;
            }
        },
        /**
         * Determine the offset for the next page; return undefined when all
         * elements have already been fetched.
         */
        getNextPageParam: (lastPage, allPages) => {
            const fetchedCount = allPages.reduce((sum, p) => sum + p.providers.length, 0);
            if (fetchedCount < lastPage.total) {
                return lastPage.nextOffset;
            }
            return undefined;
        },
    });
};