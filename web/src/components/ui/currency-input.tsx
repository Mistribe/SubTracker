import * as React from "react";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export interface CurrencyInputProps {
  value: {
    amount: number;
    currency: string;
  };
  onChange: (value: { amount: number; currency: string }) => void;
  className?: string;
  placeholder?: string;
  currencies?: Array<{ value: string; label: string }>;
  disabled?: boolean;
  error?: {
    amount?: string;
    currency?: string;
  };
}

const defaultCurrencies = [
  { value: "USD", label: "USD" },
  { value: "EUR", label: "EUR" },
  { value: "GBP", label: "GBP" },
  { value: "CAD", label: "CAD" },
  { value: "AUD", label: "AUD" },
];

export function CurrencyInput({
  value,
  onChange,
  className,
  placeholder = "Enter amount",
  currencies = defaultCurrencies,
  disabled = false,
  error,
}: CurrencyInputProps) {
  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newAmount = e.target.valueAsNumber || 0;
    onChange({ ...value, amount: newAmount });
  };

  const handleCurrencyChange = (newCurrency: string) => {
    onChange({ ...value, currency: newCurrency });
  };

  return (
    <div className="space-y-1">
      <div className={cn("flex relative", className)}>
        <Input
          type="number"
          step="0.01"
          value={value.amount || ""}
          onChange={handleAmountChange}
          placeholder={placeholder}
          className="rounded-r-none pr-2 w-full"
          disabled={disabled}
          aria-invalid={!!error?.amount}
        />
        <div className="w-24">
          <Select
            value={value.currency}
            onValueChange={handleCurrencyChange}
            disabled={disabled}
          >
            <SelectTrigger className="rounded-l-none border-l-0 h-9">
              <SelectValue placeholder="Currency" />
            </SelectTrigger>
            <SelectContent>
              {currencies.map((currency) => (
                <SelectItem key={currency.value} value={currency.value}>
                  {currency.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
      {(error?.amount || error?.currency) && (
        <p className="text-sm text-red-500">
          {error.amount || error.currency}
        </p>
      )}
    </div>
  );
}