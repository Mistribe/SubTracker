import {useState} from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog";
import {Button} from "@/components/ui/button";
import {Card, CardContent, CardFooter, CardHeader, CardTitle} from "@/components/ui/card";
import {Alert, AlertDescription} from "@/components/ui/alert";
import {AlertCircle, Pencil, Plus, Trash2} from "lucide-react";
import {useApiClient} from "@/hooks/use-api-client";
import {useQueryClient} from "@tanstack/react-query";
import {useCurrenciesQuery} from "@/hooks/currencies/useCurrenciesQuery";
import {PriceForm, type PriceFormValues} from "./PriceForm";
import Plan from "@/models/plan";
import Price from "@/models/price";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";

interface PlanDetailsDialogProps {
    isOpen: boolean;
    onClose: () => void;
    plan: Plan;
    providerId: string;
    onPlanUpdated?: () => void;
}

export function PlanDetailsDialog({isOpen, onClose, plan, providerId, onPlanUpdated}: PlanDetailsDialogProps) {
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [showAddPriceForm, setShowAddPriceForm] = useState(false);
    const [editingPrice, setEditingPrice] = useState<Price | null>(null);
    const [deletingPrice, setDeletingPrice] = useState<Price | null>(null);
    const [expandedCurrencies, setExpandedCurrencies] = useState<Record<string, boolean>>({});
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();
    const {data: currencies, isLoading: isLoadingCurrencies} = useCurrenciesQuery();
    
    // Group prices by currency
    const pricesByCurrency = plan.prices.reduce((groups, price) => {
        const currency = price.currency;
        if (!groups[currency]) {
            groups[currency] = [];
        }
        groups[currency].push(price);
        return groups;
    }, {} as Record<string, Price[]>);
    
    // Toggle expanded state for a currency
    const toggleCurrencyExpanded = (currency: string) => {
        setExpandedCurrencies(prev => ({
            ...prev,
            [currency]: !prev[currency]
        }));
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
                data.amount,
                data.currency,
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
                data.amount,
                data.currency,
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

    const handleCancelPriceForm = () => {
        setShowAddPriceForm(false);
        setEditingPrice(null);
        setError(null);
    };

    return (
        <>
            <Dialog open={isOpen} onOpenChange={onClose}>
                <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle>{plan.name}</DialogTitle>
                        {plan.description && (
                            <DialogDescription>
                                {plan.description}
                            </DialogDescription>
                        )}
                    </DialogHeader>

                    {error && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4"/>
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}

                    <div className="space-y-4 py-4">
                        <div className="flex justify-between items-center">
                            <h3 className="text-lg font-medium">Prices</h3>
                            <Button
                                size="sm"
                                onClick={() => {
                                    setShowAddPriceForm(true);
                                    setEditingPrice(null);
                                }}
                                disabled={showAddPriceForm || !!editingPrice}
                            >
                                <Plus className="h-4 w-4 mr-1"/> Add Price
                            </Button>
                        </div>
                        
                        {/* Active currencies summary */}
                        {plan.prices.length > 0 && (
                            <div className="bg-gray-50 p-3 rounded-md mb-4">
                                <div className="text-sm font-medium mb-2">Active Currencies:</div>
                                <div className="flex flex-wrap gap-2">
                                    {Object.entries(pricesByCurrency).map(([currency, prices]) => {
                                        const hasActivePrices = prices.some(p => p.isActive);
                                        return hasActivePrices ? (
                                            <div 
                                                key={currency}
                                                className="px-3 py-1 bg-white border rounded-full text-sm flex items-center gap-1"
                                            >
                                                <span className="h-2 w-2 rounded-full bg-green-500"></span>
                                                <span>{currency}</span>
                                            </div>
                                        ) : null;
                                    })}
                                </div>
                            </div>
                        )}

                        {showAddPriceForm && (
                            <PriceForm
                                existingPrices={plan.prices}
                                currencies={currencies || []}
                                isLoadingCurrencies={isLoadingCurrencies}
                                onSubmit={handlePriceFormSubmit}
                                onCancel={handleCancelPriceForm}
                                isSubmitting={isSubmitting}
                            />
                        )}

                        {plan.prices.length === 0 && !showAddPriceForm ? (
                            <div className="text-center py-4 text-gray-500">
                                No prices available for this plan.
                            </div>
                        ) : (
                            <div className="space-y-4">
                                {Object.entries(pricesByCurrency).map(([currency, prices]) => {
                                    // Sort prices: active first, then by start date (newest first)
                                    const sortedPrices = [...prices].sort((a, b) => {
                                        if (a.isActive !== b.isActive) {
                                            return a.isActive ? -1 : 1;
                                        }
                                        return b.startDate.getTime() - a.startDate.getTime();
                                    });
                                    
                                    // Split into active and inactive prices
                                    const activePrices = sortedPrices.filter(p => p.isActive);
                                    const inactivePrices = sortedPrices.filter(p => !p.isActive);
                                    const isExpanded = !!expandedCurrencies[currency];
                                    
                                    return (
                                        <div key={currency} className="border rounded-lg p-4 shadow-sm">
                                            <div className="font-medium text-lg mb-3 flex justify-between items-center border-b pb-2">
                                                <div className="flex items-center gap-2">
                                                    <span className="text-primary">{currency}</span>
                                                    <span className="text-xs text-gray-500">
                                                        ({activePrices.length} active, {inactivePrices.length} past)
                                                    </span>
                                                </div>
                                                {activePrices.length > 0 && (
                                                    <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full">
                                                        Active
                                                    </span>
                                                )}
                                            </div>
                                            
                                            {/* Active prices */}
                                            <div className="space-y-2 mb-3">
                                                {activePrices.map((price) => (
                                                    <Card key={price.id} className="overflow-hidden border-l-4 border-l-green-500 bg-green-50/30">
                                                        <CardHeader className="py-3">
                                                            <CardTitle className="text-base flex justify-between items-center">
                                                                <span className="font-medium">{price.amount} {price.currency}</span>
                                                                {editingPrice?.id === price.id ? (
                                                                    <span className="text-xs text-blue-500">Editing...</span>
                                                                ) : (
                                                                    <div className="flex space-x-1">
                                                                        <Button
                                                                            variant="ghost"
                                                                            size="icon"
                                                                            className="h-6 w-6"
                                                                            onClick={() => setEditingPrice(price)}
                                                                            disabled={!!editingPrice || showAddPriceForm}
                                                                        >
                                                                            <Pencil className="h-3 w-3"/>
                                                                        </Button>
                                                                        <Button
                                                                            variant="ghost"
                                                                            size="icon"
                                                                            className="h-6 w-6 text-red-500"
                                                                            onClick={() => setDeletingPrice(price)}
                                                                            disabled={!!editingPrice || showAddPriceForm}
                                                                        >
                                                                            <Trash2 className="h-3 w-3"/>
                                                                        </Button>
                                                                    </div>
                                                                )}
                                                            </CardTitle>
                                                        </CardHeader>
                                                        <CardContent className="py-2">
                                                            <div className="text-sm text-gray-500">
                                                                <div>From: {price.startDate.toLocaleDateString()}</div>
                                                                {price.endDate && (
                                                                    <div>To: {price.endDate.toLocaleDateString()}</div>
                                                                )}
                                                            </div>
                                                        </CardContent>
                                                        {editingPrice?.id === price.id && (
                                                            <CardFooter className="py-2 bg-gray-50">
                                                                <PriceForm
                                                                    price={price}
                                                                    existingPrices={plan.prices}
                                                                    currencies={currencies || []}
                                                                    isLoadingCurrencies={isLoadingCurrencies}
                                                                    onSubmit={handlePriceFormSubmit}
                                                                    onCancel={handleCancelPriceForm}
                                                                    isSubmitting={isSubmitting}
                                                                />
                                                            </CardFooter>
                                                        )}
                                                    </Card>
                                                ))}
                                            </div>
                                            
                                            {/* Inactive prices with toggle */}
                                            {inactivePrices.length > 0 && (
                                                <div>
                                                    <Button 
                                                        variant="outline" 
                                                        size="sm" 
                                                        onClick={() => toggleCurrencyExpanded(currency)}
                                                        className="w-full mb-2 text-gray-500"
                                                    >
                                                        {isExpanded ? "Hide" : "Show"} Past Prices ({inactivePrices.length})
                                                    </Button>
                                                    
                                                    {isExpanded && (
                                                        <div className="space-y-2 mt-2">
                                                            {inactivePrices.map((price) => (
                                                                <Card key={price.id} className="overflow-hidden border-l-4 border-l-gray-300">
                                                                    <CardHeader className="py-3">
                                                                        <CardTitle className="text-base flex justify-between items-center">
                                                                            <span>{price.amount} {price.currency}</span>
                                                                            {editingPrice?.id === price.id ? (
                                                                                <span className="text-xs text-blue-500">Editing...</span>
                                                                            ) : (
                                                                                <div className="flex space-x-1">
                                                                                    <Button
                                                                                        variant="ghost"
                                                                                        size="icon"
                                                                                        className="h-6 w-6"
                                                                                        onClick={() => setEditingPrice(price)}
                                                                                        disabled={!!editingPrice || showAddPriceForm}
                                                                                    >
                                                                                        <Pencil className="h-3 w-3"/>
                                                                                    </Button>
                                                                                    <Button
                                                                                        variant="ghost"
                                                                                        size="icon"
                                                                                        className="h-6 w-6 text-red-500"
                                                                                        onClick={() => setDeletingPrice(price)}
                                                                                        disabled={!!editingPrice || showAddPriceForm}
                                                                                    >
                                                                                        <Trash2 className="h-3 w-3"/>
                                                                                    </Button>
                                                                                </div>
                                                                            )}
                                                                        </CardTitle>
                                                                    </CardHeader>
                                                                    <CardContent className="py-2">
                                                                        <div className="text-sm text-gray-500">
                                                                            <div>From: {price.startDate.toLocaleDateString()}</div>
                                                                            {price.endDate && (
                                                                                <div>To: {price.endDate.toLocaleDateString()}</div>
                                                                            )}
                                                                        </div>
                                                                    </CardContent>
                                                                    {editingPrice?.id === price.id && (
                                                                        <CardFooter className="py-2 bg-gray-50">
                                                                            <PriceForm
                                                                                price={price}
                                                                                existingPrices={plan.prices}
                                                                                currencies={currencies || []}
                                                                                isLoadingCurrencies={isLoadingCurrencies}
                                                                                onSubmit={handlePriceFormSubmit}
                                                                                onCancel={handleCancelPriceForm}
                                                                                isSubmitting={isSubmitting}
                                                                            />
                                                                        </CardFooter>
                                                                    )}
                                                                </Card>
                                                            ))}
                                                        </div>
                                                    )}
                                                </div>
                                            )}
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>

                    <DialogFooter>
                        <Button onClick={onClose}>Close</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            <AlertDialog open={!!deletingPrice} onOpenChange={(open) => !open && setDeletingPrice(null)}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Price</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete this price? This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={handleDeletePrice}
                            disabled={isSubmitting}
                        >
                            Delete
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    );
}