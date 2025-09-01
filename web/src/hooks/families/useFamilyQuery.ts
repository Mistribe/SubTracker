import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Family from "@/models/family";
import type {FamilyModel} from "@/api/models/family";


export const useFamilyQuery = () => {
    const {apiClient} = useApiClient();

    return useQuery({
        queryKey: ['families'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            try {
                // Fetch the current user's single family; backend now supports only one family
                const result: FamilyModel | undefined = await apiClient.families.me.get();

                if (result) {
                   return Family.fromModel(result);
                }
                return undefined;
            } catch (error: unknown) {
                let status: number | undefined;
                if (typeof error === 'object' && error !== null) {
                    const e = error as Partial<{ responseStatusCode: number; statusCode: number; status: number }>;
                    status = e.responseStatusCode ?? e.statusCode ?? e.status;
                }
                if (status === 404) {
                    // No family registered for this user
                    return undefined;
                }
                throw error;
            }
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};