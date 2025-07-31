import {useEffect, useState} from "react";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import * as z from "zod";
import {useProvidersMutations} from "@/hooks/providers/useProvidersMutations";
import {useFamiliesQuery} from "@/hooks/families/useFamiliesQuery";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Textarea} from "@/components/ui/textarea";
import {Label} from "@/components/ui/label";
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
import {AlertCircle, Edit2, Loader2, Plus, Trash2} from "lucide-react";
import Provider from "@/models/provider";
import Plan from "@/models/plan";
import Price from "@/models/price";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Accordion, AccordionContent, AccordionItem, AccordionTrigger} from "@/components/ui/accordion";
import {useApiClient} from "@/hooks/use-api-client";
import {useQueryClient} from "@tanstack/react-query";

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

// Define the plan form schema
const planFormSchema = z.object({
    name: z.string().min(1, "Plan name is required"),
    description: z.string().optional(),
});

// Define the price form schema
const priceFormSchema = z.object({
    amount: z.preprocess(
        // Convert string to number and handle NaN
        (val) => {
            const parsed = parseFloat(val as string);
            return isNaN(parsed) ? undefined : parsed;
        },
        z.number()
            .min(0, "Amount must be a positive number")
            .describe("Amount is required")
    ),
    currency: z.string().min(1, "Currency is required"),
    startDate: z.string().min(1, "Start date is required"),
    endDate: z.string().optional(),
});

type FormValues = z.infer<typeof formSchema>;
type PlanFormValues = z.infer<typeof planFormSchema>;
type PriceFormValues = z.infer<typeof priceFormSchema>;

interface EditProviderFormProps {
    isOpen: boolean;
    onClose: () => void;
    provider: Provider;
}

