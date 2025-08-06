import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import {OwnerType} from "@/models/ownerType";
import type {CreateSubscriptionModel, EditableOwnerModel, EditableSubscriptionPayerModel, SubscriptionCustomPriceModel, SubscriptionFreeTrialModel, UpdateSubscriptionModel} from "@/api/models";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";

export const useSubscriptionsMutations = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // Create subscription mutation
    const createSubscriptionMutation = useMutation({
        mutationFn: async (subscriptionData: {
            friendlyName?: string,
            providerId: string,
            planId: string,
            priceId: string,
            recurrency: SubscriptionRecurrency,
            customRecurrency?: number,
            startDate: Date,
            endDate?: Date,
            ownerType?: OwnerType,
            familyId?: string,
            payer?: {
                type: OwnerType,
                familyId?: string,
                memberId?: string
            },
            serviceUsers?: string[],
            customPrice?: {
                amount: number,
                currency: string
            },
            freeTrial?: {
                startDate: Date,
                endDate: Date
            }
        }) => {
            if (!apiClient) throw new Error("API client not initialized");

            const payload: CreateSubscriptionModel = {
                friendlyName: subscriptionData.friendlyName || null,
                providerId: subscriptionData.providerId,
                planId: subscriptionData.planId,
                priceId: subscriptionData.priceId,
                recurrency: subscriptionData.recurrency,
                customRecurrency: subscriptionData.customRecurrency || null,
                startDate: subscriptionData.startDate,
                endDate: subscriptionData.endDate || null,
                serviceUsers: subscriptionData.serviceUsers || [],
            };

            // Add owner information if specified
            if (subscriptionData.ownerType) {
                const owner: EditableOwnerModel = {
                    type: subscriptionData.ownerType
                };

                // Add family ID if owner type is family and family ID is provided
                if (subscriptionData.ownerType === OwnerType.Family && subscriptionData.familyId) {
                    owner.familyId = subscriptionData.familyId;
                }

                payload.owner = owner;
            }

            // Add payer information if specified
            if (subscriptionData.payer) {
                const payer: EditableSubscriptionPayerModel = {
                    type: subscriptionData.payer.type
                };

                // Add family ID and member ID if provided
                if (subscriptionData.payer.familyId) {
                    payer.familyId = subscriptionData.payer.familyId;
                }

                if (subscriptionData.payer.memberId) {
                    payer.memberId = subscriptionData.payer.memberId;
                }

                payload.payer = payer;
            }

            // Add custom price if specified
            if (subscriptionData.customPrice) {
                const customPrice: SubscriptionCustomPriceModel = {
                    amount: subscriptionData.customPrice.amount,
                    currency: subscriptionData.customPrice.currency
                };

                payload.customPrice = customPrice;
            }

            // Add free trial if specified
            if (subscriptionData.freeTrial) {
                const freeTrial: SubscriptionFreeTrialModel = {
                    startDate: subscriptionData.freeTrial.startDate,
                    endDate: subscriptionData.freeTrial.endDate
                };

                payload.freeTrial = freeTrial;
            }

            return apiClient.subscriptions.post(payload);
        },
        onSuccess: () => {
            // Invalidate and refetch
            queryClient.invalidateQueries({queryKey: ['subscriptions']});
        }
    });

    // Delete subscription mutation
    const deleteSubscriptionMutation = useMutation({
        mutationFn: async (subscriptionId: string) => {
            if (!apiClient) throw new Error("API client not initialized");
            return apiClient.subscriptions.bySubscriptionId(subscriptionId).delete();
        },
        onSuccess: () => {
            // Invalidate and refetch
            queryClient.invalidateQueries({queryKey: ['subscriptions']});
        }
    });

    // Update subscription mutation
    const updateSubscriptionMutation = useMutation({
        mutationFn: async ({
            subscriptionId,
            subscriptionData
        }: {
            subscriptionId: string,
            subscriptionData: {
                friendlyName?: string,
                providerId: string,
                planId: string,
                priceId: string,
                recurrency: SubscriptionRecurrency,
                customRecurrency?: number,
                startDate: Date,
                endDate?: Date,
                ownerType?: OwnerType,
                familyId?: string,
                payer?: {
                    type: OwnerType,
                    familyId?: string,
                    memberId?: string
                },
                serviceUsers?: string[],
                customPrice?: {
                    amount: number,
                    currency: string
                },
                freeTrial?: {
                    startDate: Date,
                    endDate: Date
                }
            }
        }) => {
            if (!apiClient) throw new Error("API client not initialized");

            const payload: UpdateSubscriptionModel = {
                friendlyName: subscriptionData.friendlyName || null,
                providerId: subscriptionData.providerId,
                planId: subscriptionData.planId,
                priceId: subscriptionData.priceId,
                recurrency: subscriptionData.recurrency,
                customRecurrency: subscriptionData.customRecurrency || null,
                startDate: subscriptionData.startDate,
                endDate: subscriptionData.endDate || null,
                serviceUsers: subscriptionData.serviceUsers || [],
            };

            // Add owner information if specified
            if (subscriptionData.ownerType) {
                const owner: EditableOwnerModel = {
                    type: subscriptionData.ownerType
                };

                // Add family ID if owner type is family and family ID is provided
                if (subscriptionData.ownerType === OwnerType.Family && subscriptionData.familyId) {
                    owner.familyId = subscriptionData.familyId;
                }

                payload.owner = owner;
            }

            // Add payer information if specified
            if (subscriptionData.payer) {
                const payer: EditableSubscriptionPayerModel = {
                    type: subscriptionData.payer.type
                };

                // Add family ID and member ID if provided
                if (subscriptionData.payer.familyId) {
                    payer.familyId = subscriptionData.payer.familyId;
                }

                if (subscriptionData.payer.memberId) {
                    payer.memberId = subscriptionData.payer.memberId;
                }

                payload.payer = payer;
            }

            // Add custom price if specified
            if (subscriptionData.customPrice) {
                const customPrice: SubscriptionCustomPriceModel = {
                    amount: subscriptionData.customPrice.amount,
                    currency: subscriptionData.customPrice.currency
                };

                payload.customPrice = customPrice;
            }

            // Add free trial if specified
            if (subscriptionData.freeTrial) {
                const freeTrial: SubscriptionFreeTrialModel = {
                    startDate: subscriptionData.freeTrial.startDate,
                    endDate: subscriptionData.freeTrial.endDate
                };

                payload.freeTrial = freeTrial;
            }

            return apiClient.subscriptions.bySubscriptionId(subscriptionId).put(payload);
        },
        onSuccess: () => {
            // Invalidate and refetch
            queryClient.invalidateQueries({queryKey: ['subscriptions']});
        }
    });

    return {
        createSubscriptionMutation,
        updateSubscriptionMutation,
        deleteSubscriptionMutation
    };
};