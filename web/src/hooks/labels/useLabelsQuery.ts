import {useQuery} from "@tanstack/react-query";
import {useApiClient} from "@/hooks/use-api-client";
import type {LabelModel} from "@/api/models";
import Label from "@/models/label";
import {OwnerType} from "@/models/ownerType.ts";
import type {LabelsRequestBuilderGetQueryParameters} from "@/api/labels";

interface LabelsQueryOptions {
    ownerTypes?: OwnerType[];
    familyId?: string;
    offset?: number;
    limit?: number;
}

export const useLabelsQuery = (options: LabelsQueryOptions = {}) => {
    const {
        ownerTypes,
        familyId,
        offset = 0,
        limit = 10
    } = options;

    const {apiClient} = useApiClient();

    return useQuery({
        queryKey: ['labels', ownerTypes, familyId, offset, limit],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }

            const queryParameters: LabelsRequestBuilderGetQueryParameters = {
                limit: limit,
                offset: offset,
            };

            if (ownerTypes && ownerTypes.length > 0) {
                queryParameters.ownerType = ownerTypes;
            }

            if (familyId) {
                queryParameters.familyId = familyId;
            }

            const result = await apiClient?.labels.get({
                queryParameters
            });

            if (result && result.data) {
                return {
                    labels: result.data.map((model: LabelModel) => {
                        return Label.fromModel(model);
                    }),
                    length: result.data.length,
                    total: result.total ?? 0
                }
            }
            return {labels: [], length: 0, total: 0};
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });
};