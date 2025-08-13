import {useMemo} from "react";
import {useSubscriptionsQuery} from "@/hooks/subscriptions/useSubscriptionsQuery.ts";
import {useProvidersByIds} from "@/hooks/providers/useProvidersByIds";
import SummaryCards from "@/components/dashboard/SummaryCards";
import UpcomingRenewals from "@/components/dashboard/UpcomingRenewals";
import TopProviders from "@/components/dashboard/TopProviders";
import PriceEvolutionGraph from "@/components/dashboard/PriceEvolutionGraph";
import {PageHeader} from "@/components/ui/page-header";
import {useSubscriptionSummaryQuery} from "@/hooks/subscriptions/useSubscriptionSummaryQuery";
import type Subscription from "@/models/subscription.ts";
import type {ProviderSpending} from "@/models/providerSpending.ts";
import {zeroAmount} from "@/models/amount.ts";

const DashboardPage = () => {
    const {data: subscriptionsData, isLoading: isLoadingSubscriptions} = useSubscriptionsQuery();

    const allSubscriptions: Subscription[] = useMemo(() =>
            subscriptionsData?.pages.flatMap(page => page.subscriptions) || [],
        [subscriptionsData]);

    const {
        activeSubscriptions: summaryActiveSubscriptions,
        upcomingRenewals: summaryUpcomingRenewals,
        totalMonthly: summaryMonthly,
        totalYearly: summaryYearly,
        topProviders: summaryTopProviders,
        isLoading: isLoadingSummary,
    } = useSubscriptionSummaryQuery({topProviders: 5, totalMonthly: true, totalYearly: true, upcomingRenewals: 5});

    const providerIds = useMemo(() => {
        const ids = new Set<string>();
        (summaryTopProviders ?? []).forEach(tp => {
            if (tp.providerId) ids.add(tp.providerId);
        });
        (summaryUpcomingRenewals ?? []).forEach(u => {
            if (u.providerId) ids.add(u.providerId);
        });
        return Array.from(ids);
    }, [summaryTopProviders, summaryUpcomingRenewals]);

    const {providerMap, isLoading: isLoadingProvidersByIds} = useProvidersByIds(providerIds);

    const totalMonthly = summaryMonthly;
    const totalYearly = summaryYearly;
    const activeSubscriptionsCount = summaryActiveSubscriptions;

    const topProvidersData = useMemo(() => {
        if (summaryTopProviders && summaryTopProviders.length > 0) {
            return summaryTopProviders.map(tp => ({
                id: tp.providerId ?? "",
                providerName: providerMap.get(tp.providerId ?? "")?.name || (tp.providerId ?? ""),
                amount: tp.total ?? zeroAmount,
            })) as ProviderSpending[];
        }
        return [];
    }, [summaryTopProviders, providerMap]);


    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Dashboard"
            />

            {/* Summary Cards */}
            <SummaryCards
                totalMonthly={totalMonthly}
                totalYearly={totalYearly}
                activeSubscriptionsCount={activeSubscriptionsCount}
                isLoading={isLoadingSubscriptions || isLoadingSummary}
            />

            {/* Side by side: Upcoming Renewals and Top Providers */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                <UpcomingRenewals
                    summaryUpcomingRenewals={summaryUpcomingRenewals}
                    providerMap={providerMap}
                    isLoading={isLoadingSummary}
                />

                <TopProviders
                    providers={topProvidersData}
                    isLoading={isLoadingSubscriptions || isLoadingProvidersByIds || isLoadingSummary}
                />
            </div>

            <div className="mb-8">
                <PriceEvolutionGraph
                    subscriptions={allSubscriptions}
                    providerMap={providerMap}
                    isLoading={isLoadingSubscriptions || isLoadingProvidersByIds}
                />
            </div>
        </div>
    );
};

export default DashboardPage;