import { CreateFamilyDialog } from "@/components/families/CreateFamilyDialog.tsx";
import { FamilyCard } from "@/components/families/FamilyCard";
import { EmptyFamiliesState } from "@/components/families/EmptyFamiliesState";
import { FamiliesLoadingError } from "@/components/families/FamiliesLoadingError";
import { useFamiliesQuery } from "@/hooks/families/useFamiliesQuery";
import { PageHeader } from "@/components/ui/page-header.tsx";

const FamilyPage = () => {
  const {
    data: queryResponse,
    isLoading,
    error,
  } = useFamiliesQuery();

  if (isLoading || error) {
    return <FamiliesLoadingError isLoading={isLoading} error={error} />;
  }

  const families = queryResponse?.families || [];
  const hasFamily = families.length > 0;
  const family = hasFamily ? families[0] : undefined;

  return (
    <div className="container mx-auto py-6">
      <PageHeader
        title="Family"
        description="Manage your family"
        actionButton={!hasFamily ? <CreateFamilyDialog /> : undefined}
      />

      {!hasFamily ? (
        <EmptyFamiliesState />
      ) : (
        <div className="grid gap-6">
          {family && <FamilyCard key={family.id} family={family} />}
        </div>
      )}
    </div>
  );
};

export default FamilyPage;
