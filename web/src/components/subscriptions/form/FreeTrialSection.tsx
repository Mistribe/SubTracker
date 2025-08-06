import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {format} from "date-fns";
import type {FormValues} from "./SubscriptionFormSchema";
import {Toggle} from "@/components/ui/toggle";

export const FreeTrialSection = () => {
    const form = useFormContext<FormValues>();
    const hasFreeTrialPeriod = form.watch("hasFreeTrialPeriod");

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Does this subscription include a free trial?</h2>
                <p className="text-muted-foreground">Toggle below if a free trial period is included with this subscription</p>
            </div>

            <div className="flex justify-center mt-4">
                <Toggle
                    id="hasFreeTrialPeriod"
                    pressed={hasFreeTrialPeriod}
                    onPressedChange={(pressed) => form.setValue("hasFreeTrialPeriod", pressed)}
                    className="px-8 py-6 text-lg"
                    variant="outline"
                >
                    {hasFreeTrialPeriod ? "Yes, includes free trial" : "No free trial"}
                </Toggle>
            </div>

            {hasFreeTrialPeriod && (
                <div className="max-w-md mx-auto mt-6">
                    <h3 className="text-lg font-medium mb-4">Free Trial Period Details</h3>
                    <div className="grid grid-cols-1 gap-6">
                        <div>
                            <Label htmlFor="freeTrialStartDate" className="text-base mb-2 block">When does the free trial start?</Label>
                            <Input
                                id="freeTrialStartDate"
                                type="date"
                                {...form.register("freeTrialStartDate", {
                                    valueAsDate: true,
                                })}
                                defaultValue={format(new Date(), "yyyy-MM-dd")}
                                className="h-12"
                            />
                            {form.formState.errors.freeTrialStartDate && (
                                <p className="text-sm text-red-500 mt-1">
                                    {form.formState.errors.freeTrialStartDate.message}
                                </p>
                            )}
                        </div>

                        <div>
                            <Label htmlFor="freeTrialEndDate" className="text-base mb-2 block">When does the free trial end?</Label>
                            <Input
                                id="freeTrialEndDate"
                                type="date"
                                {...form.register("freeTrialEndDate", {
                                    valueAsDate: true,
                                })}
                                className="h-12"
                            />
                            {form.formState.errors.freeTrialEndDate && (
                                <p className="text-sm text-red-500 mt-1">
                                    {form.formState.errors.freeTrialEndDate.message}
                                </p>
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};