import React from "react";
import Quota from "@/models/quota.ts";
import {FeatureType} from "@/models/billing.ts";
import {Progress} from "@/components/ui/progress";
import {cn} from "@/lib/utils";

export interface QuotaUsageProps {
  /** The quota domain object */
  quota: Quota;
  /** Optional custom label to display instead of deriving from feature id */
  label?: string;
  /** Additional class names */
  className?: string;
  /** Show remaining count (only for numeric quota) */
  showRemaining?: boolean;
  /** Compact variant */
  variant?: "default" | "compact";
}

/** Maps a feature id to a friendly label (fallback to id) */
function friendlyName(feature: string): string {
  switch (feature) {
    case "custom_labels_count":
      return "Custom Labels";
    case "active_subscriptions_count":
      return "Active Subscriptions";
    case "custom_providers_count":
      return "Custom Providers";
    case "family_members_count":
      return "Family Members";
    default:
      return feature.replace(/_/g, " ").replace(/\b\w/g, c => c.toUpperCase());
  }
}

export const QuotaUsage: React.FC<QuotaUsageProps> = ({
  quota,
  label,
  className,
  showRemaining = true,
  variant = "default",
}) => {
  const isBoolean = quota.type === FeatureType.Boolean;
  const limit = quota.limit ?? 0;
  const used = quota.used ?? 0;
  const remaining = quota.remaining ?? (limit > 0 ? Math.max(limit - used, 0) : undefined);
  const percent = !isBoolean && limit > 0 ? Math.min(100, (used / limit) * 100) : 0;

  return (
    <div
      className={cn(
        "border rounded-md p-4 bg-card text-card-foreground flex flex-col gap-2",
        variant === "compact" && "p-3 gap-1",
        className
      )}
    >
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">
          {label || friendlyName(quota.feature)}
        </span>
        {isBoolean ? (
          <span className={cn("text-xs font-semibold", quota.enabled ? "text-green-600 dark:text-green-500" : "text-red-600 dark:text-red-500")}>
            {quota.enabled ? "Enabled" : "Disabled"}
          </span>
        ) : limit > 0 ? (
          <span className="text-xs text-muted-foreground">{used} / {limit}</span>
        ) : (
          <span className="text-xs text-muted-foreground">{used}</span>
        )}
      </div>

      {!isBoolean && limit > 0 && (
        <Progress value={percent} />
      )}

      {!isBoolean && showRemaining && limit > 0 && (
        <div className="text-[11px] text-muted-foreground">
          {remaining !== undefined ? `${remaining} remaining` : null}
        </div>
      )}
    </div>
  );
};

export default QuotaUsage;
