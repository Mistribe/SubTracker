import { useState } from "react";
import {
    Dialog,
    DialogContent,
    DialogFooter
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle, Trash2 } from "lucide-react";
import { useApiClient } from "@/hooks/use-api-client";
import { useQueryClient } from "@tanstack/react-query";
import { useCurrenciesQuery } from "@/hooks/currencies/useCurrenciesQuery";
import type {PriceFormValues} from "../PriceForm";
import Plan from "@/models/plan";
import Price from "@/models/price";

// Import the new components
import { PlanDetailsHeader } from "./PlanDetailsHeader";
import { CurrencyFilter } from "./CurrencyFilter";
import { PricesList } from "./PricesList";
import { DeletePriceDialog } from "./DeletePriceDialog";
import { DeletePlanDialog } from "./DeletePlanDialog";

interface PlanDetailsDialogProps {
    isOpen: boolean;
    onClose: () => void;
    plan: Plan;
    providerId: string;
    onPlanUpdated?: () => void;
}

export function PlanDetailsDialog({
    isOpen, 
    onClose, 
    plan, 
    providerId, 
    onPlanUpdated
}: PlanDetailsDialogProps) {
    // State management
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [showAddPriceForm, setShowAddPriceForm] = useState(false);
    const [editingPrice, setEditingPrice] = useState<Price | null>(null);
    const [deletingPrice, setDeletingPrice] = useState<Price | null>(null);
    const [deletingPlan, setDeletingPlan] = useState(false);
    const [selectedCurrencies, setSelectedCurrencies] = useState<string[]>([]);
    
    // Hooks
    const { apiClient } = useApiClient();
    const queryClient = useQueryClient();
    const { data: currencies, isLoading: isLoadingCurrencies } = useCurrenciesQuery();
    
    // Group prices by currency
    const pricesByCurrency = plan.prices.reduce((groups, price) => {
        const currency = price.amount.currency;
        if (!groups[currency]) {
            groups[currency] = [];
        }
        groups[currency].push(price);
        return groups;
    }, {} as Record<string, Price[]>);
    
    // Currency filter handlers
    const handleCurrencyFilterToggle = (currency: string) => {
        setSelectedCurrencies(prev => {
            if (prev.includes(currency)) {
                return prev.filter(c => c !== currency);
            } else {
                return [...prev, currency];
            }
        });
    };
    
    const handleClearFilters = () => {
        setSelectedCurrencies([]);
    };
    
    // Price form handlers
    const handleAddPriceClick = () => {
        setShowAddPriceForm(true);
        setEditingPrice(null);
    };
    
    const handleEditPriceClick = (price: Price) => {
        setEditingPrice(price);
        setShowAddPriceForm(false);
    };
    
    const handleDeletePriceClick = (price: Price) => {
        setDeletingPrice(price);
    };
    
    const handleCancelPriceForm = () => {
        setShowAddPriceForm(false);
        setEditingPrice(null);
        setError(null);
    };
    
    // Handle adding a new price
    const handleAddPrice = async (data: PriceFormValues) => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);
            setError(null);

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

            // Extract plan ID from the full plan ID if it contains a separator
            const planId = plan.id.includes('/') ? plan.id.split('/')[1] : plan.id;

            const response = await apiClient.providers.byProviderId(providerId)
                .plans.byPlanId(planId)
                .prices.post({
                    amount: data.amount,
                    currency: data.currency,
                    startDate: startDate,
                    endDate: endDate
                });

            // Create a new Price object from the response
            const newPrice = new Price(
                '',
                {
                    value: data.amount,
                    currency: data.currency,
                },
                startDate,
                endDate || null,
                new Date(),
                new Date(),
                response?.etag || ''
            );

            // Update the local plan state with the new price
            plan.prices = [...plan.prices, newPrice];

            // Invalidate the providers query to refresh the data
            await queryClient.invalidateQueries({queryKey: ['providers']});

            // Notify parent component that plan was updated
            if (onPlanUpdated) {
                onPlanUpdated();
            }

            // Hide the price form
            setShowAddPriceForm(false);
            setEditingPrice(null);
        } catch (err) {
            setError("Failed to add price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle updating an existing price
    const handleUpdatePrice = async (data: PriceFormValues) => {
        try {
            if (!apiClient || !editingPrice) return;
            setIsSubmitting(true);
            setError(null);

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

            // Extract plan ID from the full plan ID if it contains a separator
            const planId = plan.id.includes('/') ? plan.id.split('/')[1] : plan.id;

            const response = await apiClient.providers.byProviderId(providerId)
                .plans.byPlanId(planId)
                .prices.byPriceId(editingPrice.id)
                .put({
                    amount: data.amount,
                    currency: data.currency,
                    startDate: startDate,
                    endDate: endDate,
                    additionalData: {
                        etag: editingPrice.etag
                    }
                });

            // Update the local price object with the new data
            const updatedPrice = new Price(
                editingPrice.id,
                {
                    value:data.amount,
                    currency: data.currency,
                },
                startDate,
                endDate ?? null,
                response?.createdAt || new Date(),
                new Date(),
                response?.etag || ''
            );

            // Update the local plan state by replacing the edited price
            plan.prices = plan.prices.map(price =>
                price.id === editingPrice.id ? updatedPrice : price
            );

            // Invalidate the providers query to refresh the data
            await queryClient.invalidateQueries({queryKey: ['providers']});

            // Notify parent component that plan was updated
            if (onPlanUpdated) {
                onPlanUpdated();
            }

            // Hide the price form
            setEditingPrice(null);
        } catch (err) {
            setError("Failed to update price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle deleting a price
    const handleDeletePrice = async () => {
        try {
            if (!apiClient || !deletingPrice) return;
            setIsSubmitting(true);
            setError(null);

            // Extract plan ID from the full plan ID if it contains a separator
            const planId = plan.id.includes('/') ? plan.id.split('/')[1] : plan.id;

            await apiClient.providers.byProviderId(providerId)
                .plans.byPlanId(planId)
                .prices.byPriceId(deletingPrice.id)
                .delete();

            // Update the local plan state by removing the deleted price
            plan.prices = plan.prices.filter(price => price.id !== deletingPrice.id);

            // Invalidate the providers query to refresh the data
            await queryClient.invalidateQueries({queryKey: ['providers']});

            // Notify parent component that plan was updated
            if (onPlanUpdated) {
                onPlanUpdated();
            }

            // Reset the deleting price
            setDeletingPrice(null);
        } catch (err) {
            setError("Failed to delete price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handlePriceFormSubmit = (data: PriceFormValues) => {
        if (editingPrice) {
            handleUpdatePrice(data);
        } else {
            handleAddPrice(data);
        }
    };
    
    // Handle deleting a plan
    const handleDeletePlan = async () => {
        try {
            if (!apiClient) return;
            setIsSubmitting(true);
            setError(null);

            // Extract plan ID from the full plan ID if it contains a separator
            const planId = plan.id.includes('/') ? plan.id.split('/')[1] : plan.id;

            await apiClient.providers.byProviderId(providerId)
                .plans.byPlanId(planId)
                .delete();

            // Invalidate the providers query to refresh the data
            await queryClient.invalidateQueries({queryKey: ['providers']});

            // Close the dialog
            onClose();

            // Notify parent component that plan was updated
            if (onPlanUpdated) {
                onPlanUpdated();
            }
        } catch (err) {
            setError("Failed to delete plan. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
            setDeletingPlan(false);
        }
    };

    return (
        <>
            <Dialog open={isOpen} onOpenChange={onClose}>
                <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
                    <PlanDetailsHeader plan={plan} />

                    {error && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4"/>
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}

                    <div className="space-y-4 py-4">
                        {/* Currency Filter */}
                        {plan.prices.length > 0 && (
                            <CurrencyFilter 
                                pricesByCurrency={pricesByCurrency}
                                selectedCurrencies={selectedCurrencies}
                                onCurrencyFilterToggle={handleCurrencyFilterToggle}
                                onClearFilters={handleClearFilters}
                            />
                        )}

                        {/* Prices List */}
                        <PricesList 
                            pricesByCurrency={pricesByCurrency}
                            selectedCurrencies={selectedCurrencies}
                            editingPrice={editingPrice}
                            showAddPriceForm={showAddPriceForm}
                            existingPrices={plan.prices}
                            currencies={currencies || []}
                            isLoadingCurrencies={isLoadingCurrencies}
                            isSubmitting={isSubmitting}
                            onAddPrice={handleAddPriceClick}
                            onEditPrice={handleEditPriceClick}
                            onDeletePrice={handleDeletePriceClick}
                            onSubmitPriceForm={handlePriceFormSubmit}
                            onCancelPriceForm={handleCancelPriceForm}
                        />
                    </div>

                    <DialogFooter className="flex justify-between">
                        <Button 
                            variant="destructive" 
                            onClick={() => setDeletingPlan(true)}
                            disabled={isSubmitting}
                        >
                            <Trash2 className="h-4 w-4 mr-1"/> Delete Plan
                        </Button>
                        <Button onClick={onClose}>Close</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Delete Price Dialog */}
            <DeletePriceDialog 
                price={deletingPrice}
                isOpen={!!deletingPrice}
                isSubmitting={isSubmitting}
                onClose={() => setDeletingPrice(null)}
                onConfirm={handleDeletePrice}
            />
            
            {/* Delete Plan Dialog */}
            <DeletePlanDialog 
                plan={plan}
                isOpen={deletingPlan}
                isSubmitting={isSubmitting}
                onClose={() => setDeletingPlan(false)}
                onConfirm={handleDeletePlan}
            />
        </>
    );
}