import {Skeleton} from "@/components/ui/skeleton";
import {formatCurrency} from "./utils";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle
} from "@/components/ui/card";
import {Building2, DollarSign} from "lucide-react";

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
        <div>
            <Card className="overflow-hidden transition-all duration-300 hover:shadow-lg border-t-4 border-t-purple-500">
                <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                    <CardTitle className="text-xl font-medium">Top Providers by Expense</CardTitle>
                    <Building2 className="h-5 w-5 text-purple-500 animate-pulse" />
                </CardHeader>
                <CardContent>
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
                                <div key={provider.id} className="p-3 border rounded-lg bg-card transition-all duration-200 hover:bg-muted/20">
                                    <div className="flex justify-between items-start">
                                        <div className="flex items-center">
                                            <h4 className="font-medium">{provider.name}</h4>
                                        </div>
                                        <span className="font-semibold bg-gradient-to-r from-purple-500 to-purple-700 bg-clip-text text-transparent">
                                            {formatCurrency(provider.amount)}
                                        </span>
                                    </div>
                                    <p className="text-sm text-muted-foreground mt-1 flex items-center">
                                        <DollarSign className="h-3 w-3 mr-1 text-purple-500" />
                                        Annual cost
                                    </p>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <p className="text-muted-foreground">No provider spending data available.</p>
                    )}
                </CardContent>
            </Card>
        </div>
    );
};

export default TopProviders;