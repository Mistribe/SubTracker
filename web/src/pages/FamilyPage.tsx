import { CreateFamilyDialog } from "@/components/families/CreateFamilyDialog";
import { FamilyHeader } from "@/components/families/FamilyHeader";
import { FamilyMembersTable } from "@/components/families/FamilyMembersTable";
import { EmptyFamiliesState } from "@/components/families/EmptyFamiliesState";
import { FamiliesLoadingError } from "@/components/families/FamiliesLoadingError";
import { useFamiliesQuery } from "@/hooks/families/useFamiliesQuery";
import { PageHeader } from "@/components/ui/page-header";

const FamilyPage = () => {
  const {
    data: queryResponse,
    isLoading,
    error,
  } = useFamiliesQuery();

  const families = queryResponse?.families || [];
  const hasFamily = families.length > 0;
  const family = hasFamily ? families[0] : undefined;

  return (
    <div className="container mx-auto py-6">
      <PageHeader
        title="Family"
        description="Manage your family"
        actionButton={!isLoading && !error && !hasFamily ? <CreateFamilyDialog /> : undefined}
      />

      {isLoading || error ? (
        <FamiliesLoadingError isLoading={isLoading} error={error} />
      ) : !hasFamily ? (
        <EmptyFamiliesState />
      ) : (
        <div className="space-y-6">
          <FamilyHeader family={family!} />
          <FamilyMembersTable family={family!} />
        </div>
      )}
    </div>
  );
};

export default FamilyPage;
