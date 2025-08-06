import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import type {FormValues} from "./SubscriptionFormSchema";
import {Toggle} from "@/components/ui/toggle";
import {CurrencyInput} from "@/components/ui/currency-input";

export const CustomPriceSection = () => {
    const form = useFormContext<FormValues>();
    const hasCustomPrice = form.watch("customPrice") !== undefined;

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Do you want to set a custom price?</h2>
                <p className="text-muted-foreground">Toggle below if you need to specify a different price than the one
                    from the selected plan</p>
            </div>

            <div className="flex justify-center mt-4">
                <Toggle
                    id="hasCustomPrice"
                    pressed={hasCustomPrice}
                    onPressedChange={(pressed) => {
                        if (pressed) {
                            form.setValue("customPrice", {amount: 0, currency: "USD"});
                        } else {
                            form.setValue("customPrice", undefined);
                        }
                    }}
                    className="px-8 py-6 text-lg"
                    variant="outline"
                >
                    {hasCustomPrice ? "Yes, use custom price" : "No, use plan price"}
                </Toggle>
            </div>

            {hasCustomPrice && (
                <div className="max-w-md mx-auto mt-6">
                    <h3 className="text-lg font-medium mb-4">Custom Price Details</h3>
                    <div>
                        <Label htmlFor="customPrice" className="text-base mb-2 block">How much does it cost?</Label>
                        <CurrencyInput
                            value={{
                                amount: form.watch("customPrice.amount") || 0,
                                currency: form.watch("customPrice.currency") || "USD"
                            }}
                            onChange={(value) => {
                                form.setValue("customPrice.amount", value.amount);
                                form.setValue("customPrice.currency", value.currency);
                            }}
                            error={{
                                amount: form.formState.errors.customPrice?.amount?.message,
                                currency: form.formState.errors.customPrice?.currency?.message
                            }}
                            className="h-12"
                        />
                    </div>
                </div>
            )}
        </div>
    );
};