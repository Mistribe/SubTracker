import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog";
import { useFamiliesMutations } from "@/hooks/families/useFamiliesMutations";
import FamilyMember from "@/models/familyMember";
import { Loader2 } from "lucide-react";

interface FamilyMemberUnlinkAlertProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  familyId: string;
  member: FamilyMember;
  isOwner: boolean;
}

export function FamilyMemberUnlinkAlert({ open, onOpenChange, familyId, member, isOwner }: FamilyMemberUnlinkAlertProps) {
  const { revokeFamilyMemberMutation } = useFamiliesMutations();

  const onConfirm = () => {
    if (!isOwner) return;
    if (member.isYou && isOwner) return;

    revokeFamilyMemberMutation.mutate(
      {
        familyId,
        memberId: member.id,
      },
      {
        onSuccess: () => {
          onOpenChange(false);
        },
      },
    );
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Unlink account</AlertDialogTitle>
          <AlertDialogDescription>
            Are you sure you want to unlink {member.name}'s account? They will be able to re-link later via a new invitation.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={revokeFamilyMemberMutation.isPending}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            onClick={onConfirm}
            disabled={revokeFamilyMemberMutation.isPending}
          >
            {revokeFamilyMemberMutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Unlinking...
              </>
            ) : (
              "Unlink"
            )}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
