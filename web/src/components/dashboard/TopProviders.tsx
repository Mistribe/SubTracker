import {Skeleton} from "@/components/ui/skeleton";
import {Building2, DollarSign} from "lucide-react";
import Money from "@/components/ui/money.tsx";
import type TopProvider from "@/models/topProvider.ts";
import type Provider from "@/models/provider.ts";
import {formatProviderDuration} from "@/components/dashboard/utils";
import { useNavigate } from "react-router-dom";

interface TopProvidersProps {
    providers: TopProvider[];
    providerMap: Map<string, Provider>;
    isLoading: boolean;
}

const TopProviders = ({providers, providerMap, isLoading}: TopProvidersProps) => {
    const navigate = useNavigate();
    return (
        <div>
            <div className="flex flex-row items-center justify-between pb-5 space-y-0">
                <p className="text-xl font-medium">Top Providers by Expense</p>
                <Building2 className="h-5 w-5 text-purple-500 animate-pulse"/>
            </div>
            <div>
                {isLoading ? (
                    <div className="space-y-3">
                        {Array.from({length: 3}).map((_, i) => (
                            <div key={`topprov-skeleton-${i}`} className="p-3 border rounded-lg bg-muted/30">
                                <Skeleton className="h-6 w-full mb-2"/>
                                <Skeleton className="h-4 w-1/2"/>
                            </div>
                        ))}
                    </div>
                ) : providers.length > 0 ? (
                    <div className="space-y-3">
                        {providers.map((provider) => (
                            <div key={provider.providerId}
                                 className="p-3 border rounded-lg bg-card transition-all duration-300 hover:shadow-lg hover:bg-accent/50 cursor-pointer"
                                 onClick={() => provider.providerId && navigate(`/providers/${provider.providerId}`)}>
                                <div className="flex justify-between items-start">
                                    <div className="flex items-center">
                                        <h4 className="font-medium">{providerMap.get(provider.providerId ?? "")?.name || provider.providerId}</h4>
                                    </div>
                                    <span
                                        className="font-semibold bg-gradient-to-r from-purple-500 to-purple-700 bg-clip-text text-transparent">
                                            <Money amount={provider.total}/>
                                        </span>
                                </div>
                                <p className="text-sm text-muted-foreground mt-1 flex items-center">
                                    <DollarSign className="h-3 w-3 mr-1 text-purple-500"/>
                                    {formatProviderDuration(provider.duration)}
                                </p>
                            </div>
                        ))}
                    </div>
                ) : (
                    <p className="text-muted-foreground">No provider spending data available.</p>
                )}
            </div>
        </div>
    );
};

export default TopProviders;