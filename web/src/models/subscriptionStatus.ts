import {DtoSubscriptionModelStatusEnum} from "@/api";

export const SubscriptionStatus = {
    Unknow: 'Unknow',
    Active: 'Active',
    Ended: 'Ended',
    NotStarted: 'NotStarted',
} as const;

export type SubscriptionStatus = (typeof SubscriptionStatus)[keyof typeof SubscriptionStatus];

export function subscriptionStatusFromHttpApi(status: DtoSubscriptionModelStatusEnum | null | undefined): SubscriptionStatus {
    switch (status) {
        case DtoSubscriptionModelStatusEnum.Active:
            return SubscriptionStatus.Active;
        case DtoSubscriptionModelStatusEnum.Ended:
            return SubscriptionStatus.Ended;
        case DtoSubscriptionModelStatusEnum.NotStarted:
            return SubscriptionStatus.NotStarted;
        default:
            return SubscriptionStatus.Unknow;
    }
}