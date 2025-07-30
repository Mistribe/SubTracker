import {useState} from "react";
import Family from "@/models/family.ts";
import {useApiClient} from "@/hooks/use-api-client.ts";
import {useQuery} from "@tanstack/react-query";
import type {FamilyModel} from "@/api/models";
import {CreateFamilyDialog} from "@/components/families/CreateFamilyDialog.tsx";
import {FamilyCard} from "@/components/families/FamilyCard";
import {EmptyFamiliesState} from "@/components/families/EmptyFamiliesState";
import {FamiliesLoadingError} from "@/components/families/FamiliesLoadingError";

const FamiliesPage = () => {
    const {apiClient} = useApiClient();

    const [offset] = useState(0);
    const [limit] = useState(10);

    const {
        data: queryResponse,
        isLoading,
        error
    } = useQuery({
        queryKey: ['families'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }
            const result = await apiClient.families.get({
                queryParameters: {
                    offset: offset,
                    limit: limit
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

    // Show loading or error states
    if (isLoading || error) {
        return <FamiliesLoadingError isLoading={isLoading} error={error}/>;
    }

    const families = queryResponse?.families || [];

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Families</h1>
                {families.length > 0 && families.filter(x => x.isOwner).length === 0 && <CreateFamilyDialog/>}
            </div>

            {families.length === 0 ? (
                <EmptyFamiliesState/>
            ) : (
                <div className="grid gap-6">
                    {families.map((family) => (
                        <FamilyCard key={family.id} family={family}/>
                    ))}
                </div>
            )}
        </div>
    );
};

export default FamiliesPage;