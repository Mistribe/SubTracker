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

/**
 * Converts a Go-style duration string (e.g., "46019h48m14.932715s")
 * into a human-friendly format.
 *
 * Rules:
 * - >= 1 day: "Xd Yh Zm"
 * - >= 1 hour: "Xh Ym Zs"
 * - >= 1 minute: "Xm Ys"
 * - < 1 minute: "Xs" (with up to 3 decimal places if needed)
 * - Returns "-" for invalid/empty input.
 */
export const formatProviderDuration = (input: string | undefined | null): string => {
    if (!input || typeof input !== "string") return "-";

    const regex = /(-?\d+(?:\.\d+)?)(ns|us|µs|ms|s|m|h)/g;
    let totalSeconds = 0;
    let match: RegExpExecArray | null;

    while ((match = regex.exec(input)) !== null) {
        const value = parseFloat(match[1]);
        const unit = match[2];

        if (Number.isNaN(value)) continue;

        switch (unit) {
            case "h":
                totalSeconds += value * 3600;
                break;
            case "m":
                totalSeconds += value * 60;
                break;
            case "s":
                totalSeconds += value;
                break;
            case "ms":
                totalSeconds += value / 1000;
                break;
            case "us":
            case "µs":
                totalSeconds += value / 1_000_000;
                break;
            case "ns":
                totalSeconds += value / 1_000_000_000;
                break;
        }
    }

    if (!Number.isFinite(totalSeconds)) return "-";

    const sign = totalSeconds < 0 ? "-" : "";
    const abs = Math.abs(totalSeconds);

    const daySec = 86_400;
    const hourSec = 3_600;
    const minSec = 60;

    const days = Math.floor(abs / daySec);
    let remaining = abs % daySec;

    const hours = Math.floor(remaining / hourSec);
    remaining = remaining % hourSec;

    const minutes = Math.floor(remaining / minSec);
    const secondsFloat = remaining % minSec;

    // Decide formatting based on magnitude
    if (days >= 1) {
        // e.g., "1917d 11h 48m"
        return `${sign}${days}d ${hours}h ${minutes}m`;
    }

    if (hours >= 1) {
        // e.g., "3h 5m 12s"
        const seconds = Math.floor(secondsFloat);
        return `${sign}${hours}h ${minutes}m ${seconds}s`;
    }

    if (minutes >= 1) {
        // e.g., "12m 4s"
        const seconds = Math.floor(secondsFloat);
        return `${sign}${minutes}m ${seconds}s`;
    }

    // Less than a minute: keep some precision
    if (secondsFloat === 0) return `${sign}0s`;
    const secondsRounded = secondsFloat < 10 ? Math.round(secondsFloat * 1000) / 1000
                        : secondsFloat < 30 ? Math.round(secondsFloat * 10) / 10
                        : Math.round(secondsFloat);
    return `${sign}${secondsRounded}s`;
};