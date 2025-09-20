import {useInfiniteQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Subscription from "@/models/subscription";
import type {SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";
import type { DtoSubscriptionModel as SubscriptionModel } from "@/api/models/DtoSubscriptionModel";

interface SubscriptionsQueryOptions {
    limit?: number;
    search?: string;
    fromDate?: Date;
    toDate?: Date;
    providers?: string[];
    recurrencies?: SubscriptionRecurrency[];
    users?: string[];
    withInactive?: boolean;
}

/**
 * Fetches **all** subscriptions that match the given filters by requesting as
 * many pages as required. Internally relies on `useInfiniteQuery`.
 */
export const useSubscriptionsQuery = (options: SubscriptionsQueryOptions = {}) => {
    const {
        limit = 10,
        search,
        fromDate,
        toDate,
        providers,
        recurrencies,
        users,
        withInactive,
    } = options;

    const {apiClient} = useApiClient();

    return useInfiniteQuery({
        queryKey: [
            'subscriptions',
            'preferredCurrency',
            limit,
            search,
            fromDate ? fromDate.toISOString() : null,
            toDate ? toDate.toISOString() : null,
            providers ?? [],
            recurrencies ?? [],
            users ?? [],
            withInactive ?? false,
        ],
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
        initialPageParam: 0,
        queryFn: async ({pageParam = 0}) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters = {
                search: search,
                fromDate: fromDate?.toISOString(),
                toDate: toDate?.toISOString(),
                limit,
                offset: pageParam,
                providers: providers,
                recurrencies: recurrencies as string[] | undefined,
                users: users,
                withInactive: withInactive
            } as const;

            try {
                const result = await apiClient.subscriptions.subscriptionsGet(queryParameters);

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
};