export function EditProviderForm({isOpen, onClose, provider}: EditProviderFormProps) {
    const [error, setError] = useState<string | null>(null);
    const {updateProviderMutation} = useProvidersMutations();
    const {data: familiesData} = useFamiliesQuery({limit: 100});
    const families = familiesData?.families || [];
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // State for plans and prices management
    const [plans, setPlans] = useState(provider.plans);
    const [editingPlanId, setEditingPlanId] = useState<string | null>(null);
    const [editingPriceInfo, setEditingPriceInfo] = useState<{ planId: string, priceId: string } | null>(null);
    const [showAddPlanForm, setShowAddPlanForm] = useState(false);
    const [showAddPriceForm, setShowAddPriceForm] = useState<string | null>(null); // planId or null

    // Form for provider details
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
            name: provider.name,
            description: provider.description || "",
            url: provider.url || "",
            iconUrl: provider.iconUrl || "",
            pricingPageUrl: provider.pricingPageUrl || "",
            ownerType: provider.owner.type,
            familyId: provider.owner.type === OwnerType.Family ? provider.owner.familyId || undefined : undefined,
        },
    });

    // Form for plan details
    const planForm = useForm<PlanFormValues>({
        resolver: zodResolver(planFormSchema),
        defaultValues: {
            name: "",
            description: "",
        },
    });

    // Form for price details
    const priceForm = useForm<PriceFormValues>({
        resolver: zodResolver(priceFormSchema) as any, // Cast to any to resolve type incompatibility
        defaultValues: {
            amount: 0,
            currency: "USD",
            startDate: new Date().toISOString().split('T')[0],
            endDate: "",
        },
    });

    const selectedOwnerType = watch("ownerType");

    // Set familyId automatically when there's only one family and Family is selected
    useEffect(() => {
        if (selectedOwnerType === OwnerType.Family && families.length === 1) {
            setValue("familyId", families[0].id);
        }
    }, [selectedOwnerType, families, setValue]);

    // Handle provider form submission
    const onSubmit = async (data: FormValues) => {
        try {
            setError(null);
            await updateProviderMutation.mutateAsync({
                id: provider.id,
                etag: provider.etag,
                name: data.name,
                description: data.description || undefined,
                url: data.url || undefined,
                iconUrl: data.iconUrl || undefined,
                pricingPageUrl: data.pricingPageUrl || undefined,
                labels: provider.labels, // Preserve existing labels
                ownerType: provider.owner.type, // Always use the original owner type
                familyId: provider.owner.type === OwnerType.Family ? data.familyId : undefined,
            });
            reset();
            onClose();
        } catch (err) {
            setError("Failed to update provider. Please try again.");
            console.error(err);
        }
    };

    // Handle adding a new plan
    const handleAddPlan = async (data: PlanFormValues) => {
        try {
            if (!apiClient) return;

            const response = await apiClient.providers.byProviderId(provider.id).plans.post({
                name: data.name,
                description: data.description || undefined
            });

            if (response) {
                // Add the new plan to the local state
                const updatedPlans = [...plans];
                const newPlan = response as unknown as Plan;
                updatedPlans.push(newPlan);
                setPlans(updatedPlans);
                planForm.reset();
                setShowAddPlanForm(false);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to add plan. Please try again.");
            console.error(err);
        }
    };

    // Handle updating a plan
    const handleUpdatePlan = async (data: PlanFormValues) => {
        try {
            if (!apiClient || !editingPlanId) return;

            const planToUpdate = plans.find(p => p.id === editingPlanId);
            if (!planToUpdate) return;

            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(editingPlanId).put({
                name: data.name,
                description: data.description || undefined
                // etag is not needed in the UpdatePlanModel
            });

            if (response) {
                // Update the plan in the local state
                const updatedPlans = plans.map(p => {
                    if (p.id === editingPlanId) {
                        return response as unknown as Plan;
                    }
                    return p;
                });
                setPlans(updatedPlans);
                planForm.reset();
                setEditingPlanId(null);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to update plan. Please try again.");
            console.error(err);
        }
    };

    // Handle deleting a plan
    const handleDeletePlan = async (planId: string) => {
        try {
            if (!apiClient) return;

            await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).delete();

            // Remove the plan from the local state
            setPlans(plans.filter(p => p.id !== planId));
            // Invalidate the providers query to refresh the data
            queryClient.invalidateQueries({queryKey: ['providers']});
        } catch (err) {
            setError("Failed to delete plan. Please try again.");
            console.error(err);
        }
    };

    // Handle adding a new price to a plan
    const handleAddPrice = async (data: PriceFormValues) => {
        try {
            if (!apiClient || !showAddPriceForm) return;

            // Validate dates
            const startDate = new Date(data.startDate);
            if (isNaN(startDate.getTime())) {
                setError("Invalid start date. Please enter a valid date.");
                return;
            }

            let endDate = undefined;
            if (data.endDate) {
                const parsedEndDate = new Date(data.endDate);
                if (isNaN(parsedEndDate.getTime())) {
                    setError("Invalid end date. Please enter a valid date.");
                    return;
                }
                endDate = parsedEndDate;
            }

            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(showAddPriceForm).prices.post({
                amount: data.amount,
                currency: data.currency,
                startDate: startDate,
                endDate: endDate
            });

            if (response) {
                // Add the new price to the plan in the local state
                const updatedPlans = plans.map(p => {
                    if (p.id === showAddPriceForm) {
                        // Create a new plan object with the updated prices array
                        const updatedPrices = [...p.prices, response as unknown as Price];
                        // Create a deep copy of the plan to maintain all its properties
                        const updatedPlan = new Plan(
                            p.id,
                            p.name,
                            p.description,
                            updatedPrices,
                            p.createdAt,
                            p.updatedAt,
                            p.etag
                        );
                        return updatedPlan;
                    }
                    return p;
                });
                setPlans(updatedPlans);
                priceForm.reset();
                setShowAddPriceForm(null);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to add price. Please try again.");
            console.error(err);
        }
    };

    // Handle updating a price
    const handleUpdatePrice = async (data: PriceFormValues) => {
        try {
            if (!apiClient || !editingPriceInfo) return;

            const {planId, priceId} = editingPriceInfo;
            const plan = plans.find(p => p.id === planId);
            if (!plan) return;

            const priceToUpdate = plan.prices.find(p => p.id === priceId);
            if (!priceToUpdate) return;

            // Validate dates
            const startDate = new Date(data.startDate);
            if (isNaN(startDate.getTime())) {
                setError("Invalid start date. Please enter a valid date.");
                return;
            }

            let endDate = undefined;
            if (data.endDate) {
                const parsedEndDate = new Date(data.endDate);
                if (isNaN(parsedEndDate.getTime())) {
                    setError("Invalid end date. Please enter a valid date.");
                    return;
                }
                endDate = parsedEndDate;
            }

            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).prices.byPriceId(priceId).put({
                amount: data.amount,
                currency: data.currency,
                startDate: startDate,
                endDate: endDate
                // etag is not needed in the UpdatePriceModel
            });

            if (response) {
                // Update the price in the local state
                const updatedPlans = plans.map(p => {
                    if (p.id === planId) {
                        // Update the specific price in the prices array
                        const updatedPrices = p.prices.map(price => {
                            if (price.id === priceId) {
                                return response as unknown as Price;
                            }
                            return price;
                        });
                        // Create a deep copy of the plan to maintain all its properties
                        const updatedPlan = new Plan(
                            p.id,
                            p.name,
                            p.description,
                            updatedPrices,
                            p.createdAt,
                            p.updatedAt,
                            p.etag
                        );
                        return updatedPlan;
                    }
                    return p;
                });
                setPlans(updatedPlans);
                priceForm.reset();
                setEditingPriceInfo(null);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({queryKey: ['providers']});
            }
        } catch (err) {
            setError("Failed to update price. Please try again.");
            console.error(err);
        }
    };

    // Handle deleting a price
    const handleDeletePrice = async (planId: string, priceId: string) => {
        try {
            if (!apiClient) return;

            await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).prices.byPriceId(priceId).delete();

            // Remove the price from the plan in the local state
            const updatedPlans = plans.map(p => {
                if (p.id === planId) {
                    // Filter out the deleted price
                    const updatedPrices = p.prices.filter(price => price.id !== priceId);
                    // Create a deep copy of the plan to maintain all its properties
                    const updatedPlan = new Plan(
                        p.id,
                        p.name,
                        p.description,
                        updatedPrices,
                        p.createdAt,
                        p.updatedAt,
                        p.etag
                    );
                    return updatedPlan;
                }
                return p;
            });
            setPlans(updatedPlans);
            // Invalidate the providers query to refresh the data
            queryClient.invalidateQueries({queryKey: ['providers']});
        } catch (err) {
            setError("Failed to delete price. Please try again.");
            console.error(err);
        }
    };

    // Start editing a plan
    const startEditingPlan = (planId: string) => {
        const plan = plans.find(p => p.id === planId);
        if (plan) {
            planForm.setValue("name", plan.name);
            planForm.setValue("description", plan.description || "");
            setEditingPlanId(planId);
        }
    };

    // Start editing a price
    const startEditingPrice = (planId: string, priceId: string) => {
        const plan = plans.find(p => p.id === planId);
        if (plan) {
            const price = plan.prices.find(p => p.id === priceId);
            if (price) {
                priceForm.setValue("amount", price.amount);
                priceForm.setValue("currency", price.currency);
                priceForm.setValue("startDate", price.startDate.toISOString().split('T')[0]);
                priceForm.setValue("endDate", price.endDate ? price.endDate.toISOString().split('T')[0] : "");
                setEditingPriceInfo({planId, priceId});
            }
        }
    };

    const handleClose = () => {
        reset();
        planForm.reset();
        priceForm.reset();
        setError(null);
        setEditingPlanId(null);
        setEditingPriceInfo(null);
        setShowAddPlanForm(false);
        setShowAddPriceForm(null);
        onClose();
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleClose}>
            <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Edit Provider</DialogTitle>
                    <DialogDescription>
                        Update provider details, plans, and prices.
                    </DialogDescription>
                </DialogHeader>

                {error && (
                    <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4"/>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}

                <Accordion type="single" collapsible defaultValue="provider-details">
                    <AccordionItem value="provider-details">
                        <AccordionTrigger>Provider Details</AccordionTrigger>
                        <AccordionContent>
                            <form id="provider-form" onSubmit={handleSubmit(onSubmit as any)} className="space-y-4">
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
                                    <p className="text-xs text-gray-500 italic">Owner type can only be set during
                                        creation</p>
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
                        </AccordionContent>
                    </AccordionItem>

                    <AccordionItem value="plans">
                        <AccordionTrigger>Plans and Prices</AccordionTrigger>
                        <AccordionContent>
                            <div className="space-y-4">
                                <div className="flex justify-between items-center">
                                    <h3 className="text-lg font-medium">Plans</h3>
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        onClick={() => {
                                            planForm.reset();
                                            setEditingPlanId(null);
                                            setShowAddPlanForm(true);
                                        }}
                                    >
                                        <Plus className="h-4 w-4 mr-2"/>
                                        Add Plan
                                    </Button>
                                </div>

                                {/* Add/Edit Plan Form */}
                                {showAddPlanForm && (
                                    <Card className="mt-4 border-dashed border-2">
                                        <CardHeader className="pb-2">
                                            <CardTitle
                                                className="text-base">{editingPlanId ? "Edit Plan" : "Add New Plan"}</CardTitle>
                                        </CardHeader>
                                        <CardContent>
                                            <form
                                                onSubmit={planForm.handleSubmit(editingPlanId ? handleUpdatePlan : handleAddPlan)}
                                                className="space-y-4">
                                                <div className="space-y-2">
                                                    <Label htmlFor="plan-name">Name *</Label>
                                                    <Input id="plan-name" {...planForm.register("name")} />
                                                    {planForm.formState.errors.name && (
                                                        <p className="text-sm text-red-500">{planForm.formState.errors.name.message}</p>
                                                    )}
                                                </div>

                                                <div className="space-y-2">
                                                    <Label htmlFor="plan-description">Description</Label>
                                                    <Textarea
                                                        id="plan-description" {...planForm.register("description")} />
                                                </div>

                                                <div className="flex justify-end space-x-2">
                                                    <Button type="button" variant="outline" onClick={() => {
                                                        planForm.reset();
                                                        setEditingPlanId(null);
                                                        setShowAddPlanForm(false);
                                                    }}>
                                                        Cancel
                                                    </Button>
                                                    <Button type="submit" disabled={planForm.formState.isSubmitting}>
                                                        {planForm.formState.isSubmitting ? (
                                                            <>
                                                                <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                                                Saving...
                                                            </>
                                                        ) : (
                                                            editingPlanId ? "Update Plan" : "Add Plan"
                                                        )}
                                                    </Button>
                                                </div>
                                            </form>
                                        </CardContent>
                                    </Card>
                                )}

                                {/* Plans list */}
                                {plans.length === 0 && !showAddPlanForm ? (
                                    <div className="text-center py-4 text-gray-500">
                                        No plans added yet. Click "Add Plan" to create one.
                                    </div>
                                ) : (
                                    <div className="space-y-4 mt-4">
                                        {plans.map((plan) => (
                                            <Card key={plan.id}>
                                                <CardHeader className="pb-2">
                                                    <div className="flex justify-between items-start">
                                                        <CardTitle>{plan.name}</CardTitle>
                                                        <div className="flex space-x-2">
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                onClick={() => {
                                                                    startEditingPlan(plan.id);
                                                                    setShowAddPlanForm(true);
                                                                }}
                                                            >
                                                                <Edit2 className="h-4 w-4"/>
                                                            </Button>
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                onClick={() => handleDeletePlan(plan.id)}
                                                            >
                                                                <Trash2 className="h-4 w-4 text-red-500"/>
                                                            </Button>
                                                        </div>
                                                    </div>
                                                    {plan.description && (
                                                        <p className="text-sm text-gray-500">{plan.description}</p>
                                                    )}
                                                </CardHeader>
                                                <CardContent>
                                                    <div className="space-y-2">
                                                        <div className="flex justify-between items-center">
                                                            <h4 className="text-sm font-medium">Prices</h4>
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                onClick={() => {
                                                                    priceForm.reset();
                                                                    setEditingPriceInfo(null);
                                                                    setShowAddPriceForm(plan.id);
                                                                }}
                                                            >
                                                                <Plus className="h-4 w-4 mr-1"/>
                                                                Add Price
                                                            </Button>
                                                        </div>

                                                        {/* Add/Edit Price Form */}
                                                        {showAddPriceForm === plan.id && (
                                                            <Card className="mt-2 border-dashed border-2">
                                                                <CardHeader className="py-2">
                                                                    <CardTitle
                                                                        className="text-sm">{editingPriceInfo ? "Edit Price" : "Add New Price"}</CardTitle>
                                                                </CardHeader>
                                                                <CardContent className="py-2">
                                                                    <form
                                                                        onSubmit={priceForm.handleSubmit(editingPriceInfo ? handleUpdatePrice as any : handleAddPrice as any)}
                                                                        className="space-y-3">
                                                                        <div className="space-y-2">
                                                                            <Label htmlFor="price-amount"
                                                                                   className="text-xs">Amount *</Label>
                                                                            <Input
                                                                                id="price-amount"
                                                                                type="number"
                                                                                step="0.01"
                                                                                className="h-8"
                                                                                {...priceForm.register("amount", {valueAsNumber: true})}
                                                                            />
                                                                            {priceForm.formState.errors.amount && (
                                                                                <p className="text-xs text-red-500">{priceForm.formState.errors.amount.message}</p>
                                                                            )}
                                                                        </div>

                                                                        <div className="space-y-2">
                                                                            <Label htmlFor="price-currency"
                                                                                   className="text-xs">Currency
                                                                                *</Label>
                                                                            <Input id="price-currency"
                                                                                   className="h-8" {...priceForm.register("currency")} />
                                                                            {priceForm.formState.errors.currency && (
                                                                                <p className="text-xs text-red-500">{priceForm.formState.errors.currency.message}</p>
                                                                            )}
                                                                        </div>

                                                                        <div className="space-y-2">
                                                                            <Label htmlFor="price-start-date"
                                                                                   className="text-xs">Start Date
                                                                                *</Label>
                                                                            <Input id="price-start-date" type="date"
                                                                                   className="h-8" {...priceForm.register("startDate")} />
                                                                            {priceForm.formState.errors.startDate && (
                                                                                <p className="text-xs text-red-500">{priceForm.formState.errors.startDate.message}</p>
                                                                            )}
                                                                        </div>

                                                                        <div className="space-y-2">
                                                                            <Label htmlFor="price-end-date"
                                                                                   className="text-xs">End Date
                                                                                (Optional)</Label>
                                                                            <Input id="price-end-date" type="date"
                                                                                   className="h-8" {...priceForm.register("endDate")} />
                                                                        </div>

                                                                        <div className="flex justify-end space-x-2">
                                                                            <Button type="button" variant="outline"
                                                                                    size="sm" onClick={() => {
                                                                                priceForm.reset();
                                                                                setEditingPriceInfo(null);
                                                                                setShowAddPriceForm(null);
                                                                            }}>
                                                                                Cancel
                                                                            </Button>
                                                                            <Button type="submit" size="sm"
                                                                                    disabled={priceForm.formState.isSubmitting}>
                                                                                {priceForm.formState.isSubmitting ? (
                                                                                    <>
                                                                                        <Loader2
                                                                                            className="mr-1 h-3 w-3 animate-spin"/>
                                                                                        Saving...
                                                                                    </>
                                                                                ) : (
                                                                                    editingPriceInfo ? "Update" : "Add"
                                                                                )}
                                                                            </Button>
                                                                        </div>
                                                                    </form>
                                                                </CardContent>
                                                            </Card>
                                                        )}

                                                        {plan.prices.length === 0 && showAddPriceForm !== plan.id ? (
                                                            <div className="text-center py-2 text-gray-500 text-sm">
                                                                No prices added yet.
                                                            </div>
                                                        ) : (
                                                            <div className="space-y-2">
                                                                {plan.prices.map((price) => (
                                                                    <div key={price.id}
                                                                         className="flex justify-between items-center p-2 border rounded-md">
                                                                        <div>
                                                                            <div className="font-medium">
                                                                                {price.amount} {price.currency}
                                                                            </div>
                                                                            <div className="text-xs text-gray-500">
                                                                                From: {price.startDate.toLocaleDateString()}
                                                                                {price.endDate && ` To: ${price.endDate.toLocaleDateString()}`}
                                                                            </div>
                                                                        </div>
                                                                        <div className="flex space-x-1">
                                                                            <Button
                                                                                variant="ghost"
                                                                                size="sm"
                                                                                onClick={() => {
                                                                                    startEditingPrice(plan.id, price.id);
                                                                                    setShowAddPriceForm(plan.id);
                                                                                }}
                                                                            >
                                                                                <Edit2 className="h-3 w-3"/>
                                                                            </Button>
                                                                            <Button
                                                                                variant="ghost"
                                                                                size="sm"
                                                                                onClick={() => handleDeletePrice(plan.id, price.id)}
                                                                            >
                                                                                <Trash2
                                                                                    className="h-3 w-3 text-red-500"/>
                                                                            </Button>
                                                                        </div>
                                                                    </div>
                                                                ))}
                                                            </div>
                                                        )}
                                                    </div>
                                                </CardContent>
                                            </Card>
                                        ))}
                                    </div>
                                )}
                            </div>
                        </AccordionContent>
                    </AccordionItem>
                </Accordion>

                {/* No separate dialogs needed anymore as we're using inline forms */}

                <DialogFooter>
                    <Button type="button" onClick={handleClose}>
                        Close
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}