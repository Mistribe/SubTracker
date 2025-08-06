import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Input} from "@/components/ui/input";
import type {FormValues} from "./SubscriptionFormSchema";
import {Switch} from "@/components/ui/switch";
import {useEffect, useState} from "react";

export const DatesSection = () => {
    const form = useFormContext<FormValues>();
    const startDate = form.watch("startDate");
    const endDate = form.watch("endDate");
    const [hasEndDate, setHasEndDate] = useState<boolean>(!!endDate);

    // Update hasEndDate state when endDate changes
    useEffect(() => {
        setHasEndDate(!!endDate);
    }, [endDate]);

    const formatDateForInput = (date: Date | undefined | null): string => {
        if (!date) return "";
        const d = new Date(date);
        if (isNaN(d.getTime())) return "";
        return d.toISOString().split('T')[0];
    };

    const handleDateChange = (field: "startDate" | "endDate", value: string) => {
        if (value) {
            const date = new Date(value);
            if (!isNaN(date.getTime())) {
                form.setValue(field, date);
            }
        } else {
            form.setValue(field, undefined);
        }
    };

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">When does your subscription start?</h2>
                <p className="text-muted-foreground">Let us know the important dates for your subscription</p>
            </div>

            <div className="max-w-md mx-auto mt-6">
                <div className="grid grid-cols-1 gap-6">
                    <div>
                        <Label htmlFor="startDate" className="text-base mb-2 block">When did you start this
                            subscription?</Label>
                        <Input
                            id="startDate"
                            type="date"
                            value={formatDateForInput(startDate)}
                            onChange={(e) => handleDateChange("startDate", e.target.value)}
                            className="h-12"
                        />
                        {form.formState.errors.startDate && (
                            <p className="text-sm text-red-500 mt-1">{form.formState.errors.startDate.message}</p>
                        )}
                    </div>

                    <div>
                        <div className="flex items-center justify-between mb-2">
                            <Label htmlFor="hasEndDate" className="text-base">Does your subscription have an end
                                date?</Label>
                            <Switch
                                id="hasEndDate"
                                checked={hasEndDate}
                                onCheckedChange={(checked) => {
                                    setHasEndDate(checked);
                                    if (!checked) {
                                        form.setValue("endDate", undefined);
                                    }
                                }}
                            />
                        </div>

                        {hasEndDate && (
                            <>
                                <Input
                                    id="endDate"
                                    type="date"
                                    value={formatDateForInput(endDate)}
                                    onChange={(e) => handleDateChange("endDate", e.target.value)}
                                    className="h-12"
                                />
                                {form.formState.errors.endDate && (
                                    <p className="text-sm text-red-500 mt-1">{form.formState.errors.endDate.message}</p>
                                )}
                            </>
                        )}

                        <p className="text-xs text-muted-foreground mt-1">
                            {hasEndDate
                                ? "Select the end date of your subscription"
                                : "Your subscription doesn't have a specific end date"}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
};