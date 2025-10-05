import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog.tsx";
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Loader2, PlusIcon } from "lucide-react";
import { useFamiliesMutations } from "@/hooks/families/useFamiliesMutations";
import { useUser } from "@clerk/clerk-react";
import { useAccountQuotaQuery } from "@/hooks/profile/useAccountQuotaQuery";
import { useFeature, getQuotaTooltip } from "@/hooks/quotas/useFeature";
import { FeatureId } from "@/models/billing";
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";

// Define the form schema
const formSchema = z.object({
    name: z.string().min(1, "Family name is required"),
    creatorName: z.string().min(1, "Creator name is required"),
});

type FormValues = z.infer<typeof formSchema>;

interface CreateFamilyDialogProps {
    onSuccess?: () => void;
}

export function CreateFamilyDialog({ onSuccess }: CreateFamilyDialogProps) {
    const { createFamilyMutation } = useFamiliesMutations();
    const { user } = useUser();
    const [open, setOpen] = useState(false);

    // Check if the user has the Family feature enabled
    const { data: accountQuotas } = useAccountQuotaQuery();
    const { enabled: familyFeatureEnabled } = useFeature(accountQuotas, FeatureId.Family);
    const isDisabled = !familyFeatureEnabled;
    const tooltipMessage = getQuotaTooltip(familyFeatureEnabled, familyFeatureEnabled, "a family");

    let defaultCreatorName = ""
    if (user && user.firstName) {
        defaultCreatorName = user.firstName
    }
    // Initialize the form
    const form = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: "",
            creatorName: defaultCreatorName
        },
    });

    // Handle form submission
    const onSubmit = async (values: FormValues) => {
        createFamilyMutation.mutate({
            name: values.name,
            creatorName: values.creatorName
        }, {
            onSuccess: () => {
                // Close the dialog and reset the form
                setOpen(false);
                form.reset();

                // Call the onSuccess callback if provided
                if (onSuccess) {
                    onSuccess();
                }
            },
            onError: (error) => {
                console.error("Failed to create family:", error);
            }
        });
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <TooltipProvider>
                <Tooltip>
                    <TooltipTrigger asChild>
                        <span>
                            <DialogTrigger asChild>
                                <Button disabled={isDisabled}>
                                    <PlusIcon className="h-4 w-4 mr-2" />
                                    Add Family
                                </Button>
                            </DialogTrigger>
                        </span>
                    </TooltipTrigger>
                    {tooltipMessage && (
                        <TooltipContent>
                            <p>{tooltipMessage}</p>
                        </TooltipContent>
                    )}
                </Tooltip>
            </TooltipProvider>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Create a new family</DialogTitle>
                    <DialogDescription>
                        Create a family to start tracking your recurring payments.
                    </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Family Name</FormLabel>
                                    <FormControl>
                                        <Input placeholder="Enter family name" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        This is the name of your family.
                                    </FormDescription>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="creatorName"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Creator Name</FormLabel>
                                    <FormControl>
                                        <Input placeholder="Enter your name" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        Your name as the creator of this family.
                                    </FormDescription>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <DialogFooter>
                            <Button type="submit" disabled={createFamilyMutation.isPending}>
                                {createFamilyMutation.isPending ? (
                                    <>
                                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                                        Creating...
                                    </>
                                ) : (
                                    "Create Family"
                                )}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
}