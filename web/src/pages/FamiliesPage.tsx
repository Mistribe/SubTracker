import {useState} from "react";
import {CreateFamilyDialog} from "@/components/families/CreateFamilyDialog.tsx";
import {FamilyCard} from "@/components/families/FamilyCard";
import {EmptyFamiliesState} from "@/components/families/EmptyFamiliesState";
import {FamiliesLoadingError} from "@/components/families/FamiliesLoadingError";
import {useFamiliesQuery} from "@/hooks/families/useFamiliesQuery";
import {PageHeader} from "@/components/ui/page-header.tsx";

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
        <div className="container mx-auto py-6">
            <PageHeader
                title="Families"
                description="Manage your families"
                actionButton={families.length > 0 && families.filter(x => x.isOwner).length === 0 ? <CreateFamilyDialog/> : undefined}
            />

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