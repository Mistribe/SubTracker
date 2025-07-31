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
import {AlertCircle, Loader2, Plus, Trash2, Edit2} from "lucide-react";
import Provider from "@/models/provider";
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
});

// Define the plan form schema
const planFormSchema = z.object({
    name: z.string().min(1, "Plan name is required"),
    description: z.string().optional(),
});

// Define the price form schema
const priceFormSchema = z.object({
    amount: z.number().min(0, "Amount must be a positive number"),
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
    const { apiClient } = useApiClient();
    const queryClient = useQueryClient();

    // State for plans and prices management
    const [plans, setPlans] = useState(provider.plans);
    const [isAddingPlan, setIsAddingPlan] = useState(false);
    const [editingPlanId, setEditingPlanId] = useState<string | null>(null);
    const [isAddingPrice, setIsAddingPrice] = useState(false);
    const [editingPriceInfo, setEditingPriceInfo] = useState<{planId: string, priceId: string} | null>(null);
    const [currentPlanId, setCurrentPlanId] = useState<string | null>(null);

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
            familyId: provider.owner.type === OwnerType.Family ? provider.owner.familyId : undefined,
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
        resolver: zodResolver(priceFormSchema),
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
                ownerType: data.ownerType,
                familyId: data.ownerType === OwnerType.Family ? data.familyId : undefined,
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
                setPlans([...plans, response]);
                planForm.reset();
                setIsAddingPlan(false);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({ queryKey: ['providers'] });
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
                description: data.description || undefined,
                etag: planToUpdate.etag
            });
            
            if (response) {
                // Update the plan in the local state
                setPlans(plans.map(p => p.id === editingPlanId ? response : p));
                planForm.reset();
                setEditingPlanId(null);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({ queryKey: ['providers'] });
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
            queryClient.invalidateQueries({ queryKey: ['providers'] });
        } catch (err) {
            setError("Failed to delete plan. Please try again.");
            console.error(err);
        }
    };

    // Handle adding a new price to a plan
    const handleAddPrice = async (data: PriceFormValues) => {
        try {
            if (!apiClient || !currentPlanId) return;
            
            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(currentPlanId).prices.post({
                amount: data.amount,
                currency: data.currency,
                startDate: new Date(data.startDate),
                endDate: data.endDate ? new Date(data.endDate) : undefined
            });
            
            if (response) {
                // Add the new price to the plan in the local state
                setPlans(plans.map(p => {
                    if (p.id === currentPlanId) {
                        return {
                            ...p,
                            prices: [...p.prices, response]
                        };
                    }
                    return p;
                }));
                priceForm.reset();
                setIsAddingPrice(false);
                setCurrentPlanId(null);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({ queryKey: ['providers'] });
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
            
            const { planId, priceId } = editingPriceInfo;
            const plan = plans.find(p => p.id === planId);
            if (!plan) return;
            
            const priceToUpdate = plan.prices.find(p => p.id === priceId);
            if (!priceToUpdate) return;
            
            const response = await apiClient.providers.byProviderId(provider.id).plans.byPlanId(planId).prices.byPriceId(priceId).put({
                amount: data.amount,
                currency: data.currency,
                startDate: new Date(data.startDate),
                endDate: data.endDate ? new Date(data.endDate) : undefined,
                etag: priceToUpdate.etag
            });
            
            if (response) {
                // Update the price in the local state
                setPlans(plans.map(p => {
                    if (p.id === planId) {
                        return {
                            ...p,
                            prices: p.prices.map(price => price.id === priceId ? response : price)
                        };
                    }
                    return p;
                }));
                priceForm.reset();
                setEditingPriceInfo(null);
                // Invalidate the providers query to refresh the data
                queryClient.invalidateQueries({ queryKey: ['providers'] });
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
            setPlans(plans.map(p => {
                if (p.id === planId) {
                    return {
                        ...p,
                        prices: p.prices.filter(price => price.id !== priceId)
                    };
                }
                return p;
            }));
            // Invalidate the providers query to refresh the data
            queryClient.invalidateQueries({ queryKey: ['providers'] });
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
            setIsAddingPlan(true);
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
                setEditingPriceInfo({ planId, priceId });
                setIsAddingPrice(true);
                setCurrentPlanId(planId);
            }
        }
    };

    const handleClose = () => {
        reset();
        planForm.reset();
        priceForm.reset();
        setError(null);
        setIsAddingPlan(false);
        setEditingPlanId(null);
        setIsAddingPrice(false);
        setEditingPriceInfo(null);
        setCurrentPlanId(null);
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
                                    <RadioGroup
                                        defaultValue={provider.owner.type}
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
                                                defaultValue={provider.owner.familyId}
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
                                            setIsAddingPlan(true);
                                        }}
                                    >
                                        <Plus className="h-4 w-4 mr-2" />
                                        Add Plan
                                    </Button>
                                </div>

                                {/* Plans list */}
                                {plans.length === 0 ? (
                                    <div className="text-center py-4 text-gray-500">
                                        No plans added yet. Click "Add Plan" to create one.
                                    </div>
                                ) : (
                                    <div className="space-y-4">
                                        {plans.map((plan) => (
                                            <Card key={plan.id}>
                                                <CardHeader className="pb-2">
                                                    <div className="flex justify-between items-start">
                                                        <CardTitle>{plan.name}</CardTitle>
                                                        <div className="flex space-x-2">
                                                            <Button 
                                                                variant="ghost" 
                                                                size="sm"
                                                                onClick={() => startEditingPlan(plan.id)}
                                                            >
                                                                <Edit2 className="h-4 w-4" />
                                                            </Button>
                                                            <Button 
                                                                variant="ghost" 
                                                                size="sm"
                                                                onClick={() => handleDeletePlan(plan.id)}
                                                            >
                                                                <Trash2 className="h-4 w-4 text-red-500" />
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
                                                                    setIsAddingPrice(true);
                                                                    setCurrentPlanId(plan.id);
                                                                }}
                                                            >
                                                                <Plus className="h-4 w-4 mr-1" />
                                                                Add Price
                                                            </Button>
                                                        </div>
                                                        {plan.prices.length === 0 ? (
                                                            <div className="text-center py-2 text-gray-500 text-sm">
                                                                No prices added yet.
                                                            </div>
                                                        ) : (
                                                            <div className="space-y-2">
                                                                {plan.prices.map((price) => (
                                                                    <div key={price.id} className="flex justify-between items-center p-2 border rounded-md">
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
                                                                                key={`edit-${price.id}`}
                                                                                variant="ghost" 
                                                                                size="sm"
                                                                                onClick={() => startEditingPrice(plan.id, price.id)}
                                                                            >
                                                                                <Edit2 className="h-3 w-3" />
                                                                            </Button>
                                                                            <Button 
                                                                                key={`delete-${price.id}`}
                                                                                variant="ghost" 
                                                                                size="sm"
                                                                                onClick={() => handleDeletePrice(plan.id, price.id)}
                                                                            >
                                                                                <Trash2 className="h-3 w-3 text-red-500" />
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

                {/* Plan form dialog */}
                {isAddingPlan && (
                    <Dialog open={isAddingPlan} onOpenChange={(open) => {
                        if (!open) {
                            planForm.reset();
                            setEditingPlanId(null);
                            setIsAddingPlan(false);
                        }
                    }}>
                        <DialogContent className="sm:max-w-[500px]">
                            <DialogHeader>
                                <DialogTitle>{editingPlanId ? "Edit Plan" : "Add New Plan"}</DialogTitle>
                            </DialogHeader>
                            <form onSubmit={planForm.handleSubmit(editingPlanId ? handleUpdatePlan : handleAddPlan)} className="space-y-4">
                                <div className="space-y-2">
                                    <Label htmlFor="plan-name">Name *</Label>
                                    <Input id="plan-name" {...planForm.register("name")} />
                                    {planForm.formState.errors.name && (
                                        <p className="text-sm text-red-500">{planForm.formState.errors.name.message}</p>
                                    )}
                                </div>

                                <div className="space-y-2">
                                    <Label htmlFor="plan-description">Description</Label>
                                    <Textarea id="plan-description" {...planForm.register("description")} />
                                </div>

                                <DialogFooter>
                                    <Button type="button" variant="outline" onClick={() => {
                                        planForm.reset();
                                        setEditingPlanId(null);
                                        setIsAddingPlan(false);
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
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>
                )}

                {/* Price form dialog */}
                {isAddingPrice && (
                    <Dialog open={isAddingPrice} onOpenChange={(open) => {
                        if (!open) {
                            priceForm.reset();
                            setEditingPriceInfo(null);
                            setIsAddingPrice(false);
                            setCurrentPlanId(null);
                        }
                    }}>
                        <DialogContent className="sm:max-w-[500px]">
                            <DialogHeader>
                                <DialogTitle>{editingPriceInfo ? "Edit Price" : "Add New Price"}</DialogTitle>
                            </DialogHeader>
                            <form onSubmit={priceForm.handleSubmit(editingPriceInfo ? handleUpdatePrice : handleAddPrice)} className="space-y-4">
                                <div className="space-y-2">
                                    <Label htmlFor="price-amount">Amount *</Label>
                                    <Input 
                                        id="price-amount" 
                                        type="number" 
                                        step="0.01"
                                        {...priceForm.register("amount", { valueAsNumber: true })} 
                                    />
                                    {priceForm.formState.errors.amount && (
                                        <p className="text-sm text-red-500">{priceForm.formState.errors.amount.message}</p>
                                    )}
                                </div>

                                <div className="space-y-2">
                                    <Label htmlFor="price-currency">Currency *</Label>
                                    <Input id="price-currency" {...priceForm.register("currency")} />
                                    {priceForm.formState.errors.currency && (
                                        <p className="text-sm text-red-500">{priceForm.formState.errors.currency.message}</p>
                                    )}
                                </div>

                                <div className="space-y-2">
                                    <Label htmlFor="price-start-date">Start Date *</Label>
                                    <Input id="price-start-date" type="date" {...priceForm.register("startDate")} />
                                    {priceForm.formState.errors.startDate && (
                                        <p className="text-sm text-red-500">{priceForm.formState.errors.startDate.message}</p>
                                    )}
                                </div>

                                <div className="space-y-2">
                                    <Label htmlFor="price-end-date">End Date (Optional)</Label>
                                    <Input id="price-end-date" type="date" {...priceForm.register("endDate")} />
                                </div>

                                <DialogFooter>
                                    <Button type="button" variant="outline" onClick={() => {
                                        priceForm.reset();
                                        setEditingPriceInfo(null);
                                        setIsAddingPrice(false);
                                        setCurrentPlanId(null);
                                    }}>
                                        Cancel
                                    </Button>
                                    <Button type="submit" disabled={priceForm.formState.isSubmitting}>
                                        {priceForm.formState.isSubmitting ? (
                                            <>
                                                <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                                Saving...
                                            </>
                                        ) : (
                                            editingPriceInfo ? "Update Price" : "Add Price"
                                        )}
                                    </Button>
                                </DialogFooter>
                            </form>
                        </DialogContent>
                    </Dialog>
                )}

                <DialogFooter>
                    <Button type="button" variant="outline" onClick={handleClose}>
                        Cancel
                    </Button>
                    <Button form="provider-form" type="submit" disabled={isSubmitting}>
                        {isSubmitting ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                Updating...
                            </>
                        ) : (
                            "Update Provider"
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}