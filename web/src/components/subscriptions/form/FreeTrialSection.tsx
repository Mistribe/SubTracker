import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {format} from "date-fns";
import type {FormValues} from "./SubscriptionFormSchema";

export const FreeTrialSection = () => {
    const form = useFormContext<FormValues>();
    const hasFreeTrialPeriod = form.watch("hasFreeTrialPeriod");

    return (
        <div className="space-y-4">
            <div className="flex items-center space-x-2">
                <input
                    type="checkbox"
                    id="hasFreeTrialPeriod"
                    {...form.register("hasFreeTrialPeriod")}
                    className="h-4 w-4 rounded border-gray-300"
                />
                <Label htmlFor="hasFreeTrialPeriod">Has Free Trial Period</Label>
            </div>

            {hasFreeTrialPeriod && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pl-6">
                    <div>
                        <Label htmlFor="freeTrialStartDate">Free Trial Start Date</Label>
                        <Input
                            id="freeTrialStartDate"
                            type="date"
                            {...form.register("freeTrialStartDate", {
                                valueAsDate: true,
                            })}
                            defaultValue={format(new Date(), "yyyy-MM-dd")}
                        />
                        {form.formState.errors.freeTrialStartDate && (
                            <p className="text-sm text-red-500 mt-1">
                                {form.formState.errors.freeTrialStartDate.message}
                            </p>
                        )}
                    </div>

                    <div>
                        <Label htmlFor="freeTrialEndDate">Free Trial End Date</Label>
                        <Input
                            id="freeTrialEndDate"
                            type="date"
                            {...form.register("freeTrialEndDate", {
                                valueAsDate: true,
                            })}
                        />
                        {form.formState.errors.freeTrialEndDate && (
                            <p className="text-sm text-red-500 mt-1">
                                {form.formState.errors.freeTrialEndDate.message}
                            </p>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
};