import { useState } from "react";
import { Button } from "@/components/ui/button";
import Price from "@/models/price";
import { PriceCard } from "./PriceCard";
import type {PriceFormValues} from "../PriceForm";
import { Plus } from "lucide-react";

interface PricesListProps {
    pricesByCurrency: Record<string, Price[]>;
    selectedCurrencies: string[];
    editingPrice: Price | null;
    showAddPriceForm: boolean;
    existingPrices: Price[];
    currencies: string[];
    isLoadingCurrencies: boolean;
    isSubmitting: boolean;
    onAddPrice: () => void;
    onEditPrice: (price: Price) => void;
    onDeletePrice: (price: Price) => void;
    onSubmitPriceForm: (data: PriceFormValues) => void;
    onCancelPriceForm: () => void;
}

export function PricesList({
    pricesByCurrency,
    selectedCurrencies,
    editingPrice,
    showAddPriceForm,
    existingPrices,
    currencies,
    isLoadingCurrencies,
    isSubmitting,
    onAddPrice,
    onEditPrice,
    onDeletePrice,
    onSubmitPriceForm,
    onCancelPriceForm
}: PricesListProps) {
    // State to track expanded currencies (for showing inactive prices)
    const [expandedCurrencies, setExpandedCurrencies] = useState<Record<string, boolean>>({});

    // Toggle expanded state for a currency
    const toggleCurrencyExpanded = (currency: string) => {
        setExpandedCurrencies(prev => ({
            ...prev,
            [currency]: !prev[currency]
        }));
    };

    if (Object.keys(pricesByCurrency).length === 0 && !showAddPriceForm) {
        return (
            <div className="text-center py-4 text-gray-500">
                No prices available for this plan.
            </div>
        );
    }

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h3 className="text-lg font-medium">Prices</h3>
                <Button
                    size="sm"
                    onClick={onAddPrice}
                    disabled={showAddPriceForm || !!editingPrice}
                >
                    <Plus className="h-4 w-4 mr-1"/> Add Price
                </Button>
            </div>

            <div className="space-y-4">
                {Object.entries(pricesByCurrency)
                    // Filter currencies based on selection (show all if none selected)
                    .filter(([currency]) => selectedCurrencies.length === 0 || selectedCurrencies.includes(currency))
                    .map(([currency, prices]) => {
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
                                        <PriceCard
                                            key={price.id}
                                            price={price}
                                            isActive={true}
                                            isEditing={editingPrice?.id === price.id}
                                            existingPrices={existingPrices}
                                            currencies={currencies}
                                            isLoadingCurrencies={isLoadingCurrencies}
                                            isSubmitting={isSubmitting}
                                            onEdit={onEditPrice}
                                            onDelete={onDeletePrice}
                                            onSubmit={onSubmitPriceForm}
                                            onCancel={onCancelPriceForm}
                                        />
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
                                                    <PriceCard
                                                        key={price.id}
                                                        price={price}
                                                        isActive={false}
                                                        isEditing={editingPrice?.id === price.id}
                                                        existingPrices={existingPrices}
                                                        currencies={currencies}
                                                        isLoadingCurrencies={isLoadingCurrencies}
                                                        isSubmitting={isSubmitting}
                                                        onEdit={onEditPrice}
                                                        onDelete={onDeletePrice}
                                                        onSubmit={onSubmitPriceForm}
                                                        onCancel={onCancelPriceForm}
                                                    />
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                )}
                            </div>
                        );
                    })}
            </div>
        </div>
    );
}