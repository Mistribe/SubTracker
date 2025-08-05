import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useNavigate } from "react-router-dom";
import { format } from "date-fns";
import { useSubscriptionsMutations } from "@/hooks/subscriptions/useSubscriptionsMutations";
import { useAllProvidersQuery } from "@/hooks/providers/useAllProvidersQuery";
import { useFamiliesQuery } from "@/hooks/families/useFamiliesQuery";
import { PageHeader } from "@/components/ui/page-header";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import { OwnerType } from "@/models/ownerType";
import { SubscriptionRecurrency } from "@/models/subscriptionRecurrency";
import { ProviderCombobox } from "@/components/providers/ProviderCombobox";

// Define the form schema with Zod
const formSchema = z.object({
    friendlyName: z.string().optional(),
    providerId: z.string().min(1, "Provider is required"),
    planId: z.string().min(1, "Plan is required"),
    priceId: z.string().min(1, "Price is required"),
    recurrency: z.nativeEnum(SubscriptionRecurrency),
    customRecurrencyValue: z.number().positive("Value must be positive").optional(),
    customRecurrencyUnit: z.enum(["days", "months", "years"]).optional(),
    startDate: z.date(),
    endDate: z.date().optional(),
    ownerType: z.nativeEnum(OwnerType),
    familyId: z.string().optional(),
    serviceUsers: z.array(z.string()).optional(),
    customPrice: z.object({
        amount: z.number().positive("Amount must be positive"),
        currency: z.string().min(1, "Currency is required")
    }).optional(),
    hasFreeTrialPeriod: z.boolean().default(false),
    freeTrialStartDate: z.date().optional(),
    freeTrialEndDate: z.date().optional(),
});

type FormValues = z.infer<typeof formSchema>;

