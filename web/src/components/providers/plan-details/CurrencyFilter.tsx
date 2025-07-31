import Price from "@/models/price";

interface CurrencyFilterProps {
    pricesByCurrency: Record<string, Price[]>;
    selectedCurrencies: string[];
    onCurrencyFilterToggle: (currency: string) => void;
    onClearFilters: () => void;
}

export function CurrencyFilter({
    pricesByCurrency,
    selectedCurrencies,
    onCurrencyFilterToggle,
    onClearFilters
}: CurrencyFilterProps) {
    if (Object.keys(pricesByCurrency).length === 0) {
        return null;
    }

    return (
        <div className="bg-gray-50 p-3 rounded-md mb-4">
            <div className="text-sm font-medium mb-2">Filter by Active Currencies:</div>
            <div className="flex flex-wrap gap-2">
                {Object.entries(pricesByCurrency).map(([currency, prices]) => {
                    const hasActivePrices = prices.some(p => p.isActive);
                    const isSelected = selectedCurrencies.includes(currency);
                    return hasActivePrices ? (
                        <button 
                            key={currency}
                            onClick={() => onCurrencyFilterToggle(currency)}
                            className={`px-3 py-1 border rounded-full text-sm flex items-center gap-1 transition-colors ${
                                isSelected 
                                    ? "bg-green-100 border-green-300 text-green-800" 
                                    : "bg-white hover:bg-gray-100"
                            }`}
                            aria-pressed={isSelected}
                            title={`${isSelected ? 'Remove' : 'Add'} ${currency} filter`}
                        >
                            <span className={`h-2 w-2 rounded-full ${isSelected ? "bg-green-600" : "bg-green-500"}`}></span>
                            <span>{currency}</span>
                        </button>
                    ) : null;
                })}
                {selectedCurrencies.length > 0 && (
                    <button 
                        onClick={onClearFilters}
                        className="px-3 py-1 bg-gray-100 border rounded-full text-sm hover:bg-gray-200"
                        title="Clear all filters"
                    >
                        Clear filters
                    </button>
                )}
            </div>
        </div>
    );
}