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
import {Bar, CartesianGrid, ComposedChart, Line, XAxis, YAxis} from "recharts";
import {ChevronLeft, ChevronRight, RefreshCw, TrendingUp} from "lucide-react";
import type Subscription from "@/models/subscription";
import {addMonths, format, startOfMonth, subMonths} from "date-fns";
import { usePreferredCurrency } from "@/hooks/currencies/usePreferredCurrency";
import { useCurrencyRates } from "@/hooks/currencies/useCurrencyRates";
import { subscriptionMonthlyPriceInCurrency } from "@/utils/currency";

interface Provider {
    id: string;
    name: string;
}

interface PriceEvolutionGraphProps {
    subscriptions: Subscription[];
    providerMap: Map<string, Provider>;
    isLoading: boolean;
}

// Interface for subscription data by provider
interface ProviderSubscriptionData {
    providerId: string;
    providerName: string;
    monthlyPrice: number;
}

interface MonthlyData {
    month: string;
    date: Date;
    monthlySubscriptionsTotal: number;
    // Dynamic provider data will be added to this object
    [key: string]: string | Date | number | ProviderSubscriptionData[];
    providerData: ProviderSubscriptionData[];
}

const PriceEvolutionGraph = ({
                                 subscriptions,
                                 providerMap,
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
    const { preferredCurrency } = usePreferredCurrency();
    const { rates } = useCurrencyRates();

    const monthlyData = useMemo(() => {
        const data: MonthlyData[] = [];
        const currentMonth = startOfMonth(referenceDate);

        // Create data for the past 6 months and next 6 months (total 13 months including current)
        for (let i = -6; i <= 6; i++) {
            const monthDate = addMonths(currentMonth, i);
            const monthStr = format(monthDate, 'MMM yyyy');
            
            // Filter subscriptions that are active in this specific month
            const activeSubscriptions = subscriptions.filter(subscription => {
                const startDate = subscription.startDate;
                const endDate = subscription.endDate;
                
                // Check if subscription is active for this month:
                // 1. Start date is before or equal to the end of this month
                // 2. End date is undefined or strictly after the start of this month (end date is exclusive)
                const monthEnd = addMonths(monthDate, 1);
                const isActiveInMonth = 
                    startDate <= monthEnd && 
                    (!endDate || endDate > monthDate);
                    
                return isActiveInMonth;
            });
            
            // Calculate total for all active subscriptions in this month (converted to preferred currency)
            const monthlySubscriptionsTotal = activeSubscriptions.reduce((sum, subscription) => {
                return sum + subscriptionMonthlyPriceInCurrency(subscription, preferredCurrency, rates);
            }, 0);
            
            // Group subscriptions by provider
            const providerData: ProviderSubscriptionData[] = [];
            
            // Create a map to aggregate subscription costs by provider
            const providerCosts = new Map<string, number>();
            
            activeSubscriptions.forEach(subscription => {
                const providerId = subscription.providerId;
                const monthlyPrice = subscriptionMonthlyPriceInCurrency(subscription, preferredCurrency, rates);
                
                if (providerCosts.has(providerId)) {
                    providerCosts.set(providerId, providerCosts.get(providerId)! + monthlyPrice);
                } else {
                    providerCosts.set(providerId, monthlyPrice);
                }
            });
            
            // Convert the map to an array of provider data
            providerCosts.forEach((monthlyPrice, providerId) => {
                const providerName = providerMap.get(providerId)?.name || providerId;
                
                providerData.push({
                    providerId,
                    providerName,
                    monthlyPrice
                });
                
                // Also add the provider cost as a direct property on the data object
                // This will be used for the line chart
            });
            
            // Sort provider data by monthly price (descending)
            providerData.sort((a, b) => b.monthlyPrice - a.monthlyPrice);
            
            // Create the monthly data object
            const monthlyDataObj: MonthlyData = {
                month: monthStr,
                date: monthDate,
                monthlySubscriptionsTotal,
                providerData
            };
            
            // Add each provider's cost as a direct property
            providerData.forEach(provider => {
                const key = `provider_${provider.providerId}`;
                monthlyDataObj[key] = provider.monthlyPrice;
            });
            
            data.push(monthlyDataObj);
        }

        return data;
    }, [subscriptions, providerMap, referenceDate, preferredCurrency, rates]);

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
                                totalSubscriptions: {
                                    label: "Total Subscriptions",
                                    theme: {
                                        light: "#0891b2", // cyan-600
                                        dark: "#06b6d4"   // cyan-500
                                    }
                                },
                                // Add dynamic provider colors
                                ...Object.fromEntries(
                                    // Get unique providers across all months
                                    Array.from(
                                        new Set(
                                            monthlyData.flatMap(data => 
                                                data.providerData.map(p => p.providerId)
                                            )
                                        )
                                    ).map((providerId, index) => {
                                        // Find provider name from any month that has this provider
                                        const providerName = monthlyData.find(
                                            data => data.providerData.some(p => p.providerId === providerId)
                                        )?.providerData.find(p => p.providerId === providerId)?.providerName || providerId;
                                        
                                        // Generate a color based on index
                                        // Using a set of predefined colors for better visual distinction
                                        const colors = [
                                            { light: "#f97316", dark: "#fb923c" }, // orange
                                            { light: "#a855f7", dark: "#c084fc" }, // purple
                                            { light: "#ec4899", dark: "#f472b6" }, // pink
                                            { light: "#14b8a6", dark: "#2dd4bf" }, // teal
                                            { light: "#f59e0b", dark: "#fbbf24" }, // amber
                                            { light: "#8b5cf6", dark: "#a78bfa" }, // violet
                                            { light: "#ef4444", dark: "#f87171" }, // red
                                            { light: "#10b981", dark: "#34d399" }, // emerald
                                            { light: "#6366f1", dark: "#818cf8" }, // indigo
                                            { light: "#0ea5e9", dark: "#38bdf8" }  // sky
                                        ];
                                        
                                        const colorIndex = index % colors.length;
                                        
                                        return [
                                            `provider_${providerId}`,
                                            {
                                                label: providerName,
                                                theme: colors[colorIndex]
                                            }
                                        ];
                                    })
                                )
                            }}
                        >
                            <ComposedChart
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
                                    tickFormatter={(value) => new Intl.NumberFormat(undefined, { style: 'currency', currency: preferredCurrency, maximumFractionDigits: 0 }).format(value as number)}
                                />
                                <ChartTooltip
                                    content={
                                        <ChartTooltipContent
                                            formatter={(value) => (
                                                <span>{new Intl.NumberFormat(undefined, { style: 'currency', currency: preferredCurrency }).format(value as number)}</span>
                                            )}
                                        />
                                    }
                                />
                                
                                {/* Bar chart for total subscription costs */}
                                <Bar
                                    dataKey="monthlySubscriptionsTotal"
                                    name="totalSubscriptions"
                                    fill="var(--color-totalSubscriptions-bg)"
                                    barSize={20}
                                    opacity={0.7}
                                />
                                
                                {/* Line charts for each provider */}
                                {Array.from(
                                    new Set(
                                        monthlyData.flatMap(data => 
                                            data.providerData.map(p => p.providerId)
                                        )
                                    )
                                ).map(providerId => (
                                    <Line
                                        key={providerId}
                                        type="monotone"
                                        dataKey={`provider_${providerId}`}
                                        name={`provider_${providerId}`}
                                        strokeWidth={2}
                                        dot={{strokeWidth: 2}}
                                        activeDot={{r: 4}}
                                    />
                                ))}
                                
                                <ChartLegend
                                    content={<ChartLegendContent/>}
                                    verticalAlign="bottom"
                                    height={60}
                                    wrapperStyle={{paddingTop: '10px'}}
                                />
                            </ComposedChart>
                        </ChartContainer>
                    </div>
                )}
            </CardContent>
        </Card>
    );
};

export default PriceEvolutionGraph;