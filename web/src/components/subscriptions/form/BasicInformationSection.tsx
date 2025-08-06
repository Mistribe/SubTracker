import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {ProviderCombobox} from "@/components/providers/ProviderCombobox";
import {PlanCombobox} from "@/components/plans/PlanCombobox";
import {PriceCombobox} from "@/components/prices/PriceCombobox";
import {Toggle} from "@/components/ui/toggle";
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
    const planId = form.watch("planId");
    const priceId = form.watch("priceId");
    const customPrice = form.watch("customPrice");
    
    // Derive values directly
    const selectedProvider = providers.find(provider => provider.id === providerId);
    const selectedPlan = selectedProvider?.plans.find(plan => plan.id === planId);
    const hasCustomPrice = customPrice !== undefined;
    
    // Prepare values for combobox components
    const providerValue = providerId || "";
    const planValue = planId || "";
    const priceValue = priceId || "";
    
    // Prepare currency input value
    const currencyInputValue = {
        amount: customPrice?.amount || 0,
        currency: customPrice?.currency || "USD"
    };
    
    // Define handlers for form value changes
    const handleProviderChange = (value: string) => {
        // Update the provider ID
        form.setValue("providerId", value, { shouldValidate: true });
        
        // Reset dependent fields when provider changes
        form.setValue("planId", "", { shouldValidate: true });
        form.setValue("priceId", "", { shouldValidate: true });
        
        // If the selected provider doesn't have plans, force custom price
        const provider = providers.find(p => p.id === value);
        if (provider && provider.plans.length === 0) {
            form.setValue("customPrice", {amount: 0, currency: "USD"}, { shouldValidate: true });
        }
    };
    
    const handlePlanChange = (value: string) => {
        // Update the plan ID
        form.setValue("planId", value, { shouldValidate: true });
        
        // Reset price when plan changes
        form.setValue("priceId", "", { shouldValidate: true });
    };
    
    const handlePriceChange = (value: string) => {
        // Update the price ID
        form.setValue("priceId", value, { shouldValidate: true });
    };
    
    const handleCustomPriceToggle = (pressed: boolean) => {
        if (pressed) {
            form.setValue("customPrice", {amount: 0, currency: "USD"}, { shouldValidate: true });
        } else {
            form.setValue("customPrice", undefined, { shouldValidate: true });
        }
    };
    
    const handleCurrencyInputChange = (value: { amount: number, currency: string }) => {
        // Update the entire customPrice object at once
        form.setValue("customPrice", {
            amount: value.amount,
            currency: value.currency
        }, { shouldValidate: true });
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
                        <Label htmlFor="friendlyName" className="text-base mb-2 block">What would you like to call this subscription?</Label>
                        <Input
                            id="friendlyName"
                            {...form.register("friendlyName")}
                            placeholder="e.g., Netflix Family Plan"
                            className="h-12"
                        />
                        <p className="text-xs text-muted-foreground mt-1">Optional: Give your subscription a friendly name to easily identify it</p>
                        {form.formState.errors.friendlyName && (
                            <p className="text-sm text-red-500 mt-1">{form.formState.errors.friendlyName.message}</p>
                        )}
                    </div>

                    <div>
                        <Label htmlFor="providerId" className="text-base mb-2 block">Which service provider is this for?</Label>
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

                    {selectedProvider && selectedProvider.plans.length > 0 && (
                        <div>
                            <Label htmlFor="planId" className="text-base mb-2 block">Which plan are you subscribing to?</Label>
                            <PlanCombobox
                                plans={selectedProvider.plans}
                                value={planValue}
                                onChange={handlePlanChange}
                                placeholder="Select a plan"
                                emptyMessage="No plan found. Try a different search."
                            />
                            {form.formState.errors.planId && (
                                <p className="text-sm text-red-500 mt-1">{form.formState.errors.planId.message}</p>
                            )}
                        </div>
                    )}

                    {selectedPlan && selectedPlan.prices.length > 0 && (
                        <div>
                            <Label htmlFor="priceId" className="text-base mb-2 block">Which pricing option did you select?</Label>
                            <PriceCombobox
                                prices={selectedPlan.prices}
                                value={priceValue}
                                onChange={handlePriceChange}
                                placeholder="Select a price"
                                emptyMessage="No price found. Try a different search."
                            />
                            {form.formState.errors.priceId && (
                                <p className="text-sm text-red-500 mt-1">{form.formState.errors.priceId.message}</p>
                            )}
                        </div>
                    )}

                    {/* Custom price toggle - only show when provider has plans */}
                    {selectedProvider && selectedProvider.plans.length > 0 && (
                        <div className="mt-4">
                            <div className="flex items-center justify-between mb-2">
                                <Label htmlFor="hasCustomPrice" className="text-base">I have a custom price</Label>
                                <Toggle
                                    id="hasCustomPrice"
                                    pressed={hasCustomPrice}
                                    onPressedChange={handleCustomPriceToggle}
                                    className="ml-2"
                                    variant="outline"
                                >
                                    {hasCustomPrice ? "Yes" : "No"}
                                </Toggle>
                            </div>
                        </div>
                    )}

                    {/* Custom price input - show when toggle is enabled or provider has no plans, but only if a provider is selected */}
                    {selectedProvider && hasCustomPrice && (
                        <div className="mt-4">
                            <Label htmlFor="customPrice" className="text-base mb-2 block">How much does it cost?</Label>
                            <CurrencyInput
                                value={currencyInputValue}
                                onChange={handleCurrencyInputChange}
                                error={{
                                    amount: form.formState.errors.customPrice?.amount?.message,
                                    currency: form.formState.errors.customPrice?.currency?.message
                                }}
                                className="h-12"
                            />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};