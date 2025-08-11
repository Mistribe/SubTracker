import {useInfiniteQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {SubscriptionModel} from "@/api/models";
import type {SubscriptionsRequestBuilderGetQueryParameters} from "@/api/subscriptions";
import Subscription from "@/models/subscription";
import {OwnerType} from "@/models/ownerType";
import {useEffect} from "react";

interface AllSubscriptionsQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    limit?: number; // page size, API maximum is 10
}

/**
 * Fetches **all** subscriptions that match the given filters by requesting as
 * many pages as required. Internally relies on `useInfiniteQuery`.
 */
export const useAllSubscriptionsQuery = (options: AllSubscriptionsQueryOptions = {}) => {
    const {
        limit = 50,
    } = options;

    const {apiClient} = useApiClient();

    const query = useInfiniteQuery({
        queryKey: ['subscriptions', 'all', limit],
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
        initialPageParam: 0,
        queryFn: async ({pageParam = 0}) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters: SubscriptionsRequestBuilderGetQueryParameters = {
                limit,
                offset: pageParam,
            };

            try {
                const result = await apiClient.subscriptions.get({queryParameters});

                if (result && result.data) {
                    return {
                        subscriptions: result.data.map((model: SubscriptionModel) => Subscription.fromModel(model)),
                        length: result.data.length,
                        total: result.total ?? 0,
                        nextOffset: pageParam + result.data.length,
                    };
                }

                return {subscriptions: [], length: 0, total: 0, nextOffset: pageParam};
            } catch (error) {
                console.error('Failed to fetch subscriptions:', error);
                throw error;
            }
        },
        /**
         * Determine the offset for the next page; return undefined when all
         * elements have already been fetched.
         */
        getNextPageParam: (lastPage, allPages) => {
            const fetchedCount = allPages.reduce((sum, p) => sum + p.subscriptions.length, 0);
            if (fetchedCount < lastPage.total) {
                return lastPage.nextOffset;
            }
            return undefined;
        },
    });

    useEffect(() => {
        if (query.hasNextPage && !query.isFetchingNextPage) {
            query.fetchNextPage();
        }
    }, [query, query.hasNextPage, query.isFetchingNextPage]);


    return query;
};