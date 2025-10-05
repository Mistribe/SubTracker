import React, { useState } from "react";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { InfoIcon, AlertCircleIcon } from "lucide-react";
import { QuotaUsage } from "@/components/quotas/QuotaUsage";
import { QuotaUsageSkeleton } from "@/components/quotas/QuotaUsageSkeleton";
import Quota from "@/models/quota.ts";
import { FeatureId, FeatureType } from "@/models/billing.ts";
import { cn } from "@/lib/utils";

export interface QuotaButtonProps {
    /** Query hook that returns quota data - will be called when popover opens */
    useQuotaQuery: () => {
        data: Quota[] | undefined;
        isLoading: boolean;
        error: unknown;
    };
    /** Feature IDs to display from the quota data */
    featureIds: FeatureId[];
    /** Custom labels for features (optional) */
    featureLabels?: Record<FeatureId, string>;
}

/**
 * QuotaButton component displays an info icon button that opens a popover
 * with quota information. Data is fetched on-demand when the popover is opened.
 * When quota is exhausted, it shows a warning indicator.
 */
export const QuotaButton: React.FC<QuotaButtonProps> = ({
    useQuotaQuery,
    featureIds,
    featureLabels,
}) => {
    const [isOpen, setIsOpen] = useState(false);

    // Fetch data to check quota status (React Query will cache this)
    const { data: quotaData, isLoading, error } = useQuotaQuery();

    // Filter quotas by the specified feature IDs
    const quotas = quotaData?.filter((q) => featureIds.includes(q.feature)) || [];

    // Check if any quota is exhausted (remaining is 0 or used equals limit for numeric quotas)
    const hasExhaustedQuota = !isLoading && quotas.some((quota) => {
        if (quota.type === FeatureType.Boolean) return false;
        if (quota.limit === undefined || quota.limit <= 0) return false;
        const remaining = quota.remaining ?? 0;
        return remaining <= 0;
    });

    return (
        <Popover open={isOpen} onOpenChange={setIsOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant={hasExhaustedQuota ? "destructive" : "outline"}
                    size="icon"
                    className={cn(
                        hasExhaustedQuota && "animate-pulse"
                    )}
                    title={hasExhaustedQuota ? "Quota limit reached" : "View quota usage"}
                >
                    {hasExhaustedQuota ? (
                        <AlertCircleIcon className="h-4 w-4" />
                    ) : (
                        <InfoIcon className="h-4 w-4" />
                    )}
                </Button>
            </PopoverTrigger>
            <PopoverContent align="end" className="w-80">
                <div className="space-y-3">
                    <div className="font-semibold text-sm mb-3">Quota Usage</div>

                    {isLoading && (
                        <div className="space-y-2">
                            <QuotaUsageSkeleton />
                        </div>
                    )}

                    {!isLoading && !!error && (
                        <div className="text-xs text-destructive border border-destructive/30 rounded-md p-3">
                            Failed to load quota information.
                        </div>
                    )}

                    {!isLoading && !error && quotas.length === 0 && (
                        <div className="text-xs text-muted-foreground border rounded-md p-3">
                            No quota data available.
                        </div>
                    )}

                    {!isLoading && !error && quotas.length > 0 && (
                        <div className="space-y-2">
                            {quotas.map((quota) => (
                                <QuotaUsage
                                    key={quota.feature}
                                    quota={quota}
                                    label={featureLabels?.[quota.feature]}
                                    variant="compact"
                                />
                            ))}
                        </div>
                    )}
                </div>
            </PopoverContent>
        </Popover>
    );
};

export default QuotaButton;
