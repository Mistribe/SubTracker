import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import currencyCodes from "currency-codes";
import getSymbolFromCurrency from "currency-symbol-map";
import type {UpdatePreferredCurrencyModel} from "@/api/models/user";

interface ProfileQueryOptions {
    /**
     * Whether the query should be enabled (default: true)
     */
    enabled?: boolean;
}

/**
 * Currency object with code, name, and symbol
 */
interface Currency {
    code: string;
    name: string;
    symbol: string;
}

const fallbackCurrencyCodes = ["USD", "EUR"]

/**
 * Hook for managing user profile data
 * Provides functions to fetch and update profile information
 *
 * @param options - Options for the profile query
 * @returns Object containing profile data, loading states, and update functions
 */
export const useProfileManagement = (options: ProfileQueryOptions = {}) => {
    const {enabled = true} = options;
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // Query to fetch preferred currency from the backend
    const preferredCurrencyQuery = useQuery({
        queryKey: ['preferredCurrency'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            try {
                return await apiClient.users.preferred.currency.get();
            } catch (error) {
                console.error('Error fetching preferred currency:', error);
                throw error;
            }
        },
        enabled: !!apiClient && enabled,
        // Don't refetch on window focus to avoid unnecessary API calls
        refetchOnWindowFocus: false,
    });

    const availableCurrencyQuery = useQuery({
        queryKey: ['availableCurrencies'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            try {
                // Get currency codes from the API
                let response = await apiClient.currencies.supported.get();

                if (!response || response.length === 0) {
                    console.warn('No supported currencies found');
                    response = fallbackCurrencyCodes;
                }
                // Transform currency codes to objects with code, name, and symbol
                const currencies: Currency[] = response.map((currencyCode: string) => {
                    try {
                        const currencyDetails = currencyCodes.code(currencyCode);
                        const symbol = getSymbolFromCurrency(currencyCode) || '';

                        return {
                            code: currencyCode,
                            name: currencyDetails?.currency || currencyCode,
                            symbol: symbol
                        };
                    } catch (error) {
                        // Fallback for unsupported currency codes
                        console.warn(`Currency code not found: ${currencyCode}`, error);
                        return {
                            code: currencyCode,
                            name: currencyCode,
                            symbol: ''
                        };
                    }
                });

                return currencies;
            } catch (error) {
                console.error('Error fetching available currencies:', error);
                throw error;
            }
        },
        enabled: !!apiClient,
        refetchOnWindowFocus: false,
    })

    // Mutation to update profile data, specifically the preferred currency
    const updatePreferredCurrencyMutation = useMutation({
        mutationFn: async (currency: string) => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            try {
                const payload: UpdatePreferredCurrencyModel = {
                    currency,
                    additionalData: {}
                };

                return await apiClient.users.preferred.currency.put(payload);
            } catch (error) {
                console.error('Error updating preferred currency:', error);
                throw error;
            }
        },
        onSuccess: async () => {
            // Invalidate and refetch profile and preferred currency data to ensure the UI is up to date
            await queryClient.invalidateQueries({ queryKey: ['profile'] });
            await queryClient.invalidateQueries({ queryKey: ['preferredCurrency'] });

            // Also invalidate the global preferred currency hook used across the app
            await queryClient.invalidateQueries({ queryKey: ['user', 'preferred', 'currency'] });

            // Invalidate subscription summary so totals are recomputed with the new preferred currency
            await queryClient.invalidateQueries({ queryKey: ['subscriptions', 'summary'] });
        }
    });


    return {
        preferredCurrency: preferredCurrencyQuery.data,
        isLoadingPreferredCurrency: preferredCurrencyQuery.isLoading,
        isErrorPreferredCurrency: preferredCurrencyQuery.isError,
        errorPreferredCurrency: preferredCurrencyQuery.error,
        availableCurrencies: availableCurrencyQuery.data,
        isLoadingAvailableCurrencies: availableCurrencyQuery.isLoading,
        isErrorAvailableCurrencies: availableCurrencyQuery.isError,
        errorAvailableCurrencies: availableCurrencyQuery.error,
        isUpdating: updatePreferredCurrencyMutation.isPending
    };
};