import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import {OwnerType} from "@/models/ownerType";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";
import {PayerType} from "@/models/payerType.ts";
import type { DtoCreateSubscriptionRequest as CreateSubscriptionModel } from "@/api/models/DtoCreateSubscriptionRequest";
import { DtoCreateSubscriptionRequestOwnerEnum } from "@/api/models/DtoCreateSubscriptionRequest";
import type { DtoUpdateSubscriptionRequest as UpdateSubscriptionModel } from "@/api/models/DtoUpdateSubscriptionRequest";
import { DtoUpdateSubscriptionRequestOwnerEnum } from "@/api/models/DtoUpdateSubscriptionRequest";
import type { DtoEditableSubscriptionPayerModel as EditableSubscriptionPayerModel } from "@/api/models/DtoEditableSubscriptionPayerModel";
import { DtoEditableSubscriptionPayerModelTypeEnum } from "@/api/models/DtoEditableSubscriptionPayerModel";
import type { DtoSubscriptionFreeTrialModel as SubscriptionFreeTrialModel } from "@/api/models/DtoSubscriptionFreeTrialModel";
import type { DtoAmountModel as AmountModel } from "@/api/models/DtoAmountModel";

export const useSubscriptionsMutations = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // Create subscription mutation
    const createSubscriptionMutation = useMutation({
        mutationFn: async (subscriptionData: {
            subscriptionId?: string
            friendlyName?: string,
            providerId?: string,
            providerKey?: string,
            recurrency: SubscriptionRecurrency,
            customRecurrency?: number,
            startDate: Date,
            endDate?: Date,
            ownerType?: OwnerType,
            payer?: {
                type: PayerType,
                memberId?: string
            },
            familyUsers?: string[],
            price?: {
                amount: number,
                currency: string
            },
            freeTrial?: {
                startDate: Date,
                endDate: Date
            }
        }) => {
            if (!apiClient) throw new Error("API client not initialized");

            let owner: DtoCreateSubscriptionRequestOwnerEnum;
            switch (subscriptionData.ownerType ?? OwnerType.Personal) {
                case OwnerType.Family:
                    owner = DtoCreateSubscriptionRequestOwnerEnum.Family;
                    break;
                case OwnerType.System:
                    owner = DtoCreateSubscriptionRequestOwnerEnum.System;
                    break;
                case OwnerType.Personal:
                default:
                    owner = DtoCreateSubscriptionRequestOwnerEnum.Personal;
                    break;
            }

            const payload: CreateSubscriptionModel = {
                id: subscriptionData.subscriptionId,
                friendlyName: subscriptionData.friendlyName,
                recurrency: subscriptionData.recurrency as unknown as string,
                customRecurrency: subscriptionData.customRecurrency,
                startDate: subscriptionData.startDate,
                endDate: subscriptionData.endDate,
                familyUsers: subscriptionData.familyUsers,
                owner,
            };
            if (subscriptionData.providerKey) {
                payload.providerKey = subscriptionData.providerKey;
            } else {
                payload.providerId = subscriptionData.providerId;
            }

            // Add payer information if specified
            if (subscriptionData.payer) {
                let payerType: DtoEditableSubscriptionPayerModelTypeEnum;
                switch (subscriptionData.payer.type) {
                    case PayerType.FamilyMember:
                        payerType = DtoEditableSubscriptionPayerModelTypeEnum.FamilyMember;
                        break;
                    case PayerType.Family:
                    default:
                        payerType = DtoEditableSubscriptionPayerModelTypeEnum.Family;
                        break;
                }
                const payer: EditableSubscriptionPayerModel = {
                    type: payerType,
                    ...(subscriptionData.payer.memberId ? { memberId: subscriptionData.payer.memberId } : {}),
                };

                payload.payer = payer;
            }

            // Add custom price if specified
            if (subscriptionData.price) {
                const customPrice: AmountModel = {
                    value: subscriptionData.price.amount,
                    currency: subscriptionData.price.currency
                };

                payload.price = customPrice;
            }

            // Add free trial if specified
            if (subscriptionData.freeTrial) {
                const freeTrial: SubscriptionFreeTrialModel = {
                    startDate: subscriptionData.freeTrial.startDate,
                    endDate: subscriptionData.freeTrial.endDate
                };

                payload.freeTrial = freeTrial;
            }

            return apiClient.subscriptions.subscriptionsPost({ dtoCreateSubscriptionRequest: payload });
        },
        onSuccess: async () => {
            // Invalidate and refetch
            await queryClient.invalidateQueries({queryKey: ['subscriptions']});
        }
    });

    // Delete subscription mutation
    const deleteSubscriptionMutation = useMutation({
        mutationFn: async (subscriptionId: string) => {
            if (!apiClient) throw new Error("API client not initialized");
            return apiClient.subscriptions.subscriptionsSubscriptionIdDelete({ subscriptionId });
        },
        onSuccess: async () => {
            // Invalidate and refetch
            await queryClient.invalidateQueries({queryKey: ['subscriptions']});
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
                providerKey?: string,
                planId?: string,
                priceId?: string,
                recurrency: SubscriptionRecurrency,
                customRecurrency?: number,
                startDate: Date,
                endDate?: Date,
                ownerType?: OwnerType,
                familyId?: string,
                payer?: {
                    type: PayerType,
                    familyId?: string,
                    memberId?: string
                },
                serviceUsers?: string[],
                price?: {
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

            let owner: DtoUpdateSubscriptionRequestOwnerEnum;
            switch (subscriptionData.ownerType ?? OwnerType.Personal) {
                case OwnerType.Family:
                    owner = DtoUpdateSubscriptionRequestOwnerEnum.Family;
                    break;
                case OwnerType.System:
                    owner = DtoUpdateSubscriptionRequestOwnerEnum.System;
                    break;
                case OwnerType.Personal:
                default:
                    owner = DtoUpdateSubscriptionRequestOwnerEnum.Personal;
                    break;
            }

            const payload: UpdateSubscriptionModel = {
                friendlyName: subscriptionData.friendlyName,
                // Keep providerId for transitional API typing. Prefer providerKey when provided by caller.
                providerId: subscriptionData.providerId,
                ...(subscriptionData.providerKey ? { providerKey: subscriptionData.providerKey } : {}),
                recurrency: subscriptionData.recurrency as unknown as string,
                customRecurrency: subscriptionData.customRecurrency,
                startDate: subscriptionData.startDate,
                endDate: subscriptionData.endDate,
                serviceUsers: subscriptionData.serviceUsers,
                owner,
            };

            // Add payer information if specified
            if (subscriptionData.payer) {
                let payerType: DtoEditableSubscriptionPayerModelTypeEnum;
                switch (subscriptionData.payer.type) {
                    case PayerType.FamilyMember:
                        payerType = DtoEditableSubscriptionPayerModelTypeEnum.FamilyMember;
                        break;
                    case PayerType.Family:
                    default:
                        payerType = DtoEditableSubscriptionPayerModelTypeEnum.Family;
                        break;
                }
                const payer: EditableSubscriptionPayerModel = {
                    type: payerType,
                    ...(subscriptionData.payer.memberId ? { memberId: subscriptionData.payer.memberId } : {}),
                };

                payload.payer = payer;
            }

            // Add custom price if specified
            if (subscriptionData.price) {
                const customPrice: AmountModel = {
                    value: subscriptionData.price.amount,
                    currency: subscriptionData.price.currency
                };

                payload.price = customPrice;
            }

            // Add free trial if specified
            if (subscriptionData.freeTrial) {
                const freeTrial: SubscriptionFreeTrialModel = {
                    startDate: subscriptionData.freeTrial.startDate,
                    endDate: subscriptionData.freeTrial.endDate
                };

                payload.freeTrial = freeTrial;
            }

            return apiClient.subscriptions.subscriptionsSubscriptionIdPut({ subscriptionId, dtoUpdateSubscriptionRequest: payload });
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