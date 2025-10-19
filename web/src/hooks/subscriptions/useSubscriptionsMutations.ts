import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import {OwnerType} from "@/models/ownerType";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";
import {PayerType} from "@/models/payerType.ts";
import type { DtoCreateSubscriptionRequest as CreateSubscriptionModel } from "@/api/models/DtoCreateSubscriptionRequest";
import type { DtoUpdateSubscriptionRequest as UpdateSubscriptionModel } from "@/api/models/DtoUpdateSubscriptionRequest";
import type { DtoEditableOwnerModel as EditableOwnerModel } from "@/api/models/DtoEditableOwnerModel";
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
            friendlyName?: string,
            providerId: string,
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
            familyUsers?: string[],
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

            const owner: EditableOwnerModel = {
                type: subscriptionData.ownerType ?? OwnerType.Personal,
            };
            if ((subscriptionData.ownerType ?? OwnerType.Personal) === OwnerType.Family && subscriptionData.familyId) {
                owner.familyId = subscriptionData.familyId;
            }

            const payload: CreateSubscriptionModel = {
                friendlyName: subscriptionData.friendlyName,
                providerId: subscriptionData.providerId,
                recurrency: subscriptionData.recurrency as unknown as string,
                customRecurrency: subscriptionData.customRecurrency,
                startDate: subscriptionData.startDate,
                endDate: subscriptionData.endDate,
                familyUsers: subscriptionData.familyUsers,
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
                const familyId = subscriptionData.payer.familyId ?? subscriptionData.familyId;
                if (!familyId) {
                    throw new Error("Payer requires a familyId");
                }
                const payer: EditableSubscriptionPayerModel = {
                    type: payerType,
                    familyId,
                    ...(subscriptionData.payer.memberId ? { memberId: subscriptionData.payer.memberId } : {}),
                };

                payload.payer = payer;
            }

            // Add custom price if specified
            if (subscriptionData.customPrice) {
                const customPrice: AmountModel = {
                    value: subscriptionData.customPrice.amount,
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

            const owner: EditableOwnerModel = {
                type: subscriptionData.ownerType ?? OwnerType.Personal,
            };
            if ((subscriptionData.ownerType ?? OwnerType.Personal) === OwnerType.Family && subscriptionData.familyId) {
                owner.familyId = subscriptionData.familyId;
            }

            const payload: UpdateSubscriptionModel = {
                friendlyName: subscriptionData.friendlyName,
                providerId: subscriptionData.providerId,
                planId: subscriptionData.planId,
                priceId: subscriptionData.priceId,
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
                const familyId = subscriptionData.payer.familyId ?? subscriptionData.familyId;
                if (!familyId) {
                    throw new Error("Payer requires a familyId");
                }
                const payer: EditableSubscriptionPayerModel = {
                    type: payerType,
                    familyId,
                    ...(subscriptionData.payer.memberId ? { memberId: subscriptionData.payer.memberId } : {}),
                };

                payload.payer = payer;
            }

            // Add custom price if specified
            if (subscriptionData.customPrice) {
                const customPrice: AmountModel = {
                    value: subscriptionData.customPrice.amount,
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