import {useMemo} from "react";
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
        activePersonal: summaryActivePersonal,
        activeFamily: summaryActiveFamily,
        upcomingRenewals: summaryUpcomingRenewals,
        totalMonthly: summaryMonthly,
        totalYearly: summaryYearly,
        totalLastMonth: summaryLastMonth,
        totalLastYear: summaryLastYear,
        personalMonthly: summaryPersonalMonthly,
        personalYearly: summaryPersonalYearly,
        personalLastMonth: summaryPersonalLastMonth,
        personalLastYear: summaryPersonalLastYear,
        familyMonthly: summaryFamilyMonthly,
        familyYearly: summaryFamilyYearly,
        familyLastMonth: summaryFamilyLastMonth,
        familyLastYear: summaryFamilyLastYear,
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


    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Dashboard"
            />

            <SummaryCards
                totalMonthly={summaryMonthly}
                totalYearly={summaryYearly}
                totalLastMonth={summaryLastMonth}
                totalLastYear={summaryLastYear}
                personalMonthly={summaryPersonalMonthly}
                personalYearly={summaryPersonalYearly}
                personalLastMonth={summaryPersonalLastMonth}
                personalLastYear={summaryPersonalLastYear}
                familyMonthly={summaryFamilyMonthly}
                familyYearly={summaryFamilyYearly}
                familyLastMonth={summaryFamilyLastMonth}
                familyLastYear={summaryFamilyLastYear}
                activeSubscriptionsCount={summaryActiveSubscriptions}
                activePersonal={summaryActivePersonal}
                activeFamily={summaryActiveFamily}
                isLoading={ isLoadingSummary}
            />

            { !isLoadingSummary && (summaryActiveSubscriptions ?? 0) === 0 ? (
                <EmptySubscriptionsState />
            ) : (
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
                    <UpcomingRenewals
                        summaryUpcomingRenewals={summaryUpcomingRenewals}
                        isLoading={isLoadingSummary}
                    />

                    <TopProviders
                        providers={summaryTopProviders}
                        isLoading={isLoadingSummary}
                    />

                    <TopLabels
                        labels={summaryTopLabels ?? []}
                        isLoading={isLoadingSummary}
                    />
                </div>
            )}
        </div>
    );
};

export default DashboardPage;