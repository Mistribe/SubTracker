import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";

// Format currency
export const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(amount);
};

// Format recurrency
export const formatRecurrency = (recurrency: SubscriptionRecurrency, customRecurrency: number | undefined) => {
    if (recurrency === SubscriptionRecurrency.Custom && customRecurrency) {
        return `Every ${customRecurrency} days`;
    }

    switch (recurrency) {
        case SubscriptionRecurrency.Monthly:
            return 'Monthly';
        case SubscriptionRecurrency.Quarterly:
            return 'Quarterly';
        case SubscriptionRecurrency.HalfYearly:
            return 'Half Yearly';
        case SubscriptionRecurrency.Yearly:
            return 'Yearly';
        case SubscriptionRecurrency.OneTime:
            return 'One Time';
        default:
            return 'Unknown';
    }
};