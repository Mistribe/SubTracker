import { useEffect, useMemo, useRef } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import type { CurrencyGetRateResponse } from "@/api/models";

export interface UseCurrencyRatesOptions {
  enabled?: boolean;
  /** Optional conversion date to fetch rates for (defaults to now). */
  date?: Date;
}

export interface CurrencyRatesData {
  /** Map of currency code (e.g., "USD") to rate (number) relative to base. */
  rates: Record<string, number>;
  /** Timestamp from backend parsed as Date if available. */
  timestamp: Date | null;
  /** Raw response in case consumers need extended data. */
  raw?: CurrencyGetRateResponse;
}

function parseRates(response?: CurrencyGetRateResponse | null): CurrencyRatesData {
  const timestamp = response?.timestamp ? new Date(response.timestamp) : null;
  // Kiota AdditionalDataHolder stores unknown fields in `additionalData`
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const ratesAdditional: Record<string, unknown> = (response?.rates as any)?.additionalData ?? {};
  const rates: Record<string, number> = {};
  for (const [k, v] of Object.entries(ratesAdditional)) {
    const num = typeof v === "number" ? v : Number(v);
    if (!Number.isNaN(num)) rates[k] = num;
  }
  return { rates, timestamp, raw: response ?? undefined };
}

function isOlderThanOneDay(date: Date | null): boolean {
  if (!date) return true;
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  return diffMs > 24 * 60 * 60 * 1000; // 24 hours
}

export function useCurrencyRates(options: UseCurrencyRatesOptions = {}) {
  const { enabled = true, date } = options;
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();

  const currentKey = ["currency", "rates", date ? date.toISOString() : "current"] as const;

  // Fetch current rates
  const query = useQuery<ReturnType<typeof parseRates>>({
    queryKey: currentKey,
    enabled: !!apiClient && enabled,
    refetchOnWindowFocus: false,
    staleTime: 60 * 60 * 1000, // 1 hour
    queryFn: async () => {
      if (!apiClient) throw new Error("API client not initialized");
      const res = await apiClient.currencies.rates.get({
        queryParameters: date ? { date: date.toISOString() } : undefined,
      });
      const parsed = parseRates(res);
      // If server returned no rates, ask it to refresh and use the refreshed values
      if (!parsed || Object.keys(parsed.rates).length === 0) {
        const refreshed = await apiClient.currencies.rates.refresh.post();
        return parseRates(refreshed as CurrencyGetRateResponse);
      }
      return parsed;
    },
  });

  // Manual refresh endpoint now returns refreshed values
  const refreshMutation = useMutation({
    mutationKey: ["currency", "rates", "refresh"],
    mutationFn: async () => {
      if (!apiClient) throw new Error("API client not initialized");
      const refreshed = await apiClient.currencies.rates.refresh.post();
      return parseRates(refreshed as CurrencyGetRateResponse);
    },
    onSuccess: (data) => {
      // Directly set the latest data into the current query cache
      queryClient.setQueryData(currentKey, data);
    },
  });

  // Auto refresh if the fetched rates are older than one day (only once per mount)
  const didAutoRefresh = useRef(false);
  useEffect(() => {
    if (!enabled || !apiClient) return;
    const data = query.data;
    if (!didAutoRefresh.current && data && isOlderThanOneDay(data.timestamp)) {
      didAutoRefresh.current = true;
      // Fire and forget; consumers can also call refresh manually
      refreshMutation.mutate();
    }
  }, [apiClient, enabled, query.data, refreshMutation]);

  const currencyCodes = useMemo(() => Object.keys(query.data?.rates ?? {}), [query.data]);

  return {
    // data
    data: query.data,
    rates: query.data?.rates ?? {},
    timestamp: query.data?.timestamp ?? null,
    currencyCodes,

    // states
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isRefetching: query.isRefetching,
    error: query.error as unknown,

    // actions
    refetch: query.refetch,
    refresh: refreshMutation.mutateAsync,
    isRefreshing: refreshMutation.isPending,
  } as const;
}
