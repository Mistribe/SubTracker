import { FamilyHeader } from "@/components/families/FamilyHeader";
import { FamilyMembersTable } from "@/components/families/FamilyMembersTable";
import { EmptyFamiliesState } from "@/components/families/EmptyFamiliesState";
import { FamiliesLoadingError } from "@/components/families/FamiliesLoadingError";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";
import { PageHeader } from "@/components/ui/page-header";
import { useFamilyQuotaQuery } from "@/hooks/families/useFamilyQuotaQuery.ts";
import { QuotaButton } from "@/components/quotas/QuotaButton";
import { FeatureId } from "@/models/billing.ts";
import { AddFamilyMemberDialog } from "@/components/families/AddFamilyMemberDialog";

const FamilyPage = () => {
  const { data: family, isLoading, error } = useFamilyQuery();

  // Determine if user has a family first
  const hasFamily = !!family;

  return (
    <div className="container mx-auto py-6">
      <PageHeader
        title="Family"
        description="Manage your family"
        quotaButton={
          hasFamily ? (
            <QuotaButton
              useQuotaQuery={() => useFamilyQuotaQuery(hasFamily)}
              featureIds={[FeatureId.FamilyMembersCount]}
              featureLabels={{
                [FeatureId.FamilyMembersCount]: "Family Members"
              }}
            />
          ) : undefined
        }
        actionButton={
          hasFamily ? (
            <AddFamilyMemberDialog familyId={family!.id} />
          ) : undefined
        }
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
