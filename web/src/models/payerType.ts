export const PayerType = {
    Family: 'family',
    FamilyMember: 'family_member'
} as const;

export type PayerType = (typeof PayerType)[keyof typeof PayerType];