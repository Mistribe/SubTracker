import type Subscription from "@/models/subscription";

function fromCurrencyToUSD(amount: number, currency: string, rates: Map<string, number>): number | undefined {
    if (currency === "USD") return amount;

    const currencyRate = rates.get(currency);
    if (currencyRate === undefined) {
        console.log(`Currency ${currency} to USD rate not found`);
        return undefined;
    }
    return amount * currencyRate;
}

function fromUSDtoCurrency(amount: number, currency: string, rates: Map<string, number>): number | undefined {
    if (currency === "USD") return amount;
    const currencyRate = rates.get(currency);
    if (currencyRate === undefined) {
        console.log(`Currency USD to ${currency} rate not found`);
        return undefined;
    }

    return amount / currencyRate;
}

/**
 * Convert amount from one currency to another using a rates map.
 * The rates map is expected to be relative to a common base (same assumption as Money component).
 * Returns null if conversion can't be performed due to missing rates.
 */
export function convertAmount(
    amount: number,
    from: string,
    to: string,
    rates: Map<string, number>
): number | null {
    if (!Number.isFinite(amount)) return 0;
    if (!from || !to) return null;
    if (from === to) return amount;
    const usdValue = fromCurrencyToUSD(amount, from, rates);
    if (usdValue === undefined) return null;
    const toValue = fromUSDtoCurrency(usdValue, to, rates);
    if (toValue === undefined) return null;
    return toValue;
}

/**
 * Compute the subscription's monthly price expressed in the target currency.
 * Uses the subscription's getMonthlyPrice() (which is in the subscription currency)
 * then converts it to the target currency via rates.
 */
export function subscriptionMonthlyPriceInCurrency(
    subscription: Subscription,
    targetCurrency: string,
    rates: Map<string, number>
): number {
    const monthly = subscription.getMonthlyPrice();
    const fromCurrency = subscription.price?.currency;
    if (!fromCurrency) return 0;
    if (fromCurrency === targetCurrency) return monthly;
    const converted = convertAmount(monthly, fromCurrency, targetCurrency, rates);
    return converted ?? 0;
}

/**
 * Compute the subscription's yearly price expressed in the target currency.
 * Uses the subscription's getYearlyPrice() (which is in the subscription currency)
 * then converts it to the target currency via rates.
 */
export function subscriptionYearlyPriceInCurrency(
    subscription: Subscription,
    targetCurrency: string,
    rates: Map<string, number>
): number {
    const yearly = subscription.getYearlyPrice();
    const fromCurrency = subscription.price?.currency;
    if (!fromCurrency) return 0;
    if (fromCurrency === targetCurrency) return yearly;
    const converted = convertAmount(yearly, fromCurrency, targetCurrency, rates);
    return converted ?? 0;
}
