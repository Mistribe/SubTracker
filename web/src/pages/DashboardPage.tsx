import {useMemo} from "react";
import {useSubscriptionsQuery} from "@/hooks/subscriptions/useSubscriptionsQuery.ts";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import SummaryCards from "@/components/dashboard/SummaryCards";
import UpcomingRenewals from "@/components/dashboard/UpcomingRenewals";
import TopProviders from "@/components/dashboard/TopProviders";
import PriceEvolutionGraph from "@/components/dashboard/PriceEvolutionGraph";
import {PageHeader} from "@/components/ui/page-header";
import {usePreferredCurrency} from "@/hooks/currencies/usePreferredCurrency";
import {useCurrencyRates} from "@/hooks/currencies/useCurrencyRates";
import {convertAmount} from "@/utils/currency";
import {useSubscriptionSummaryQuery} from "@/hooks/subscriptions/useSubscriptionSummaryQuery";
import type Provider from "@/models/provider.ts";
import type {SubscriptionWithNextRenewal} from "@/models/subscriptionWithNextRenewal.ts";
import type {ProviderSpending} from "@/models/providerSpending.ts";
import type Subscription from "@/models/subscription.ts";

const DashboardPage = () => {
    const {data: subscriptionsData, isLoading: isLoadingSubscriptions} = useSubscriptionsQuery();
    const {data: providersData, isLoading: isLoadingProviders} = useAllProvidersQuery();

    const allSubscriptions: Subscription[] = useMemo(() =>
            subscriptionsData?.pages.flatMap(page => page.subscriptions) || [],
        [subscriptionsData]);

    const providerMap = useMemo(() => {
        const map = new Map<string, Provider>();
        providersData?.pages.flatMap(page => page.providers || []).forEach(provider => {
            map.set(provider.id, provider);
        });
        return map;
    }, [providersData]);

    const {preferredCurrency} = usePreferredCurrency();
    const {rates, isLoading: isLoadingRates} = useCurrencyRates();
    const {
        totalMonthly: summaryMonthly,
        totalYearly: summaryYearly,
        topProviders: summaryTopProviders,
        isLoading: isLoadingSummary,
    } = useSubscriptionSummaryQuery({topProviders: 5, totalMonthly: true, totalYearly: true, upcomingRenewals: 5});
    
    const totalMonthly = summaryMonthly;
    const totalYearly = summaryYearly;

    const activeSubscriptionsCount = useMemo(() => {
        return allSubscriptions.filter(sub => sub.isActive).length;
    }, [allSubscriptions]);

    const totalsCurrency = preferredCurrency;

    const subscriptionsWithNextRenewal = useMemo(() => {
        return allSubscriptions
            .filter((sub: Subscription) => sub.isActive)
            .map((sub: Subscription) => {
                return {
                    subscription: sub,
                    nextRenewalDate: sub.getNextRenewalDate()
                } as SubscriptionWithNextRenewal;
            })
            .filter((sub: SubscriptionWithNextRenewal) => sub.nextRenewalDate)
            .sort((a, b) => a.nextRenewalDate!.getTime() - b.nextRenewalDate!.getTime());
    }, [allSubscriptions]);

    const topUpcomingRenewals = useMemo(() => {
        return subscriptionsWithNextRenewal.slice(0, 5);
    }, [subscriptionsWithNextRenewal]);

    // Calculate spending by provider (convert to preferred currency yearly)
    const providerSpending = useMemo(() => {
        const spending = new Map<string, ProviderSpending>();

        allSubscriptions.forEach((sub: Subscription) => {
            const providerId = sub.providerId;
            const providerName = providerMap.get(providerId)?.name || providerId;
            const subscriptionAmount = convertAmount(sub.getTotalAmount(),
                sub.getCurrency(),
                preferredCurrency,
                rates);
            if (!subscriptionAmount) {
                console.log(`Error converting amount for subscription ${sub.id}(${providerName}): ${sub.getAmount()} ${sub.getCurrency()} -> ${preferredCurrency}`);
                return;
            }

            if (spending.has(providerId)) {
                spending.set(providerId, {
                    id: providerId,
                    name: providerName,
                    amount: spending.get(providerId)!.amount + subscriptionAmount,
                    currency: preferredCurrency
                });
            } else {
                spending.set(providerId, {
                    id: providerId,
                    name: providerName,
                    amount: subscriptionAmount,
                    currency: preferredCurrency
                });
            }
        });

        return Array.from(spending.values())
            .sort((a, b) => b.amount - a.amount)
            .slice(0, 5);
    }, [allSubscriptions, providerMap, preferredCurrency, rates]);

    const topProvidersData = useMemo(() => {
        if (summaryTopProviders && summaryTopProviders.length > 0) {
            return summaryTopProviders.map(tp => ({
                id: tp.providerId ?? "",
                name: providerMap.get(tp.providerId ?? "")?.name || (tp.providerId ?? ""),
                amount: tp.total ?? 0,
                currency: preferredCurrency,
            }));
        }
        return providerSpending;
    }, [summaryTopProviders, providerMap, preferredCurrency, providerSpending]);


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
                isLoading={isLoadingSubscriptions || isLoadingRates || isLoadingSummary}
                totalsCurrency={totalsCurrency}
            />

            {/* Side by side: Upcoming Renewals and Top Providers */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                <UpcomingRenewals
                    upcomingRenewals={topUpcomingRenewals}
                    providerMap={providerMap}
                    isLoading={isLoadingSubscriptions}
                />

                <TopProviders
                    providers={topProvidersData}
                    isLoading={isLoadingSubscriptions || isLoadingProviders || isLoadingSummary}
                />
            </div>

            <div className="mb-8">
                <PriceEvolutionGraph
                    subscriptions={allSubscriptions}
                    providerMap={providerMap}
                    isLoading={isLoadingSubscriptions || isLoadingProviders}
                />
            </div>
        </div>
    );
};

export default DashboardPage;