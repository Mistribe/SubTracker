import { FamilyHeader } from "@/components/families/FamilyHeader";
import { FamilyMembersTable } from "@/components/families/FamilyMembersTable";
import { EmptyFamiliesState } from "@/components/families/EmptyFamiliesState";
import { FamiliesLoadingError } from "@/components/families/FamiliesLoadingError";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";
import { PageHeader } from "@/components/ui/page-header";
import { useFamilyQuotaQuery } from "@/hooks/families/useFamilyQuotaQuery.ts";
import { QuotaUsage } from "@/components/quotas/QuotaUsage";
import { QuotaUsageSkeleton } from "@/components/quotas/QuotaUsageSkeleton";
import { FeatureId } from "@/models/billing.ts";

const FamilyPage = () => {
  const { data: family, isLoading, error } = useFamilyQuery();

  // Determine if user has a family first
  const hasFamily = !!family;

  // Fetch family quota only if a family exists
  const { data: familyQuotaData, isLoading: isFamilyQuotaLoading, error: familyQuotaError } = useFamilyQuotaQuery(hasFamily);
  const membersCountQuota = hasFamily ? familyQuotaData?.find(q => q.feature === FeatureId.FamilyMembersCount) : undefined;

  const quotaSection = hasFamily && (
    <div className="mt-4 max-w-xs">
      {isFamilyQuotaLoading && <QuotaUsageSkeleton />}
      {!isFamilyQuotaLoading && membersCountQuota && (
        <QuotaUsage quota={membersCountQuota} label="Family Members" />
      )}
      {!isFamilyQuotaLoading && !membersCountQuota && !familyQuotaError && (
        <div className="text-xs text-muted-foreground border rounded-md p-3">No quota data available.</div>
      )}
      {familyQuotaError && (
        <div className="text-xs text-destructive border border-destructive/30 rounded-md p-3">Failed to load quota.</div>
      )}
    </div>
  );

  return (
    <div className="container mx-auto py-6">
      <PageHeader
        title="Family"
        description="Manage your family"
      />
      {quotaSection}
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
