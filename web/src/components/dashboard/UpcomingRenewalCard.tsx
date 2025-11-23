import {useNavigate} from "react-router-dom";
import {format} from "date-fns";
import {Calendar} from "lucide-react";
import {Money} from "@/components/ui/money";
import {Skeleton} from "@/components/ui/skeleton";
import type UpcomingRenewal from "@/models/upcomingRenewals.ts";
import {useSubscription} from "@/hooks/subscriptions/useSubscriptionQuery.ts";
import {useProviderQuery} from "@/hooks/providers/useProviderQuery.ts";

interface RenewalCardProps {
    item: UpcomingRenewal;
}

const UpcomingRenewalCard = ({item}: RenewalCardProps) => {
    const navigate = useNavigate();

    const {data: subscription, isLoading: isLoadingSubscription} = useSubscription(item.subscriptionId);
    const {data: provider, isLoading: isLoadingProvider} = useProviderQuery(item.providerId);

    const isLoading = isLoadingSubscription || isLoadingProvider;

    if (isLoading) {
        return (
            <div className="p-3 border rounded-lg bg-card">
                <div className="flex justify-between items-start">
                    <div className="flex items-center gap-2 flex-1">
                        <Skeleton className="h-10 w-10 rounded-md"/>
                        <Skeleton className="h-4 w-32"/>
                    </div>
                    <Skeleton className="h-6 w-16"/>
                </div>
                <div className="mt-2 flex items-center gap-1">
                    <Skeleton className="h-3 w-3"/>
                    <Skeleton className="h-4 w-40"/>
                </div>
            </div>
        );
    }

    return (
        <div
            className="p-3 border rounded-lg bg-card transition-all duration-300 hover:shadow-lg hover:bg-accent/50 cursor-pointer"
            onClick={() => item.providerId && navigate(`/providers/${item.providerId}`)}
        >
            <div className="flex justify-between items-start">
                <h4 className="flex items-center gap-2 font-medium">
                    {provider?.iconUrl  && (
                        <img
                            src={provider.iconUrl}
                            alt={provider.name}
                            className="h-10 w-10 object-contain rounded-md"
                        />
                    )}
                    {subscription && provider && (
                        <span className="text-xs text-muted-foreground">({subscription.friendlyName ?? provider.name})</span>
                    )}
                </h4>
                <span
                    className="font-semibold bg-gradient-to-r from-cyan-500 to-cyan-700 bg-clip-text text-transparent">
          <Money amount={item.total}/>
        </span>
            </div>
            <div className="mt-2 text-sm flex items-center">
                <Calendar className="h-3 w-3 mr-1 text-cyan-500"/>
                <span className="font-medium">Next renewal:</span>{" "}
                {item.at ? format(item.at, 'MMM d, yyyy') : ""}
            </div>
        </div>
    );
};

export default UpcomingRenewalCard;