import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogTitle,
    DialogTrigger
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
import { Checkbox } from "@/components/ui/checkbox.tsx";
import { useState } from "react";
import { Button } from "@/components/ui/button.tsx";
import { Loader2, PlusIcon } from "lucide-react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useFamiliesMutations } from "@/hooks/families/useFamiliesMutations";
import { useFamilyQuotaQuery } from "@/hooks/families/useFamilyQuotaQuery";
import { useQuotaLimit, getQuotaTooltip } from "@/hooks/quotas/useFeature";
import { FeatureId } from "@/models/billing";
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";

const formSchema = z.object({
    name: z.string().min(1, "Family member name is required"),
    isKid: z.boolean(),
});

type FormValues = z.infer<typeof formSchema>;

interface AddFamilyMemberDialogProps {
    familyId: string;
    onSuccess?: () => void;
}

export function AddFamilyMemberDialog({ familyId, onSuccess }: AddFamilyMemberDialogProps) {
    const { addFamilyMemberMutation } = useFamiliesMutations();
    const [open, setOpen] = useState(false);

    // Check family members quota
    const { data: familyQuotas } = useFamilyQuotaQuery(true);
    const { enabled, canAdd, used, limit, remaining } = useQuotaLimit(
        familyQuotas,
        FeatureId.FamilyMembersCount
    );
    const isDisabled = !enabled || !canAdd;
    const tooltipMessage = getQuotaTooltip(enabled, canAdd, "family members");

    const form = useForm<FormValues>(
        {
            resolver: zodResolver(formSchema),
            defaultValues: {
                name: "",
                isKid: false,
            },
        });

    const onSubmit = async (values: FormValues) => {
        addFamilyMemberMutation.mutate({
            familyId,
            name: values.name,
            isKid: values.isKid
        }, {
            onSuccess: () => {
                setOpen(false);
                form.reset();

                if (onSuccess) {
                    onSuccess();
                }
            },
            onError: (error) => {
                console.error("Failed to add family member:", error);
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
                                    Add Member
                                    {enabled && limit !== undefined && (
                                        <span className="ml-2 text-xs opacity-70">
                                            ({used}/{limit})
                                        </span>
                                    )}
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
                <DialogTitle>Add a new family member</DialogTitle>
                <DialogDescription>Add a new family member to start attributing member to
                    subscription</DialogDescription>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Name</FormLabel>
                                    <FormControl>
                                        <Input placeholder="Enter member name" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="isKid"
                            render={({ field }) => (
                                <FormItem
                                    className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                        />
                                    </FormControl>
                                    <div className="space-y-1 leading-none">
                                        <FormLabel>Kid Account</FormLabel>
                                        <FormDescription>
                                            Is this member a kid?
                                        </FormDescription>
                                    </div>
                                </FormItem>
                            )}
                        />
                        <DialogFooter>
                            <Button type="submit" disabled={addFamilyMemberMutation.isPending}>
                                {addFamilyMemberMutation.isPending ? (
                                    <>
                                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                                        Adding...
                                    </>
                                ) : (
                                    "Add Family member"
                                )}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    )
}