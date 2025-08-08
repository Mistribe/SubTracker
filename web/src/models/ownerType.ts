export type OwnerTypeValues = 'system' | 'personal' | 'family';

export const OwnerType = {
    System: 'system',
    Personal: 'personal',
    Family: 'family',
} as const;

export type OwnerType = (typeof OwnerType)[keyof typeof OwnerType];

export function isValidOwnerType(type: string): type is OwnerType {
    return Object.values(OwnerType).includes(type as OwnerType);
}
