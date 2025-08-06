import { useFormContext } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { FormValues } from "./SubscriptionFormSchema";

export const CustomPriceSection = () => {
    const form = useFormContext<FormValues>();
    const hasCustomPrice = form.watch("customPrice") !== undefined;
    
    return (
        <div className="space-y-4">
            <div className="flex items-center space-x-2">
                <input
                    type="checkbox"
                    id="hasCustomPrice"
                    checked={hasCustomPrice}
                    onChange={(e) => {
                        if (e.target.checked) {
                            form.setValue("customPrice", {amount: 0, currency: "USD"});
                        } else {
                            form.setValue("customPrice", undefined);
                        }
                    }}
                    className="h-4 w-4 rounded border-gray-300"
                />
                <Label htmlFor="hasCustomPrice">Add Custom Price</Label>
            </div>

            {hasCustomPrice && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pl-6">
                    <div>
                        <Label htmlFor="customPriceAmount">Amount</Label>
                        <Input
                            id="customPriceAmount"
                            type="number"
                            step="0.01"
                            {...form.register("customPrice.amount", {valueAsNumber: true})}
                            placeholder="Enter amount"
                        />
                        {form.formState.errors.customPrice?.amount && (
                            <p className="text-sm text-red-500 mt-1">
                                {form.formState.errors.customPrice?.amount?.message}
                            </p>
                        )}
                    </div>

                    <div>
                        <Label htmlFor="customPriceCurrency">Currency</Label>
                        <Select
                            onValueChange={(value) => form.setValue("customPrice.currency", value)}
                            defaultValue="USD"
                            value={form.watch("customPrice.currency") || "USD"}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Select currency"/>
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="USD">USD</SelectItem>
                                <SelectItem value="EUR">EUR</SelectItem>
                                <SelectItem value="GBP">GBP</SelectItem>
                                <SelectItem value="CAD">CAD</SelectItem>
                                <SelectItem value="AUD">AUD</SelectItem>
                            </SelectContent>
                        </Select>
                        {form.formState.errors.customPrice?.currency && (
                            <p className="text-sm text-red-500 mt-1">
                                {form.formState.errors.customPrice?.currency?.message}
                            </p>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
};