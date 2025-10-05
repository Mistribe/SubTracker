import { DtoFamilyMemberModelTypeEnum } from "@/api/models/DtoFamilyMemberModel";

export const FamilyMemberType = {
    Owner: 'owner',
    Adult: 'adult',
    Kid: 'kid'
} as const;

export type FamilyMemberType = (typeof FamilyMemberType)[keyof typeof FamilyMemberType];

export function fromHttpApi(memberType: DtoFamilyMemberModelTypeEnum | null | undefined): FamilyMemberType {
    switch (memberType) {
        case DtoFamilyMemberModelTypeEnum.Owner:
            return FamilyMemberType.Owner;
        case DtoFamilyMemberModelTypeEnum.Adult:
            return FamilyMemberType.Adult;
        default:
            return FamilyMemberType.Kid;
    }
}