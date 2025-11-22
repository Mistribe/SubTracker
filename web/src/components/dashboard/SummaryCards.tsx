import {Skeleton} from "@/components/ui/skeleton";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Calendar, CreditCard, TrendingUp, ArrowUpRight, ArrowDownRight, User, Users} from "lucide-react";
import Money from "@/components/ui/money";
import type {Amount} from "@/models/amount.ts";
import React, { useMemo } from "react";

interface SummaryCardsProps {
    totalMonthly: Amount;
    totalYearly: Amount;
    totalLastMonth: Amount;
    totalLastYear: Amount;
    personalMonthly?: Amount;
    personalYearly?: Amount;
    personalLastMonth?: Amount;
    personalLastYear?: Amount;
    familyMonthly?: Amount;
    familyYearly?: Amount;
    familyLastMonth?: Amount;
    familyLastYear?: Amount;
    activeSubscriptionsCount: number;
    activePersonal: number;
    activeFamily: number;
    isLoading: boolean;
}

const SummaryCards = ({
                          totalMonthly,
                          totalYearly,
                          totalLastMonth,
                          totalLastYear,
                          personalMonthly,
                          personalYearly,
                          personalLastMonth,
                          personalLastYear,
                          familyMonthly,
                          familyYearly,
                          familyLastMonth,
                          familyLastYear,
                          activeSubscriptionsCount,
                          activePersonal,
                          activeFamily,
                          isLoading,
                      }: SummaryCardsProps) => {
    const computeChange = (current: Amount | undefined, previous: Amount | undefined): number | null => {
        const prevVal = previous?.value;
        const currVal = current?.value;
        if (prevVal == null || currVal == null || prevVal === 0) return null;
        const delta = currVal - prevVal;
        return (delta / prevVal) * 100;
    };

    const monthlyChange = useMemo(() => computeChange(totalMonthly, totalLastMonth), [totalMonthly, totalLastMonth]);
    const yearlyChange = useMemo(() => computeChange(totalYearly, totalLastYear), [totalYearly, totalLastYear]);
    const personalMonthlyChange = useMemo(() => computeChange(personalMonthly, personalLastMonth), [personalMonthly, personalLastMonth]);
    const personalYearlyChange = useMemo(() => computeChange(personalYearly, personalLastYear), [personalYearly, personalLastYear]);
    const familyMonthlyChange = useMemo(() => computeChange(familyMonthly, familyLastMonth), [familyMonthly, familyLastMonth]);
    const familyYearlyChange = useMemo(() => computeChange(familyYearly, familyLastYear), [familyYearly, familyLastYear]);

    const renderTrendIndicator = (change: number | null) => {
        if (change === null) {
            return <span className="ml-1 text-sm text-muted-foreground">—</span>;
        }
        return (
            <span
                className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                    change > 0
                        ? "bg-red-100 text-red-700 dark:bg-red-500/10 dark:text-red-400"
                        : change < 0
                            ? "bg-green-100 text-green-700 dark:bg-green-500/10 dark:text-green-400"
                            : "bg-muted text-muted-foreground"
                }`}
            >
                {change > 0 ? (
                    <ArrowUpRight className="h-3 w-3 mr-1" />
                ) : change < 0 ? (
                    <ArrowDownRight className="h-3 w-3 mr-1" />
                ) : null}
                {Math.abs(change).toFixed(1)}%
            </span>
        );
    };

    const renderBreakdownRow = (
        label: string,
        amount: Amount | undefined,
        change: number | null,
        icon: React.ReactNode,
        testId: string
    ) => {
        if (!amount) return null;
        return (
            <div className="flex items-center justify-between py-2 border-b last:border-b-0" data-testid={testId}>
                <div className="flex items-center gap-2">
                    {icon}
                    <span className="text-sm font-medium text-muted-foreground">{label}</span>
                </div>
                <div className="flex items-center gap-2">
                    <Money amount={amount} className="text-lg font-semibold" />
                    {renderTrendIndicator(change)}
                </div>
            </div>
        );
    };

    const renderCountBreakdownRow = (
        label: string,
        count: number,
        icon: React.ReactNode,
        testId: string
    ) => {
        return (
            <div className="flex items-center justify-between py-2 border-b last:border-b-0" data-testid={testId}>
                <div className="flex items-center gap-2">
                    {icon}
                    <span className="text-sm font-medium text-muted-foreground">{label}</span>
                </div>
                <span className="text-lg font-semibold">{count}</span>
            </div>
        );
    };

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
                            <div className="space-y-2">
                                <Skeleton className="h-10 w-full"/>
                                <Skeleton className="h-10 w-full"/>
                                <Skeleton className="h-10 w-full"/>
                            </div>
                        ) : (
                            <div className="space-y-1">
                                {renderBreakdownRow(
                                    "Personal",
                                    personalMonthly,
                                    personalMonthlyChange,
                                    <User className="h-4 w-4 text-blue-500" />,
                                    "monthly-personal"
                                )}
                                {renderBreakdownRow(
                                    "Family",
                                    familyMonthly,
                                    familyMonthlyChange,
                                    <Users className="h-4 w-4 text-purple-500" />,
                                    "monthly-family"
                                )}
                                {renderBreakdownRow(
                                    "Total",
                                    totalMonthly,
                                    monthlyChange,
                                    <CreditCard className="h-4 w-4 text-blue-700" />,
                                    "monthly-total"
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
                            <div className="space-y-2">
                                <Skeleton className="h-10 w-full"/>
                                <Skeleton className="h-10 w-full"/>
                                <Skeleton className="h-10 w-full"/>
                            </div>
                        ) : (
                            <div className="space-y-1">
                                {renderBreakdownRow(
                                    "Personal",
                                    personalYearly,
                                    personalYearlyChange,
                                    <User className="h-4 w-4 text-blue-500" />,
                                    "yearly-personal"
                                )}
                                {renderBreakdownRow(
                                    "Family",
                                    familyYearly,
                                    familyYearlyChange,
                                    <Users className="h-4 w-4 text-purple-500" />,
                                    "yearly-family"
                                )}
                                {renderBreakdownRow(
                                    "Total",
                                    totalYearly,
                                    yearlyChange,
                                    <TrendingUp className="h-4 w-4 text-purple-700" />,
                                    "yearly-total"
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
                            <div className="space-y-2">
                                <Skeleton className="h-10 w-full"/>
                                <Skeleton className="h-10 w-full"/>
                                <Skeleton className="h-10 w-full"/>
                            </div>
                        ) : (
                            <div className="space-y-1">
                                {renderCountBreakdownRow(
                                    "Personal",
                                    activePersonal,
                                    <User className="h-4 w-4 text-blue-500" />,
                                    "active-personal"
                                )}
                                {renderCountBreakdownRow(
                                    "Family",
                                    activeFamily,
                                    <Users className="h-4 w-4 text-purple-500" />,
                                    "active-family"
                                )}
                                {renderCountBreakdownRow(
                                    "Total",
                                    activeSubscriptionsCount,
                                    <Calendar className="h-4 w-4 text-green-700" />,
                                    "active-total"
                                )}
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};

export default SummaryCards;