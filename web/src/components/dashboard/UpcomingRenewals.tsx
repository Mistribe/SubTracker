import {Skeleton} from "@/components/ui/skeleton";
import {format} from "date-fns";
import {formatRecurrency} from "./utils";
import { Money } from "@/components/ui/money";
import type {SubscriptionWithNextRenewal} from "@/models/subscriptionWithNextRenewal";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Calendar, Clock} from "lucide-react";

interface Provider {
    id: string;
    name: string;
}

interface UpcomingRenewalsProps {
    upcomingRenewals: SubscriptionWithNextRenewal[];
    providerMap: Map<string, Provider>;
    isLoading: boolean;
}

const UpcomingRenewals = ({
                              upcomingRenewals,
                              providerMap,
                              isLoading
                          }: UpcomingRenewalsProps) => {
    return (
        <div>
            <Card className="overflow-hidden transition-all duration-300 hover:shadow-lg border-t-4 border-t-cyan-500">
                <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                    <CardTitle className="text-xl font-medium">Upcoming Renewals</CardTitle>
                    <Calendar className="h-5 w-5 text-cyan-500 animate-pulse"/>
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <div className="space-y-3">
                            {Array.from({length: 3}).map((_, i) => (
                                <div key={`upcoming-skeleton-${i}`} className="p-3 border rounded-lg bg-muted/30">
                                    <Skeleton className="h-6 w-full mb-2"/>
                                    <Skeleton className="h-4 w-1/2"/>
                                </div>
                            ))}
                        </div>
                    ) : upcomingRenewals.length > 0 ? (
                        <div className="space-y-3">
                            {upcomingRenewals.map((item, idx) => (
                                <div
                                    key={
                                        item.subscription.id ??
                                        `${item.subscription.providerId}-${format(item.nextRenewalDate, 'yyyy-MM-dd')}-${idx}`
                                    }
                                    className="p-3 border rounded-lg bg-card transition-all duration-200 hover:bg-muted/20"
                                >
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <h4 className="font-medium">
                                                {item.subscription.friendlyName ||
                                                    providerMap.get(item.subscription.providerId)?.name ||
                                                    item.subscription.providerId}
                                            </h4>
                                            <p className="text-sm text-muted-foreground">
                                                {formatRecurrency(item.subscription.recurrency, item.subscription.customRecurrency)}
                                            </p>
                                        </div>
                                        {item.subscription.customPrice && (
                                            <span className="font-semibold bg-gradient-to-r from-cyan-500 to-cyan-700 bg-clip-text text-transparent">
                                                <Money amount={item.subscription.customPrice.amount} currency={item.subscription.customPrice.currency} />
                                            </span>
                                        )}
                                    </div>
                                    <div className="mt-2 text-sm flex items-center">
                                        <Clock className="h-3 w-3 mr-1 text-cyan-500"/>
                                        <span className="font-medium">Next renewal:</span>{" "}
                                        {format(item.nextRenewalDate, 'MMM d, yyyy')}
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <p className="text-muted-foreground">No upcoming renewals found.</p>
                    )}
                </CardContent>
            </Card>
        </div>
    );
};

export default UpcomingRenewals;