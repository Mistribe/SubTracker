import {useEffect} from "react";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import * as z from "zod";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Textarea} from "@/components/ui/textarea";
import {Label} from "@/components/ui/label";
import {OwnerType} from "@/models/ownerType";
import {Loader2} from "lucide-react";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import Provider from "@/models/provider";
import Family from "@/models/family";

// Define the form schema with Zod
const formSchema = z.object({
    name: z.string().min(1, "Name is required"),
    description: z.string().optional(),
    url: z.string().url("Must be a valid URL").optional().or(z.literal("")),
    iconUrl: z.string().url("Must be a valid URL").optional().or(z.literal("")),
    pricingPageUrl: z.string().url("Must be a valid URL").optional().or(z.literal("")),
    ownerType: z.enum([OwnerType.Personal, OwnerType.Family, OwnerType.System]),
    familyId: z.string().optional(),
}).refine(data => {
    // If owner type is Family, familyId is required
    return data.ownerType !== OwnerType.Family || (data.familyId && data.familyId.length > 0);
}, {
    message: "Family ID is required when owner type is Family",
    path: ["familyId"],
});

export type ProviderFormValues = z.infer<typeof formSchema>;

interface ProviderDetailsFormProps {
    provider: Provider;
    families: Family[];
    onSubmit: (data: ProviderFormValues) => Promise<void>;
    isSubmitting: boolean;
}

export function ProviderDetailsForm({provider, families, onSubmit, isSubmitting}: ProviderDetailsFormProps) {
    const {
        register,
        handleSubmit,
        watch,
        setValue,
        formState: {errors},
    } = useForm<ProviderFormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            name: provider.name,
            description: provider.description || "",
            url: provider.url || "",
            iconUrl: provider.iconUrl || "",
            pricingPageUrl: provider.pricingPageUrl || "",
            ownerType: provider.owner.type,
            familyId: provider.owner.type === OwnerType.Family ? provider.owner.familyId || undefined : undefined,
        },
    });

    const selectedOwnerType = watch("ownerType");

    // Set familyId automatically when there's only one family and Family is selected
    useEffect(() => {
        if (selectedOwnerType === OwnerType.Family && families.length === 1) {
            setValue("familyId", families[0].id);
        }
    }, [selectedOwnerType, families, setValue]);

    return (
        <form id="provider-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">
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
                <div className="p-2 border rounded-md text-sm">
                    {provider.owner.type === OwnerType.Personal ? "Personal" :
                        provider.owner.type === OwnerType.Family ? "Family" : "System"}
                    <input
                        type="hidden"
                        {...register("ownerType")}
                        value={provider.owner.type}
                    />
                </div>
                <p className="text-xs text-gray-500 italic">Owner type can only be set during creation</p>
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
                            defaultValue={provider.owner.familyId || undefined}
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
                    ) : null}
                    {errors.familyId && (
                        <p className="text-sm text-red-500">{errors.familyId.message}</p>
                    )}
                </div>
            )}

            <div className="flex justify-end mt-4">
                <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting ? (
                        <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                            Updating...
                        </>
                    ) : (
                        "Update Provider"
                    )}
                </Button>
            </div>
        </form>
    );
}