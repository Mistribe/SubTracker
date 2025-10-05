import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";

export interface UsePreferredCurrencyOptions {
  enabled?: boolean;
}

export function usePreferredCurrency(options: UsePreferredCurrencyOptions = {}) {
  const { enabled = true } = options;
  const { apiClient } = useApiClient();

  const query = useQuery({
    queryKey: ["user", "preferred", "currency"],
    enabled: !!apiClient && enabled,
    queryFn: async () => {
      if (!apiClient) throw new Error("API client not initialized");
      const res = await apiClient.accounts.accountsPreferredCurrencyGet();
      const currency = res?.currency ?? "USD";
      return currency;
    },
    staleTime: 24 * 60 * 60 * 1000,
    refetchOnWindowFocus: false,
  });

  return {
    preferredCurrency: query.data ?? "USD",
    isLoading: query.isLoading,
    error: query.error as unknown,
    refetch: query.refetch,
  } as const;
}
