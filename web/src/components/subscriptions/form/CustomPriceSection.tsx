import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import type {FormValues} from "./SubscriptionFormSchema";
import {Toggle} from "@/components/ui/toggle";
import {CurrencyInput} from "@/components/ui/currency-input";

export const priceSection = () => {
    const form = useFormContext<FormValues>();

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Set the subscription price</h2>
                <p className="text-muted-foreground">Enter the price for this subscription</p>
            </div>

            <div className="max-w-md mx-auto mt-6">
                <div>
                    <Label htmlFor="price" className="text-base mb-2 block">How much does it cost? *</Label>
                    <CurrencyInput
                        value={{
                            amount: form.watch("price.amount") || 0,
                            currency: form.watch("price.currency") || "USD"
                        }}
                        onChange={(value) => {
                            form.setValue("price.amount", value.amount, {shouldValidate: true});
                            form.setValue("price.currency", value.currency, {shouldValidate: true});
                        }}
                        error={{
                            amount: form.formState.errors.price?.amount?.message,
                            currency: form.formState.errors.price?.currency?.message
                        }}
                        className="h-12"
                    />
                    <p className="text-xs text-muted-foreground mt-1">Enter the subscription price</p>
                    {(form.formState.errors.price?.amount || form.formState.errors.price?.currency) && (
                        <div className="mt-2">
                            {form.formState.errors.price?.amount && (
                                <p className="text-sm text-red-500">{form.formState.errors.price.amount.message}</p>
                            )}
                            {form.formState.errors.price?.currency && (
                                <p className="text-sm text-red-500">{form.formState.errors.price.currency.message}</p>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};