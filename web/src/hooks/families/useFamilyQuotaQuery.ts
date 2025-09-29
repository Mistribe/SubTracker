import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client.ts";
import Quota from "@/models/quota.ts";

// Pass `enabled` to control whether the quota is fetched (e.g., only if user has a family)
export const useFamilyQuotaQuery = (enabled: boolean = true) => {
    const {apiClient} = useApiClient();


    return useQuery({
        queryKey: ['familyQuota'],
        queryFn: async () => {
            if (!apiClient) throw new Error("API client not initialized");
            const res = await apiClient.families.familyQuotaUsageGet();
            return res.map((item) => Quota.fromModel(item));
        },
        enabled: !!apiClient && enabled,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};