const CreateSubscriptionPage = () => {
    const navigate = useNavigate();
    const [error, setError] = useState<string | null>(null);
    const { createSubscriptionMutation } = useSubscriptionsMutations();
    const { data: providersData } = useAllProvidersQuery();
    const { data: familiesData } = useFamiliesQuery({ limit: 100 });

    const providers = providersData?.pages.flatMap(page => page.providers) || [];
    const families = familiesData?.families || [];

    const {
        register,
        handleSubmit,
        watch,
        setValue,
        formState: { errors, isSubmitting },
    } = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            recurrency: SubscriptionRecurrency.Monthly,
            startDate: new Date(),
            ownerType: OwnerType.Personal,
            hasFreeTrialPeriod: false,
            serviceUsers: [],
            customRecurrencyValue: 1,
            customRecurrencyUnit: "days",
        },
    });

    const selectedOwnerType = watch("ownerType");
    const selectedRecurrency = watch("recurrency");
    const hasFreeTrialPeriod = watch("hasFreeTrialPeriod");
    const hasCustomPrice = watch("customPrice") !== undefined;

    // Function to convert custom recurrency to days
    const convertTodays = (value: number, unit: string): number => {
        switch (unit) {
            case "days":
                return value;
            case "months":
                return value * 30; // Approximate days in a month
            case "years":
                return value * 365; // Approximate days in a year
            default:
                return value;
        }
    };

    const onSubmit = async (data: FormValues) => {
        try {
            setError(null);
            
            // Calculate custom recurrency in days if needed
            let customRecurrencyInDays: number | undefined = undefined;
            if (data.recurrency === "custom" && data.customRecurrencyValue && data.customRecurrencyUnit) {
                customRecurrencyInDays = convertTodays(data.customRecurrencyValue, data.customRecurrencyUnit);
            }
            
            await createSubscriptionMutation.mutateAsync({
                friendlyName: data.friendlyName,
                providerId: data.providerId,
                planId: data.planId,
                priceId: data.priceId,
                recurrency: data.recurrency,
                customRecurrency: customRecurrencyInDays,
                startDate: data.startDate,
                endDate: data.endDate,
                ownerType: data.ownerType,
                familyId: data.ownerType === OwnerType.Family ? data.familyId : undefined,
                serviceUsers: data.serviceUsers,
                customPrice: data.customPrice,
                freeTrial: data.hasFreeTrialPeriod && data.freeTrialStartDate && data.freeTrialEndDate
                    ? {
                        startDate: data.freeTrialStartDate,
                        endDate: data.freeTrialEndDate,
                    }
                    : undefined,
            });
            navigate("/subscriptions");
        } catch (err) {
            setError("Failed to create subscription. Please try again.");
            console.error(err);
        }
    };

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Add New Subscription"
                description="Create a new subscription by filling out the form below."
                actionButton={
                    <Button onClick={() => navigate("/subscriptions")}>
                        Back to Subscriptions
                    </Button>
                }
            />

            <div className="max-w-3xl mx-auto mt-8">
                {error && (
                    <Alert variant="destructive" className="mb-6">
                        <AlertCircle className="h-4 w-4" />
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                    {/* Basic Information */}
                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold">Basic Information</h2>

                        <div className="grid grid-cols-1 gap-4">
                            <div>
                                <Label htmlFor="friendlyName">Friendly Name (Optional)</Label>
                                <Input
                                    id="friendlyName"
                                    {...register("friendlyName")}
                                    placeholder="e.g., Netflix Family Plan"
                                />
                                {errors.friendlyName && (
                                    <p className="text-sm text-red-500 mt-1">{errors.friendlyName.message}</p>
                                )}
                            </div>

                            <div>
                                <Label htmlFor="providerId">Provider</Label>
                                <ProviderCombobox
                                    providers={providers}
                                    value={watch("providerId") || ""}
                                    onChange={(value) => setValue("providerId", value)}
                                    placeholder="Select a provider"
                                    emptyMessage="No provider found. Try a different search."
                                />
                                {errors.providerId && (
                                    <p className="text-sm text-red-500 mt-1">{errors.providerId.message}</p>
                                )}
                            </div>

                            <div>
                                <Label htmlFor="planId">Plan ID</Label>
                                <Input
                                    id="planId"
                                    {...register("planId")}
                                    placeholder="Enter plan ID"
                                />
                                {errors.planId && (
                                    <p className="text-sm text-red-500 mt-1">{errors.planId.message}</p>
                                )}
                            </div>

                            <div>
                                <Label htmlFor="priceId">Price ID</Label>
                                <Input
                                    id="priceId"
                                    {...register("priceId")}
                                    placeholder="Enter price ID"
                                />
                                {errors.priceId && (
                                    <p className="text-sm text-red-500 mt-1">{errors.priceId.message}</p>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Recurrency */}
                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold">Recurrency</h2>

                        <div className="flex flex-col">
                            <Label htmlFor="recurrency">Recurrency Type</Label>
                            <div className="flex gap-4 items-start">
                                <div className={selectedRecurrency === "custom" ? "w-1/2" : "w-full"}>
                                    <Select
                                        onValueChange={(value) => setValue("recurrency", value as SubscriptionRecurrency)}
                                        defaultValue={SubscriptionRecurrency.Monthly}
                                    >
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select recurrency" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value={SubscriptionRecurrency.Monthly}>Monthly</SelectItem>
                                            <SelectItem value={SubscriptionRecurrency.Quarterly}>Quarterly</SelectItem>
                                            <SelectItem value={SubscriptionRecurrency.HalfYearly}>Half Yearly</SelectItem>
                                            <SelectItem value={SubscriptionRecurrency.Yearly}>Yearly</SelectItem>
                                            <SelectItem value={SubscriptionRecurrency.OneTime}>One Time</SelectItem>
                                            <SelectItem value="custom">Custom</SelectItem>
                                        </SelectContent>
                                    </Select>
                                </div>
                                
                                {selectedRecurrency === "custom" && (
                                    <div className="w-1/2 flex flex-col gap-1">
                                        <div className="flex gap-2">
                                            <div className="w-1/2">
                                                <Label htmlFor="customRecurrencyValue" className="sr-only">Value</Label>
                                                <Input
                                                    id="customRecurrencyValue"
                                                    type="number"
                                                    min="1"
                                                    {...register("customRecurrencyValue", { valueAsNumber: true })}
                                                    placeholder="Value"
                                                />
                                                {errors.customRecurrencyValue && (
                                                    <p className="text-sm text-red-500 mt-1">{errors.customRecurrencyValue.message}</p>
                                                )}
                                            </div>
                                            <div className="w-1/2">
                                                <Label htmlFor="customRecurrencyUnit" className="sr-only">Unit</Label>
                                                <Select
                                                    onValueChange={(value) => setValue("customRecurrencyUnit", value)}
                                                    defaultValue="days"
                                                >
                                                    <SelectTrigger id="customRecurrencyUnit">
                                                        <SelectValue placeholder="Unit" />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        <SelectItem value="days">Days</SelectItem>
                                                        <SelectItem value="months">Months</SelectItem>
                                                        <SelectItem value="years">Years</SelectItem>
                                                    </SelectContent>
                                                </Select>
                                                {errors.customRecurrencyUnit && (
                                                    <p className="text-sm text-red-500 mt-1">{errors.customRecurrencyUnit.message}</p>
                                                )}
                                            </div>
                                        </div>
                                        <p className="text-xs text-muted-foreground">
                                            Specify how often the subscription recurs. This will be converted to days (30 days per month, 365 days per year).
                                        </p>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Dates */}
                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold">Dates</h2>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <Label htmlFor="startDate">Start Date</Label>
                                <Input
                                    id="startDate"
                                    type="date"
                                    {...register("startDate", {
                                        valueAsDate: true,
                                    })}
                                    defaultValue={format(new Date(), "yyyy-MM-dd")}
                                />
                                {errors.startDate && (
                                    <p className="text-sm text-red-500 mt-1">{errors.startDate.message}</p>
                                )}
                            </div>

                            <div>
                                <Label htmlFor="endDate">End Date (Optional)</Label>
                                <Input
                                    id="endDate"
                                    type="date"
                                    {...register("endDate", {
                                        valueAsDate: true,
                                    })}
                                />
                                {errors.endDate && (
                                    <p className="text-sm text-red-500 mt-1">{errors.endDate.message}</p>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Ownership */}
                    <div className="space-y-4">
                        <h2 className="text-xl font-semibold">Ownership</h2>

                        <div>
                            <Label>Owner Type</Label>
                            <RadioGroup
                                defaultValue={OwnerType.Personal}
                                onValueChange={(value) => setValue("ownerType", value as OwnerType)}
                                className="grid grid-cols-2 gap-4 mt-2"
                            >
                                <div className="flex items-center space-x-2">
                                    <RadioGroupItem value={OwnerType.Personal} id="personal" />
                                    <Label htmlFor="personal">Personal</Label>
                                </div>
                                <div className="flex items-center space-x-2">
                                    <RadioGroupItem value={OwnerType.Family} id="family" />
                                    <Label htmlFor="family">Family</Label>
                                </div>
                            </RadioGroup>
                        </div>

                        {selectedOwnerType === OwnerType.Family && (
                            <div>
                                <Label htmlFor="familyId">Family</Label>
                                <Select
                                    onValueChange={(value) => setValue("familyId", value)}
                                    defaultValue=""
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select a family" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {families.map((family) => (
                                            <SelectItem key={family.id} value={family.id}>
                                                {family.name}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                                {errors.familyId && (
                                    <p className="text-sm text-red-500 mt-1">{errors.familyId.message}</p>
                                )}
                            </div>
                        )}
                    </div>

                    {/* Custom Price */}
                    <div className="space-y-4">
                        <div className="flex items-center space-x-2">
                            <input
                                type="checkbox"
                                id="hasCustomPrice"
                                checked={hasCustomPrice}
                                onChange={(e) => {
                                    if (e.target.checked) {
                                        setValue("customPrice", { amount: 0, currency: "USD" });
                                    } else {
                                        setValue("customPrice", undefined);
                                    }
                                }}
                                className="h-4 w-4 rounded border-gray-300"
                            />
                            <Label htmlFor="hasCustomPrice">Add Custom Price</Label>
                        </div>

                        {hasCustomPrice && (
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pl-6">
                                <div>
                                    <Label htmlFor="customPriceAmount">Amount</Label>
                                    <Input
                                        id="customPriceAmount"
                                        type="number"
                                        step="0.01"
                                        {...register("customPrice.amount", { valueAsNumber: true })}
                                        placeholder="Enter amount"
                                    />
                                    {errors.customPrice?.amount && (
                                        <p className="text-sm text-red-500 mt-1">{errors.customPrice.amount.message}</p>
                                    )}
                                </div>

                                <div>
                                    <Label htmlFor="customPriceCurrency">Currency</Label>
                                    <Select
                                        onValueChange={(value) => setValue("customPrice.currency", value)}
                                        defaultValue="USD"
                                    >
                                        <SelectTrigger>
                                            <SelectValue placeholder="Select currency" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="USD">USD</SelectItem>
                                            <SelectItem value="EUR">EUR</SelectItem>
                                            <SelectItem value="GBP">GBP</SelectItem>
                                            <SelectItem value="CAD">CAD</SelectItem>
                                            <SelectItem value="AUD">AUD</SelectItem>
                                        </SelectContent>
                                    </Select>
                                    {errors.customPrice?.currency && (
                                        <p className="text-sm text-red-500 mt-1">{errors.customPrice.currency.message}</p>
                                    )}
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Free Trial */}
                    <div className="space-y-4">
                        <div className="flex items-center space-x-2">
                            <input
                                type="checkbox"
                                id="hasFreeTrialPeriod"
                                {...register("hasFreeTrialPeriod")}
                                className="h-4 w-4 rounded border-gray-300"
                            />
                            <Label htmlFor="hasFreeTrialPeriod">Has Free Trial Period</Label>
                        </div>

                        {hasFreeTrialPeriod && (
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pl-6">
                                <div>
                                    <Label htmlFor="freeTrialStartDate">Free Trial Start Date</Label>
                                    <Input
                                        id="freeTrialStartDate"
                                        type="date"
                                        {...register("freeTrialStartDate", {
                                            valueAsDate: true,
                                        })}
                                        defaultValue={format(new Date(), "yyyy-MM-dd")}
                                    />
                                    {errors.freeTrialStartDate && (
                                        <p className="text-sm text-red-500 mt-1">{errors.freeTrialStartDate.message}</p>
                                    )}
                                </div>

                                <div>
                                    <Label htmlFor="freeTrialEndDate">Free Trial End Date</Label>
                                    <Input
                                        id="freeTrialEndDate"
                                        type="date"
                                        {...register("freeTrialEndDate", {
                                            valueAsDate: true,
                                        })}
                                    />
                                    {errors.freeTrialEndDate && (
                                        <p className="text-sm text-red-500 mt-1">{errors.freeTrialEndDate.message}</p>
                                    )}
                                </div>
                            </div>
                        )}
                    </div>

                    {/* Submit Button */}
                    <div className="flex justify-end">
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Creating...
                                </>
                            ) : (
                                "Create Subscription"
                            )}
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default CreateSubscriptionPage;