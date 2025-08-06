import {useEffect, useState} from "react";
import {FormProvider, useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {useNavigate, useParams} from "react-router-dom";
import {useSubscriptionsMutations} from "@/hooks/subscriptions/useSubscriptionsMutations";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import {useFamiliesQuery} from "@/hooks/families/useFamiliesQuery";
import {useAllSubscriptionsQuery} from "@/hooks/subscriptions/useAllSubscriptionsQuery";
import {PageHeader} from "@/components/ui/page-header";
import {Button} from "@/components/ui/button";
import {Alert, AlertDescription} from "@/components/ui/alert";
import {AlertCircle, Loader2} from "lucide-react";
import {OwnerType} from "@/models/ownerType";
import {convertToDays, formSchema, type FormValues} from "@/components/subscriptions/form/SubscriptionFormSchema";
import {BasicInformationSection} from "@/components/subscriptions/form/BasicInformationSection";
import {RecurrencySection} from "@/components/subscriptions/form/RecurrencySection";
import {DatesSection} from "@/components/subscriptions/form/DatesSection";
import {OwnershipSection} from "@/components/subscriptions/form/OwnershipSection";
import {CustomPriceSection} from "@/components/subscriptions/form/CustomPriceSection";
import {FreeTrialSection} from "@/components/subscriptions/form/FreeTrialSection";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";

const SubscriptionFormPage = () => {
    const navigate = useNavigate();
    const {subscriptionId} = useParams<{ subscriptionId: string }>();
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const {createSubscriptionMutation, updateSubscriptionMutation} = useSubscriptionsMutations();
    const {data: providersData} = useAllProvidersQuery();
    const {data: familiesData} = useFamiliesQuery({limit: 100});
    const {data: subscriptionsData} = useAllSubscriptionsQuery();

    const isEditMode = !!subscriptionId;
    const providers = providersData?.pages.flatMap(page => page.providers) || [];
    const families = familiesData?.families || [];

    // Find the subscription to edit if in edit mode
    const subscriptionToEdit = isEditMode
        ? subscriptionsData?.pages.flatMap(page => page.subscriptions).find(sub => sub.id === subscriptionId)
        : undefined;

    const form = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            recurrency: SubscriptionRecurrency.Monthly,
            startDate: new Date(),
            ownerType: OwnerType.Personal,
            friendlyName: undefined,
            hasFreeTrialPeriod: false,
            serviceUsers: [],
            customRecurrencyValue: 1,
            customRecurrencyUnit: "days",
            customPrice: {
                amount: 10,
                currency: "USD"
            },
            endDate: undefined,
            familyId: undefined,
            freeTrialEndDate: new Date(),
            freeTrialStartDate: new Date(),
            planId: undefined,
            priceId: undefined,
            providerId: undefined
        },
    });

    // Populate form with subscription data when in edit mode
    useEffect(() => {
        if (isEditMode && subscriptionToEdit) {
            setIsLoading(true);

            // Set basic information
            form.setValue("friendlyName", subscriptionToEdit.friendlyName || undefined);
            form.setValue("providerId", subscriptionToEdit.providerId);
            form.setValue("planId", subscriptionToEdit.planId);
            form.setValue("priceId", subscriptionToEdit.priceId);

            // Set recurrency
            form.setValue("recurrency", subscriptionToEdit.recurrency);
            if (subscriptionToEdit.recurrency === "custom" && subscriptionToEdit.customRecurrency) {
                form.setValue("customRecurrencyValue", subscriptionToEdit.customRecurrency);
                form.setValue("customRecurrencyUnit", "days");
            }

            // Set dates
            form.setValue("startDate", subscriptionToEdit.startDate);
            if (subscriptionToEdit.endDate) {
                form.setValue("endDate", subscriptionToEdit.endDate);
            }

            // Set owner type and family ID if applicable
            form.setValue("ownerType", subscriptionToEdit.owner.type);
            if (subscriptionToEdit.owner.type === OwnerType.Family && subscriptionToEdit.owner.familyId) {
                form.setValue("familyId", subscriptionToEdit.owner.familyId);
            }

            // Set service users
            if (subscriptionToEdit.serviceUsers.length > 0) {
                form.setValue("serviceUsers", subscriptionToEdit.serviceUsers);
            }

            // Set custom price if available
            if (subscriptionToEdit.customPrice) {
                form.setValue("customPrice", {
                    amount: subscriptionToEdit.customPrice.amount,
                    currency: subscriptionToEdit.customPrice.currency
                });
            }

            // Set free trial if available
            if (subscriptionToEdit.freeTrial) {
                form.setValue("hasFreeTrialPeriod", true);
                form.setValue("freeTrialStartDate", subscriptionToEdit.freeTrial.startDate);
                form.setValue("freeTrialEndDate", subscriptionToEdit.freeTrial.endDate);
            }

            setIsLoading(false);
        }
    }, [isEditMode, subscriptionToEdit, form]);

    const onSubmit = async (data: FormValues) => {
        try {
            setError(null);

            // Calculate custom recurrency in days if needed
            let customRecurrencyInDays: number | undefined = undefined;
            if (data.recurrency === "custom" && data.customRecurrencyValue && data.customRecurrencyUnit) {
                customRecurrencyInDays = convertToDays(data.customRecurrencyValue, data.customRecurrencyUnit);
            }

            const subscriptionData = {
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
            };

            if (isEditMode && subscriptionId) {
                // Update existing subscription
                await updateSubscriptionMutation.mutateAsync({
                    subscriptionId,
                    subscriptionData
                });
            } else {
                // Create new subscription
                await createSubscriptionMutation.mutateAsync(subscriptionData);
            }

            navigate("/subscriptions");
        } catch (err) {
            setError(isEditMode
                ? "Failed to update subscription. Please try again."
                : "Failed to create subscription. Please try again.");
            console.error(err);
        }
    };

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title={isEditMode ? "Edit Subscription" : "Add New Subscription"}
                description={isEditMode
                    ? "Edit your subscription by modifying the form below."
                    : "Create a new subscription by filling out the form below."}
                actionButton={
                    <Button onClick={() => navigate("/subscriptions")}>
                        Back to Subscriptions
                    </Button>
                }
            />

            {isLoading && (
                <div className="flex justify-center my-8">
                    <div className="flex items-center space-x-2">
                        <Loader2 className="h-6 w-6 animate-spin text-primary"/>
                        <p>Loading subscription data...</p>
                    </div>
                </div>
            )}

            <div className="max-w-3xl mx-auto mt-8">
                {error && (
                    <Alert variant="destructive" className="mb-6">
                        <AlertCircle className="h-4 w-4"/>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}
                <FormProvider {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                        {/* Form Sections */}
                        <BasicInformationSection providers={providers}/>
                        <RecurrencySection/>
                        <DatesSection/>
                        <OwnershipSection families={families}/>
                        <CustomPriceSection/>
                        <FreeTrialSection/>

                        {/* Submit Button */}
                        <div className="flex justify-end">
                            <Button type="submit" disabled={form.formState.isSubmitting}>
                                {form.formState.isSubmitting ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                        {isEditMode ? "Updating..." : "Creating..."}
                                    </>
                                ) : (
                                    isEditMode ? "Update Subscription" : "Create Subscription"
                                )}
                            </Button>
                        </div>
                    </form>
                </FormProvider>
            </div>
        </div>
    );
};

export default SubscriptionFormPage;