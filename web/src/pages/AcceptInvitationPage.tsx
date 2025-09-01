import { useMemo } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, CheckCircle2, AlertTriangle } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";

const AcceptInvitationPage = () => {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();
  const familiesQuery = useFamilyQuery();

  const familyId = params.get("fid") || "";
  const memberId = params.get("mid") || "";
  const code = params.get("c") || "";

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

  const leaveMutation = useMutation({
    mutationFn: async () => {
      if (!apiClient) throw new Error("API client not initialized");
      const currentFamily = familiesQuery.data;
      if (!currentFamily) throw new Error("No family to leave");
      const you = currentFamily.members.find(m => m.isYou);
      if (!you) throw new Error("Could not determine your member entry");

      // If you are the owner and the only member, deleting the family is equivalent to leaving
      if (currentFamily.isOwner && currentFamily.members.length === 1) {
        return apiClient.families.byFamilyId(currentFamily.id).delete();
      }
      // Otherwise remove yourself from members
      return apiClient.families
        .byFamilyId(currentFamily.id)
        .members
        .byFamilyMemberId(you.id)
        .delete();
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["families"] });
    },
  });

  const hasFamily = !!familiesQuery.data;
  const currentFamily = familiesQuery.data || undefined;
  const isInTargetFamily = hasFamily && currentFamily?.id === familyId;

  return (
    <div className="container mx-auto py-10">
      <div className="mx-auto max-w-md">
        <Card>
          <CardHeader>
            <CardTitle>Accept Family Invitation</CardTitle>
            <CardDescription>
              Join a family to share subscriptions, labels and insights with people you trust.
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

            {isParamsValid && (
              <div className="space-y-3 text-sm text-muted-foreground">
                <p className="text-foreground font-medium">What is a family?</p>
                <ul className="list-disc pl-5 space-y-1">
                  <li>Share your subscriptions list with your household.</li>
                  <li>Collaborate on labels and keep everything organized.</li>
                  <li>Get a shared view of costs and renewal dates.</li>
                </ul>
              </div>
            )}

            {familiesQuery.isLoading && (
              <div className="flex items-center gap-2 text-sm">
                <Loader2 className="h-4 w-4 animate-spin" />
                Checking your current family...
              </div>
            )}

            {isInTargetFamily && (
              <div className="flex items-start gap-3 rounded-md border border-emerald-300 bg-emerald-50 p-3 text-emerald-900">
                <CheckCircle2 className="h-5 w-5 mt-0.5" />
                <div className="space-y-1">
                  <div className="font-medium">You're already a member of this family</div>
                  <div className="text-sm">
                    {currentFamily?.name ? (
                      <>You already belong to "{currentFamily.name}". No action is required.</>
                    ) : (
                      <>You already belong to this family. No action is required.</>
                    )}
                  </div>
                </div>
              </div>
            )}

            {hasFamily && !isInTargetFamily && (
              <div className="flex items-start gap-3 rounded-md border border-amber-300 bg-amber-50 p-3 text-amber-900">
                <AlertTriangle className="h-5 w-5 mt-0.5" />
                <div className="space-y-1">
                  <div className="font-medium">You're already in a family</div>
                  <div className="text-sm">
                    {currentFamily?.name ? (
                      <>You are currently a member of "{currentFamily.name}". To join this new family, you need to leave your current family first.</>
                    ) : (
                      <>To join this new family, you need to leave your current family first.</>
                    )}
                  </div>
                </div>
              </div>
            )}

            {acceptMutation.isPending && (
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
              <div className="flex items-start gap-3 rounded-md border border-red-300 bg-red-50 p-3 text-red-900">
                <AlertTriangle className="h-5 w-5 mt-0.5" />
                <div>
                  <div className="font-medium">Could not accept the invitation</div>
                  <div className="text-sm">The invitation may be invalid or expired. Please try again or request a new link.</div>
                </div>
              </div>
            )}

            <div className="flex flex-wrap gap-2 pt-2">
              <Button variant="outline" onClick={() => navigate("/family")}>
                {isInTargetFamily ? "Go to family" : "Decline"}
              </Button>
              {!isInTargetFamily && (
                <Button onClick={() => acceptMutation.mutate()} disabled={!isParamsValid || hasFamily || acceptMutation.isPending}>
                  {acceptMutation.isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  Accept invitation
                </Button>
              )}
              {hasFamily && !isInTargetFamily && (
                <Button variant="destructive" onClick={() => leaveMutation.mutate()} disabled={leaveMutation.isPending}>
                  {leaveMutation.isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  Leave current family
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default AcceptInvitationPage;
