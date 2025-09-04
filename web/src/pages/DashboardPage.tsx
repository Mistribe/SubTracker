import {useMemo} from "react";
import {useProvidersByIds} from "@/hooks/providers/useProvidersByIds";
import SummaryCards from "@/components/dashboard/SummaryCards";
import UpcomingRenewals from "@/components/dashboard/UpcomingRenewals";
import TopProviders from "@/components/dashboard/TopProviders";
import TopLabels from "@/components/dashboard/TopLabels";
import EmptySubscriptionsState from "@/components/dashboard/EmptySubscriptionsState";
import {PageHeader} from "@/components/ui/page-header";
import {useSubscriptionSummaryQuery} from "@/hooks/subscriptions/useSubscriptionSummaryQuery";

const DashboardPage = () => {
    const {
        activeSubscriptions: summaryActiveSubscriptions,
        upcomingRenewals: summaryUpcomingRenewals,
        totalMonthly: summaryMonthly,
        totalYearly: summaryYearly,
        totalLastMonth: summaryLastMonth,
        totalLastYear: summaryLastYear,
        topProviders: summaryTopProviders,
        topLabels: summaryTopLabels,
        isLoading: isLoadingSummary,
    } = useSubscriptionSummaryQuery({topProviders: 5, topLabels: 5, totalMonthly: true, totalYearly: true, upcomingRenewals: 5});

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

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Dashboard"
            />

            { !isLoadingSummary && (summaryActiveSubscriptions ?? 0) === 0 ? (
                <EmptySubscriptionsState />
            ) : (
                <>
                    <SummaryCards
                        totalMonthly={summaryMonthly}
                        totalYearly={summaryYearly}
                        totalLastMonth={summaryLastMonth}
                        totalLastYear={summaryLastYear}
                        activeSubscriptionsCount={summaryActiveSubscriptions}
                        isLoading={ isLoadingSummary}
                    />

                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
                        <UpcomingRenewals
                            summaryUpcomingRenewals={summaryUpcomingRenewals}
                            providerMap={providerMap}
                            isLoading={isLoadingSummary}
                        />

                        <TopProviders
                            providers={summaryTopProviders}
                            providerMap={providerMap}
                            isLoading={isLoadingProvidersByIds || isLoadingSummary}
                        />

                        <TopLabels
                            labels={summaryTopLabels ?? []}
                            isLoading={isLoadingSummary}
                        />
                    </div>
                </>
            )}
        </div>
    );
};

export default DashboardPage;