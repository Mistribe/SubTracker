import { useMemo } from 'react';
import { FeatureId, FeatureType } from '@/models/billing';
import Quota from '@/models/quota';

/**
 * Hook to check if a specific feature is enabled
 * @param quotas - Array of quota objects from a quota query
 * @param featureId - The feature ID to check
 * @returns Object with `enabled` boolean and `quota` object (if applicable)
 */
export const useFeature = (quotas: Quota[] | undefined, featureId: FeatureId) => {
    return useMemo(() => {
        if (!quotas) {
            return { enabled: false, quota: undefined };
        }

        const quota = quotas.find((q) => q.feature === featureId);

        if (!quota) {
            // If quota not found, assume feature is disabled
            return { enabled: false, quota: undefined };
        }

        return { enabled: quota.enabled, quota };
    }, [quotas, featureId]);
};

/**
 * Hook to check if a quota-based feature has reached its limit
 * @param quotas - Array of quota objects from a quota query
 * @param featureId - The feature ID to check (must be a quota type feature)
 * @returns Object with:
 *   - `enabled`: whether the feature is enabled
 *   - `canAdd`: whether user can add more items (quota not reached)
 *   - `isAtLimit`: whether the quota limit has been reached
 *   - `used`: current usage count
 *   - `limit`: maximum allowed count
 *   - `remaining`: remaining quota
 *   - `quota`: the full quota object
 */
export const useQuotaLimit = (quotas: Quota[] | undefined, featureId: FeatureId) => {
    return useMemo(() => {
        if (!quotas) {
            return {
                enabled: false,
                canAdd: false,
                isAtLimit: true,
                used: 0,
                limit: 0,
                remaining: 0,
                quota: undefined,
            };
        }

        const quota = quotas.find((q) => q.feature === featureId);

        if (!quota) {
            return {
                enabled: false,
                canAdd: false,
                isAtLimit: true,
                used: 0,
                limit: 0,
                remaining: 0,
                quota: undefined,
            };
        }

        // For quota-type features, check if limit/used/remaining are present
        // If they are, the feature exists and we should check the quota, not just enabled flag
        if (quota.type === FeatureType.Quota) {
            const hasQuotaData = quota.limit !== undefined || quota.used !== undefined;

            if (hasQuotaData) {
                // Feature exists with quota data
                const remaining = quota.remaining ?? 0;
                const canAdd = remaining > 0;
                const isAtLimit = !canAdd;

                return {
                    enabled: true, // Feature exists (has quota data)
                    canAdd,
                    isAtLimit,
                    used: quota.used ?? 0,
                    limit: quota.limit ?? 0,
                    remaining,
                    quota,
                };
            }

            // No quota data and not enabled = feature not available
            if (!quota.enabled) {
                return {
                    enabled: false,
                    canAdd: false,
                    isAtLimit: true,
                    used: 0,
                    limit: 0,
                    remaining: 0,
                    quota,
                };
            }
        }

        // For boolean-type features, rely on enabled flag
        if (!quota.enabled) {
            return {
                enabled: false,
                canAdd: false,
                isAtLimit: true,
                used: quota.used ?? 0,
                limit: quota.limit ?? 0,
                remaining: quota.remaining ?? 0,
                quota,
            };
        }

        // Boolean feature is enabled (always can add)
        return {
            enabled: true,
            canAdd: true,
            isAtLimit: false,
            used: 0,
            limit: undefined,
            remaining: undefined,
            quota,
        };
    }, [quotas, featureId]);
};

/**
 * Helper to generate a tooltip message for disabled buttons based on quota status
 * @param enabled - Whether the feature is enabled
 * @param canAdd - Whether user can add more items
 * @param featureName - Human-readable name of the feature (e.g., "family members", "subscriptions")
 * @returns A tooltip message or undefined if button should be enabled
 */
export const getQuotaTooltip = (
    enabled: boolean,
    canAdd: boolean,
    featureName: string
): string | undefined => {
    if (!enabled) {
        return `This feature is not available on your current plan. Upgrade to create ${featureName}.`;
    }
    if (!canAdd) {
        return `You've reached the maximum number of ${featureName} for your plan. Upgrade to add more.`;
    }
    return undefined;
};
