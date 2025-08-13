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

    const {apiClient} = useApiClient();

    const trimmedSearch = (search ?? '').trim();

    return useInfiniteQuery({
        queryKey: ['providers', 'all', limit, trimmedSearch],
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

            if (trimmedSearch) {
                queryParameters.search = trimmedSearch;
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