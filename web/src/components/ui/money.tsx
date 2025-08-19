import React, {useMemo} from "react";
import type {Amount} from "@/models/amount.ts";

export interface MoneyProps {
    // Supports new API: pass an Amount instance, and optional source Amount.
    // Backward compatibility: accepts a number with a separate currency prop.
    amount: Amount | undefined | null;
    className?: string;
    /**
     * Minimum fraction digits when formatting. Defaults to 2.
     */
    minimumFractionDigits?: number;
    /**
     * Maximum fraction digits when formatting. Defaults to 2.
     */
    maximumFractionDigits?: number;
}

function format(amount: number | undefined, currency: string | undefined, min = 2, max = 2) {
    if (amount === undefined || currency === undefined) {
        return "Unknown value";
    }
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

export const Money: React.FC<MoneyProps> = ({
                                                amount,
                                                className,
                                                minimumFractionDigits = 2,
                                                maximumFractionDigits = 2,
                                            }) => {
    const {displayMain, displayOriginal} = useMemo(() => {
        const main = format(
            amount?.value,
            amount?.currency,
            minimumFractionDigits,
            maximumFractionDigits,
        );

        const showSource = !!amount?.source && amount.currency !== amount?.currency;
        const original = showSource
            ? format(
                amount.value,
                amount.currency,
                minimumFractionDigits,
                maximumFractionDigits,
            )
            : null;

        return {displayMain: main, displayOriginal: original};
    }, [amount,  minimumFractionDigits, maximumFractionDigits]);

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
