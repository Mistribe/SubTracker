import { useQuery } from '@tanstack/react-query';
import { useApiClient } from '@/hooks/use-api-client.ts';
import Quota from '@/models/quota.ts';

// Query hook to retrieve quota usage related to the authenticated account
export const useAccountQuotaQuery = () => {
  const { apiClient } = useApiClient();

  return useQuery({
    queryKey: ['accountQuota'],
    queryFn: async () => {
      if (!apiClient) throw new Error('API client not initialized');
      const res = await apiClient.accounts.accountsQuotaUsageGet();
      return res.map(item => Quota.fromModel(item));
    },
    enabled: !!apiClient,
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: true,
  });
};

