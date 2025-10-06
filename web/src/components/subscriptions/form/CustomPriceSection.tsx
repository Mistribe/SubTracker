import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import type {FormValues} from "./SubscriptionFormSchema";
import {Toggle} from "@/components/ui/toggle";
import {CurrencyInput} from "@/components/ui/currency-input";

export const CustomPriceSection = () => {
    const form = useFormContext<FormValues>();

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Set the subscription price</h2>
                <p className="text-muted-foreground">Enter the price for this subscription</p>
            </div>

            <div className="max-w-md mx-auto mt-6">
                <div>
                    <Label htmlFor="customPrice" className="text-base mb-2 block">How much does it cost? *</Label>
                    <CurrencyInput
                        value={{
                            amount: form.watch("customPrice.amount") || 0,
                            currency: form.watch("customPrice.currency") || "USD"
                        }}
                        onChange={(value) => {
                            form.setValue("customPrice.amount", value.amount, {shouldValidate: true});
                            form.setValue("customPrice.currency", value.currency, {shouldValidate: true});
                        }}
                        error={{
                            amount: form.formState.errors.customPrice?.amount?.message,
                            currency: form.formState.errors.customPrice?.currency?.message
                        }}
                        className="h-12"
                    />
                    <p className="text-xs text-muted-foreground mt-1">Enter the subscription price</p>
                    {(form.formState.errors.customPrice?.amount || form.formState.errors.customPrice?.currency) && (
                        <div className="mt-2">
                            {form.formState.errors.customPrice?.amount && (
                                <p className="text-sm text-red-500">{form.formState.errors.customPrice.amount.message}</p>
                            )}
                            {form.formState.errors.customPrice?.currency && (
                                <p className="text-sm text-red-500">{form.formState.errors.customPrice.currency.message}</p>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};