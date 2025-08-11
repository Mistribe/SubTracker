import {useMemo} from "react";
import {useAllSubscriptionsQuery} from "@/hooks/subscriptions/useAllSubscriptionsQuery";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import SummaryCards from "@/components/dashboard/SummaryCards";
import UpcomingRenewals from "@/components/dashboard/UpcomingRenewals";
import TopProviders from "@/components/dashboard/TopProviders";
import PriceEvolutionGraph from "@/components/dashboard/PriceEvolutionGraph";
import {PageHeader} from "@/components/ui/page-header";
import {usePreferredCurrency} from "@/hooks/currencies/usePreferredCurrency";
import {useCurrencyRates} from "@/hooks/currencies/useCurrencyRates";
import {convertAmount, subscriptionMonthlyPriceInCurrency, subscriptionYearlyPriceInCurrency} from "@/utils/currency";
import type Provider from "@/models/provider.ts";
import type {SubscriptionWithNextRenewal} from "@/models/subscriptionWithNextRenewal.ts";
import type {ProviderSpending} from "@/models/providerSpending.ts";
import type Subscription from "@/models/subscription.ts";

const DashboardPage = () => {
    const {data: subscriptionsData, isLoading: isLoadingSubscriptions} = useAllSubscriptionsQuery();
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

    const totalMonthly = useMemo(() => {
        return allSubscriptions
            .filter(sub => sub.isActive)
            .reduce((sum, sub) => {
                return sum + subscriptionMonthlyPriceInCurrency(sub, preferredCurrency, rates);
            }, 0);
    }, [allSubscriptions, preferredCurrency, rates]);

    const totalYearly = useMemo(() => {
        return allSubscriptions
            .filter(sub => sub.isActive)
            .reduce((sum, sub) => {
                return sum + subscriptionYearlyPriceInCurrency(sub, preferredCurrency, rates);
            }, 0);
    }, [allSubscriptions, preferredCurrency, rates]);

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
                isLoading={isLoadingSubscriptions || isLoadingRates}
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
                    providers={providerSpending}
                    isLoading={isLoadingSubscriptions || isLoadingProviders}
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