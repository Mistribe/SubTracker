import React, { useMemo } from "react";
import { usePreferredCurrency } from "@/hooks/currencies/usePreferredCurrency";
import { useCurrencyRates } from "@/hooks/currencies/useCurrencyRates";

export interface MoneyProps {
  amount: number;
  currency: string; // source currency (e.g., "USD")
  className?: string;
  /**
   * Optional override for preferred currency (primarily for testing or special cases).
   */
  preferredCurrencyOverride?: string;
  /**
   * Minimum fraction digits when formatting. Defaults to 2.
   */
  minimumFractionDigits?: number;
  /**
   * Maximum fraction digits when formatting. Defaults to 2.
   */
  maximumFractionDigits?: number;
}

function format(amount: number, currency: string, min = 2, max = 2) {
  try {
    return new Intl.NumberFormat(undefined, {
      style: "currency",
      currency,
      minimumFractionDigits: min,
      maximumFractionDigits: max,
    }).format(amount);
  } catch {
    // Fallback if currency code is invalid for the current runtime
    return `${amount.toFixed(Math.max(min, Math.min(max, 2)))} ${currency}`;
  }
}

function convertAmount(amount: number, from: string, to: string, rates: Record<string, number>): number | null {
  if (from === to) return amount;
  const fromRate = rates[from];
  const toRate = rates[to];
  if (typeof fromRate !== "number" || typeof toRate !== "number") return null;
  // Convert via base: amount_in_base = amount / fromRate; amount_in_to = amount_in_base * toRate
  const amountInBase = amount / fromRate;
  return amountInBase * toRate;
}

export const Money: React.FC<MoneyProps> = ({
  amount,
  currency,
  className,
  preferredCurrencyOverride,
  minimumFractionDigits = 2,
  maximumFractionDigits = 2,
}) => {
  const { preferredCurrency } = usePreferredCurrency();
  const { rates } = useCurrencyRates();

  const targetCurrency = preferredCurrencyOverride ?? preferredCurrency;

  const { displayMain, displayOriginal } = useMemo(() => {
    if (!currency || !targetCurrency) {
      return {
        displayMain: format(amount, currency || targetCurrency || "USD", minimumFractionDigits, maximumFractionDigits),
        displayOriginal: null as string | null,
      };
    }

    if (currency === targetCurrency) {
      return {
        displayMain: format(amount, targetCurrency, minimumFractionDigits, maximumFractionDigits),
        displayOriginal: null as string | null,
      };
    }

    const converted = convertAmount(amount, currency, targetCurrency, rates);

    if (converted == null) {
      // Cannot convert due to missing rates; show original only
      return {
        displayMain: format(amount, currency, minimumFractionDigits, maximumFractionDigits),
        displayOriginal: null as string | null,
      };
    }

    return {
      displayMain: format(converted, targetCurrency, minimumFractionDigits, maximumFractionDigits),
      displayOriginal: format(amount, currency, minimumFractionDigits, maximumFractionDigits),
    };
  }, [amount, currency, targetCurrency, rates, minimumFractionDigits, maximumFractionDigits]);

  return (
    <span className={className}>
      {displayMain}
      {displayOriginal && (
        <span className="text-muted-foreground"> ({displayOriginal})</span>
      )}
    </span>
  );
};

export default Money;
