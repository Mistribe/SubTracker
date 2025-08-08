import {Skeleton} from "@/components/ui/skeleton";
import {formatCurrency} from "./utils";

interface Provider {
    id: string;
    name: string;
    amount: number;
}

interface TopProvidersProps {
    providers: Provider[];
    isLoading: boolean;
}

const TopProviders = ({providers, isLoading}: TopProvidersProps) => {
    return (
        <div className="mb-8">
            <h3 className="text-xl font-semibold mb-4">Top Providers by Expense</h3>
            {isLoading ? (
                <div className="space-y-2">
                    {Array.from({length: 5}).map((_, i) => (
                        <div key={`topprov-skeleton-${i}`} className="p-4 border rounded-lg">
                            <Skeleton className="h-6 w-full mb-2"/>
                            <Skeleton className="h-4 w-1/2"/>
                        </div>
                    ))}
                </div>
            ) : providers.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {providers.map((provider) => (
                        <div key={provider.id} className="p-4 border rounded-lg">
                            <div className="flex justify-between items-start">
                                <h4 className="font-medium">{provider.name}</h4>
                                <span className="font-semibold">{formatCurrency(provider.amount)}</span>
                            </div>
                            <p className="text-sm text-muted-foreground mt-1">Annual cost</p>
                        </div>
                    ))}
                </div>
            ) : (
                <p className="text-muted-foreground">No provider spending data available.</p>
            )}
        </div>
    );
};

export default TopProviders;