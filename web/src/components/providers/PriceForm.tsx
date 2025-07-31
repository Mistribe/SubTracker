import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import Price from "@/models/price";
import getSymbolFromCurrency from "currency-symbol-map";

// Define the price form schema
const priceFormSchema = z.object({
    amount: z.preprocess(
        // Convert string to number and handle NaN
        (val) => {
            const parsed = parseFloat(val as string);
            return isNaN(parsed) ? undefined : parsed;
        },
        z.number()
            .min(0, "Amount must be a positive number")
            .describe("Amount is required")
    ),
    currency: z.string().min(1, "Currency is required"),
    startDate: z.string().min(1, "Start date is required"),
    endDate: z.string().optional(),
});

export type PriceFormValues = z.infer<typeof priceFormSchema>;

interface PriceFormProps {
    price?: Price; // Optional - if provided, we're editing an existing price
    currencies: string[];
    isLoadingCurrencies: boolean;
    onSubmit: (data: PriceFormValues) => Promise<void>;
    onCancel: () => void;
    isSubmitting: boolean;
}

export function PriceForm({ 
    price, 
    currencies, 
    isLoadingCurrencies, 
    onSubmit, 
    onCancel, 
    isSubmitting 
}: PriceFormProps) {
    const {
        register,
        handleSubmit,
        setValue,
        watch,
        formState: { errors },
    } = useForm<PriceFormValues>({
        resolver: zodResolver(priceFormSchema) as any, // Cast to any to resolve type incompatibility
        defaultValues: {
            amount: price?.amount || 0,
            currency: price?.currency || "USD",
            startDate: price?.startDate 
                ? price.startDate.toISOString().split('T')[0] 
                : new Date().toISOString().split('T')[0],
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
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
                    <div className="space-y-2">
                        <Label htmlFor="price-amount" className="text-xs">Amount and Currency *</Label>
                        <div className="flex space-x-2">
                            <div className="flex-grow">
                                <Input
                                    id="price-amount"
                                    type="number"
                                    step="0.01"
                                    className="h-8"
                                    {...register("amount", {valueAsNumber: true})}
                                />
                                {errors.amount && (
                                    <p className="text-xs text-red-500">{errors.amount.message}</p>
                                )}
                            </div>
                            <div className="w-1/3">
                                {isLoadingCurrencies ? (
                                    <Input 
                                        id="price-currency" 
                                        className="h-8" 
                                        {...register("currency")} 
                                    />
                                ) : (
                                    <Select
                                        onValueChange={(value) => setValue("currency", value)}
                                        defaultValue={currentCurrency}
                                    >
                                        <SelectTrigger className="h-8">
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
                            className="h-8" 
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
                            className="h-8" 
                            {...register("endDate")} 
                        />
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