import { Skeleton } from "@/components/ui/skeleton";
import { formatCurrency } from "./utils";

interface SummaryCardsProps {
    totalMonthly: number;
    totalYearly: number;
    activeSubscriptionsCount: number;
    isLoading: boolean;
}

const SummaryCards = ({
    totalMonthly,
    totalYearly,
    activeSubscriptionsCount,
    isLoading
}: SummaryCardsProps) => {
    return (
        <div className="mb-8">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                {/* Monthly Expenses */}
                <div className="p-6 border rounded-lg bg-card">
                    <h3 className="text-xl font-medium mb-2">Monthly Expenses</h3>
                    {isLoading ? (
                        <Skeleton className="h-8 w-24" />
                    ) : (
                        <p className="text-3xl font-bold">{formatCurrency(totalMonthly)}</p>
                    )}
                </div>
                
                {/* Yearly Expenses */}
                <div className="p-6 border rounded-lg bg-card">
                    <h3 className="text-xl font-medium mb-2">Yearly Expenses</h3>
                    {isLoading ? (
                        <Skeleton className="h-8 w-24" />
                    ) : (
                        <p className="text-3xl font-bold">{formatCurrency(totalYearly)}</p>
                    )}
                </div>
                
                {/* Total Annual Cost */}
                <div className="p-6 border rounded-lg bg-card">
                    <h3 className="text-xl font-medium mb-2">Total Annual Cost</h3>
                    {isLoading ? (
                        <Skeleton className="h-8 w-24" />
                    ) : (
                        <p className="text-3xl font-bold">{formatCurrency(totalMonthly * 12 + totalYearly)}</p>
                    )}
                </div>
                
                {/* Active Subscriptions */}
                <div className="p-6 border rounded-lg bg-card">
                    <h3 className="text-xl font-medium mb-2">Active Subscriptions</h3>
                    {isLoading ? (
                        <Skeleton className="h-8 w-24" />
                    ) : (
                        <p className="text-3xl font-bold">{activeSubscriptionsCount}</p>
                    )}
                </div>
            </div>
        </div>
    );
};

export default SummaryCards;