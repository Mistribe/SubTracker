import {useMemo, useState} from "react";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {
    ChartContainer,
    ChartLegend,
    ChartLegendContent,
    ChartTooltip,
    ChartTooltipContent
} from "@/components/ui/chart";
import {CartesianGrid, Line, LineChart, XAxis, YAxis} from "recharts";
import {formatCurrency} from "./utils";
import {ChevronLeft, ChevronRight, RefreshCw, TrendingUp} from "lucide-react";
import type Subscription from "@/models/subscription";
import {addMonths, format, startOfMonth, subMonths} from "date-fns";

interface PriceEvolutionGraphProps {
    subscriptions: Subscription[];
    isLoading: boolean;
}

interface MonthlyData {
    month: string;
    date: Date;
    monthlySubscriptionsTotal: number;
}

const PriceEvolutionGraph = ({
                                 subscriptions,
                                 isLoading
                             }: PriceEvolutionGraphProps) => {
    // State to track the reference date for time navigation
    const [referenceDate, setReferenceDate] = useState<Date>(() => new Date());

    // Handler to move backward in time (3 months at a time)
    const handleMoveBackward = () => {
        setReferenceDate(prevDate => subMonths(prevDate, 3));
    };

    // Handler to move forward in time (3 months at a time)
    const handleMoveForward = () => {
        setReferenceDate(prevDate => addMonths(prevDate, 3));
    };

    // Handler to reset to current date
    const handleResetToToday = () => {
        setReferenceDate(new Date());
    };

    // Generate monthly data for the past 6 months and next 6 months from the reference date
    const monthlyData = useMemo(() => {
        const data: MonthlyData[] = [];
        const currentMonth = startOfMonth(referenceDate);

        // Create data for the past 6 months and next 6 months (total 13 months including current)
        for (let i = -6; i <= 6; i++) {
            const monthDate = addMonths(currentMonth, i);
            const monthStr = format(monthDate, 'MMM yyyy');
            
            // Calculate total for subscriptions that are active in this specific month
            // This shows how subscription costs change over time based on start/end dates
            const monthlySubscriptionsTotal = subscriptions
                .filter(subscription => {
                    const startDate = subscription.startDate;
                    const endDate = subscription.endDate;
                    
                    // Check if subscription is active for this month:
                    // 1. Start date is before or equal to the end of this month
                    // 2. End date is undefined or after the start of this month
                    const monthEnd = addMonths(monthDate, 1);
                    const isActiveInMonth = 
                        startDate <= monthEnd && 
                        (!endDate || endDate >= monthDate);
                        
                    return isActiveInMonth;
                })
                .reduce((sum, subscription) => {
                    return sum + subscription.getMonthlyPrice();
                }, 0);

            data.push({
                month: monthStr,
                date: monthDate,
                monthlySubscriptionsTotal
            });
        }

        return data;
    }, [subscriptions, referenceDate]);

    // Calculate the visible date range for display
    const dateRangeStart = format(monthlyData[0].date, 'MMM yyyy');
    const dateRangeEnd = format(monthlyData[monthlyData.length - 1].date, 'MMM yyyy');

    return (
        <Card className="overflow-hidden transition-all duration-300 hover:shadow-lg border-t-4 border-t-green-500">
            <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                <CardTitle className="text-xl font-medium">Monthly Subscription Costs</CardTitle>
                <TrendingUp className="h-5 w-5 text-green-500"/>
            </CardHeader>
            <CardContent>
                {isLoading ? (
                    <div className="h-96 w-full flex items-center justify-center">
                        <div className="animate-pulse text-muted-foreground">Loading chart data...</div>
                    </div>
                ) : (
                    <div className="w-full">
                        {/* Time Navigation Controls */}
                        <div className="flex items-center justify-between mb-4">
                            <div className="text-sm text-muted-foreground">
                                Showing data from <span className="font-medium">{dateRangeStart}</span> to <span
                                className="font-medium">{dateRangeEnd}</span>
                            </div>
                            <div className="flex items-center gap-2">
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={handleMoveBackward}
                                    className="h-8 px-2"
                                >
                                    <ChevronLeft className="h-4 w-4 mr-1"/>
                                    <span>Back</span>
                                </Button>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={handleResetToToday}
                                    className="h-8 px-2"
                                >
                                    <RefreshCw className="h-4 w-4 mr-1"/>
                                    <span>Today</span>
                                </Button>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={handleMoveForward}
                                    className="h-8 px-2"
                                >
                                    <span>Forward</span>
                                    <ChevronRight className="h-4 w-4 ml-1"/>
                                </Button>
                            </div>
                        </div>

                        <ChartContainer
                            config={{
                                monthlySubscriptions: {
                                    label: "Monthly Subscriptions",
                                    theme: {
                                        light: "#0891b2", // cyan-600
                                        dark: "#06b6d4"   // cyan-500
                                    }
                                }
                            }}
                        >
                            <LineChart
                                data={monthlyData}
                                margin={{top: 10, right: 30, left: 20, bottom: 30}}
                            >
                                <CartesianGrid strokeDasharray="3 3"/>
                                <XAxis
                                    dataKey="month"
                                    tick={{fontSize: 12}}
                                />
                                <YAxis
                                    tick={{fontSize: 12}}
                                    tickFormatter={(value) => formatCurrency(value).replace('.00', '')}
                                />
                                <ChartTooltip
                                    content={
                                        <ChartTooltipContent
                                            formatter={(value) => (
                                                <span>{formatCurrency(value as number)}</span>
                                            )}
                                        />
                                    }
                                />
                                <Line
                                    type="monotone"
                                    dataKey="monthlySubscriptionsTotal"
                                    name="monthlySubscriptions"
                                    strokeWidth={2}
                                    dot={{strokeWidth: 2}}
                                    activeDot={{r: 6}}
                                />
                                <ChartLegend
                                    content={<ChartLegendContent/>}
                                    verticalAlign="bottom"
                                    height={36}
                                    wrapperStyle={{paddingTop: '10px'}}
                                />
                            </LineChart>
                        </ChartContainer>
                    </div>
                )}
            </CardContent>
        </Card>
    );
};

export default PriceEvolutionGraph;