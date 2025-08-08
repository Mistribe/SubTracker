import {Button} from "@/components/ui/button";
import {Skeleton} from "@/components/ui/skeleton";
import Subscription from "@/models/subscription";
import {format} from "date-fns";
import {formatCurrency, formatRecurrency} from "./utils";

interface SubscriptionWithNextRenewal extends Subscription {
    nextRenewalDate: Date;
}

interface Provider {
    id: string;
    name: string;
}

interface SubscriptionsTableProps {
    subscriptions: Subscription[];
    subscriptionsWithNextRenewal: SubscriptionWithNextRenewal[];
    providerMap: Map<string, Provider>;
    isLoading: boolean;
}

const SubscriptionsTable = ({
                                subscriptions,
                                subscriptionsWithNextRenewal,
                                providerMap,
                                isLoading
                            }: SubscriptionsTableProps) => {
    return (
        <div>
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-xl font-semibold">Your Subscriptions</h3>
                <Button onClick={() => window.location.href = "/subscriptions/create"}>
                    Add New
                </Button>
            </div>

            {isLoading ? (
                <div className="space-y-2">
                    {Array.from({length: 3}).map((_, i) => (
                        <div key={`skeleton-${i}`} className="p-4 border rounded-lg">
                            <Skeleton className="h-6 w-full mb-2"/>
                            <Skeleton className="h-4 w-3/4"/>
                        </div>
                    ))}
                </div>
            ) : subscriptions.length > 0 ? (
                <div className="overflow-x-auto">
                    <table className="w-full border-collapse">
                        <thead>
                        <tr className="bg-muted">
                            <th className="p-3 text-left">Name</th>
                            <th className="p-3 text-left">Amount</th>
                            <th className="p-3 text-left">Frequency</th>
                            <th className="p-3 text-left">Next Payment</th>
                            <th className="p-3 text-left">Actions</th>
                        </tr>
                        </thead>
                        <tbody>
                        {subscriptions.map((subscription, idx) => (
                            <tr key={subscription.id ?? `${subscription.providerId}-${idx}`} className="border-b">
                                <td className="p-3">
                                    {subscription.friendlyName ||
                                        providerMap.get(subscription.providerId)?.name ||
                                        subscription.providerId}
                                </td>
                                <td className="p-3">
                                    {subscription.customPrice ?
                                        formatCurrency(subscription.customPrice.amount) :
                                        'N/A'}
                                </td>
                                <td className="p-3">
                                    {formatRecurrency(subscription.recurrency, subscription.customRecurrency)}
                                </td>
                                <td className="p-3">
                                    {subscription.isActive ?
                                        format(
                                            subscriptionsWithNextRenewal.find(s => s.id === subscription.id)?.nextRenewalDate ||
                                            new Date(),
                                            'MMM d, yyyy'
                                        ) :
                                        'Ended'}
                                </td>
                                <td className="p-3">
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="text-blue-500 hover:underline mr-2"
                                        onClick={() => window.location.href = `/subscriptions/edit/${subscription.id}`}
                                    >
                                        Edit
                                    </Button>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="text-red-500 hover:underline"
                                    >
                                        Delete
                                    </Button>
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            ) : (
                <p className="text-muted-foreground">No subscriptions found. Add your first subscription to get
                    started.</p>
            )}
        </div>
    );
};

export default SubscriptionsTable;