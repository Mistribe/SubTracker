import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {CreateFamilyMemberModel, CreateFamilyModel, UpdateFamilyMemberModel} from "@/api/models";
import {FamilyMemberType} from "@/models/familyMemberType.ts";

export const useFamiliesMutations = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // Create family mutation
    const createFamilyMutation = useMutation({
        mutationFn: async (familyData: { name: string, creatorName: string }) => {
            const payload: CreateFamilyModel = {
                name: familyData.name,
                creatorName: familyData.creatorName
            };

            return apiClient?.families.post(payload);
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

            return apiClient?.families.byFamilyId(familyId).members.post(payload);
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Delete family mutation
    const deleteFamilyMutation = useMutation({
        mutationFn: async (familyId: string) => {
            return apiClient?.families.byFamilyId(familyId).delete();
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Remove family member mutation
    const removeFamilyMemberMutation = useMutation({
        mutationFn: async ({familyId, memberId}: { familyId: string, memberId: string }) => {
            return apiClient?.families.byFamilyId(familyId).members.byFamilyMemberId(memberId).delete();
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Create label for family mutation
    const createFamilyLabelMutation = useMutation({
        mutationFn: async ({familyId, name, color}: { familyId: string, name: string, color: string }) => {
            return apiClient?.labels.post({
                name,
                color,
                owner: {
                    type: 'family',
                    familyId: familyId
                }
            });
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

            return apiClient?.families.byFamilyId(familyId).members.byId(memberId).put(payload);
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({queryKey: ['families']});
        }
    });

    // Invite family member mutation (email or link generation)
    const inviteFamilyMemberMutation = useMutation({
        mutationFn: async ({ familyId, familyMemberId, email }: { familyId: string, familyMemberId: string, email?: string }) => {
            return apiClient?.families.byFamilyId(familyId).invite.post({
                familyMemberId,
                email: email || undefined,
            });
        },
        onSuccess: async () => {
            // Invitation does not change members list immediately; optional invalidate
            await queryClient.invalidateQueries({ queryKey: ['families'] });
        }
    });

    // Revoke (unlink) linked account for a family member
    const revokeFamilyMemberMutation = useMutation({
        mutationFn: async ({ familyId, memberId }: { familyId: string, memberId: string }) => {
            return apiClient?.families
                .byFamilyId(familyId)
                .members
                .byFamilyMemberId(memberId)
                .revoke
                .post({});
        },
        onSuccess: async () => {
            await queryClient.invalidateQueries({ queryKey: ['families'] });
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
        revokeFamilyMemberMutation
    };
};