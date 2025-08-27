import type {SubmitHandler} from "react-hook-form";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import * as z from "zod";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Label} from "@/components/ui/label";
import {Loader2} from "lucide-react";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import Price from "@/models/price";
import getSymbolFromCurrency from "currency-symbol-map";

// Define the type for form values
export type PriceFormValues = {
    amount: number;
    currency: string;
    startDate: string;
    endDate?: string;
};

interface PriceFormProps {
    price?: Price; // Optional - if provided, we're editing an existing price
    existingPrices?: Price[]; // All prices in the current plan for overlap validation
    currencies: string[];
    isLoadingCurrencies: boolean;
    onSubmit: SubmitHandler<PriceFormValues>;
    onCancel: () => void;
    isSubmitting: boolean;
}

// Helper function to check if date ranges overlap
const datesOverlap = (
    startDate1: Date,
    endDate1: Date | null,
    startDate2: Date,
    endDate2: Date | null
): boolean => {
    // If either range has no end date, it extends indefinitely
    if (endDate1 === null && endDate2 === null) {
        // Both ranges extend indefinitely, they must overlap
        return true;
    }

    if (endDate1 === null) {
        // Range 1 extends indefinitely, check if it starts before range 2 ends
        return startDate1 <= (endDate2 || new Date(8640000000000000));
    }

    if (endDate2 === null) {
        // Range 2 extends indefinitely, check if it starts before range 1 ends
        return startDate2 <= endDate1;
    }

    // Both ranges have end dates, check standard overlap
    return startDate1 <= endDate2 && startDate2 <= endDate1;
};

// Define the price form schema with custom refinement for date overlap validation
const createPriceFormSchema = (existingPrices: Price[] = [], currentPriceId?: string) => {
    return z.object({
        amount: z.number().min(0, "Amount must be a positive number").describe("Amount must be a valid number"),
        currency: z.string().min(1, "Currency is required"),
        startDate: z.string().min(1, "Start date is required"),
        endDate: z.string().optional(),
    })
        // First refinement: Validate end date is after start date
        .refine((data) => {
            if (!data.startDate || !data.endDate) {
                return true; // Skip validation if either date is missing
            }

            const startDate = new Date(data.startDate);
            const endDate = new Date(data.endDate);

            return startDate <= endDate;
        }, {
            message: "End date must be after start date",
            path: ["endDate"] // Show error on the end date field
        })
        // Second refinement: Check for overlapping date ranges
        .refine((data) => {
            // Skip validation if no currency or start date
            if (!data.currency || !data.startDate) {
                return true;
            }

            const newStartDate = new Date(data.startDate);
            const newEndDate = data.endDate ? new Date(data.endDate) : null;

            // Check for overlaps with existing prices of the same currency
            const overlappingPrice = existingPrices.find(price => {
                // Skip the current price being edited
                if (currentPriceId && price.id === currentPriceId) {
                    return false;
                }

                // Only check prices with the same currency
                if (price.amount.currency !== data.currency) {
                    return false;
                }

                // Check if date ranges overlap
                return datesOverlap(
                    newStartDate,
                    newEndDate,
                    price.startDate,
                    price.endDate
                );
            });

            return !overlappingPrice;
        }, {
            message: "Date range overlaps with an existing price for the same currency",
            path: ["startDate"] // Show error on the start date field
        });
};


