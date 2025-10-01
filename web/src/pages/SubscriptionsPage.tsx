import { useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useSubscriptionsQuery } from "@/hooks/subscriptions/useSubscriptionsQuery.ts";
import { useProvidersByIds } from "@/hooks/providers/useProvidersByIds";
import { useSubscriptionsMutations } from "@/hooks/subscriptions/useSubscriptionsMutations";
import { PageHeader } from "@/components/ui/page-header";
import { Button } from "@/components/ui/button";
import { PlusIcon } from "lucide-react";
import Subscription from "@/models/subscription";
import { DeleteSubscriptionDialog } from "@/components/subscriptions/DeleteSubscriptionDialog";
import { SubscriptionsTable } from "@/components/subscriptions/ui/SubscriptionsTable";
import { SubscriptionsTableSkeleton } from "@/components/subscriptions/ui/SubscriptionsTableSkeleton";
import { SubscriptionsErrorState } from "@/components/subscriptions/ui/SubscriptionsErrorState";
import { SubscriptionsEmptyState } from "@/components/subscriptions/ui/SubscriptionsEmptyState";
import {
    SubscriptionsFilters,
    type SubscriptionsFiltersValues
} from "@/components/subscriptions/ui/SubscriptionsFilters";
import { SubscriptionRecurrency } from "@/models/subscriptionRecurrency.ts";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";
import { useSubscriptionsQuotaQuery } from "@/hooks/subscriptions/useSubscriptionsQuotaQuery.ts";
import { QuotaButton } from "@/components/quotas/QuotaButton";
import { FeatureId } from "@/models/billing.ts";
import { useQuotaLimit, getQuotaTooltip } from "@/hooks/quotas/useFeature";
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";

const SubscriptionsPage = () => {
    const navigate = useNavigate();
    const [searchText, setSearchText] = useState("");
    const [subscriptionToDelete, setSubscriptionToDelete] = useState<Subscription | null>(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    // Filters state for drawer (draft values)
    const [isFilterOpen, setIsFilterOpen] = useState(false);
    const [fromDate, setFromDate] = useState<Date | undefined>(undefined);
    const [toDate, setToDate] = useState<Date | undefined>(undefined);
    const [providersFilter, setProvidersFilter] = useState<string[]>([]);
    const [recurrenciesFilter, setRecurrenciesFilter] = useState<SubscriptionRecurrency[]>([]);
    const [users, setUsers] = useState<string[]>([]);
    const [withInactive, setWithInactive] = useState<boolean>(false);

    // Family and family members for users multi-select (single family context)
    const { data: familyData } = useFamilyQuery();

    // Check subscriptions quota
    const { data: subscriptionsQuota } = useSubscriptionsQuotaQuery();
    const { enabled: subsEnabled, canAdd: canAddSubs, used: subsUsed, limit: subsLimit } = useQuotaLimit(
        subscriptionsQuota,
        FeatureId.ActiveSubscriptionsCount
    );
    const isAddDisabled = !subsEnabled || !canAddSubs;
    const addSubTooltip = getQuotaTooltip(subsEnabled, canAddSubs, "subscriptions");

    const usersOptions = useMemo(() => {
        const members = familyData?.members ?? [];
        const labelByName = new Map<string, string>();
        for (const m of members) {
            if (!m.name) continue;
            if (!labelByName.has(m.name)) {
                labelByName.set(m.name, m.name);
            }
        }
        return Array.from(labelByName.entries()).map(([value, label]) => ({ value, label }));
    }, [familyData]);


    const userNameMap = useMemo(() => {
        const members = familyData?.members ?? [];
        const map = new Map<string, string>();
        for (const m of members) {
            map.set(m.id, m.name);
        }
        return map;
    }, [familyData]);

    const handleClearFilters = () => {
        setFromDate(undefined);
        setToDate(undefined);
        setProvidersFilter([]);
        setRecurrenciesFilter([]);
        setUsers([]);
        setWithInactive(false);
        setIsFilterOpen(false);
    };


    const sentinelRef = useRef<HTMLDivElement | null>(null);

    const { deleteSubscriptionMutation } = useSubscriptionsMutations();

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
        limit: 10,
        search: searchText || undefined,
        fromDate,
        toDate,
        providers: providersFilter.length ? providersFilter : undefined,
        recurrencies: recurrenciesFilter.length ? recurrenciesFilter : undefined,
        users: users.length ? users : undefined,
        withInactive: withInactive,
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
    const { providerMap } = useProvidersByIds(providerIds);


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
            { root: null, rootMargin: '400px', threshold: 0 }
        );

        observer.observe(el);
        return () => {
            observer.unobserve(el);
            observer.disconnect();
        };
    }, [hasNextPage, isFetchingNextPage, fetchNextPage]);


    const handleApplyFilters = (values: SubscriptionsFiltersValues) => {
        setFromDate(values.fromDate);
        setToDate(values.toDate);
        setProvidersFilter(values.providers ?? []);
        setRecurrenciesFilter(values.recurrenciesFilter ?? []);
        setUsers(values.users ?? []);
        setWithInactive(values.withInactive ?? false);
        setIsFilterOpen(false);
    };

    const openFilter = () => {
        setIsFilterOpen(true);
    };

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Subscriptions"
                description="Manage your subscriptions"
                searchText={searchText}
                onSearchChange={setSearchText}
                quotaButton={
                    <QuotaButton
                        useQuotaQuery={useSubscriptionsQuotaQuery}
                        featureIds={[FeatureId.ActiveSubscriptionsCount]}
                        featureLabels={{
                            [FeatureId.ActiveSubscriptionsCount]: "Active Subscriptions"
                        }}
                    />
                }
                actionButton={
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <span>
                                    <Button
                                        onClick={() => navigate("/subscriptions/create")}
                                        disabled={isAddDisabled}
                                    >
                                        <PlusIcon className="mr-2 h-4 w-4" />
                                        Add Subscription
                                        {subsEnabled && subsLimit !== undefined && (
                                            <span className="ml-2 text-xs opacity-70">
                                                ({subsUsed}/{subsLimit})
                                            </span>
                                        )}
                                    </Button>
                                </span>
                            </TooltipTrigger>
                            {addSubTooltip && (
                                <TooltipContent>
                                    <p>{addSubTooltip}</p>
                                </TooltipContent>
                            )}
                        </Tooltip>
                    </TooltipProvider>
                }
                onFilter={openFilter}
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
                            userNameMap={userNameMap}
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
                <div ref={sentinelRef} className="h-1" aria-hidden />
            )}

            <SubscriptionsFilters
                open={isFilterOpen}
                onOpenChange={setIsFilterOpen}
                fromDate={fromDate}
                toDate={toDate}
                providerMap={providerMap}
                providersFilter={providersFilter}
                recurrenciesFilter={recurrenciesFilter}
                users={users}
                usersOptions={usersOptions}
                withInactive={withInactive}
                onClear={handleClearFilters}
                onApply={handleApplyFilters}
            />

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

