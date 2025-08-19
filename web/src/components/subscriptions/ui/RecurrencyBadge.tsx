import { Badge } from "@/components/ui/badge";
import { type SubscriptionRecurrency } from "@/models/subscriptionRecurrency";
import { formatRecurrency } from "@/components/dashboard/utils";

export type RecurrencyBadgeProps = {
  recurrency: SubscriptionRecurrency;
  customRecurrency?: number;
};

export function RecurrencyBadge({ recurrency, customRecurrency }: RecurrencyBadgeProps) {
  const label = formatRecurrency(recurrency, customRecurrency);
  return (
    <Badge variant="outline" className="px-1.5">
      {label}
    </Badge>
  );
}

export default RecurrencyBadge;
