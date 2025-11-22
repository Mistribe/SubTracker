import {useEffect, useState} from "react";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import * as z from "zod";
import {useProvidersMutations} from "@/hooks/providers/useProvidersMutations";
import {useFamilyQuery} from "@/hooks/families/useFamilyQuery.ts";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Textarea} from "@/components/ui/textarea";
import {Label} from "@/components/ui/label";
import {RadioGroup, RadioGroupItem} from "@/components/ui/radio-group";
import {OwnerType} from "@/models/ownerType";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue,} from "@/components/ui/select";
import {Alert, AlertDescription} from "@/components/ui/alert";
import {AlertCircle, Loader2} from "lucide-react";

// Define the form schema with Zod
const formSchema = z.object({
    name: z.string().min(1, "Name is required"),
    description: z.string().optional(),
    url: z.string().url("Must be a valid URL").optional().or(z.literal("")),
    iconUrl: z.string().url("Must be a valid URL").optional().or(z.literal("")),
    pricingPageUrl: z.string().url("Must be a valid URL").optional().or(z.literal("")),
    ownerType: z.enum([OwnerType.Personal, OwnerType.Family, OwnerType.System]),
    familyId: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;

interface AddProviderFormProps {
    isOpen: boolean;
    onClose: () => void;
    familyId?: string;
}

export function AddProviderForm({isOpen, onClose, familyId}: AddProviderFormProps) {
    const [error, setError] = useState<string | null>(null);
    const {createProviderMutation} = useProvidersMutations();
    const {data: familyData} = useFamilyQuery();
    const families = familyData ? [familyData] : [];

    const {
        register,
        handleSubmit,
        reset,
        watch,
        setValue,
        formState: {errors, isSubmitting},
    } = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: "",
            description: "",
            url: "",
            iconUrl: "",
            pricingPageUrl: "",
            ownerType: familyId ? OwnerType.Family : OwnerType.Personal,
            familyId: familyId || undefined,
        },
    });

    const selectedOwnerType = watch("ownerType");

    // Set familyId automatically when there's only one family and Family is selected
    useEffect(() => {
        if (selectedOwnerType === OwnerType.Family && families.length === 1) {
            setValue("familyId", families[0].id);
        }
    }, [selectedOwnerType, families, setValue]);

    const onSubmit = async (data: FormValues) => {
        try {
            setError(null);
            await createProviderMutation.mutateAsync({
                name: data.name,
                description: data.description || undefined,
                url: data.url || undefined,
                iconUrl: data.iconUrl || undefined,
                pricingPageUrl: data.pricingPageUrl || undefined,
                ownerType: data.ownerType,
            });
            reset();
            onClose();
        } catch (err) {
            setError("Failed to create provider. Please try again.");
            console.error(err);
        }
    };

    const handleClose = () => {
        reset();
        setError(null);
        onClose();
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleClose}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Add New Provider</DialogTitle>
                    <DialogDescription>
                        Create a new service provider. Fill in the details below.
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    {error && (
                        <Alert variant="destructive">
                            <AlertCircle className="h-4 w-4"/>
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}

                    <div className="space-y-2">
                        <Label htmlFor="name">Name *</Label>
                        <Input id="name" {...register("name")} />
                        {errors.name && (
                            <p className="text-sm text-red-500">{errors.name.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="description">Description</Label>
                        <Textarea id="description" {...register("description")} />
                        {errors.description && (
                            <p className="text-sm text-red-500">{errors.description.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="url">Website URL</Label>
                        <Input id="url" {...register("url")} placeholder="https://"/>
                        {errors.url && (
                            <p className="text-sm text-red-500">{errors.url.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="iconUrl">Icon URL</Label>
                        <Input id="iconUrl" {...register("iconUrl")} placeholder="https://"/>
                        {errors.iconUrl && (
                            <p className="text-sm text-red-500">{errors.iconUrl.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="pricingPageUrl">Pricing Page URL</Label>
                        <Input id="pricingPageUrl" {...register("pricingPageUrl")} placeholder="https://"/>
                        {errors.pricingPageUrl && (
                            <p className="text-sm text-red-500">{errors.pricingPageUrl.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label>Owner Type</Label>
                        <RadioGroup
                            defaultValue={familyId ? OwnerType.Family : OwnerType.Personal}
                            onValueChange={(value) => setValue("ownerType", value as OwnerType)}
                            value={watch("ownerType")}
                            className="flex flex-col space-y-1"
                        >
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="personal" id="personal"/>
                                <Label htmlFor="personal">Personal</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="family" id="family"/>
                                <Label htmlFor="family">Family</Label>
                            </div>
                            {/* System providers can only be created by administrators */}
                        </RadioGroup>
                        {errors.ownerType && (
                            <p className="text-sm text-red-500">{errors.ownerType.message}</p>
                        )}
                    </div>

                    {selectedOwnerType === OwnerType.Family && (
                        <div className="space-y-2">
                            <Label htmlFor="familyId">Family</Label>
                            {families.length > 1 ? (
                                <Select
                                    onValueChange={(value) => setValue("familyId", value)}
                                    defaultValue={familyId}
                                >
                                    <SelectTrigger className="w-full">
                                        <SelectValue placeholder="Select a family"/>
                                    </SelectTrigger>
                                    <SelectContent>
                                        {families.map((family) => (
                                            <SelectItem key={family.id} value={family.id}>
                                                {family.name}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            ) : families.length === 1 ? (
                                <div className="p-2 border rounded-md">
                                    {families[0].name}
                                    <input
                                        type="hidden"
                                        {...register("familyId")}
                                        value={families[0].id}
                                    />
                                </div>
                            ) : !familyId ? (
                                <Input id="familyId" {...register("familyId")} />
                            ) : null}
                            {errors.familyId && (
                                <p className="text-sm text-red-500">{errors.familyId.message}</p>
                            )}
                        </div>
                    )}

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={handleClose}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                    Creating...
                                </>
                            ) : (
                                "Create Provider"
                            )}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}