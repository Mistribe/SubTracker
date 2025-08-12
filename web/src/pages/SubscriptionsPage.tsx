import {useMemo, useState} from "react";
import {useNavigate} from "react-router-dom";
import {useSubscriptionsQuery} from "@/hooks/subscriptions/useSubscriptionsQuery.ts";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import {useSubscriptionsMutations} from "@/hooks/subscriptions/useSubscriptionsMutations";
import {PageHeader} from "@/components/ui/page-header";
import {Skeleton} from "@/components/ui/skeleton";
import {Badge} from "@/components/ui/badge";
import {Button} from "@/components/ui/button";
import {CalendarIcon, CreditCardIcon, PencilIcon, PlusIcon, TagIcon, TrashIcon, UsersIcon} from "lucide-react";
import {format} from "date-fns";
import Subscription from "@/models/subscription";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Money} from "@/components/ui/money";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";
import {DeleteSubscriptionDialog} from "@/components/subscriptions/DeleteSubscriptionDialog";

const SubscriptionsPage = () => {
    const navigate = useNavigate();
    const [searchText, setSearchText] = useState("");
    const [subscriptionToDelete, setSubscriptionToDelete] = useState<Subscription | null>(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    const {deleteSubscriptionMutation} = useSubscriptionsMutations();

    // Handler to open the delete dialog
    const handleDeleteClick = (subscription: Subscription) => {
        setSubscriptionToDelete(subscription);
        setIsDeleteDialogOpen(true);
    };

    // Handler to close the delete dialog
    const handleDeleteDialogClose = () => {
        setIsDeleteDialogOpen(false);
    };

    // Handler to confirm deletion
    const handleDeleteConfirm = async () => {
        if (!subscriptionToDelete) return;

        try {
            setIsDeleting(true);
            await deleteSubscriptionMutation.mutateAsync(subscriptionToDelete.id);
            setIsDeleteDialogOpen(false);
        } catch (error) {
            console.error("Error deleting subscription:", error);
        } finally {
            setIsDeleting(false);
            setSubscriptionToDelete(null);
        }
    };

    // Query all subscriptions using the dedicated hook
    const {
        data,
        isLoading,
        isError
    } = useSubscriptionsQuery();

    // Query all providers to resolve provider names from IDs
    const {
        data: providersData
    } = useAllProvidersQuery();

    // Flatten all subscriptions from all pages
    const allSubscriptions = data?.pages.flatMap(page => page.subscriptions) || [];

    // Flatten all providers from all pages and create a mapping from ID to provider
    const allProviders = providersData?.pages.flatMap(page => page.providers) || [];
    const providerMap = useMemo(() => {
        const map = new Map();
        allProviders.forEach(provider => {
            map.set(provider.id, provider);
        });
        return map;
    }, [allProviders]);

    // Filter subscriptions based on search text
    const filteredSubscriptions = useMemo(() => {
        if (searchText === "") return allSubscriptions;

        const searchLower = searchText.toLowerCase();

        return allSubscriptions.filter(subscription => {
            // Filter by subscription friendly name
            if (subscription.friendlyName && subscription.friendlyName.toLowerCase().includes(searchLower)) {
                return true;
            }

            // Filter by provider ID, plan ID, or price ID
            if (subscription.providerId.toLowerCase().includes(searchLower) ||
                subscription.planId.toLowerCase().includes(searchLower) ||
                subscription.priceId.toLowerCase().includes(searchLower)) {
                return true;
            }

            return false;
        });
    }, [allSubscriptions, searchText]);

    // Use Money component for currency display (conversion handled in UI).

    // Function to format recurrency
    const formatRecurrency = (recurrency: SubscriptionRecurrency, customRecurrency: number | undefined) => {
        if (recurrency === 'custom' && customRecurrency) {
            return `Every ${customRecurrency} days`;
        }

        switch (recurrency) {
            case SubscriptionRecurrency.Monthly:
                return 'Monthly';
            case SubscriptionRecurrency.Quarterly:
                return 'Quarterly';
            case SubscriptionRecurrency.HalfYearly:
                return 'Half Yearly';
            case SubscriptionRecurrency.Yearly:
                return 'Yearly';
            case SubscriptionRecurrency.OneTime:
                return 'OneTime';
            default:
                return recurrency;
        }
    };

    // Render subscriptions table
    const renderSubscriptionsTable = (subscriptions: Subscription[]) => {
        return (
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Provider</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Price</TableHead>
                        <TableHead>Recurrency</TableHead>
                        <TableHead>Dates</TableHead>
                        <TableHead>Users</TableHead>
                        <TableHead>Plan</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead className="w-[50px]">Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {subscriptions.map((subscription) => (
                        <TableRow key={subscription.id}>
                            <TableCell>
                                <div className="flex items-center">
                                    {providerMap.get(subscription.providerId)?.iconUrl ? (
                                        <img
                                            src={providerMap.get(subscription.providerId)?.iconUrl || ''}
                                            alt={`${providerMap.get(subscription.providerId)?.name} logo`}
                                            className="mr-2 h-5 w-5 object-contain"
                                        />
                                    ) : (
                                        <CreditCardIcon className="mr-2 h-4 w-4 text-muted-foreground"/>
                                    )}
                                    <span>
                                        {providerMap.get(subscription.providerId)?.name || subscription.providerId}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell className="font-medium">
                                {subscription.friendlyName || (providerMap.get(subscription.providerId)?.name || subscription.providerId)}
                            </TableCell>
                            <TableCell>
                                {subscription.customPrice && (
                                    <Badge variant="outline">
                                        <Money amount={subscription.customPrice.amount}
                                               currency={subscription.customPrice.currency}/>
                                    </Badge>
                                )}
                            </TableCell>
                            <TableCell>
                                {formatRecurrency(subscription.recurrency, subscription.customRecurrency)}
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center">
                                    <CalendarIcon className="mr-2 h-4 w-4 text-muted-foreground"/>
                                    <span>
                                        {format(subscription.startDate, 'MMM d, yyyy')}
                                        {subscription.endDate && <br/>}
                                        {subscription.endDate && `Ends: ${format(subscription.endDate, 'MMM d, yyyy')}`}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell>
                                {subscription.serviceUsers.length > 0 && (
                                    <div className="flex items-center">
                                        <UsersIcon className="mr-2 h-4 w-4 text-muted-foreground"/>
                                        <span>{subscription.serviceUsers.length}</span>
                                    </div>
                                )}
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center">
                                    <TagIcon className="mr-1 h-4 w-4 text-muted-foreground"/>
                                    <span>{subscription.planId.substring(0, 8)}</span>
                                </div>
                            </TableCell>
                            <TableCell>
                                {subscription.freeTrial ? "Free Trial" : ""}
                                {subscription.freeTrial && subscription.isActive ? " - " : ""}
                                {subscription.isActive
                                    ? <Badge variant="outline" className="bg-green-50 text-green-700">Active</Badge>
                                    : <Badge variant="outline" className="bg-red-50 text-red-700">Ended</Badge>}
                            </TableCell>
                            <TableCell>
                                <div className="flex space-x-1">
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        onClick={() => navigate(`/subscriptions/edit/${subscription.id}`)}
                                        title="Edit subscription"
                                    >
                                        <PencilIcon className="h-4 w-4 text-blue-500"/>
                                    </Button>
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        onClick={() => handleDeleteClick(subscription)}
                                        title="Delete subscription"
                                    >
                                        <TrashIcon className="h-4 w-4 text-red-500"/>
                                    </Button>
                                </div>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        );
    };

    // Render loading skeleton for table
    const renderTableSkeleton = () => {
        return (
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                        <TableHead><Skeleton className="h-4 w-20"/></TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {Array(8).fill(0).map((_, index) => (
                        <TableRow key={index}>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                            <TableCell><Skeleton className="h-4 w-full"/></TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        );
    };

    // Render error state
    const renderError = () => {
        return (
            <div className="flex flex-col items-center justify-center py-12">
                <h3 className="text-xl font-semibold mb-2">Error Loading Subscriptions</h3>
                <p className="text-muted-foreground mb-6">
                    There was a problem loading your subscriptions. Please try again later.
                </p>
            </div>
        );
    };

    // Render empty state
    const renderEmptyState = () => {
        return (
            <div className="flex flex-col items-center justify-center py-12">
                <h3 className="text-xl font-semibold mb-2">No Subscriptions Found</h3>
                <p className="text-muted-foreground mb-6">
                    You don't have any subscriptions yet.
                </p>
            </div>
        );
    };

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Subscriptions"
                description="Manage your subscriptions"
                searchText={searchText}
                onSearchChange={setSearchText}
                actionButton={
                    <Button onClick={() => navigate("/subscriptions/create")}>
                        <PlusIcon className="mr-2 h-4 w-4"/>
                        Add Subscription
                    </Button>
                }
            />

            {isLoading ? (
                renderTableSkeleton()
            ) : isError ? (
                renderError()
            ) : (
                <>
                    {filteredSubscriptions.length > 0 ? (
                        renderSubscriptionsTable(filteredSubscriptions)
                    ) : searchText !== "" ? (
                        <div className="text-center mt-8">
                            <p className="text-muted-foreground">No subscriptions match your search criteria.</p>
                        </div>
                    ) : (
                        renderEmptyState()
                    )}
                </>
            )}

            {/* Delete Subscription Dialog */}
            {subscriptionToDelete && (
                <DeleteSubscriptionDialog
                    subscription={subscriptionToDelete}
                    isOpen={isDeleteDialogOpen}
                    isSubmitting={isDeleting}
                    onClose={handleDeleteDialogClose}
                    onConfirm={handleDeleteConfirm}
                />
            )}
        </div>
    );
};

export default SubscriptionsPage;