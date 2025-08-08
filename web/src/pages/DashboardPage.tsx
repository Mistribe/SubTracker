import {useMemo} from "react";
import {useAllSubscriptionsQuery} from "@/hooks/subscriptions/useAllSubscriptionsQuery";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";
import {addDays, addMonths, addYears} from "date-fns";
import SummaryCards from "@/components/dashboard/SummaryCards";
import UpcomingRenewals from "@/components/dashboard/UpcomingRenewals";
import TopProviders from "@/components/dashboard/TopProviders";
import PriceEvolutionGraph from "@/components/dashboard/PriceEvolutionGraph";
import {PageHeader} from "@/components/ui/page-header";

const DashboardPage = () => {
    // Fetch all subscriptions
    const {data: subscriptionsData, isLoading: isLoadingSubscriptions} = useAllSubscriptionsQuery();

    // Fetch all providers
    const {data: providersData, isLoading: isLoadingProviders} = useAllProvidersQuery();

    // Flatten all subscriptions from all pages
    const allSubscriptions = useMemo(() =>
            subscriptionsData?.pages.flatMap(page => page.subscriptions) || [],
        [subscriptionsData]
    );

    // Flatten all providers from all pages and create a mapping from ID to provider
    const providerMap = useMemo(() => {
        const map = new Map();
        providersData?.pages.flatMap(page => page.providers || []).forEach(provider => {
            map.set(provider.id, provider);
        });
        return map;
    }, [providersData]);

    // Calculate monthly expenses
    const totalMonthly = useMemo(() => {
        return allSubscriptions
            .filter(sub => sub.isActive)
            .reduce((sum, sub) => {
                return sum + sub.getMonthlyPrice()
            }, 0);
    }, [allSubscriptions]);

    // Calculate yearly expenses
    const totalYearly = useMemo(() => {
        return allSubscriptions
            .filter(sub => sub.isActive)
            .reduce((sum, sub) => {
                return sum + sub.getYearlyPrice()
            }, 0);
    }, [allSubscriptions]);

    // Count active subscriptions
    const activeSubscriptionsCount = useMemo(() => {
        return allSubscriptions.filter(sub => sub.isActive).length;
    }, [allSubscriptions]);

    // Calculate next renewal date for each subscription
    const subscriptionsWithNextRenewal = useMemo(() => {
        return allSubscriptions
            .filter(sub => sub.isActive)
            .map(sub => {
                let nextRenewalDate;
                const today = new Date();

                // Calculate next renewal based on recurrency
                switch (sub.recurrency) {
                    case SubscriptionRecurrency.Monthly:
                        // Start from subscription start date and add months until we get a future date
                        nextRenewalDate = new Date(sub.startDate);
                        while (nextRenewalDate < today) {
                            nextRenewalDate = addMonths(nextRenewalDate, 1);
                        }
                        break;
                    case SubscriptionRecurrency.Quarterly:
                        nextRenewalDate = new Date(sub.startDate);
                        while (nextRenewalDate < today) {
                            nextRenewalDate = addMonths(nextRenewalDate, 3);
                        }
                        break;
                    case SubscriptionRecurrency.HalfYearly:
                        nextRenewalDate = new Date(sub.startDate);
                        while (nextRenewalDate < today) {
                            nextRenewalDate = addMonths(nextRenewalDate, 6);
                        }
                        break;
                    case SubscriptionRecurrency.Yearly:
                        nextRenewalDate = new Date(sub.startDate);
                        while (nextRenewalDate < today) {
                            nextRenewalDate = addYears(nextRenewalDate, 1);
                        }
                        break;
                    case SubscriptionRecurrency.Custom:
                        if (sub.customRecurrency) {
                            nextRenewalDate = new Date(sub.startDate);
                            while (nextRenewalDate < today) {
                                nextRenewalDate = addDays(nextRenewalDate, sub.customRecurrency);
                            }
                        } else {
                            nextRenewalDate = sub.endDate || today;
                        }
                        break;
                    default:
                        nextRenewalDate = sub.endDate || today;
                }

                return {
                    subscription: sub,
                    nextRenewalDate
                };
            })
            .sort((a, b) => a.nextRenewalDate.getTime() - b.nextRenewalDate.getTime());
    }, [allSubscriptions]);

    // Get top 5 upcoming renewals
    const topUpcomingRenewals = useMemo(() => {
        return subscriptionsWithNextRenewal.slice(0, 5);
    }, [subscriptionsWithNextRenewal]);

    // Calculate spending by provider
    const providerSpending = useMemo(() => {
        const spending = new Map();

        allSubscriptions.forEach(sub => {
            if (!sub.customPrice || !sub.isActive) return;

            const providerId = sub.providerId;
            const providerName = providerMap.get(providerId)?.name || providerId;
            const amount = sub.customPrice.amount;

            // Convert all subscriptions to yearly cost for fair comparison
            let yearlyAmount = amount;
            switch (sub.recurrency) {
                case SubscriptionRecurrency.Monthly:
                    yearlyAmount = amount * 12;
                    break;
                case SubscriptionRecurrency.Quarterly:
                    yearlyAmount = amount * 4;
                    break;
                case SubscriptionRecurrency.HalfYearly:
                    yearlyAmount = amount * 2;
                    break;
                case SubscriptionRecurrency.Custom:
                    if (sub.customRecurrency) {
                        yearlyAmount = amount * (365 / sub.customRecurrency);
                    }
                    break;
            }

            if (spending.has(providerId)) {
                spending.set(providerId, {
                    id: providerId,
                    name: providerName,
                    amount: spending.get(providerId).amount + yearlyAmount
                });
            } else {
                spending.set(providerId, {
                    id: providerId,
                    name: providerName,
                    amount: yearlyAmount
                });
            }
        });

        return Array.from(spending.values())
            .sort((a, b) => b.amount - a.amount)
            .slice(0, 5);
    }, [allSubscriptions, providerMap]);


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
                isLoading={isLoadingSubscriptions}
            />

            {/* Side by side: Upcoming Renewals and Top Providers */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                {/* Top 5 Upcoming Renewals */}
                <UpcomingRenewals
                    upcomingRenewals={topUpcomingRenewals}
                    providerMap={providerMap}
                    isLoading={isLoadingSubscriptions}
                />

                {/* Top 5 Providers by Expense */}
                <TopProviders
                    providers={providerSpending}
                    isLoading={isLoadingSubscriptions || isLoadingProviders}
                />
            </div>

            {/* Monthly Price Evolution Graph */}
            <div className="mb-8">
                <PriceEvolutionGraph
                    subscriptions={allSubscriptions}
                    isLoading={isLoadingSubscriptions || isLoadingProviders}
                />
            </div>
        </div>
    );
};

export default DashboardPage;