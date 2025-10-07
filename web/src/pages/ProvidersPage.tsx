import { useEffect, useRef, useState } from "react";
import { useAllProvidersQuery } from "@/hooks/providers/useAllProvidersQuery";
import { AddProviderForm } from "@/components/providers/AddProviderForm";
import { EditProviderForm } from "@/components/providers/EditProviderForm";
import { ProviderCard } from "@/components/providers/ui/ProviderCard";
import { ProviderCardSkeletonGrid } from "@/components/providers/ui/ProviderCardSkeleton";
import { ErrorState } from "@/components/providers/ui/ErrorState";
import { NoProviders } from "@/components/providers/ui/EmptyStates";
import { PageHeader } from "@/components/ui/page-header";
import { Button } from "@/components/ui/button";
import Provider from "@/models/provider";
import { PlusIcon } from "lucide-react";
import { useProvidersQuotaQuery } from "@/hooks/providers/useProvidersQuotaQuery.ts";
import { QuotaButton } from "@/components/quotas/QuotaButton";
import { FeatureId } from "@/models/billing.ts";
import { useQuotaLimit, getQuotaTooltip } from "@/hooks/quotas/useFeature";
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";

const ProvidersPage = () => {
    const [isAddingProvider, setIsAddingProvider] = useState(false);
    const [editingProvider, setEditingProvider] = useState<Provider | null>(null);
    const [searchText, setSearchText] = useState("");

    // Check providers quota
    const { data: providersQuota } = useProvidersQuotaQuery();
    const { enabled: providersEnabled, canAdd: canAddProviders, used: providersUsed, limit: providersLimit } = useQuotaLimit(
        providersQuota,
        FeatureId.CustomProvidersCount
    );
    const isDisabled = !providersEnabled || !canAddProviders;
    const tooltipMessage = getQuotaTooltip(providersEnabled, canAddProviders, "custom providers");

    // Query all providers using the dedicated hook
    const {
        data,
        isLoading,
        isError,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
    } = useAllProvidersQuery({ search: searchText });

    // Flatten all providers from all pages
    const allProviders = data?.pages.flatMap(page => page.providers) || [];

    const loadMoreRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (!loadMoreRef.current) return;
        const node = loadMoreRef.current;
        const observer = new IntersectionObserver((entries) => {
            const entry = entries[0];
            if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
                fetchNextPage();
            }
        }, { rootMargin: "200px" });
        observer.observe(node);
        return () => observer.disconnect();
    }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Providers"
                description="Manage your providers"
                searchText={searchText}
                onSearchChange={setSearchText}
                quotaButton={
                    <QuotaButton
                        useQuotaQuery={useProvidersQuotaQuery}
                        featureIds={[FeatureId.CustomProvidersCount]}
                        featureLabels={{
                            [FeatureId.CustomProvidersCount]: "Providers Used"
                        }}
                    />
                }
                actionButton={
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <span>
                                    <Button
                                        onClick={() => setIsAddingProvider(true)}
                                        disabled={isDisabled}
                                    >
                                        <PlusIcon className="mr-2 h-4 w-4" />
                                        Add Provider
                                        {providersEnabled && providersLimit !== undefined && (
                                            <span className="ml-2 text-xs opacity-70">
                                                ({providersUsed}/{providersLimit})
                                            </span>
                                        )}
                                    </Button>
                                </span>
                            </TooltipTrigger>
                            {tooltipMessage && (
                                <TooltipContent>
                                    <p>{tooltipMessage}</p>
                                </TooltipContent>
                            )}
                        </Tooltip>
                    </TooltipProvider>
                }
            />

            {isLoading ? (
                <ProviderCardSkeletonGrid />
            ) : isError ? (
                <ErrorState />
            ) : (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                        {allProviders.map((provider) => (
                            <ProviderCard
                                key={provider.id}
                                provider={provider}
                                onEdit={setEditingProvider}
                            />
                        ))}
                    </div>

                    {searchText.trim() && allProviders.length === 0 && (
                        <div className="text-center mt-8">
                            <p className="text-muted-foreground">No providers match your search criteria.</p>
                        </div>
                    )}

                    <NoProviders
                        show={!searchText.trim() && allProviders.length === 0}
                        onAddProvider={() => setIsAddingProvider(true)}
                    />

                    {/* Infinite scroll sentinel - only show when there are results */}
                    {allProviders.length > 0 && (
                        <div ref={loadMoreRef} className="flex justify-center items-center py-4">
                            {isFetchingNextPage && (
                                <div className="flex items-center gap-2 text-muted-foreground text-sm">
                                    <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none"></circle>
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
                                    </svg>
                                    Loading more providers...
                                </div>
                            )}
                            {!hasNextPage && (
                                <div className="text-xs text-muted-foreground">You have reached the end.</div>
                            )}
                        </div>
                    )}
                </>
            )}

            {/* Provider form modals */}
            <AddProviderForm
                isOpen={isAddingProvider}
                onClose={() => setIsAddingProvider(false)}
            />

            {editingProvider && (
                <EditProviderForm
                    isOpen={!!editingProvider}
                    onClose={() => setEditingProvider(null)}
                    provider={editingProvider}
                />
            )}
        </div>
    );
};

export default ProvidersPage;