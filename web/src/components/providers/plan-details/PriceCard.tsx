import { Button } from "@/components/ui/button";
import { Card, CardFooter } from "@/components/ui/card";
import { Pencil, Trash2 } from "lucide-react";
import Price from "@/models/price";
import { PriceForm, type PriceFormValues } from "../PriceForm";
import { Money } from "@/components/ui/money";

interface PriceCardProps {
    price: Price;
    isActive: boolean;
    isEditing: boolean;
    existingPrices: Price[];
    currencies: string[];
    isLoadingCurrencies: boolean;
    isSubmitting: boolean;
    onEdit: (price: Price) => void;
    onDelete: (price: Price) => void;
    onSubmit: (data: PriceFormValues) => void;
    onCancel: () => void;
}

export function PriceCard({
    price,
    isActive,
    isEditing,
    existingPrices,
    currencies,
    isLoadingCurrencies,
    isSubmitting,
    onEdit,
    onDelete,
    onSubmit,
    onCancel
}: PriceCardProps) {
    // Determine card styling based on active status
    const cardClassName = isActive 
        ? "overflow-hidden border-l-4 border-l-green-500 bg-green-50/30" 
        : "overflow-hidden border-l-4 border-l-gray-300";

    return (
        <Card key={price.id} className={cardClassName}>
            <div className="p-2 flex items-center justify-between">
                <div>
                    <div className={isActive ? "font-medium" : ""}>
                        <span>
                            {/* Display in preferred currency, with original in parentheses when converted */}
                            <Money amount={price.amount} currency={price.currency} />
                        </span>
                    </div>
                    <div className="text-xs text-gray-500">
                        {price.endDate 
                            ? `${price.startDate.toLocaleDateString()} - ${price.endDate.toLocaleDateString()}`
                            : `From ${price.startDate.toLocaleDateString()}`}
                    </div>
                </div>
                {isEditing ? (
                    <span className="text-xs text-blue-500">Editing...</span>
                ) : (
                    <div className="flex space-x-1">
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6"
                            onClick={() => onEdit(price)}
                        >
                            <Pencil className="h-3 w-3"/>
                        </Button>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6 text-red-500"
                            onClick={() => onDelete(price)}
                        >
                            <Trash2 className="h-3 w-3"/>
                        </Button>
                    </div>
                )}
            </div>
            {isEditing && (
                <CardFooter className="py-2 bg-gray-50">
                    <PriceForm
                        price={price}
                        existingPrices={existingPrices}
                        currencies={currencies}
                        isLoadingCurrencies={isLoadingCurrencies}
                        onSubmit={onSubmit}
                        onCancel={onCancel}
                        isSubmitting={isSubmitting}
                    />
                </CardFooter>
            )}
        </Card>
    );
}