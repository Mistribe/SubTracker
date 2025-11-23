import Money from "@/components/ui/money.tsx";
import {Clock} from "lucide-react";
import {formatProviderDuration} from "@/components/dashboard/utils.ts";
import type TopProvider from "@/models/topProvider.ts";
import {useProviderQuery} from "@/hooks/providers/useProviderQuery.ts";
import {useNavigate} from "react-router-dom";
import {Skeleton} from "@/components/ui/skeleton";

interface TopProviderCardProps {
    item: TopProvider;
}

const TopProviderCard = ({item}: TopProviderCardProps) => {
    const navigate = useNavigate();
    const {data: provider, isLoading} = useProviderQuery(item.providerId);

    if (isLoading) {
        return (
            <div className="p-3 border rounded-lg bg-card">
                <div className="flex justify-between items-start">
                    <div className="flex items-center gap-2 flex-1">
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

    return <div key={item.providerId}
                className="p-3 border rounded-lg bg-card transition-all duration-300 hover:shadow-lg hover:bg-accent/50 cursor-pointer"
                onClick={() => item.providerId && navigate(`/providers/${item.providerId}`)}>
        <div className="flex justify-between items-start">
            <div className="flex items-center">
                <h4 className="font-medium">{provider?.name ?? item.providerId ?? "Unknown provider"}</h4>
            </div>
            <span
                className="font-semibold bg-gradient-to-r from-purple-500 to-purple-700 bg-clip-text text-transparent">
                                            <Money amount={item.total}/>
                                        </span>
        </div>
        <p className="text-sm text-muted-foreground mt-1 flex items-center">
            <Clock className="h-3 w-3 mr-1 text-purple-500"/>
            {formatProviderDuration(item.duration)}
        </p>
    </div>
}

export default TopProviderCard;