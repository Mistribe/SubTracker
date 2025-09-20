import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";

interface CurrenciesQueryOptions {
  enabled?: boolean;
}

export const useCurrenciesQuery = (options: CurrenciesQueryOptions = {}) => {
  const { enabled = true } = options;
  const { apiClient } = useApiClient();

  return useQuery({
    queryKey: ['currencies', 'supported'],
    queryFn: async () => {
      if (!apiClient) {
        throw new Error('API client not initialized');
      }

      const result = await apiClient.currencies.currenciesSupportedGet();
      
      return result || [];
    },
    enabled: !!apiClient && enabled,
    staleTime: 24 * 60 * 60 * 1000, // 24 hours - currencies don't change often
    refetchOnWindowFocus: false,
  });
};