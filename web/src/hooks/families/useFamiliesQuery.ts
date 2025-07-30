import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {FamilyModel} from "@/api/models";
import Family from "@/models/family";

interface FamiliesQueryOptions {
    offset?: number;
    limit?: number;
}

export const useFamiliesQuery = (options: FamiliesQueryOptions = {}) => {
    const {
        offset = 0,
        limit = 10
    } = options;
    
    const {apiClient} = useApiClient();

    return useQuery({
        queryKey: ['families', offset, limit],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }
            
            const result = await apiClient.families.get({
                queryParameters: {
                    offset,
                    limit
                }
            });
            
            if (result && result.data) {
                return {
                    families: result.data.map((model: FamilyModel) => {
                        return Family.fromModel(model);
                    }),
                    length: result.data.length,
                    total: result.total ?? 0
                }
            }
            return {families: [], length: 0, total: 0};
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};