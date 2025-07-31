import {useState} from "react";
import {useProvidersMutations} from "@/hooks/providers/useProvidersMutations";
import {useFamiliesQuery} from "@/hooks/families/useFamiliesQuery";
import {useCurrenciesQuery} from "@/hooks/currencies/useCurrenciesQuery";
import {Button} from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import {Alert, AlertDescription} from "@/components/ui/alert";
import {AlertCircle} from "lucide-react";
import Provider from "@/models/provider";
import Plan from "@/models/plan";
import Price from "@/models/price";
import {Accordion, AccordionContent, AccordionItem, AccordionTrigger} from "@/components/ui/accordion";
import {useApiClient} from "@/hooks/use-api-client";
import {useQueryClient} from "@tanstack/react-query";
import {ProviderDetailsForm, type ProviderFormValues} from "./ProviderDetailsForm";
import {PlansList} from "./PlansList";
import type {PlanFormValues} from "./PlanForm";
import type {PriceFormValues} from "./PriceForm";
import {OwnerType} from "@/models/ownerType.ts";

interface EditProviderFormProps {
    isOpen: boolean;
    onClose: () => void;
    provider: Provider;
}

export function EditProviderForm({isOpen, onClose, provider}: EditProviderFormProps) {
    const [error, setError] = useState<string | null>(null);
    const {updateProviderMutation} = useProvidersMutations();
    const {data: familiesData} = useFamiliesQuery({limit: 100});
    const {data: currencies, isLoading: isLoadingCurrencies} = useCurrenciesQuery();
    const families = familiesData?.families || [];
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // State for plans and prices management
    const [plans, setPlans] = useState(provider.plans);
    const [isSubmitting, setIsSubmitting] = useState(false);

    // Handle provider form submission
    const handleProviderSubmit = async (data: ProviderFormValues) => {
        try {
            setError(null);
            setIsSubmitting(true);
            await updateProviderMutation.mutateAsync({
                id: provider.id,
                etag: provider.etag,
                name: data.name,
                description: data.description || undefined,
                url: data.url || undefined,
                iconUrl: data.iconUrl || undefined,
                pricingPageUrl: data.pricingPageUrl || undefined,
                labels: provider.labels, // Preserve existing labels
                ownerType: provider.owner.type, // Always use the original owner type
                familyId: provider.owner.type === OwnerType.Family ? data.familyId : undefined,
            });
            onClose();
        } catch (err) {
            setError("Failed to update provider. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle adding a new plan
    const handleAddPlan = async (data: PlanFormValues) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);

            const response = await apiClient.providers.byProviderId(provider.id).plans.post({
                name: data.name,
                description: data.description || undefined
            });

            if (response) {
                // Add the new plan to the local state
                const updatedPlans = [...plans];
                const newPlan = response as unknown as Plan;
                updatedPlans.push(newPlan);
                setPlans(updatedPlans);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to add plan. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle updating a plan
    const handleUpdatePlan = async (planId: string, data: PlanFormValues) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);

            const planToUpdate = plans.find(p => p.id === planId);
            if (!planToUpdate) return;

            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).put({
                name: data.name,
                description: data.description || undefined
                // etag is not needed in the UpdatePlanModel
            });

            if (response) {
                // Update the plan in the local state
                const updatedPlans = plans.map(p => {
                    if (p.id === planId) {
                        return response as unknown as Plan;
                    }
                    return p;
                });
                setPlans(updatedPlans);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to update plan. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle deleting a plan
    const handleDeletePlan = async (planId: string) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);

            await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).delete();

            // Remove the plan from the local state
            setPlans(plans.filter(p => p.id !== planId));
            // Invalidate the providers query to refresh the data
            queryClient.invalidateQueries({queryKey: ['providers']});
        } catch (err) {
            setError("Failed to delete plan. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle adding a new price to a plan
    const handleAddPrice = async (planId: string, data: PriceFormValues) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);

            // Validate dates
            const startDate = new Date(data.startDate);
            if (isNaN(startDate.getTime())) {
                setError("Invalid start date. Please enter a valid date.");
                return;
            }

            let endDate = undefined;
            if (data.endDate) {
                const parsedEndDate = new Date(data.endDate);
                if (isNaN(parsedEndDate.getTime())) {
                    setError("Invalid end date. Please enter a valid date.");
                    return;
                }
                endDate = parsedEndDate;
            }

            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).prices.post({
                amount: data.amount,
                currency: data.currency,
                startDate: startDate,
                endDate: endDate
            });

            if (response) {
                // Add the new price to the plan in the local state
                const updatedPlans = plans.map(p => {
                    if (p.id === planId) {
                        // Create a new plan object with the updated prices array
                        const updatedPrices = [...p.prices, response as unknown as Price];
                        // Create a deep copy of the plan to maintain all its properties
                        const updatedPlan = new Plan(
                            p.id,
                            p.name,
                            p.description,
                            updatedPrices,
                            p.createdAt,
                            p.updatedAt,
                            p.etag
                        );
                        return updatedPlan;
                    }
                    return p;
                });
                setPlans(updatedPlans);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to add price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle updating a price
    const handleUpdatePrice = async (planId: string, priceId: string, data: PriceFormValues) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);

            const plan = plans.find(p => p.id === planId);
            if (!plan) return;

            const priceToUpdate = plan.prices.find(p => p.id === priceId);
            if (!priceToUpdate) return;

            // Validate dates
            const startDate = new Date(data.startDate);
            if (isNaN(startDate.getTime())) {
                setError("Invalid start date. Please enter a valid date.");
                return;
            }

            let endDate = undefined;
            if (data.endDate) {
                const parsedEndDate = new Date(data.endDate);
                if (isNaN(parsedEndDate.getTime())) {
                    setError("Invalid end date. Please enter a valid date.");
                    return;
                }
                endDate = parsedEndDate;
            }

            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).prices.byPriceId(priceId).put({
                amount: data.amount,
                currency: data.currency,
                startDate: startDate,
                endDate: endDate
                // etag is not needed in the UpdatePriceModel
            });

            if (response) {
                // Update the price in the local state
                const updatedPlans = plans.map(p => {
                    if (p.id === planId) {
                        // Update the specific price in the prices array
                        const updatedPrices = p.prices.map(price => {
                            if (price.id === priceId) {
                                return response as unknown as Price;
                            }
                            return price;
                        });
                        // Create a deep copy of the plan to maintain all its properties
                        const updatedPlan = new Plan(
                            p.id,
                            p.name,
                            p.description,
                            updatedPrices,
                            p.createdAt,
                            p.updatedAt,
                            p.etag
                        );
                        return updatedPlan;
                    }
                    return p;
                });
                setPlans(updatedPlans);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to update price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle deleting a price
    const handleDeletePrice = async (planId: string, priceId: string) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);

            await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).prices.byPriceId(priceId).delete();

            // Remove the price from the plan in the local state
            const updatedPlans = plans.map(p => {
                if (p.id === planId) {
                    // Filter out the deleted price
                    const updatedPrices = p.prices.filter(price => price.id !== priceId);
                    // Create a deep copy of the plan to maintain all its properties
                    const updatedPlan = new Plan(
                        p.id,
                        p.name,
                        p.description,
                        updatedPrices,
                        p.createdAt,
                        p.updatedAt,
                        p.etag
                    );
                    return updatedPlan;
                }
                return p;
            });
            setPlans(updatedPlans);
            // Invalidate the providers query to refresh the data
            queryClient.invalidateQueries({queryKey: ['providers']});
        } catch (err) {
            setError("Failed to delete price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setError(null);
        onClose();
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleClose}>
            <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Edit Provider</DialogTitle>
                    <DialogDescription>
                        Update provider details, plans, and prices.
                    </DialogDescription>
                </DialogHeader>

                {error && (
                    <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4"/>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}

                <Accordion type="single" collapsible defaultValue="provider-details">
                    <AccordionItem value="provider-details">
                        <AccordionTrigger>Provider Details</AccordionTrigger>
                        <AccordionContent>
                            <ProviderDetailsForm
                                provider={provider}
                                families={families}
                                onSubmit={handleProviderSubmit}
                                isSubmitting={isSubmitting}
                            />
                        </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="plans">
                        <AccordionTrigger>Plans and Prices</AccordionTrigger>
                        <AccordionContent>
                            <PlansList
                                plans={plans}
                                currencies={currencies || []}
                                isLoadingCurrencies={isLoadingCurrencies}
                                onAddPlan={handleAddPlan}
                                onUpdatePlan={handleUpdatePlan}
                                onDeletePlan={handleDeletePlan}
                                onAddPrice={handleAddPrice}
                                onUpdatePrice={handleUpdatePrice}
                                onDeletePrice={handleDeletePrice}
                                isSubmitting={isSubmitting}
                            />
                        </AccordionContent>
                    </AccordionItem>
                </Accordion>

                <DialogFooter>
                    <Button type="button" onClick={handleClose}>
                        Close
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}