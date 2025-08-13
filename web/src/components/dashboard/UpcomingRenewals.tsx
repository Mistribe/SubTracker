import {Skeleton} from "@/components/ui/skeleton";
import {format} from "date-fns";
import {Money} from "@/components/ui/money";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card";
import {Calendar, Clock} from "lucide-react";
import type UpcomingRenewal from "@/models/upcomingRenewals.ts";

interface Provider {
    id: string;
    name: string;
}

interface UpcomingRenewalsProps {
    summaryUpcomingRenewals?: UpcomingRenewal[];
    providerMap: Map<string, Provider>;
    isLoading: boolean;
}

const UpcomingRenewals = ({
                              summaryUpcomingRenewals = [],
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
                    ) : summaryUpcomingRenewals.length > 0 ? (
                        <div className="space-y-3">
                            {summaryUpcomingRenewals.map((item, idx) => (
                                <div
                                    key={`${item.providerId}-${item.at?.toString() ?? ""}-${idx}`}
                                    className="p-3 border rounded-lg bg-card transition-all duration-200 hover:bg-muted/20"
                                >
                                    <div className="flex justify-between items-start">
                                        <h4 className="flex items-center font-medium">
                                            {providerMap.get(item.providerId ?? "")?.name || item.providerId}
                                        </h4>
                                        <span
                                            className="font-semibold bg-gradient-to-r from-cyan-500 to-cyan-700 bg-clip-text text-transparent">
                                                <Money amount={item.total}/>
                                            </span>
                                    </div>
                                    <div className="mt-2 text-sm flex items-center">
                                        <Clock className="h-3 w-3 mr-1 text-cyan-500"/>
                                        <span className="font-medium">Next renewal:</span>{" "}
                                        {item.at ? format(item.at, 'MMM d, yyyy') : ""}
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