export function PriceForm({
                              price,
                              existingPrices = [],
                              currencies,
                              isLoadingCurrencies,
                              onSubmit,
                              onCancel,
                              isSubmitting
                          }: PriceFormProps) {
    // Create a schema with validation for overlapping date ranges
    const validationSchema = createPriceFormSchema(existingPrices, price?.id);

    // Helper function to find the latest end date for a currency
    const getLatestEndDateForCurrency = (currency: string): string => {
        if (!existingPrices.length) return new Date().toISOString().split('T')[0];

        // Filter prices by currency and sort by end date (most recent first)
        const currencyPrices = existingPrices
            .filter(p => p.amount.currency === currency)
            .sort((a, b) => {
                // Handle null end dates (active prices)
                if (!a.endDate && !b.endDate) return 0;
                if (!a.endDate) return -1; // a is active, should come first
                if (!b.endDate) return 1;  // b is active, should come first

                // Both have end dates, compare them
                return b.endDate.getTime() - a.endDate.getTime();
            });

        // If no prices for this currency, return today
        if (currencyPrices.length === 0) return new Date().toISOString().split('T')[0];

        // Get the most recent price
        const latestPrice = currencyPrices[0];

        // If it has an end date, use that (plus one day), otherwise use today
        if (latestPrice.endDate) {
            const nextDay = new Date(latestPrice.endDate);
            nextDay.setDate(nextDay.getDate() + 1); // Add one day to avoid overlap
            return nextDay.toISOString().split('T')[0];
        }

        return new Date().toISOString().split('T')[0];
    };

    // Get default start date based on editing status and currency
    const getDefaultStartDate = (): string => {
        if (price?.startDate) {
            // If editing, use the existing start date
            return price.startDate.toISOString().split('T')[0];
        }

        // For new prices, use the latest end date for the default currency
        return getLatestEndDateForCurrency(price?.amount?.currency || "USD");
    };

    const {
        register,
        handleSubmit,
        setValue,
        watch,
        formState: {errors},
    } = useForm<PriceFormValues>({
        resolver: zodResolver(validationSchema),
        defaultValues: {
            amount: price?.amount?.value || 0,
            currency: price?.amount?.currency || "USD",
            startDate: getDefaultStartDate(),
            endDate: price?.endDate
                ? price.endDate.toISOString().split('T')[0]
                : "",
        },
    });

    const isEditing = !!price;
    const currentCurrency = watch("currency");
    const currencySymbol = getSymbolFromCurrency(currentCurrency) || '';

    return (
        <Card className="mt-2 border-dashed border-2">
            <CardHeader className="py-2">
                <CardTitle className="text-sm">{isEditing ? "Edit Price" : "Add New Price"}</CardTitle>
            </CardHeader>
            <CardContent className="py-2">
                <form onSubmit={handleSubmit(
                    (data) => onSubmit(data),
                    (errors) => {
                        console.error("Form validation failed:", errors);
                        // Form validation failed, errors object contains validation errors
                    }
                )} className="space-y-3">
                    <div className="space-y-2">
                        <Label htmlFor="price-amount" className="text-xs">Amount and Currency *</Label>
                        <div className="flex space-x-2">
                            <div className="flex-grow">
                                <Input
                                    id="price-amount"
                                    type="number"
                                    step="0.01"
                                    className={`h-8 ${errors.amount ? "border-red-500" : ""}`}
                                    {...register("amount", {valueAsNumber: true})}
                                />
                                {errors.amount && (
                                    <p className="text-xs text-red-500">{errors.amount.message}</p>
                                )}
                            </div>
                            <div className="w-1/3">
                                {isLoadingCurrencies ? (
                                    <div>
                                        <Input
                                            id="price-currency"
                                            className={`h-8 ${errors.currency ? "border-red-500" : ""}`}
                                            disabled
                                            value={currentCurrency}
                                            placeholder="Loading currencies..."
                                        />
                                        <input
                                            type="hidden"
                                            {...register("currency")}
                                            value={currentCurrency}
                                        />
                                    </div>
                                ) : (
                                    <Select
                                        onValueChange={(value) => {
                                            setValue("currency", value);

                                            // Only update start date for new prices (not when editing)
                                            if (!isEditing) {
                                                // Set start date based on the selected currency
                                                setValue("startDate", getLatestEndDateForCurrency(value));
                                            }
                                        }}
                                        defaultValue={currentCurrency}
                                    >
                                        <SelectTrigger className={`h-8 ${errors.currency ? "border-red-500" : ""}`}>
                                            <SelectValue placeholder="Select currency">
                                                {currentCurrency && `${currencySymbol} ${currentCurrency}`}
                                            </SelectValue>
                                        </SelectTrigger>
                                        <SelectContent>
                                            {currencies?.map((currency) => (
                                                <SelectItem key={currency} value={currency}>
                                                    {getSymbolFromCurrency(currency) || ''} {currency}
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                )}
                                {errors.currency && (
                                    <p className="text-xs text-red-500">{errors.currency.message}</p>
                                )}
                            </div>
                        </div>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="price-start-date" className="text-xs">Start Date *</Label>
                        <Input
                            id="price-start-date"
                            type="date"
                            className={`h-8 ${errors.startDate ? "border-red-500" : ""}`}
                            {...register("startDate")}
                        />
                        {errors.startDate && (
                            <p className="text-xs text-red-500">{errors.startDate.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="price-end-date" className="text-xs">End Date (Optional)</Label>
                        <Input
                            id="price-end-date"
                            type="date"
                            className={`h-8 ${errors.endDate ? "border-red-500" : ""}`}
                            {...register("endDate")}
                        />
                        {errors.endDate && (
                            <p className="text-xs text-red-500">{errors.endDate.message}</p>
                        )}
                    </div>

                    <div className="flex justify-end space-x-2">
                        <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={onCancel}
                        >
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            size="sm"
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? (
                                <>
                                    <Loader2 className="mr-1 h-3 w-3 animate-spin"/>
                                    Saving...
                                </>
                            ) : (
                                isEditing ? "Update" : "Add"
                            )}
                        </Button>
                    </div>
                </form>
            </CardContent>
        </Card>
    );
}