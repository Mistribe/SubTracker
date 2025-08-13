import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {SummaryRequestBuilderGetQueryParameters} from "@/api/summary";
import type Summary from "@/models/summary.ts";
import {type Amount, zeroAmount} from "@/models/amount.ts";

export interface UseSubscriptionSummaryQueryOptions {
    /** Number of top providers to return */
    topProviders?: number;
    /** Whether to include monthly total costs */
    totalMonthly?: boolean;
    /** Whether to include yearly total costs */
    totalYearly?: boolean;
    /** Number of upcoming renewals to return */
    upcomingRenewals?: number;
    /** Enable/disable the query */
    enabled?: boolean;
}

/**
 * Fetch subscription summary from backend.
 * Wraps the generated apiClient.summary.get endpoint with React Query.
 */
export function useSubscriptionSummaryQuery(options: UseSubscriptionSummaryQueryOptions = {}) {
    const {
        topProviders = 5,
        totalMonthly = true,
        totalYearly = true,
        upcomingRenewals = 5,
        enabled = true,
    } = options;

    const {apiClient} = useApiClient();

    const query = useQuery<Summary | undefined>({
        queryKey: [
            "subscriptions",
            "summary",
            topProviders,
            totalMonthly,
            totalYearly,
            upcomingRenewals,
        ],
        enabled: !!apiClient && enabled,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
        queryFn: async () => {
            if (!apiClient) throw new Error("API client not initialized");

            const queryParameters: SummaryRequestBuilderGetQueryParameters = {
                topProviders,
                totalMonthly,
                totalYearly,
                upcomingRenewals,
            };
            const response = await apiClient.subscriptions.summary.get({queryParameters});

            return {
                activeSubscriptions: response?.active ?? 0,
                totalMonthly: response?.totalMonthly as Amount ?? zeroAmount,
                totalYearly: response?.totalYearly as Amount ?? zeroAmount,
                topProviders: response?.topProviders ?? [],
                upcomingRenewals: response?.upcomingRenewals ?? [],
            } as Summary;
        },
    });

    return {
        // raw response
        data: query.data,

        // convenience accessors (with safe fallbacks)
        activeSubscriptions: query.data?.activeSubscriptions ?? 0,
        totalMonthly: query.data?.totalMonthly ?? zeroAmount,
        totalYearly: query.data?.totalYearly ?? zeroAmount,
        topProviders: query.data?.topProviders ?? [],
        upcomingRenewals: query.data?.upcomingRenewals ?? [],

        // states
        isLoading: query.isLoading,
        isFetching: query.isFetching,
        isRefetching: query.isRefetching,
        error: query.error as unknown,

        // actions
        refetch: query.refetch,
    } as const;
}
