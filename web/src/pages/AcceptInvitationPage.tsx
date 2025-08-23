import { useEffect, useMemo } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, CheckCircle2, AlertTriangle } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";

const AcceptInvitationPage = () => {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { apiClient } = useApiClient();

  const familyId = params.get("familyId") || "";
  const memberId = params.get("memberId") || "";
  const code = params.get("code") || "";

  const isParamsValid = useMemo(() => !!familyId && !!code, [familyId, code]);

  const acceptMutation = useMutation({
    mutationFn: async () => {
      if (!apiClient) throw new Error("API client not initialized");
      // memberId is optional from the API model, but we pass it if present
      return apiClient.families
        .byFamilyId(familyId)
        .accept
        .post({ invitationCode: code, familyMemberId: memberId || undefined });
    },
    onSuccess: () => {
      // Navigate to the Family page after a short delay for UX
      setTimeout(() => navigate("/family"), 500);
    },
  });

  useEffect(() => {
    if (apiClient && isParamsValid && !acceptMutation.isPending && !acceptMutation.isSuccess) {
      acceptMutation.mutate();
    }
  }, [apiClient, isParamsValid, acceptMutation]);

  return (
    <div className="container mx-auto py-10">
      <div className="mx-auto max-w-md">
        <Card>
          <CardHeader>
            <CardTitle>Accept Family Invitation</CardTitle>
            <CardDescription>
              Join the family by validating your invitation.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {!isParamsValid && (
              <div className="flex items-start gap-3 rounded-md border border-amber-300 bg-amber-50 p-3 text-amber-900">
                <AlertTriangle className="h-5 w-5 mt-0.5" />
                <div>
                  <div className="font-medium">Missing information</div>
                  <div className="text-sm">The invitation link is invalid. Please ask the owner for a new link.</div>
                </div>
              </div>
            )}

            {isParamsValid && acceptMutation.isIdle && (
              <div className="text-sm text-muted-foreground">
                Preparing acceptance...
              </div>
            )}

            {isParamsValid && acceptMutation.isPending && (
              <div className="flex items-center gap-2 text-sm">
                <Loader2 className="h-4 w-4 animate-spin" />
                Accepting your invitation...
              </div>
            )}

            {acceptMutation.isSuccess && (
              <div className="flex items-start gap-3 rounded-md border border-emerald-300 bg-emerald-50 p-3 text-emerald-900">
                <CheckCircle2 className="h-5 w-5 mt-0.5" />
                <div>
                  <div className="font-medium">You're in!</div>
                  <div className="text-sm">Redirecting to your family...</div>
                </div>
              </div>
            )}

            {acceptMutation.isError && (
              <div className="flex flex-col gap-3">
                <div className="flex items-start gap-3 rounded-md border border-red-300 bg-red-50 p-3 text-red-900">
                  <AlertTriangle className="h-5 w-5 mt-0.5" />
                  <div>
                    <div className="font-medium">Could not accept the invitation</div>
                    <div className="text-sm">The invitation may be invalid or expired. Please try again or request a new link.</div>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button variant="outline" onClick={() => navigate("/family")}>Go to Family</Button>
                  <Button onClick={() => acceptMutation.mutate()} disabled={acceptMutation.isPending}>
                    {acceptMutation.isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                    Retry
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default AcceptInvitationPage;
