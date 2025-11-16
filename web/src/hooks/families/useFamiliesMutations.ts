import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import {FamilyMemberType} from "@/models/familyMemberType.ts";
import type {
    DtoCreateFamilyMemberRequest as CreateFamilyMemberModel,
    DtoCreateFamilyRequest as CreateFamilyModel,
    DtoUpdateFamilyMemberRequest as UpdateFamilyMemberModel
} from "@/api/models";
import {DtoCreateLabelRequestOwnerEnum} from "@/api/models";
import {useFamilyQuery} from "@/hooks/families/useFamilyQuery.ts";

export const useFamiliesMutations = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();
    const familiesQuery = useFamilyQuery();

    // Create family mutation
    const createFamilyMutation = useMutation({
        mutationFn: async (familyData: { name: string, creatorName: string }) => {
            const payload: CreateFamilyModel = {
                name: familyData.name,
                creatorName: familyData.creatorName
            };

            return apiClient?.families.familyPost({dtoCreateFamilyRequest: payload});
        },
        onSuccess: async () => {
            // Invalidate and refetch
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Add family member mutation
    const addFamilyMemberMutation = useMutation({
        mutationFn: async ({familyId, name, isKid}: { familyId: string, name: string, isKid?: boolean }) => {
            const payload: CreateFamilyMemberModel = {
                name,
                type: isKid ? FamilyMemberType.Kid : FamilyMemberType.Adult
            };

            return apiClient?.families.familyFamilyIdMembersPost({familyId, dtoCreateFamilyMemberRequest: payload});
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Delete family mutation
    const deleteFamilyMutation = useMutation({
        mutationFn: async (familyId: string) => {
            return apiClient?.families.familyFamilyIdDelete({familyId});
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Remove family member mutation
    const removeFamilyMemberMutation = useMutation({
        mutationFn: async ({familyId, memberId}: { familyId: string, memberId: string }) => {
            return apiClient?.families.familyFamilyIdMembersFamilyMemberIdDelete({familyId, familyMemberId: memberId});
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Create label for family mutation
    const createFamilyLabelMutation = useMutation({
        mutationFn: async ({name, color}: { familyId: string, name: string, color: string }) => {
            const payload = {
                name,
                color,
                owner: DtoCreateLabelRequestOwnerEnum.Family
            };
            return apiClient?.labels.labelsPost({dtoCreateLabelRequest: payload});
        },
        onSuccess: async () => {
            // Invalidate both families and labels queries
            await queryClient.invalidateQueries({queryKey: ['families']});
            await queryClient.invalidateQueries({queryKey: ['labels']});
        }
    });

    // Update family member mutation
    const updateFamilyMemberMutation = useMutation({
        mutationFn: async ({familyId, memberId, name, isKid}: {
            familyId: string,
            memberId: string,
            name: string,
            isKid: boolean
        }) => {
            const payload: Partial<UpdateFamilyMemberModel> = {
                name,
                type: isKid ? FamilyMemberType.Kid : FamilyMemberType.Adult
            };

            return apiClient?.families.familyFamilyIdMembersFamilyMemberIdPut({
                familyId,
                familyMemberId: memberId,
                dtoUpdateFamilyMemberRequest: payload as UpdateFamilyMemberModel
            });
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Invite family member mutation (email or link generation)
    const inviteFamilyMemberMutation = useMutation({
        mutationFn: async ({familyId, familyMemberId, email}: {
            familyId: string,
            familyMemberId: string,
            email?: string
        }) => {
            return apiClient?.families.familyFamilyIdInvitePost({
                familyId,
                dtoFamilyInviteRequest: {
                    familyMemberId,
                    email: email || undefined,
                },
            });
        },
        onSuccess: async () => {
            // Invitation does not change members list immediately; optional invalidate
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Revoke (unlink) linked account for a family member
    const revokeFamilyMemberMutation = useMutation({
        mutationFn: async ({familyId, memberId}: { familyId: string, memberId: string }) => {
            return apiClient?.families.familyFamilyIdMembersFamilyMemberIdRevokePost({
                familyId,
                familyMemberId: memberId,
                body: {}
            });
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Accept family invitation mutation
    const acceptFamilyInvitationMutation = useMutation({
        mutationFn: async ({familyId, invitationCode, familyMemberId}: {
            familyId: string,
            invitationCode: string,
            familyMemberId?: string
        }) => {
            return apiClient?.families.familyFamilyIdAcceptPost({
                familyId,
                dtoFamilyAcceptInvitationRequest: {
                    invitationCode,
                    familyMemberId: familyMemberId ?? '',
                },
            });
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Leave current family mutation
    const leaveFamilyMutation = useMutation({
        mutationFn: async () => {
            if (!apiClient) throw new Error("API client not initialized");
            const currentFamily = familiesQuery.data;
            if (!currentFamily) throw new Error("No family to leave");
            const you = currentFamily.members.find(m => m.isYou);
            if (!you) throw new Error("Could not determine your member entry");

            if (currentFamily.isOwner && currentFamily.members.length === 1) {
                return apiClient.families.familyFamilyIdDelete({familyId: currentFamily.id});
            }

            return apiClient.families.familyFamilyIdMembersFamilyMemberIdDelete({
                familyId: currentFamily.id,
                familyMemberId: you.id
            });
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    return {
        createFamilyMutation,
        addFamilyMemberMutation,
        deleteFamilyMutation,
        removeFamilyMemberMutation,
        createFamilyLabelMutation,
        updateFamilyMemberMutation,
        inviteFamilyMemberMutation,
        revokeFamilyMemberMutation,
        acceptFamilyInvitationMutation,
        leaveFamilyMutation
    };
};