export type FeatureIdValues =
    "unknown"
    | "subscriptions"
    | "active_subscriptions_count"
    | "custom_labels"
    | "custom_labels_count"
    | "custom_providers"
    | "custom_providers_count"
    | "family"
    | "family_members_count";

export const FeatureId = {
    Unknown: 'unknown',
    Subscriptions: 'subscriptions',
    ActiveSubscriptionsCount: 'active_subscriptions_count',
    CustomLabels: 'custom_labels',
    CustomLabelsCount: 'custom_labels_count',
    CustomProviders: 'custom_providers',
    CustomProvidersCount: 'custom_providers_count',
    Family: 'family',
    FamilyMembersCount: 'family_members_count',
}

export type FeatureId = (typeof FeatureId)[keyof typeof FeatureId];

export function isValidFeatureId(id: string): id is FeatureId {
    return Object.values(FeatureId).includes(id as FeatureId);
}

export type FeatureTypeValues = 'boolean' | 'quota' | 'unknown';

export const FeatureType = {Boolean: 'boolean', Quota: 'quota', Unknown: 'unknown'};
export type FeatureType = (typeof FeatureType)[keyof typeof FeatureType];

export function isValidFeatureType(type: string): type is FeatureTypeValues {
    return Object.values(FeatureType).includes(type as FeatureTypeValues);
}