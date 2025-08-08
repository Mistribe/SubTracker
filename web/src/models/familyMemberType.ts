import {type FamilyMemberModel_type, FamilyMemberModel_typeObject} from "@/api/models";

export const FamilyMemberType = {
    Owner: 'owner',
    Adult: 'adult',
    Kid: 'kid'
} as const;

export type FamilyMemberType = (typeof FamilyMemberType)[keyof typeof FamilyMemberType];

export function fromHttpApi(memberType: FamilyMemberModel_type | null | undefined): FamilyMemberType {
    switch (memberType) {
        case FamilyMemberModel_typeObject.Owner:
            return FamilyMemberType.Owner;
        case FamilyMemberModel_typeObject.Adult:
            return FamilyMemberType.Adult;
        default:
            return FamilyMemberType.Kid;
    }
}