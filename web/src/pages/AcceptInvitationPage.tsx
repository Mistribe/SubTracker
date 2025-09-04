import { useMemo } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2, CheckCircle2, AlertTriangle } from "lucide-react";
import { useFamilyQuery } from "@/hooks/families/useFamilyQuery.ts";
import { useFamiliesMutations } from "@/hooks/families/useFamiliesMutations";

const AcceptInvitationPage = () => {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const familiesQuery = useFamilyQuery();

  const familyId = params.get("fid") || "";
  const memberId = params.get("mui") || "";
  const code = params.get("c") || "";

  const isParamsValid = useMemo(() => !!familyId && !!code, [familyId, code]);

  const { acceptFamilyInvitationMutation, leaveFamilyMutation } = useFamiliesMutations();


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

            {acceptFamilyInvitationMutation.isPending && (
              <div className="flex items-center gap-2 text-sm">
                <Loader2 className="h-4 w-4 animate-spin" />
                Accepting your invitation...
              </div>
            )}

            {acceptFamilyInvitationMutation.isSuccess && (
              <div className="flex items-start gap-3 rounded-md border border-emerald-300 bg-emerald-50 p-3 text-emerald-900">
                <CheckCircle2 className="h-5 w-5 mt-0.5" />
                <div>
                  <div className="font-medium">You're in!</div>
                  <div className="text-sm">Redirecting to your family...</div>
                </div>
              </div>
            )}

            {acceptFamilyInvitationMutation.isError && (
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
                <Button
                  onClick={() =>
                    acceptFamilyInvitationMutation.mutate(
                      { familyId, invitationCode: code, familyMemberId: memberId || undefined },
                      {
                        onSuccess: () => {
                          setTimeout(() => navigate("/family"), 500);
                        },
                      }
                    )
                  }
                  disabled={!isParamsValid || hasFamily || acceptFamilyInvitationMutation.isPending}
                >
                  {acceptFamilyInvitationMutation.isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  Accept invitation
                </Button>
              )}
              {hasFamily && !isInTargetFamily && (
                <Button variant="destructive" onClick={() => leaveFamilyMutation.mutate()} disabled={leaveFamilyMutation.isPending}>
                  {leaveFamilyMutation.isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
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
