import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import type {FormValues} from "./SubscriptionFormSchema";
import {Toggle} from "@/components/ui/toggle";

export const CustomPriceSection = () => {
    const form = useFormContext<FormValues>();
    const hasCustomPrice = form.watch("customPrice") !== undefined;

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Do you want to set a custom price?</h2>
                <p className="text-muted-foreground">Toggle below if you need to specify a different price than the one from the selected plan</p>
            </div>

            <div className="flex justify-center mt-4">
                <Toggle
                    id="hasCustomPrice"
                    pressed={hasCustomPrice}
                    onPressedChange={(pressed) => {
                        if (pressed) {
                            form.setValue("customPrice", {amount: 0, currency: "USD"});
                        } else {
                            form.setValue("customPrice", undefined);
                        }
                    }}
                    className="px-8 py-6 text-lg"
                    variant="outline"
                >
                    {hasCustomPrice ? "Yes, use custom price" : "No, use plan price"}
                </Toggle>
            </div>

            {hasCustomPrice && (
                <div className="max-w-md mx-auto mt-6">
                    <h3 className="text-lg font-medium mb-4">Custom Price Details</h3>
                    <div className="grid grid-cols-1 gap-6">
                        <div>
                            <Label htmlFor="customPriceAmount" className="text-base mb-2 block">How much does it cost?</Label>
                            <Input
                                id="customPriceAmount"
                                type="number"
                                step="0.01"
                                {...form.register("customPrice.amount", {valueAsNumber: true})}
                                placeholder="Enter amount"
                                className="h-12"
                            />
                            {form.formState.errors.customPrice?.amount && (
                                <p className="text-sm text-red-500 mt-1">
                                    {form.formState.errors.customPrice?.amount?.message}
                                </p>
                            )}
                        </div>

                        <div>
                            <Label htmlFor="customPriceCurrency" className="text-base mb-2 block">Which currency?</Label>
                            <Select
                                onValueChange={(value) => form.setValue("customPrice.currency", value)}
                                defaultValue="USD"
                                value={form.watch("customPrice.currency") || "USD"}
                            >
                                <SelectTrigger className="h-12">
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
                </div>
            )}
        </div>
    );
};