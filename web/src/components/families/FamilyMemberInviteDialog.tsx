import { useMemo, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
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
  const [copied, setCopied] = useState(false);
  const [method, setMethod] = useState<"email" | "link">("email");

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

  const handleCopy = async () => {
    if (!inviteLink) return;
    try {
      await navigator.clipboard.writeText(inviteLink);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch (e) {
      // ignore clipboard errors silently
    }
  };

  const handleOpenChange = (o: boolean) => {
    onOpenChange(o);
    if (!o) {
      setInviteCode(null);
      setCopied(false);
      setMethod("email");
      form.reset();
    }
  };

  const handleMethodChange = (val: string) => {
    if (val === method) return;
    setMethod(val as "email" | "link");
    setInviteCode(null);
    setCopied(false);
    form.reset();
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Invite {member.name}</DialogTitle>
          <DialogDescription>
            Choose how you want to invite this member.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <Tabs value={method} onValueChange={handleMethodChange}>
            <TabsList className="grid grid-cols-2 w-full">
              <TabsTrigger value="email">Email</TabsTrigger>
              <TabsTrigger value="link">Link</TabsTrigger>
            </TabsList>
            <TabsContent value="email" className="mt-4">
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
                  <div>
                    <Button type="submit" disabled={inviteFamilyMemberMutation.isPending}>
                      {inviteFamilyMemberMutation.isPending && (
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      )}
                      Send email invite
                    </Button>
                  </div>
                </form>
              </Form>
            </TabsContent>
            <TabsContent value="link" className="mt-4">
              <div className="flex gap-2">
                <Button type="button" onClick={onGenerateLink} disabled={inviteFamilyMemberMutation.isPending}>
                  {inviteFamilyMemberMutation.isPending && (
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  )}
                  Generate link
                </Button>
              </div>
            </TabsContent>
          </Tabs>

          {inviteCode && method === "link" && (
            <div className="space-y-4">
              <div className="space-y-2">
                <div className="text-sm text-muted-foreground">Invitation link</div>
                <div className="flex gap-2">
                  <Input readOnly value={inviteLink} />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={handleCopy}
                    className={copied ? "animate-pulse" : undefined}
                  >
                    {copied ? "Copied" : "Copy"}
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
