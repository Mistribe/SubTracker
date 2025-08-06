import { useFormContext } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { FormValues } from "./SubscriptionFormSchema";
import { SubscriptionRecurrency } from "@/models/subscriptionRecurrency";

export const RecurrencySection = () => {
    const form = useFormContext<FormValues>();
    const selectedRecurrency = form.watch("recurrency");
    
    return (
        <div className="space-y-4">
            <h2 className="text-xl font-semibold">Recurrency</h2>

            <div className="flex flex-col">
                <Label htmlFor="recurrency">Recurrency Type</Label>
                <div className="flex gap-4 items-start">
                    <div className={selectedRecurrency === "custom" ? "w-1/2" : "w-full"}>
                        <Select
                            onValueChange={(value) => form.setValue("recurrency", value as SubscriptionRecurrency)}
                            defaultValue={SubscriptionRecurrency.Monthly}
                            value={selectedRecurrency}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Select recurrency"/>
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value={SubscriptionRecurrency.Monthly}>Monthly</SelectItem>
                                <SelectItem value={SubscriptionRecurrency.Quarterly}>Quarterly</SelectItem>
                                <SelectItem value={SubscriptionRecurrency.HalfYearly}>Half Yearly</SelectItem>
                                <SelectItem value={SubscriptionRecurrency.Yearly}>Yearly</SelectItem>
                                <SelectItem value={SubscriptionRecurrency.OneTime}>One Time</SelectItem>
                                <SelectItem value="custom">Custom</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {selectedRecurrency === "custom" && (
                        <div className="w-1/2 flex flex-col gap-1">
                            <div className="flex gap-2">
                                <div className="w-1/2">
                                    <Label htmlFor="customRecurrencyValue" className="sr-only">Value</Label>
                                    <Input
                                        id="customRecurrencyValue"
                                        type="number"
                                        min="1"
                                        {...form.register("customRecurrencyValue", {valueAsNumber: true})}
                                        placeholder="Value"
                                    />
                                    {form.formState.errors.customRecurrencyValue && (
                                        <p className="text-sm text-red-500 mt-1">
                                            {form.formState.errors.customRecurrencyValue.message}
                                        </p>
                                    )}
                                </div>
                                <div className="w-1/2">
                                    <Label htmlFor="customRecurrencyUnit" className="sr-only">Unit</Label>
                                    <Select
                                        onValueChange={(value) => form.setValue("customRecurrencyUnit", value)}
                                        defaultValue="days"
                                        value={form.watch("customRecurrencyUnit")}
                                    >
                                        <SelectTrigger id="customRecurrencyUnit">
                                            <SelectValue placeholder="Unit"/>
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="days">Days</SelectItem>
                                            <SelectItem value="months">Months</SelectItem>
                                            <SelectItem value="years">Years</SelectItem>
                                        </SelectContent>
                                    </Select>
                                    {form.formState.errors.customRecurrencyUnit && (
                                        <p className="text-sm text-red-500 mt-1">
                                            {form.formState.errors.customRecurrencyUnit.message}
                                        </p>
                                    )}
                                </div>
                            </div>
                            <p className="text-xs text-muted-foreground">
                                Specify how often the subscription recurs. This will be converted to days
                                (30 days per month, 365 days per year).
                            </p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};