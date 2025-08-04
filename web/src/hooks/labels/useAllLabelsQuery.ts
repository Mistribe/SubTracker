import {useInfiniteQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {LabelModel} from "@/api/models";
import type {LabelsRequestBuilderGetQueryParameters} from "@/api/labels";
import Label from "@/models/label";
import {OwnerType} from "@/models/ownerType";

interface AllLabelsQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    limit?: number; // page size, API maximum is 10
}

/**
 * Fetches **all** labels that match the given filters by requesting as
 * many pages as required. Internally relies on `useInfiniteQuery`.
 */
export const useAllLabelsQuery = (options: AllLabelsQueryOptions = {}) => {
    const {
        ownerTypes,
        familyId,
        limit = 10,
    } = options;

    const {apiClient} = useApiClient();

    return useInfiniteQuery({
        queryKey: ['labels', 'all', ownerTypes, familyId, limit],
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

            if (ownerTypes && ownerTypes.length > 0) {
                // Convert OwnerType enum values to strings for the API
                queryParameters.ownerType = ownerTypes.map(type => type.toString());
            }

            if (familyId) {
                queryParameters.familyId = familyId;
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
