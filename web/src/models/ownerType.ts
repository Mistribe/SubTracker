export type OwnerTypeValues = 'system' | 'personal' | 'family';

export enum OwnerType {
    System = 'system',
    Personal = 'personal',
    Family = 'family',
}

export function isValidOwnerType(type: string): type is OwnerType {
    return Object.values(OwnerType).includes(type as OwnerType);
}
