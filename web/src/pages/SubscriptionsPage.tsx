import {useEffect, useMemo, useRef, useState} from "react";
import {useNavigate} from "react-router-dom";
import {useSubscriptionsQuery} from "@/hooks/subscriptions/useSubscriptionsQuery.ts";
import {useProvidersByIds} from "@/hooks/providers/useProvidersByIds";
import {useSubscriptionsMutations} from "@/hooks/subscriptions/useSubscriptionsMutations";
import {PageHeader} from "@/components/ui/page-header";
import {Button} from "@/components/ui/button";
import {PlusIcon} from "lucide-react";
import Subscription from "@/models/subscription";
import {DeleteSubscriptionDialog} from "@/components/subscriptions/DeleteSubscriptionDialog";
import { SubscriptionsTable } from "@/components/subscriptions/ui/SubscriptionsTable";
import { SubscriptionsTableSkeleton } from "@/components/subscriptions/ui/SubscriptionsTableSkeleton";
import { SubscriptionsErrorState } from "@/components/subscriptions/ui/SubscriptionsErrorState";
import { SubscriptionsEmptyState } from "@/components/subscriptions/ui/SubscriptionsEmptyState";
import type { SortingState } from "@tanstack/react-table";

const SubscriptionsPage = () => {
    const navigate = useNavigate();
    const [searchText, setSearchText] = useState("");
    const [subscriptionToDelete, setSubscriptionToDelete] = useState<Subscription | null>(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [sorting, setSorting] = useState<SortingState>([]);

    const sentinelRef = useRef<HTMLDivElement | null>(null);

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
        isError,
        hasNextPage,
        fetchNextPage,
        isFetchingNextPage,
    } = useSubscriptionsQuery({ 
        search: searchText || undefined,
        sortBy: sorting.length > 0 ? (sorting[0].id as 'provider' | 'name' | 'price' | 'recurrency' | 'dates' | 'status') : undefined,
        sortOrder: sorting.length > 0 ? (sorting[0].desc ? 'desc' : 'asc') : undefined,
    });

    // Flatten all subscriptions from all pages
    const allSubscriptions = data?.pages.flatMap(page => page.subscriptions) || [];

    // Build a unique list of provider IDs from subscriptions
    const providerIds = useMemo(() => {
        const ids = new Set<string>();
        allSubscriptions.forEach(s => {
            if (s.providerId) ids.add(s.providerId);
        });
        return Array.from(ids);
    }, [allSubscriptions]);

    // Fetch providers individually and build a map (cached by React Query)
    const {providerMap} = useProvidersByIds(providerIds);


    // Infinite scroll: observe sentinel and fetch next page when visible
    useEffect(() => {
        const el = sentinelRef.current;
        if (!el) return;

        const observer = new IntersectionObserver(
            (entries) => {
                const [entry] = entries;
                if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
                    fetchNextPage();
                }
            },
            {root: null, rootMargin: '400px', threshold: 0}
        );

        observer.observe(el);
        return () => {
            observer.unobserve(el);
            observer.disconnect();
        };
    }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

    // Use Money component for currency display (conversion handled in UI).





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
                <SubscriptionsTableSkeleton />
            ) : isError ? (
                <SubscriptionsErrorState />
            ) : (
                <>
                    {allSubscriptions.length > 0 ? (
                        <SubscriptionsTable
                            subscriptions={allSubscriptions}
                            providerMap={providerMap}
                            onEdit={(s) => navigate(`/subscriptions/edit/${s.id}`)}
                            onDelete={handleDeleteClick}
                            isFetchingNextPage={isFetchingNextPage}
                            sorting={sorting}
                            onSortingChange={setSorting}
                        />
                    ) : searchText !== "" ? (
                        <div className="text-center mt-8">
                            <p className="text-muted-foreground">No subscriptions match your search criteria.</p>
                        </div>
                    ) : (
                        <SubscriptionsEmptyState />
                    )}
                </>
            )}

            {hasNextPage && (
                <div ref={sentinelRef} className="h-1" aria-hidden/>
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