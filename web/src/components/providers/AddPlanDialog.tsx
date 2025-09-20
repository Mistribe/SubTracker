import { useState } from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogFooter
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { useApiClient } from "@/hooks/use-api-client";
import { useQueryClient } from "@tanstack/react-query";
import { useCurrenciesQuery } from "@/hooks/currencies/useCurrenciesQuery";
import { PlanForm, type PlanFormValues } from "./PlanForm";
import { PriceForm, type PriceFormValues } from "./PriceForm";
import Plan from "@/models/plan";

interface AddPlanDialogProps {
    isOpen: boolean;
    onClose: () => void;
    providerId: string;
}

export function AddPlanDialog({ isOpen, onClose, providerId }: AddPlanDialogProps) {
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [createdPlan, setCreatedPlan] = useState<Plan | null>(null);
    const [showPriceForm, setShowPriceForm] = useState(false);
    const { apiClient } = useApiClient();
    const queryClient = useQueryClient();
    const { data: currencies, isLoading: isLoadingCurrencies } = useCurrenciesQuery();

    // Handle adding a new plan
    const handleAddPlan = async (data: PlanFormValues) => {
        try {
            setIsSubmitting(true);
            setError(null);

            // The current OpenAPI client does not expose endpoints to create plans.
            // Display an informative error for now to avoid build errors.
            setError("Adding plans is not supported with the current API version.");
            return;
        } catch (err) {
            setError("Failed to add plan. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    // Handle adding a price to the newly created plan
    const handleAddPrice = async (data: PriceFormValues) => {
        try {
            setIsSubmitting(true);
            setError(null);

            // The current OpenAPI client does not expose endpoints to manage plan prices.
            setError("Adding prices to plans is not supported with the current API version.");
            return;
        } catch (err) {
            setError("Failed to add price. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setError(null);
        setCreatedPlan(null);
        setShowPriceForm(false);
        onClose();
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleClose}>
            <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{showPriceForm ? "Add Price to Plan" : "Add New Plan"}</DialogTitle>
                    <DialogDescription>
                        {showPriceForm 
                            ? `Add pricing details for ${createdPlan?.name}`
                            : "Create a new plan for this provider."}
                    </DialogDescription>
                </DialogHeader>

                {error && (
                    <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4" />
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}

                {!showPriceForm ? (
                    <PlanForm
                        onSubmit={handleAddPlan}
                        onCancel={handleClose}
                        isSubmitting={isSubmitting}
                    />
                ) : createdPlan && (
                    <div className="py-4">
                        <PriceForm
                            existingPrices={[]}
                            currencies={currencies || []}
                            isLoadingCurrencies={isLoadingCurrencies}
                            onSubmit={handleAddPrice}
                            onCancel={handleClose}
                            isSubmitting={isSubmitting}
                        />
                    </div>
                )}

                <DialogFooter>
                    <Button type="button" variant="outline" onClick={handleClose}>
                        {showPriceForm ? "Skip Adding Price" : "Cancel"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}