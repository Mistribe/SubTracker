import {Badge} from "@/components/ui/badge";
import {Button} from "@/components/ui/button";
import {Money} from "@/components/ui/money";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import Provider from "@/models/provider";
import Subscription from "@/models/subscription";
import {RecurrencyBadge} from "@/components/subscriptions/ui/RecurrencyBadge";
import {
    CalendarIcon,
    CircleCheckIcon,
    CircleXIcon,
    CreditCardIcon,
    EllipsisVerticalIcon,
    Loader2,
    UsersIcon,
} from "lucide-react";
import {format} from "date-fns";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger
} from "@/components/ui/dropdown-menu.tsx";

export type SubscriptionsTableProps = {
    subscriptions: Subscription[];
    providerMap: Map<string, Provider>;
    onEdit: (subscription: Subscription) => void;
    onDelete: (subscription: Subscription) => void;
    isFetchingNextPage?: boolean;
};


export function SubscriptionsTable({
                                       subscriptions,
                                       providerMap,
                                       onEdit,
                                       onDelete,
                                       isFetchingNextPage,
                                   }: SubscriptionsTableProps) {
    return (
        <div className="overflow-hidden rounded-lg border">
            <Table>
                <TableHeader className="bg-muted sticky">
                    <TableRow>
                        <TableHead>Provider</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>
                            Price
                        </TableHead>
                        <TableHead>Recurrency</TableHead>
                        <TableHead>Dates</TableHead>
                        <TableHead>Users</TableHead>
                        <TableHead>
                            Status
                        </TableHead>
                        <TableHead className="w-[50px]"></TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {subscriptions.map((subscription) => (
                        <TableRow key={subscription.id}>
                            <TableCell>
                                <div className="flex items-center">
                                    {providerMap.get(subscription.providerId)?.iconUrl ? (
                                        <img
                                            src={providerMap.get(subscription.providerId)?.iconUrl || ""}
                                            alt={`${providerMap.get(subscription.providerId)?.name} logo`}
                                            className="mr-2 h-5 w-5 object-contain"
                                        />
                                    ) : (
                                        <CreditCardIcon className="mr-2 h-4 w-4 text-muted-foreground"/>
                                    )}
                                    <span className="font-bold">
                                        {providerMap.get(subscription.providerId)?.name || subscription.providerId}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell className="font-medium">
                                {subscription.friendlyName ||
                                    (providerMap.get(subscription.providerId)?.name ||
                                        subscription.providerId)}
                            </TableCell>
                            <TableCell>
                                {subscription.customPrice && (
                                    <Money amount={subscription.customPrice}/>
                                )}
                            </TableCell>
                            <TableCell>
                                <RecurrencyBadge recurrency={subscription.recurrency}
                                                 customRecurrency={subscription.customRecurrency}/>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center">
                                    <CalendarIcon className="mr-2 h-4 w-4 text-muted-foreground"/>
                                    <span>
                  {format(subscription.startDate, "MMM d, yyyy")}
                                        {subscription.endDate && <br/>}
                                        {subscription.endDate &&
                                            `Ends: ${format(subscription.endDate, "MMM d, yyyy")}`}
                </span>
                                </div>
                            </TableCell>
                            <TableCell>
                                {subscription.serviceUsers.length > 0 && (
                                    <div className="flex items-center">
                                        <UsersIcon className="mr-2 h-4 w-4 text-muted-foreground"/>
                                        <span>{subscription.serviceUsers.length}</span>
                                    </div>
                                )}
                            </TableCell>
                            <TableCell>
                                {subscription.freeTrial ? "Free Trial" : ""}
                                {subscription.freeTrial && subscription.isActive ? " - " : ""}
                                {subscription.isActive ? (
                                    <Badge variant="outline" className="text-muted-foreground px-1.5">
                                        <CircleCheckIcon className="fill-green-500 dark:fill-green-400"/> Active
                                    </Badge>
                                ) : (
                                    <Badge variant="outline" className="text-muted-foreground px-1.5">
                                        <CircleXIcon className="fill-red-500 dark:fill-red-400"/>Ended
                                    </Badge>
                                )}
                            </TableCell>
                            <TableCell>
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button
                                            variant="ghost"
                                            className="data-[state=open]:bg-muted text-muted-foreground flex size-8"
                                            size="icon"
                                        >
                                            <EllipsisVerticalIcon/>
                                            <span className="sr-only">Open menu</span>
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end" className="w-32">
                                        <DropdownMenuItem onClick={() => onEdit(subscription)}>
                                            Edit
                                        </DropdownMenuItem>
                                        <DropdownMenuItem>Make a copy</DropdownMenuItem>
                                        <DropdownMenuSeparator/>
                                        <DropdownMenuItem variant="destructive"
                                                          onClick={() => onDelete(subscription)}>
                                            Delete
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </TableCell>
                        </TableRow>
                    ))}
                    {isFetchingNextPage && (
                        <TableRow>
                            <TableCell colSpan={9} className="py-4 text-center">
                                <div className="flex items-center justify-center gap-2 text-muted-foreground">
                                    <Loader2 className="h-4 w-4 animate-spin"/>
                                    Loading more subscriptions...
                                </div>
                            </TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </div>
    );
}

export default SubscriptionsTable;
