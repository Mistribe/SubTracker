import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Family from "@/models/family";
import type {DtoFamilyModel} from "@/api";

export const useFamilyQuery = () => {
    const {apiClient} = useApiClient();

    return useQuery<Family | null>({
        queryKey: ['families'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            try {
                // Fetch the current user's family response; backend now returns a wrapper with optional family
                const family: DtoFamilyModel = await apiClient.families.familyGet();
                if (family) {
                   return Family.fromModel(family);
                }
                return null;
            } catch (error: unknown) {
                let status: number | undefined;
                if (typeof error === 'object' && error !== null) {
                    const e = error as Partial<{ responseStatusCode: number; statusCode: number; status: number }>;
                    status = e.responseStatusCode ?? e.statusCode ?? e.status;
                }
                if (status === 404) {
                    // No family registered for this user
                    return null;
                }
                throw error;
            }
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};