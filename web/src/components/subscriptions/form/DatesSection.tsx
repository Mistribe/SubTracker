import { useFormContext } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { format } from "date-fns";
import { FormValues } from "./SubscriptionFormSchema";

export const DatesSection = () => {
    const form = useFormContext<FormValues>();
    
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
                        defaultValue={format(new Date(), "yyyy-MM-dd")}
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
                    />
                    {form.formState.errors.endDate && (
                        <p className="text-sm text-red-500 mt-1">{form.formState.errors.endDate.message}</p>
                    )}
                </div>
            </div>
        </div>
    );
};