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
import {Input} from "@/components/ui/input.tsx";
import {Checkbox} from "@/components/ui/checkbox.tsx";
import {useState} from "react";
import {Button} from "@/components/ui/button.tsx";
import {Loader2, PlusIcon} from "lucide-react";
import {zodResolver} from "@hookform/resolvers/zod";
import {useForm} from "react-hook-form";
import {z} from "zod";
import {useApiClient} from "@/hooks/use-api-client.ts";
import type {CreateFamilyMemberModel} from "@/api/models";

const formSchema = z.object({
    name: z.string().min(1, "Family member name is required"),
    isKid: z.boolean(),
    email: z.email().or(z.literal("")).optional()
});

type FormValues = z.infer<typeof formSchema>;

interface AddFamilyMemberDialogProps {
    familyId: string;
    onSuccess?: () => void;
}

export function AddFamilyMemberDialog({familyId, onSuccess}: AddFamilyMemberDialogProps) {
    const {apiClient} = useApiClient();
    const [open, setOpen] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const form = useForm<FormValues>(
        {
            resolver: zodResolver(formSchema),
            defaultValues: {
                name: "",
                isKid: false,
                email: "",
            },
        });

    const onSubmit = async (values: FormValues) => {
        if (!apiClient) {
            return;
        }

        setIsSubmitting(true);

        try {
            const familyMember: CreateFamilyMemberModel = {
                name: values.name,
                email: values.email,
                isKid: values.isKid,
            }

            await apiClient.families.byFamilyId(familyId).members.post(familyMember);
            setOpen(false);
            form.reset();

            if (onSuccess) {
                onSuccess();
            }
        } catch (error) {
            console.error("Failed to add family member:", error);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button>
                    <PlusIcon className="h-4 w-4 mr-2"/>
                    Add Member
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <DialogTitle>Add a new family member</DialogTitle>
                <DialogDescription>Add a new family member to start attributing member to
                    subscription</DialogDescription>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Name</FormLabel>
                                    <FormControl>
                                        <Input placeholder="Enter member name" {...field} />
                                    </FormControl>
                                    <FormMessage/>
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="isKid"
                            render={({field}) => (
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
                        <FormField
                            control={form.control}
                            name="email"
                            render={({field}) => (
                                <FormItem>
                                    <FormLabel>Email (Optional)</FormLabel>
                                    <FormControl>
                                        <Input type="email" placeholder="Enter email address" {...field} />
                                    </FormControl>
                                    <FormDescription>
                                        If email is set, the user will receive an invitation to join this family.
                                    </FormDescription>
                                    <FormMessage/>
                                </FormItem>
                            )}
                        />
                        <DialogFooter>
                            <Button type="submit" disabled={isSubmitting}>
                                {isSubmitting ? (
                                    <>
                                        <Loader2 className="h-4 w-4 mr-2 animate-spin"/>
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