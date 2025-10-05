import { useQuery } from '@tanstack/react-query';
import { useApiClient } from '@/hooks/use-api-client.ts';
import Quota from '@/models/quota.ts';

// Query hook to retrieve quota usage related to subscriptions
export const useSubscriptionsQuotaQuery = () => {
  const { apiClient } = useApiClient();

  return useQuery({
    queryKey: ['subscriptionsQuota'],
    queryFn: async () => {
      if (!apiClient) throw new Error('API client not initialized');
      const res = await apiClient.subscriptions.subscriptionsQuotaUsageGet();
      return res.map(item => Quota.fromModel(item));
    },
    enabled: !!apiClient,
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: true,
  });
};

