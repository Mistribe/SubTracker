import { DtoSubscriptionModelRecurrencyEnum } from "@/api/models/DtoSubscriptionModel";

export const SubscriptionRecurrency = {
    Unknown: 'unknown',
    OneTime: 'one_time',
    Monthly: 'monthly',
    Quarterly: 'quarterly',
    HalfYearly: 'half_yearly',
    Yearly: 'yearly',
    Custom: 'custom'
} as const;

export type SubscriptionRecurrency = (typeof SubscriptionRecurrency)[keyof typeof SubscriptionRecurrency];

export function fromHttpApi(recurrency: DtoSubscriptionModelRecurrencyEnum | null | undefined): SubscriptionRecurrency {
    if (!recurrency) {
        return SubscriptionRecurrency.Unknown;
    }
    switch (recurrency) {
        case DtoSubscriptionModelRecurrencyEnum.Unknown:
            return SubscriptionRecurrency.Unknown;
        case DtoSubscriptionModelRecurrencyEnum.OneTime:
            return SubscriptionRecurrency.OneTime;
        case DtoSubscriptionModelRecurrencyEnum.Monthly:
            return SubscriptionRecurrency.Monthly;
        case DtoSubscriptionModelRecurrencyEnum.Quarterly:
            return SubscriptionRecurrency.Quarterly;
        case DtoSubscriptionModelRecurrencyEnum.HalfYearly:
            return SubscriptionRecurrency.HalfYearly;
        case DtoSubscriptionModelRecurrencyEnum.Yearly:
            return SubscriptionRecurrency.Yearly;
        default:
            return SubscriptionRecurrency.Custom;
    }
}