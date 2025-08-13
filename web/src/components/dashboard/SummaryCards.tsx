import {Skeleton} from "@/components/ui/skeleton";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Calendar, CreditCard, TrendingUp} from "lucide-react";
import Money from "@/components/ui/money";
import type {Amount} from "@/models/amount.ts";

interface SummaryCardsProps {
    totalMonthly: Amount;
    totalYearly: Amount;
    activeSubscriptionsCount: number;
    isLoading: boolean;
}

const SummaryCards = ({
                          totalMonthly,
                          totalYearly,
                          activeSubscriptionsCount,
                          isLoading,
                      }: SummaryCardsProps) => {
    return (
        <div className="mb-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Monthly Expenses */}
                <Card
                    className="overflow-hidden transition-all duration-300 hover:shadow-lg hover:-translate-y-1 border-l-4 border-l-blue-500">
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-xl font-medium">Monthly Expenses</CardTitle>
                        <CreditCard className="h-5 w-5 text-blue-500 animate-pulse"/>
                    </CardHeader>
                    <CardContent>
                        {isLoading ? (
                            <Skeleton className="h-10 w-28"/>
                        ) : (
                            <div className="flex items-center">
                                <Money
                                    amount={totalMonthly}
                                    className="text-3xl font-bold bg-gradient-to-r from-blue-500 to-blue-700 bg-clip-text text-transparent"
                                />
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Yearly Expenses */}
                <Card
                    className="overflow-hidden transition-all duration-300 hover:shadow-lg hover:-translate-y-1 border-l-4 border-l-purple-500">
                    <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                        <CardTitle className="text-xl font-medium">Yearly Expenses</CardTitle>
                        <TrendingUp className="h-5 w-5 text-purple-500 animate-pulse"/>
                    </CardHeader>
                    <CardContent>
                        {isLoading ? (
                            <Skeleton className="h-10 w-28"/>
                        ) : (
                            <div className="flex items-center">
                                <Money
                                    amount={totalYearly}
                                    className="text-3xl font-bold bg-gradient-to-r from-purple-500 to-purple-700 bg-clip-text text-transparent"
                                />
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Active Subscriptions */}
                <Card
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