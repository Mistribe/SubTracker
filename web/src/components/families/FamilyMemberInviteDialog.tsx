import { useMemo, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Loader2 } from "lucide-react";
import FamilyMember from "@/models/familyMember";
import { useFamiliesMutations } from "@/hooks/families/useFamiliesMutations";

interface FamilyMemberInviteDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  familyId: string;
  member: FamilyMember;
}

export function FamilyMemberInviteDialog({ open, onOpenChange, familyId, member }: FamilyMemberInviteDialogProps) {
  const { inviteFamilyMemberMutation } = useFamiliesMutations();

  const [inviteCode, setInviteCode] = useState<string | null>(null);

  const inviteLink = useMemo(() => {
    if (!inviteCode) return "";
    try {
      const origin = window.location.origin;
      const url = new URL(`${origin}/invite/accept`);
      url.searchParams.set("fid", familyId);
      url.searchParams.set("mui", member.id);
      url.searchParams.set("c", inviteCode);
      return url.toString();
    } catch {
      return "";
    }
  }, [inviteCode, familyId, member.id]);

  const emailSchema = z.object({
    email: z.string().email("Please enter a valid email address"),
  });

  const form = useForm<z.infer<typeof emailSchema>>({
    resolver: zodResolver(emailSchema),
    defaultValues: { email: "" },
  });

  const onSendEmailInvite = (values: z.infer<typeof emailSchema>) => {
    inviteFamilyMemberMutation.mutate(
      {
        familyId,
        familyMemberId: member.id,
        email: values.email,
      },
      {
        onSuccess: (resp?: { code?: string | null }) => {
          const code = resp?.code ?? null;
          setInviteCode(code);
        },
      },
    );
  };

  const onGenerateLink = () => {
    inviteFamilyMemberMutation.mutate(
      {
        familyId,
        familyMemberId: member.id,
      },
      {
        onSuccess: (resp?: { code?: string | null }) => {
          const code = resp?.code ?? null;
          setInviteCode(code);
        },
      },
    );
  };

  const handleOpenChange = (o: boolean) => {
    onOpenChange(o);
    if (!o) {
      setInviteCode(null);
      form.reset();
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Invite {member.name}</DialogTitle>
          <DialogDescription>
            Invite this member by email or generate an invitation code to share.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSendEmailInvite)} className="space-y-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input placeholder="name@example.com" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className="flex gap-2">
                <Button type="submit" disabled={inviteFamilyMemberMutation.isPending}>
                  {inviteFamilyMemberMutation.isPending && (
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  )}
                  Send email invite
                </Button>
                <Button type="button" variant="outline" onClick={onGenerateLink} disabled={inviteFamilyMemberMutation.isPending}>
                  Generate link
                </Button>
              </div>
            </form>
          </Form>

          {inviteCode && (
            <div className="space-y-4">
              <div className="space-y-2">
                <div className="text-sm text-muted-foreground">Invitation code</div>
                <div className="flex gap-2">
                  <Input readOnly value={inviteCode} />
                  <Button type="button" variant="outline" onClick={() => navigator.clipboard.writeText(inviteCode!)}>
                    Copy
                  </Button>
                </div>
              </div>
              <div className="space-y-2">
                <div className="text-sm text-muted-foreground">Invitation link</div>
                <div className="flex gap-2">
                  <Input readOnly value={inviteLink} />
                  <Button type="button" variant="outline" onClick={() => inviteLink && navigator.clipboard.writeText(inviteLink)}>
                    Copy
                  </Button>
                </div>
              </div>
            </div>
          )}
        </div>
        <DialogFooter>
          <Button variant="ghost" onClick={() => handleOpenChange(false)}>Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
