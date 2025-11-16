import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Provider from "@/models/provider";
import {OwnerType} from "@/models/ownerType";
import type { DtoCreateProviderRequest as CreateProviderModel } from "@/api/models/DtoCreateProviderRequest";
import { DtoCreateProviderRequestOwnerEnum } from "@/api/models/DtoCreateProviderRequest";
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
        }) => {
            let owner: DtoCreateProviderRequestOwnerEnum;
            switch (providerData.ownerType ?? OwnerType.Personal) {
                case OwnerType.Family:
                    owner = DtoCreateProviderRequestOwnerEnum.Family;
                    break;
                case OwnerType.System:
                    owner = DtoCreateProviderRequestOwnerEnum.System;
                    break;
                case OwnerType.Personal:
                default:
                    owner = DtoCreateProviderRequestOwnerEnum.Personal;
                    break;
            }

            const payload: CreateProviderModel = {
                name: providerData.name,
                description: providerData.description,
                iconUrl: providerData.iconUrl,
                url: providerData.url,
                pricingPageUrl: providerData.pricingPageUrl,
                labels: providerData.labels || [],
                owner,
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
            
            // Note: Generated API does not expose If-Match header; backend update is optimistic here.
            return apiClient.providers.providersProviderIdPut({ providerId: providerData.id, dtoUpdateProviderRequest: payload });
        },
        onSuccess: async (updated) => {
            // Optimistically update any cached infinite providers lists so UI reflects changes immediately
            // Match all queries that start with ['providers', 'all'] (any limit/search)
            queryClient.setQueriesData({ queryKey: ['providers', 'all'] }, (oldData: any) => {
                if (!oldData || !oldData.pages) return oldData;
                const updatedId = updated?.id;
                const updatedName = updated?.name;
                const updatedDesc = updated?.description;
                const updatedUrl = updated?.url;
                const updatedIcon = updated?.iconUrl;
                const updatedPricingUrl = updated?.pricingPageUrl;
                return {
                    ...oldData,
                    pages: oldData.pages.map((page: any) => ({
                        ...page,
                        providers: page.providers.map((p: any) =>
                            p.id === updatedId
                                ? { ...p, name: updatedName, description: updatedDesc, url: updatedUrl, iconUrl: updatedIcon, pricingPageUrl: updatedPricingUrl }
                                : p
                        ),
                    })),
                };
            });

            // Invalidate queries but avoid refetching active ones immediately to keep UI stable for e2e flows
            await queryClient.invalidateQueries({ queryKey: ['providers'], refetchType: 'inactive' as any });
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