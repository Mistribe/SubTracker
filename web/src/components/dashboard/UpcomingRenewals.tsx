import {Skeleton} from "@/components/ui/skeleton";
import {Calendar} from "lucide-react";
import type UpcomingRenewal from "@/models/upcomingRenewals.ts";
import UpcomingRenewalCard from "./UpcomingRenewalCard.tsx";

interface UpcomingRenewalsProps {
    summaryUpcomingRenewals?: UpcomingRenewal[];
    isLoading: boolean;
}

const UpcomingRenewals = ({
                               summaryUpcomingRenewals = [],
                               isLoading
                           }: UpcomingRenewalsProps) => {
    return (
        <div>
            <div>
                <div className="flex flex-row items-center justify-between pb-5 space-y-0">
                    <p className="text-xl font-medium">Upcoming Renewals</p>
                    <Calendar className="h-5 w-5 text-cyan-500 animate-pulse"/>
                </div>
                <div>
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
                                <UpcomingRenewalCard key={`${item.providerId}-${item.at?.toString() ?? ""}-${idx}`} item={item}/>
                            ))}
                        </div>
                    ) : (
                        <p className="text-muted-foreground">No upcoming renewals found.</p>
                    )}
                </div>
            </div>
        </div>
    );
};

export default UpcomingRenewals;