import {type SubscriptionModel_recurrency, SubscriptionModel_recurrencyObject} from "@/api/models";

export enum SubscriptionRecurrency {
    Unknown = 'unknown',
    OneTime = 'one_time',
    Monthly = 'monthly',
    Quarterly = 'quarterly',
    HalfYearly = 'half_yearly',
    Yearly = 'yearly',
    Custom = 'custom'
}

export function fromHttpApi(recurrency: SubscriptionModel_recurrency | null | undefined): SubscriptionRecurrency {
    if (!recurrency) {
        return SubscriptionRecurrency.Unknown;
    }
    switch (recurrency) {
        case SubscriptionModel_recurrencyObject.Unknown:
            return SubscriptionRecurrency.Unknown;
        case SubscriptionModel_recurrencyObject.One_time:
            return SubscriptionRecurrency.OneTime;
        case SubscriptionModel_recurrencyObject.Monthly:
            return SubscriptionRecurrency.Monthly;
        case SubscriptionModel_recurrencyObject.Quarterly:
            return SubscriptionRecurrency.Quarterly;
        case SubscriptionModel_recurrencyObject.Half_yearly:
            return SubscriptionRecurrency.HalfYearly;
        case SubscriptionModel_recurrencyObject.Yearly:
            return SubscriptionRecurrency.Yearly;
        default:
            return SubscriptionRecurrency.Custom;
    }
}