import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Loader2 } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import Plan from "@/models/plan";

// Define the plan form schema
const planFormSchema = z.object({
    name: z.string().min(1, "Plan name is required"),
    description: z.string().optional(),
});

export type PlanFormValues = z.infer<typeof planFormSchema>;

interface PlanFormProps {
    plan?: Plan; // Optional - if provided, we're editing an existing plan
    onSubmit: (data: PlanFormValues) => Promise<void>;
    onCancel: () => void;
    isSubmitting: boolean;
}

export function PlanForm({ plan, onSubmit, onCancel, isSubmitting }: PlanFormProps) {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<PlanFormValues>({
        resolver: zodResolver(planFormSchema),
        defaultValues: {
            name: plan?.name || "",
            description: plan?.description || "",
        },
    });

    const isEditing = !!plan;

    return (
        <Card className="mt-4 border-dashed border-2">
            <CardHeader className="pb-2">
                <CardTitle className="text-base">{isEditing ? "Edit Plan" : "Add New Plan"}</CardTitle>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="plan-name">Name *</Label>
                        <Input id="plan-name" {...register("name")} />
                        {errors.name && (
                            <p className="text-sm text-red-500">{errors.name.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="plan-description">Description</Label>
                        <Textarea id="plan-description" {...register("description")} />
                    </div>

                    <div className="flex justify-end space-x-2">
                        <Button type="button" variant="outline" onClick={onCancel}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                    Saving...
                                </>
                            ) : (
                                isEditing ? "Update Plan" : "Add Plan"
                            )}
                        </Button>
                    </div>
                </form>
            </CardContent>
        </Card>
    );
}