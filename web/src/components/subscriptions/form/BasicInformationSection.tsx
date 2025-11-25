import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {ProviderCombobox} from "@/components/providers/ProviderCombobox";
import {CurrencyInput} from "@/components/ui/currency-input";
import type {FormValues} from "./SubscriptionFormSchema";
import Provider from "@/models/provider";

interface BasicInformationSectionProps {
    providers: Provider[];
}

export const BasicInformationSection = ({providers}: BasicInformationSectionProps) => {
    const form = useFormContext<FormValues>();

    // Get form values directly
    const providerId = form.watch("providerId");
    const price = form.watch("price");

    // Prepare values for combobox components
    const providerValue = providerId || "";

    // Prepare currency input value
    const currencyInputValue = {
        amount: price?.amount || 0,
        currency: price?.currency || "USD"
    };

    // Define handlers for form value changes
    const handleProviderChange = (value: string) => {
        // Update the provider ID
        form.setValue("providerId", value, {shouldValidate: true});

        // Ensure price has a valid default value since it's now required
        if (!form.getValues("price") || form.getValues("price.amount") <= 0) {
            form.setValue("price", {amount: 10, currency: "USD"}, {shouldValidate: true});
        }
    };

    const handleCurrencyInputChange = (value: { amount: number, currency: string }) => {
        // Update the entire price object at once
        form.setValue("price", {
            amount: value.amount,
            currency: value.currency
        }, {shouldValidate: true});
    };

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Tell us about your subscription</h2>
                <p className="text-muted-foreground">Let's start with the basic details of your subscription</p>
            </div>

            <div className="max-w-md mx-auto mt-6">
                <div className="grid grid-cols-1 gap-6">
                    <div>
                        <Label htmlFor="providerId" className="text-base mb-2 block">Which service provider is this
                            for?</Label>
                        <ProviderCombobox
                            providers={providers}
                            value={providerValue}
                            onChange={handleProviderChange}
                            placeholder="Select a provider"
                            emptyMessage="No provider found. Try a different search."
                        />
                        {form.formState.errors.providerId && (
                            <p className="text-sm text-red-500 mt-1">{form.formState.errors.providerId.message}</p>
                        )}
                    </div>
                    <div>
                        <Label htmlFor="friendlyName" className="text-base mb-2 block">What would you like to call this
                            subscription?</Label>
                        <Input
                            id="friendlyName"
                            {...form.register("friendlyName")}
                            placeholder="e.g., Netflix Family Plan"
                            className="h-12"
                        />
                        <p className="text-xs text-muted-foreground mt-1">Optional: Give your subscription a friendly
                            name to easily identify it</p>
                        {form.formState.errors.friendlyName && (
                            <p className="text-sm text-red-500 mt-1">{form.formState.errors.friendlyName.message}</p>
                        )}
                    </div>

                   

                    {/* Custom price input - now always required */}
                    <div className="mt-4">
                        <Label htmlFor="price" className="text-base mb-2 block">How much does it cost? *</Label>
                        <CurrencyInput
                            value={currencyInputValue}
                            onChange={handleCurrencyInputChange}
                            error={{
                                amount: form.formState.errors.price?.amount?.message,
                                currency: form.formState.errors.price?.currency?.message
                            }}
                            className="h-12"
                        />
                        <p className="text-xs text-muted-foreground mt-1">Enter the subscription price</p>
                    </div>
                </div>
            </div>
        </div>
    );
};