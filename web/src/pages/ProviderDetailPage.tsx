import { useEffect, useMemo, useRef } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { PageHeader } from "@/components/ui/page-header";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useProvidersByIds } from "@/hooks/providers/useProvidersByIds";
import { useSubscriptionsQuery } from "@/hooks/subscriptions/useSubscriptionsQuery";
import { SubscriptionsTable } from "@/components/subscriptions/ui/SubscriptionsTable";
import { SubscriptionsTableSkeleton } from "@/components/subscriptions/ui/SubscriptionsTableSkeleton";
import { SubscriptionsErrorState } from "@/components/subscriptions/ui/SubscriptionsErrorState";
import { SubscriptionsEmptyState } from "@/components/subscriptions/ui/SubscriptionsEmptyState";
import { useSubscriptionsMutations } from "@/hooks/subscriptions/useSubscriptionsMutations";
import { useLabelQuery } from "@/hooks/labels/useLabelQuery";
import { argbToRgba } from "@/components/ui/utils/color-utils.ts";
import type Subscription from "@/models/subscription";
import type Provider from "@/models/provider";

const LabelPill = ({ labelId }: { labelId: string }) => {
  const { data: label, isLoading } = useLabelQuery(labelId);
  return (
    <Badge
      variant="outline"
      className="text-xs py-0.5"
      style={{ backgroundColor: label?.color ? argbToRgba(label.color) : undefined }}
    >
      {label ? label.name : isLoading ? "..." : labelId}
    </Badge>
  );
};

const ProviderDetailPage = () => {
  const { providerId } = useParams<{ providerId: string }>();
  const navigate = useNavigate();

  const { providerMap, isLoading: isProviderLoading } = useProvidersByIds(
    providerId ? [providerId] : []
  );
  const provider = providerId ? providerMap.get(providerId) : undefined;

  const {
    data: subsData,
    isLoading: isSubsLoading,
    isError: isSubsError,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useSubscriptionsQuery({ providers: providerId ? [providerId] : [], withInactive: true, limit: 10 });

  const allSubscriptions = useMemo(() => subsData?.pages.flatMap(p => p.subscriptions) ?? [], [subsData]);

  const providerMapForTable = useMemo(() => {
    const map = new Map<string, Provider>();
    if (provider) map.set(provider.id, provider);
    return map;
  }, [provider]);

  // Infinite scroll sentinel for subscriptions list
  const sentinelRef = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    if (!sentinelRef.current) return;
    const node = sentinelRef.current;
    const observer = new IntersectionObserver((entries) => {
      const entry = entries[0];
      if (entry.isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage();
      }
    }, { rootMargin: "200px" });
    observer.observe(node);
    return () => observer.disconnect();
  }, [fetchNextPage, hasNextPage, isFetchingNextPage]);

  const { deleteSubscriptionMutation } = useSubscriptionsMutations();

  const handleEditSubscription = (sub: Subscription) => {
    navigate(`/subscriptions/edit/${sub.id}`);
  };

  const handleDeleteSubscription = async (sub: Subscription) => {
    try {
      await deleteSubscriptionMutation.mutateAsync(sub.id);
      // No explicit refetch needed; rely on cache updates or manual refetch if required in future
    } catch (e) {
      console.error("Failed to delete subscription", e);
    }
  };

  return (
    <div className="container mx-auto py-6">
      <PageHeader
        title={provider?.name || "Provider"}
        description="Provider details and subscriptions"
        actionButton={
          <Button variant="secondary" onClick={() => navigate(-1)}>
            Back
          </Button>
        }
      />

      <Card className="mb-6">
        <CardContent className="p-4">
          <div className="flex items-start gap-4">
            {provider?.iconUrl && (
              <img
                src={provider.iconUrl}
                alt={`${provider.name} logo`}
                className="h-16 w-16 object-contain"
              />
            )}
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <h2 className="text-xl font-semibold">{provider?.name ?? (isProviderLoading ? "Loading..." : "Unknown provider")}</h2>
              </div>
              {provider?.description && (
                <p className="text-muted-foreground mt-1">{provider.description}</p>
              )}
              {provider?.labels?.length ? (
                <div className="flex flex-wrap gap-1 mt-2">
                  {provider.labels.map((id) => (
                    <LabelPill key={id} labelId={id} />
                  ))}
                </div>
              ) : null}
              <div className="flex flex-wrap gap-2 mt-3">
                {provider?.url && (
                  <a href={provider.url} target="_blank" rel="noopener noreferrer">
                    <Button variant="outline" size="sm" onClick={(e) => e.stopPropagation()}>Website</Button>
                  </a>
                )}
                {provider?.pricingPageUrl && (
                  <a href={provider.pricingPageUrl} target="_blank" rel="noopener noreferrer">
                    <Button variant="outline" size="sm" onClick={(e) => e.stopPropagation()}>Pricing</Button>
                  </a>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Subscriptions List */}
      {isSubsLoading ? (
        <SubscriptionsTableSkeleton />
      ) : isSubsError ? (
        <SubscriptionsErrorState />
      ) : allSubscriptions.length === 0 ? (
        <SubscriptionsEmptyState />
      ) : (
        <>
          <SubscriptionsTable
            subscriptions={allSubscriptions}
            providerMap={providerMapForTable}
            onEdit={handleEditSubscription}
            onDelete={handleDeleteSubscription}
            isFetchingNextPage={isFetchingNextPage}
          />
          <div ref={sentinelRef} className="py-4 flex justify-center items-center">
            {!hasNextPage && allSubscriptions.length > 0 && (
              <div className="text-xs text-muted-foreground">You have reached the end.</div>
            )}
          </div>
        </>
      )}
    </div>
  );
};

export default ProviderDetailPage;
