import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {ProviderCombobox} from "@/components/providers/ProviderCombobox";
import {PlanCombobox} from "@/components/plans/PlanCombobox";
import {PriceCombobox} from "@/components/prices/PriceCombobox";
import type {FormValues} from "./SubscriptionFormSchema";
import Provider from "@/models/provider";
import {useEffect} from "react";

interface BasicInformationSectionProps {
    providers: Provider[];
}

export const BasicInformationSection = ({providers}: BasicInformationSectionProps) => {
    const form = useFormContext<FormValues>();
    
    // Use formState to trigger re-renders when values change
    // This is more efficient than using multiple watch calls
    const { dirtyFields } = form.formState;
    
    // Get the current values without setting up watchers
    const providerId = form.getValues("providerId");
    const selectedProvider = providers.find(provider => provider.id === providerId);
    
    // Get the selected plan
    const planId = form.getValues("planId");
    const selectedPlan = selectedProvider?.plans.find(plan => plan.id === planId);
    
    // Use a single watch for the provider ID to reset dependent fields
    const watchedProviderId = form.watch("providerId");
    useEffect(() => {
        if (watchedProviderId) {
            // Use a reference to the form that won't change between renders
            const formRef = form;
            formRef.setValue("planId", "");
            formRef.setValue("priceId", "");
        }
    }, [watchedProviderId]);
    
    // Use a single watch for the plan ID to reset price
    const watchedPlanId = form.watch("planId");
    useEffect(() => {
        if (watchedPlanId) {
            // Use a reference to the form that won't change between renders
            const formRef = form;
            formRef.setValue("priceId", "");
        }
    }, [watchedPlanId]);

    return (
        <div className="space-y-4">
            <h2 className="text-xl font-semibold">Basic Information</h2>

            <div className="grid grid-cols-1 gap-4">
                <div>
                    <Label htmlFor="friendlyName">Friendly Name (Optional)</Label>
                    <Input
                        id="friendlyName"
                        {...form.register("friendlyName")}
                        placeholder="e.g., Netflix Family Plan"
                    />
                    {form.formState.errors.friendlyName && (
                        <p className="text-sm text-red-500 mt-1">{form.formState.errors.friendlyName.message}</p>
                    )}
                </div>

                <div>
                    <Label htmlFor="providerId">Provider</Label>
                    <ProviderCombobox
                        providers={providers}
                        value={form.getValues("providerId") || ""}
                        onChange={(value) => form.setValue("providerId", value)}
                        placeholder="Select a provider"
                        emptyMessage="No provider found. Try a different search."
                    />
                    {form.formState.errors.providerId && (
                        <p className="text-sm text-red-500 mt-1">{form.formState.errors.providerId.message}</p>
                    )}
                </div>

                {selectedProvider && selectedProvider.plans.length > 0 && (
                    <div>
                        <Label htmlFor="planId">Plan</Label>
                        <PlanCombobox
                            plans={selectedProvider.plans}
                            value={form.getValues("planId") || ""}
                            onChange={(value) => form.setValue("planId", value)}
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
                        <Label htmlFor="priceId">Price</Label>
                        <PriceCombobox
                            prices={selectedPlan.prices}
                            value={form.getValues("priceId") || ""}
                            onChange={(value) => form.setValue("priceId", value)}
                            placeholder="Select a price"
                            emptyMessage="No price found. Try a different search."
                        />
                        {form.formState.errors.priceId && (
                            <p className="text-sm text-red-500 mt-1">{form.formState.errors.priceId.message}</p>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};