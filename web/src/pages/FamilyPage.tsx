import { FamilyHeader } from "@/components/families/FamilyHeader";
import { FamilyMembersTable } from "@/components/families/FamilyMembersTable";
import { EmptyFamiliesState } from "@/components/families/EmptyFamiliesState";
import { FamiliesLoadingError } from "@/components/families/FamiliesLoadingError";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";
import { PageHeader } from "@/components/ui/page-header";

const FamilyPage = () => {
  const {
    data: family,
    isLoading,
    error,
  } = useFamilyQuery();

  const hasFamily = !!family;

  return (
    <div className="container mx-auto py-6">
      <PageHeader
        title="Family"
        description="Manage your family"
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
