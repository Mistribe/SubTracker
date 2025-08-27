import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import Family from "@/models/family";
import type {FamilyModel} from "@/api/models/family";

interface FamiliesQueryOptions {
    offset?: number;
    limit?: number;
}

export const useFamiliesQuery = (options: FamiliesQueryOptions = {}) => {
    // Keep options to maintain backward compatibility with callers, but they are unused now.
    void options; // mark as used to satisfy lint rules
    const {apiClient} = useApiClient();

    return useQuery({
        queryKey: ['families'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            // Fetch the current user's single family; backend now supports only one family
            const result: FamilyModel | undefined = await apiClient.families.me.get()

            if (result) {
                const family = Family.fromModel(result);
                return {
                    families: [family],
                    length: 1,
                    total: 1
                };
            }
            return {families: [], length: 0, total: 0};
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};