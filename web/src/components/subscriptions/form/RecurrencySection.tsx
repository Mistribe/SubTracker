import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {ToggleGroup, ToggleGroupItem} from "@/components/ui/toggle-group";
import type {FormValues} from "./SubscriptionFormSchema";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";

type CustomRecurrencyUnit = "days" | "months" | "years";

export const RecurrencySection = () => {
    const form = useFormContext<FormValues>();
    const selectedRecurrency = form.watch("recurrency");

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">How often do you pay for this subscription?</h2>
                <p className="text-muted-foreground">Select the billing frequency for your subscription</p>
            </div>

            <div className="max-w-md mx-auto mt-6">
                <div className="flex flex-col gap-6">
                    <div>
                        <Label htmlFor="recurrency" className="text-base mb-2 block">Billing frequency</Label>
                        <ToggleGroup
                            type="single"
                            defaultValue={SubscriptionRecurrency.Monthly}
                            value={selectedRecurrency}
                            onValueChange={(value) => value && form.setValue("recurrency", value as SubscriptionRecurrency)}
                            variant="outline"
                            className="w-full"
                        >
                            <ToggleGroupItem value={SubscriptionRecurrency.Monthly} className="flex-1 justify-center py-3">
                                Monthly
                            </ToggleGroupItem>
                            <ToggleGroupItem value={SubscriptionRecurrency.Quarterly} className="flex-1 justify-center py-3">
                                Quarterly
                            </ToggleGroupItem>
                            <ToggleGroupItem value={SubscriptionRecurrency.HalfYearly} className="flex-1 justify-center py-3">
                                Half Yearly
                            </ToggleGroupItem>
                            <ToggleGroupItem value={SubscriptionRecurrency.Yearly} className="flex-1 justify-center py-3">
                                Yearly
                            </ToggleGroupItem>
                            <ToggleGroupItem value={SubscriptionRecurrency.OneTime} className="flex-1 justify-center py-3">
                                One Time
                            </ToggleGroupItem>
                            <ToggleGroupItem value="custom" className="flex-1 justify-center py-3">
                                Custom
                            </ToggleGroupItem>
                        </ToggleGroup>
                    </div>

                    {selectedRecurrency === "custom" && (
                        <div className="bg-muted/30 p-4 rounded-lg">
                            <h3 className="text-lg font-medium mb-3">Custom Billing Frequency</h3>
                            <p className="text-sm text-muted-foreground mb-4">
                                Specify exactly how often this subscription recurs
                            </p>
                            <div className="flex gap-4">
                                <div className="w-1/2">
                                    <Label htmlFor="customRecurrencyValue" className="text-base mb-2 block">Every</Label>
                                    <Input
                                        id="customRecurrencyValue"
                                        type="number"
                                        min="1"
                                        {...form.register("customRecurrencyValue", {valueAsNumber: true})}
                                        placeholder="Value"
                                        className="h-12"
                                    />
                                    {form.formState.errors.customRecurrencyValue && (
                                        <p className="text-sm text-red-500 mt-1">
                                            {form.formState.errors.customRecurrencyValue.message}
                                        </p>
                                    )}
                                </div>
                                <div className="w-1/2">
                                    <Label htmlFor="customRecurrencyUnit" className="text-base mb-2 block">Unit</Label>
                                    <Select
                                        onValueChange={(value) => form.setValue("customRecurrencyUnit", value as CustomRecurrencyUnit)}
                                        defaultValue="days"
                                        value={form.watch("customRecurrencyUnit")}
                                    >
                                        <SelectTrigger id="customRecurrencyUnit" className="h-12">
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
                            <p className="text-xs text-muted-foreground mt-2">
                                This will be converted to days (30 days per month, 365 days per year).
                            </p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};