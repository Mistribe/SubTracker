import {Skeleton} from "@/components/ui/skeleton";
import {Building2, Clock} from "lucide-react";
import Money from "@/components/ui/money.tsx";
import type TopProvider from "@/models/topProvider.ts";
import type Provider from "@/models/provider.ts";
import {formatProviderDuration} from "@/components/dashboard/utils";
import { useNavigate } from "react-router-dom";
import TopProviderCard from "@/components/dashboard/TopProviderCard.tsx";

interface TopProvidersProps {
    providers: TopProvider[];
    isLoading: boolean;
}

const TopProviders = ({providers, isLoading}: TopProvidersProps) => {
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
                            <TopProviderCard item={provider} />
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