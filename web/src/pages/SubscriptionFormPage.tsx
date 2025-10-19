import {useEffect, useState} from "react";
import {FormProvider, useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {useNavigate, useParams} from "react-router-dom";
import {useSubscriptionsMutations} from "@/hooks/subscriptions/useSubscriptionsMutations";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import {useFamilyQuery} from "@/hooks/families/useFamilyQuery.ts";
import {useSubscriptionsQuery} from "@/hooks/subscriptions/useSubscriptionsQuery.ts";
import {PageHeader} from "@/components/ui/page-header";
import {Button} from "@/components/ui/button";
import {Alert, AlertDescription} from "@/components/ui/alert";
import {AlertCircle, ChevronLeft, ChevronRight, Loader2} from "lucide-react";
import {OwnerType} from "@/models/ownerType";
import {convertToDays, formSchema, type FormValues} from "@/components/subscriptions/form/SubscriptionFormSchema";
import {BasicInformationSection} from "@/components/subscriptions/form/BasicInformationSection";
import {RecurrencySection} from "@/components/subscriptions/form/RecurrencySection";
import {DatesSection} from "@/components/subscriptions/form/DatesSection";
import {OwnershipSection} from "@/components/subscriptions/form/OwnershipSection";
import {FreeTrialSection} from "@/components/subscriptions/form/FreeTrialSection";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";

const SubscriptionFormPage = () => {
    const navigate = useNavigate();
    const {subscriptionId} = useParams<{ subscriptionId: string }>();
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [currentStep, setCurrentStep] = useState(0);
    const {createSubscriptionMutation, updateSubscriptionMutation} = useSubscriptionsMutations();
    const {data: providersData} = useAllProvidersQuery();
    const {data: familyData} = useFamilyQuery();
    const {data: subscriptionsData} = useSubscriptionsQuery();

    const isEditMode = !!subscriptionId;
    const providers = providersData?.pages.flatMap(page => page.providers) || [];

    // Find the subscription to edit if in edit mode
    const subscriptionToEdit = isEditMode
        ? subscriptionsData?.pages.flatMap(page => page.subscriptions).find(sub => sub.id === subscriptionId)
        : undefined;

    // Define the steps for the wizard
    const steps = [
        {title: "Basic Information", component: BasicInformationSection},
        {title: "Recurrency", component: RecurrencySection},
        {title: "Dates", component: DatesSection},
        {title: "Ownership", component: OwnershipSection},
        {title: "Free Trial", component: FreeTrialSection},
    ];

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
            }, // customPrice is now required
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
                    amount: subscriptionToEdit.customPrice.value,
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
                planId: data.planId === "" ? undefined : data.planId,
                priceId: data.priceId === "" ? undefined : data.priceId,
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

            <div className="max-w-3xl mx-auto mt-8 bg-white dark:bg-black rounded-lg shadow-md p-6 border border-gray-200 dark:border-gray-700">
                {error && (
                    <Alert variant="destructive" className="mb-6">
                        <AlertCircle className="h-4 w-4"/>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}
                <FormProvider {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit, (errors) => {
                        // Identify which section contains errors and redirect to that section
                        const errorFields = Object.keys(errors);

                        // Map fields to their respective sections
                        const fieldToSectionMap: Record<string, number> = {
                            // Basic Information (Section 0)
                            'providerId': 0,
                            'planId': 0,
                            'priceId': 0,
                            'friendlyName': 0,

                            // Recurrency (Section 1)
                            'recurrency': 1,
                            'customRecurrencyValue': 1,
                            'customRecurrencyUnit': 1,

                            // Dates (Section 2)
                            'startDate': 2,
                            'endDate': 2,

                            // Ownership (Section 3)
                            'ownerType': 3,
                            'familyId': 3,
                            'serviceUsers': 3,

                            // Free Trial (Section 4)
                            'hasFreeTrialPeriod': 4,
                            'freeTrialStartDate': 4,
                            'freeTrialEndDate': 4,
                        };

                        // Find the earliest section with errors
                        let earliestSectionWithError = steps.length - 1;
                        for (const field of errorFields) {
                            const sectionIndex = fieldToSectionMap[field];
                            if (sectionIndex !== undefined && sectionIndex < earliestSectionWithError) {
                                earliestSectionWithError = sectionIndex;
                            }
                        }

                        // Set the current step to the section with errors
                        setCurrentStep(earliestSectionWithError);

                        // Create a detailed error message showing which fields have errors
                        const fieldErrorMessages = errorFields.map(field => {
                            const error = errors[field as keyof typeof errors];
                            const fieldName = field.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase());
                            return `${fieldName}: ${error?.message || 'Invalid value'}`;
                        });

                        const detailedErrorMessage = fieldErrorMessages.length > 0 
                            ? `Please correct the following errors: ${fieldErrorMessages.join(', ')}`
                            : "Please correct the errors in the form before submitting.";
                        
                        setError(detailedErrorMessage);
                    })} className="space-y-6">
                        {/* Step Indicator */}
                        <div className="mb-6">
                            <div className="flex justify-between items-center mb-2">
                                {steps.map((step, index) => (
                                    <div
                                        key={index}
                                        className={`flex flex-col items-center cursor-pointer ${
                                            index === currentStep
                                                ? "text-primary font-medium"
                                                : index < currentStep
                                                    ? "text-primary/70"
                                                    : "text-muted-foreground"
                                        }`}
                                        onClick={() => setCurrentStep(index)}
                                    >
                                        <div
                                            className={`w-8 h-8 rounded-full flex items-center justify-center mb-1 ${
                                                index === currentStep
                                                    ? "bg-primary text-primary-foreground"
                                                    : index < currentStep
                                                        ? "bg-primary/20 text-primary"
                                                        : "bg-muted text-muted-foreground"
                                            }`}
                                        >
                                            {index + 1}
                                        </div>
                                        <span className="text-xs">{step.title}</span>
                                    </div>
                                ))}
                            </div>
                            <div className="w-full bg-muted h-1 rounded-full overflow-hidden">
                                <div
                                    className="bg-primary h-full transition-all duration-300 ease-in-out"
                                    style={{width: `${((currentStep + 1) / steps.length) * 100}%`}}
                                />
                            </div>
                        </div>

                        {/* Current Step Component */}
                        <div className="py-4">
                            {currentStep === 0 && <BasicInformationSection providers={providers}/>}
                            {currentStep === 1 && <RecurrencySection/>}
                            {currentStep === 2 && <DatesSection/>}
                            {currentStep === 3 && <OwnershipSection family={familyData}/>}
                            {currentStep === 4 && <FreeTrialSection/>}
                        </div>

                        {/* Navigation Buttons */}
                        <div className="flex justify-between pt-4">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => setCurrentStep(prev => Math.max(0, prev - 1))}
                                disabled={currentStep === 0}
                            >
                                <ChevronLeft className="mr-2 h-4 w-4"/>
                                Previous
                            </Button>

                            {currentStep < steps.length - 1 ? (
                                <Button
                                    type="button"
                                    onClick={async () => {
                                        // Validate current step fields before proceeding
                                        let isValid = true;

                                        // Step 0: Basic Information
                                        if (currentStep === 0) {
                                            const result = await form.trigger(['providerId']);
                                            isValid = result;
                                        }
                                        // Step 1: Recurrency
                                        else if (currentStep === 1) {
                                            const result = await form.trigger(['recurrency']);
                                            isValid = result;

                                            // If custom recurrency, validate custom fields
                                            if (form.getValues('recurrency') === 'custom') {
                                                const customResult = await form.trigger(['customRecurrencyValue', 'customRecurrencyUnit']);
                                                isValid = isValid && customResult;
                                            }
                                        }
                                        // Step 2: Dates
                                        else if (currentStep === 2) {
                                            const result = await form.trigger(['startDate']);
                                            isValid = result;
                                        }
                                        // Step 3: Ownership
                                        else if (currentStep === 3) {
                                            const result = await form.trigger(['ownerType']);
                                            isValid = result;

                                            // If family ownership, validate family ID
                                            if (form.getValues('ownerType') === 'family') {
                                                const familyResult = await form.trigger(['familyId']);
                                                isValid = isValid && familyResult;
                                            }
                                        }
                                        // Step 4: Free Trial
                                        else if (currentStep === 4) {
                                            const result = await form.trigger(['hasFreeTrialPeriod']);
                                            isValid = result;

                                            // If has free trial period, validate start and end dates
                                            if (form.getValues('hasFreeTrialPeriod')) {
                                                const trialResult = await form.trigger(['freeTrialStartDate', 'freeTrialEndDate']);
                                                isValid = isValid && trialResult;
                                            }
                                        }

                                        if (isValid) {
                                            setCurrentStep(prev => Math.min(steps.length - 1, prev + 1));
                                        }
                                    }}
                                >
                                    Next
                                    <ChevronRight className="ml-2 h-4 w-4"/>
                                </Button>
                            ) : (
                                <Button
                                    type="submit"
                                    disabled={form.formState.isSubmitting}
                                >
                                    {form.formState.isSubmitting ? (
                                        <>
                                            <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                            {isEditMode ? "Updating..." : "Creating..."}
                                        </>
                                    ) : (
                                        isEditMode ? "Update Subscription" : "Create Subscription"
                                    )}
                                </Button>
                            )}
                        </div>
                    </form>
                </FormProvider>
            </div>
        </div>
    );
};

export default SubscriptionFormPage;