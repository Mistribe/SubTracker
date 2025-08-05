import { Skeleton } from "@/components/ui/skeleton";
import { format } from "date-fns";
import { formatCurrency, formatRecurrency } from "./utils";
import Subscription from "@/models/subscription";

interface SubscriptionWithNextRenewal extends Subscription {
    nextRenewalDate: Date;
}

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
        <div className="mb-8">
            <h3 className="text-xl font-semibold mb-4">Upcoming Renewals</h3>
            {isLoading ? (
                <div className="space-y-2">
                    {[...Array(5)].map((_, i) => (
                        <div key={i} className="p-4 border rounded-lg">
                            <Skeleton className="h-6 w-full mb-2" />
                            <Skeleton className="h-4 w-1/2" />
                        </div>
                    ))}
                </div>
            ) : upcomingRenewals.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {upcomingRenewals.map((sub) => (
                        <div key={sub.id} className="p-4 border rounded-lg">
                            <div className="flex justify-between items-start">
                                <div>
                                    <h4 className="font-medium">
                                        {sub.friendlyName || providerMap.get(sub.providerId)?.name || sub.providerId}
                                    </h4>
                                    <p className="text-sm text-muted-foreground">
                                        {formatRecurrency(sub.recurrency, sub.customRecurrency)}
                                    </p>
                                </div>
                                {sub.customPrice && (
                                    <span className="font-semibold">
                                        {formatCurrency(sub.customPrice.amount)}
                                    </span>
                                )}
                            </div>
                            <div className="mt-2 text-sm">
                                <span className="font-medium">Next renewal:</span> {format(sub.nextRenewalDate, 'MMM d, yyyy')}
                            </div>
                        </div>
                    ))}
                </div>
            ) : (
                <p className="text-muted-foreground">No upcoming renewals found.</p>
            )}
        </div>
    );
};

export default UpcomingRenewals;