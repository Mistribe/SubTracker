import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Provider from "@/models/provider";
import {OwnerType} from "@/models/ownerType";
import type {CreateProviderModel, UpdateProviderModel} from "@/api/models";

export const useProvidersMutations = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // Create provider mutation
    const createProviderMutation = useMutation({
        mutationFn: async (providerData: {
            name: string,
            description?: string,
            iconUrl?: string,
            url?: string,
            pricingPageUrl?: string,
            labels?: string[],
            ownerType?: OwnerType,
            familyId?: string
        }) => {
            const payload: CreateProviderModel = {
                name: providerData.name,
                description: providerData.description,
                iconUrl: providerData.iconUrl,
                url: providerData.url,
                pricingPageUrl: providerData.pricingPageUrl,
                labels: providerData.labels || []
            };

            // Add owner information if specified
            if (providerData.ownerType) {
                payload.owner = {
                    type: providerData.ownerType
                };

                // Add family ID if owner type is family and family ID is provided
                if (providerData.ownerType === OwnerType.Family && providerData.familyId) {
                    payload.owner.familyId = providerData.familyId;
                }
            }

            return apiClient?.providers.post(payload);
        },
        onSuccess: () => {
            // Invalidate and refetch
            queryClient.invalidateQueries({queryKey: ['providers']});
        }
    });

    // Update provider mutation
    const updateProviderMutation = useMutation({
        mutationFn: async (providerData: {
            id: string,
            etag: string,
            name: string,
            description?: string,
            iconUrl?: string,
            url?: string,
            pricingPageUrl?: string,
            labels?: string[],
            ownerType?: OwnerType,
            familyId?: string
        }) => {
            if (!apiClient) throw new Error("API client not initialized");

            const payload: UpdateProviderModel = {
                name: providerData.name,
                description: providerData.description,
                iconUrl: providerData.iconUrl,
                url: providerData.url,
                pricingPageUrl: providerData.pricingPageUrl,
                labels: providerData.labels || [],
            };
            
            // Add etag to the additional data for optimistic concurrency control
            payload.additionalData = {
                ...payload.additionalData,
                etag: providerData.etag
            };
            
            return apiClient.providers.byProviderId(providerData.id).put(payload);
        },
        onSuccess: () => {
            // Invalidate and refetch
            queryClient.invalidateQueries({queryKey: ['providers']});
        }
    });

    // Helper function to check if a provider can be edited or deleted
    const canModifyProvider = (provider: Provider): boolean => {
        return !provider.owner.isSystem;
    };

    return {
        createProviderMutation,
        updateProviderMutation,
        canModifyProvider
    };
};