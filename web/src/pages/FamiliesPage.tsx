import {useState} from "react";
import {CreateFamilyDialog} from "@/components/families/CreateFamilyDialog.tsx";
import {FamilyCard} from "@/components/families/FamilyCard";
import {EmptyFamiliesState} from "@/components/families/EmptyFamiliesState";
import {FamiliesLoadingError} from "@/components/families/FamiliesLoadingError";
import {useFamiliesQuery} from "@/hooks/families/useFamiliesQuery";

const FamiliesPage = () => {
    const [offset] = useState(0);
    const [limit] = useState(10);

    const {
        data: queryResponse,
        isLoading,
        error
    } = useFamiliesQuery({
        offset,
        limit
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