import { useMemo } from "react";
import { useAllSubscriptionsQuery } from "@/hooks/subscriptions/useAllSubscriptionsQuery";
import { useAllProvidersQuery } from "@/hooks/providers/useAllProvidersQuery";
import { SubscriptionRecurrency } from "@/models/subscriptionRecurrency";
import { addMonths, addYears, addDays } from "date-fns";
import SummaryCards from "@/components/dashboard/SummaryCards";
import UpcomingRenewals from "@/components/dashboard/UpcomingRenewals";
import TopProviders from "@/components/dashboard/TopProviders";
import SubscriptionsTable from "@/components/dashboard/SubscriptionsTable";

const DashboardPage = () => {
    // Fetch all subscriptions
    const { data: subscriptionsData, isLoading: isLoadingSubscriptions } = useAllSubscriptionsQuery();
    
    // Fetch all providers
    const { data: providersData, isLoading: isLoadingProviders } = useAllProvidersQuery();
    
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
            .filter(sub => sub.recurrency === SubscriptionRecurrency.Monthly && sub.isActive)
            .reduce((sum, sub) => {
                if (sub.customPrice) {
                    return sum + sub.customPrice.amount;
                }
                return sum;
            }, 0);
    }, [allSubscriptions]);
    
    // Calculate yearly expenses
    const totalYearly = useMemo(() => {
        return allSubscriptions
            .filter(sub => sub.recurrency === SubscriptionRecurrency.Yearly && sub.isActive)
            .reduce((sum, sub) => {
                if (sub.customPrice) {
                    return sum + sub.customPrice.amount;
                }
                return sum;
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
                    ...sub,
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
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Dashboard</h1>
            </div>

            {/* Summary Cards */}
            <SummaryCards 
                totalMonthly={totalMonthly}
                totalYearly={totalYearly}
                activeSubscriptionsCount={activeSubscriptionsCount}
                isLoading={isLoadingSubscriptions}
            />

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

            {/* All Subscriptions */}
            <SubscriptionsTable 
                subscriptions={allSubscriptions}
                subscriptionsWithNextRenewal={subscriptionsWithNextRenewal}
                providerMap={providerMap}
                isLoading={isLoadingSubscriptions}
            />
        </div>
    );
};

export default DashboardPage;