import {useQuery} from '@tanstack/react-query';
import {useApiClient} from "@/hooks/use-api-client.ts";
import type {ApiClient} from "@/lib/api-client.ts";
import Subscription from "@/models/subscription.ts";

// Mock API function - replace with your actual API call
const getSubscription = async (apiClient: ApiClient, id: string | undefined) => {
    if (!apiClient || !id) return undefined;
    const result = await apiClient.subscriptions.subscriptionsSubscriptionIdGet({
        subscriptionId: id
    })

    return Subscription.fromModel(result);
};

export const useSubscription = (subscriptionId: string | undefined) => {
    const {apiClient} = useApiClient();

    return useQuery({
        queryKey: ['subscription', subscriptionId],
        queryFn: () => getSubscription(apiClient, subscriptionId),
        // Only run the query if an ID is provided
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};