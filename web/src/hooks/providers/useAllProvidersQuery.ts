import { useInfiniteQuery, type InfiniteData } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import Provider from "@/models/provider";
import { OwnerType } from "@/models/ownerType";
import type { DtoProviderModel as ProviderModel } from "@/api/models/DtoProviderModel";

// Types for the infinite query pages
type ProvidersPage = {
    providers: Provider[];
    length: number;
    total: number;
    nextOffset: number;
};

interface AllProvidersQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    limit?: number; // page size, API maximum is 10
    search?: string;
}

/**
 * Fetches **all** providers that match the given filters by requesting as
 * many pages as required. Internally relies on `useInfiniteQuery`.
 */
export const useAllProvidersQuery = (options: AllProvidersQueryOptions = {}) => {
    const {
        limit = 10,
        search,
    } = options;

    const { apiClient } = useApiClient();

    const trimmedSearch = (search ?? '').trim();

    return useInfiniteQuery<ProvidersPage, Error, InfiniteData<ProvidersPage>, ['providers', 'all', number, string], number>({
        queryKey: ['providers', 'all', limit, trimmedSearch],
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
        initialPageParam: 0,
        // Keep showing the previous list while a new search is loading to avoid empty UI flashes
        // This stabilizes e2e flows that search immediately after edits
        placeholderData: (prev) => prev as any,
        queryFn: async ({ pageParam = 0 }) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters: any = {
                limit,
                offset: pageParam,
            };

            if (trimmedSearch) {
                queryParameters.search = trimmedSearch;
            }

            try {
                const result = await apiClient.providers.providersGet(queryParameters);

                // Handle null or empty data array
                if (result && result.data && Array.isArray(result.data)) {
                    return {
                        providers: result.data.map((model: ProviderModel) => Provider.fromModel(model)),
                        length: result.data.length,
                        total: result.total ?? 0,
                        nextOffset: pageParam + result.data.length,
                    };
                }

                // Backend returned no data (null or empty) - this is a valid empty result, not an error
                return { providers: [], length: 0, total: 0, nextOffset: pageParam };
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