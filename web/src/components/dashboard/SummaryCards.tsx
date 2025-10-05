import {Skeleton} from "@/components/ui/skeleton";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Calendar, CreditCard, TrendingUp, ArrowUpRight, ArrowDownRight} from "lucide-react";
import Money from "@/components/ui/money";
import type {Amount} from "@/models/amount.ts";
import { useMemo } from "react";

interface SummaryCardsProps {
    totalMonthly: Amount;
    totalYearly: Amount;
    totalLastMonth: Amount;
    totalLastYear: Amount;
    activeSubscriptionsCount: number;
    isLoading: boolean;
}

const SummaryCards = ({
                          totalMonthly,
                          totalYearly,
                          totalLastMonth,
                          totalLastYear,
                          activeSubscriptionsCount,
                          isLoading,
                      }: SummaryCardsProps) => {
    const computeChange = (current: Amount, previous: Amount): number | null => {
        const prevVal = previous?.value;
        const currVal = current?.value;
        if (prevVal == null || currVal == null || prevVal === 0) return null;
        const delta = currVal - prevVal;
        return (delta / prevVal) * 100;
    };

    const monthlyChange = useMemo(() => computeChange(totalMonthly, totalLastMonth), [totalMonthly, totalLastMonth]);
    const yearlyChange = useMemo(() => computeChange(totalYearly, totalLastYear), [totalYearly, totalLastYear]);

    return (
        <div className="mb-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6" data-testid="summary-cards">
                {/* Monthly Expenses */}
                <Card data-testid="summary-card-monthly"
                    className="overflow-hidden transition-all duration-300 hover:shadow-lg hover:-translate-y-1 border-l-4 border-l-blue-500">
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-xl font-medium">Monthly Expenses</CardTitle>
                        <CreditCard className="h-5 w-5 text-blue-500 animate-pulse"/>
                    </CardHeader>
                    <CardContent>
                        {isLoading ? (
                            <Skeleton className="h-10 w-28"/>
                        ) : (
                            <div className="flex items-center gap-3">
                                <Money
                                    amount={totalMonthly}
                                    className="text-3xl font-bold bg-gradient-to-r from-blue-500 to-blue-700 bg-clip-text text-transparent"
                                />
                                {monthlyChange === null ? (
                                    <span className="ml-1 text-sm text-muted-foreground">—</span>
                                ) : (
                                    <span
                                        className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${monthlyChange > 0
                                            ? "bg-red-100 text-red-700 dark:bg-red-500/10 dark:text-red-400"
                                            : monthlyChange < 0
                                                ? "bg-green-100 text-green-700 dark:bg-green-500/10 dark:text-green-400"
                                                : "bg-muted text-muted-foreground"}
                                        `}
                                    >
                                        {monthlyChange > 0 ? (
                                            <ArrowUpRight className="h-3 w-3 mr-1" />
                                        ) : monthlyChange < 0 ? (
                                            <ArrowDownRight className="h-3 w-3 mr-1" />
                                        ) : null}
                                        {Math.abs(monthlyChange).toFixed(1)}%
                                    </span>
                                )}
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Yearly Expenses */}
                <Card data-testid="summary-card-yearly"
                    className="overflow-hidden transition-all duration-300 hover:shadow-lg hover:-translate-y-1 border-l-4 border-l-purple-500">
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-xl font-medium">Yearly Expenses</CardTitle>
                        <TrendingUp className="h-5 w-5 text-purple-500 animate-pulse"/>
                    </CardHeader>
                    <CardContent>
                        {isLoading ? (
                            <Skeleton className="h-10 w-28"/>
                        ) : (
                            <div className="flex items-center gap-3">
                                <Money
                                    amount={totalYearly}
                                    className="text-3xl font-bold bg-gradient-to-r from-purple-500 to-purple-700 bg-clip-text text-transparent"
                                />
                                {yearlyChange === null ? (
                                    <span className="ml-1 text-sm text-muted-foreground">—</span>
                                ) : (
                                    <span
                                        className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${yearlyChange > 0
                                            ? "bg-red-100 text-red-700 dark:bg-red-500/10 dark:text-red-400"
                                            : yearlyChange < 0
                                                ? "bg-green-100 text-green-700 dark:bg-green-500/10 dark:text-green-400"
                                                : "bg-muted text-muted-foreground"}
                                        `}
                                    >
                                        {yearlyChange > 0 ? (
                                            <ArrowUpRight className="h-3 w-3 mr-1" />
                                        ) : yearlyChange < 0 ? (
                                            <ArrowDownRight className="h-3 w-3 mr-1" />
                                        ) : null}
                                        {Math.abs(yearlyChange).toFixed(1)}%
                                    </span>
                                )}
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Active Subscriptions */}
                <Card data-testid="summary-card-active"
                    className="overflow-hidden transition-all duration-300 hover:shadow-lg hover:-translate-y-1 border-l-4 border-l-green-500">
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-xl font-medium">Active Subscriptions</CardTitle>
                        <Calendar className="h-5 w-5 text-green-500 animate-pulse"/>
                    </CardHeader>
                    <CardContent>
                        {isLoading ? (
                            <Skeleton className="h-10 w-28"/>
                        ) : (
                            <div className="flex items-center">
                                <p className="text-3xl font-bold bg-gradient-to-r from-green-500 to-green-700 bg-clip-text text-transparent">
                                    {activeSubscriptionsCount}
                                </p>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};

export default SummaryCards;