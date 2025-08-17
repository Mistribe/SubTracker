import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Money } from "@/components/ui/money";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import Provider from "@/models/provider";
import Subscription from "@/models/subscription";
import { SubscriptionRecurrency } from "@/models/subscriptionRecurrency.ts";
import { CalendarIcon, CreditCardIcon, Loader2, PencilIcon, TagIcon, TrashIcon, UsersIcon } from "lucide-react";
import { format } from "date-fns";
import React from "react";

export type SubscriptionsTableProps = {
  subscriptions: Subscription[];
  providerMap: Map<string, Provider>;
  onEdit: (subscription: Subscription) => void;
  onDelete: (subscription: Subscription) => void;
  isFetchingNextPage?: boolean;
};

const formatRecurrency = (
  recurrency: SubscriptionRecurrency,
  customRecurrency: number | undefined
) => {
  if (recurrency === "custom" && customRecurrency) {
    return `Every ${customRecurrency} days`;
  }
  switch (recurrency) {
    case SubscriptionRecurrency.Monthly:
      return "Monthly";
    case SubscriptionRecurrency.Quarterly:
      return "Quarterly";
    case SubscriptionRecurrency.HalfYearly:
      return "Half Yearly";
    case SubscriptionRecurrency.Yearly:
      return "Yearly";
    case SubscriptionRecurrency.OneTime:
      return "OneTime";
    default:
      return recurrency as string;
  }
};

export function SubscriptionsTable({
  subscriptions,
  providerMap,
  onEdit,
  onDelete,
  isFetchingNextPage,
}: SubscriptionsTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Provider</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>Price</TableHead>
          <TableHead>Recurrency</TableHead>
          <TableHead>Dates</TableHead>
          <TableHead>Users</TableHead>
          <TableHead>Plan</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="w-[50px]">Actions</TableHead>
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
                  <CreditCardIcon className="mr-2 h-4 w-4 text-muted-foreground" />
                )}
                <span>
                  {providerMap.get(subscription.providerId)?.name ||
                    subscription.providerId}
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
                <Badge variant="outline">
                  <Money amount={subscription.customPrice} />
                </Badge>
              )}
            </TableCell>
            <TableCell>
              {formatRecurrency(
                subscription.recurrency,
                subscription.customRecurrency
              )}
            </TableCell>
            <TableCell>
              <div className="flex items-center">
                <CalendarIcon className="mr-2 h-4 w-4 text-muted-foreground" />
                <span>
                  {format(subscription.startDate, "MMM d, yyyy")}
                  {subscription.endDate && <br />}
                  {subscription.endDate &&
                    `Ends: ${format(subscription.endDate, "MMM d, yyyy")}`}
                </span>
              </div>
            </TableCell>
            <TableCell>
              {subscription.serviceUsers.length > 0 && (
                <div className="flex items-center">
                  <UsersIcon className="mr-2 h-4 w-4 text-muted-foreground" />
                  <span>{subscription.serviceUsers.length}</span>
                </div>
              )}
            </TableCell>
            <TableCell>
              <div className="flex items-center">
                <TagIcon className="mr-1 h-4 w-4 text-muted-foreground" />
                <span>{subscription.planId.substring(0, 8)}</span>
              </div>
            </TableCell>
            <TableCell>
              {subscription.freeTrial ? "Free Trial" : ""}
              {subscription.freeTrial && subscription.isActive ? " - " : ""}
              {subscription.isActive ? (
                <Badge variant="outline" className="bg-green-50 text-green-700">
                  Active
                </Badge>
              ) : (
                <Badge variant="outline" className="bg-red-50 text-red-700">
                  Ended
                </Badge>
              )}
            </TableCell>
            <TableCell>
              <div className="flex space-x-1">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => onEdit(subscription)}
                  title="Edit subscription"
                >
                  <PencilIcon className="h-4 w-4 text-blue-500" />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => onDelete(subscription)}
                  title="Delete subscription"
                >
                  <TrashIcon className="h-4 w-4 text-red-500" />
                </Button>
              </div>
            </TableCell>
          </TableRow>
        ))}
        {isFetchingNextPage && (
          <TableRow>
            <TableCell colSpan={9} className="py-4 text-center">
              <div className="flex items-center justify-center gap-2 text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" />
                Loading more subscriptions...
              </div>
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
}

export default SubscriptionsTable;
