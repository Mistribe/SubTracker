import {useInfiniteQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {LabelsRequestBuilderGetQueryParameters} from "@/api/labels";
import Label from "@/models/label";
import {OwnerType} from "@/models/ownerType";
import type {LabelModel} from "@/api/models/label";

interface AllLabelsQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    limit?: number; // page size, API maximum is 10
    search?: string;
}

/**
 * Fetches **all** labels that match the given filters by requesting as
 * many pages as required. Internally relies on `useInfiniteQuery`.
 */
export const useAllLabelsQuery = (options: AllLabelsQueryOptions = {}) => {
    const {
        limit = 10,
        search,
    } = options;

    const {apiClient} = useApiClient();

    return useInfiniteQuery({
        queryKey: ['labels', 'all', limit, (search ?? '').trim()],
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000,
        initialPageParam: 0,
        queryFn: async ({pageParam = 0}) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters: LabelsRequestBuilderGetQueryParameters = {
                limit,
                offset: pageParam,
            };

            const trimmedSearch = (search ?? '').trim();
            if (trimmedSearch) {
                // Use backend-side search when a search term is provided
                queryParameters.search = trimmedSearch;
            }

            try {
                const result = await apiClient.labels.get({queryParameters});

                if (result && result.data) {
                    return {
                        labels: result.data.map((model: LabelModel) => Label.fromModel(model)),
                        length: result.data.length,
                        total: result.total ?? 0,
                        nextOffset: pageParam + result.data.length,
                    };
                }

                return {labels: [], length: 0, total: 0, nextOffset: pageParam};
            } catch (error) {
                console.error('Failed to fetch labels:', error);
                throw error;
            }
        },
        /**
         * Determine the offset for the next page; return undefined when all
         * elements have already been fetched.
         */
        getNextPageParam: (lastPage, allPages) => {
            const fetchedCount = allPages.reduce((sum, p) => sum + p.labels.length, 0);
            if (fetchedCount < lastPage.total) {
                return lastPage.nextOffset;
            }
            return undefined;
        },
    });
};
