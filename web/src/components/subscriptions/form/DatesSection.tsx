import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import {format, isValid} from "date-fns";
import type {FormValues} from "./SubscriptionFormSchema";

export const DatesSection = () => {
    const form = useFormContext<FormValues>();
    const startDate = form.watch("startDate");
    const endDate = form.watch("endDate");

    const formatDateForInput = (date: Date | null | undefined): string => {
        if (!date) return "";

        // Handle case where date might be a string
        const dateObj = date instanceof Date ? date : new Date(date);

        // Check if the date is valid
        if (!isValid(dateObj)) return "";

        return format(dateObj, "yyyy-MM-dd");
    };

    return (
        <div className="space-y-4">
            <h2 className="text-xl font-semibold">Dates</h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <Label htmlFor="startDate">Start Date</Label>
                    <Input
                        id="startDate"
                        type="date"
                        {...form.register("startDate", {
                            valueAsDate: true,
                        })}
                        value={formatDateForInput(startDate)}
                    />
                    {form.formState.errors.startDate && (
                        <p className="text-sm text-red-500 mt-1">{form.formState.errors.startDate.message}</p>
                    )}
                </div>

                <div>
                    <Label htmlFor="endDate">End Date (Optional)</Label>
                    <Input
                        id="endDate"
                        type="date"
                        {...form.register("endDate", {
                            valueAsDate: true,
                        })}
                        value={formatDateForInput(endDate)}
                    />
                    {form.formState.errors.endDate && (
                        <p className="text-sm text-red-500 mt-1">{form.formState.errors.endDate.message}</p>
                    )}
                </div>
            </div>
        </div>
    );
};