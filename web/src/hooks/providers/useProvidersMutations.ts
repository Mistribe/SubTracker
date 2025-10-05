import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Provider from "@/models/provider";
import {OwnerType} from "@/models/ownerType";
import type { DtoCreateProviderRequest as CreateProviderModel } from "@/api/models/DtoCreateProviderRequest";
import type { DtoUpdateProviderRequest as UpdateProviderModel } from "@/api/models/DtoUpdateProviderRequest";

export const useProvidersMutations = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

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
                labels: providerData.labels || [],
                owner: {
                    type: providerData.ownerType ?? OwnerType.Personal,
                    ...(providerData.ownerType === OwnerType.Family && providerData.familyId ? { familyId: providerData.familyId } : {}),
                },
            };

            return apiClient?.providers.providersPost({ dtoCreateProviderRequest: payload });
        },
        onSuccess: async () => {
            // Invalidate and refetch
            await queryClient.invalidateQueries({queryKey: ['providers']});
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
                updatedAt: new Date(),
            };
            
            return apiClient.providers.providersProviderIdPut({ providerId: providerData.id, dtoUpdateProviderRequest: payload });
        },
        onSuccess: async () => {
            // Invalidate and refetch
            await queryClient.invalidateQueries({queryKey: ['providers']});
        }
    });

    // Delete provider mutation
    const deleteProviderMutation = useMutation({
        mutationFn: async (providerId: string) => {
            if (!apiClient) throw new Error("API client not initialized");
            return apiClient.providers.providersProviderIdDelete({ providerId });
        },
        onSuccess: async () => {
            // Invalidate and refetch
            await queryClient.invalidateQueries({queryKey: ['providers']});
        }
    });

    // Helper function to check if a provider can be edited or deleted
    const canModifyProvider = (provider: Provider): boolean => {
        return !provider.owner.isSystem;
    };

    // Helper function to check if a provider can be deleted (only family or personal)
    const canDeleteProvider = (provider: Provider): boolean => {
        return !provider.owner.isSystem && (provider.owner.isFamily || provider.owner.isPersonal);
    };

    return {
        createProviderMutation,
        updateProviderMutation,
        deleteProviderMutation,
        canModifyProvider,
        canDeleteProvider
    };
};