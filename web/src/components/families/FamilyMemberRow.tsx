import {useMemo, useState} from "react";
import {TableCell, TableRow} from "@/components/ui/table";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Checkbox} from "@/components/ui/checkbox";
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form.tsx";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import {CheckIcon, Loader2, MoreVertical, UserCheck, UserX, XIcon} from "lucide-react";
import FamilyMember from "@/models/familyMember.ts";
import {useFamiliesMutations} from "@/hooks/families/useFamiliesMutations";
import {FamilyMemberType} from "@/models/familyMemberType.ts";

interface FamilyMemberRowProps {
    member: FamilyMember;
    familyId: string;
    isOwner: boolean;
}

export const FamilyMemberRow = ({member, familyId, isOwner}: FamilyMemberRowProps) => {
    const {updateFamilyMemberMutation, removeFamilyMemberMutation, inviteFamilyMemberMutation, revokeFamilyMemberMutation} = useFamiliesMutations();

    const [editingMemberId, setEditingMemberId] = useState<string | null>(null);
    const [editedMemberName, setEditedMemberName] = useState<string>("");
    const [editedMemberIsKid, setEditedMemberIsKid] = useState<boolean>(false);

    // Invite dialog state
    const [inviteOpen, setInviteOpen] = useState(false);
    const [inviteCode, setInviteCode] = useState<string | null>(null);
    const [unlinkOpen, setUnlinkOpen] = useState(false);

    const inviteLink = useMemo(() => {
        if (!inviteCode) return "";
        try {
            const origin = window.location.origin;
            const url = new URL(`${origin}/invite/accept`);
            url.searchParams.set("familyId", familyId);
            url.searchParams.set("memberId", member.id);
            url.searchParams.set("code", inviteCode);
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
        inviteFamilyMemberMutation.mutate({
            familyId,
            familyMemberId: member.id,
            email: values.email,
        }, {
            onSuccess: (resp?: { code?: string | null }) => {
                const code = resp?.code ?? null;
                setInviteCode(code);
            },
        });
    };

    const onGenerateLink = () => {
        inviteFamilyMemberMutation.mutate({
            familyId,
            familyMemberId: member.id,
        }, {
            onSuccess: (resp?: { code?: string | null }) => {
                const code = resp?.code ?? null;
                setInviteCode(code);
            },
        });
    };

    // Function to start editing a family member
    const startEditingMember = (member: FamilyMember) => {
        setEditingMemberId(member.id);
        setEditedMemberName(member.name);
        setEditedMemberIsKid(member.type === FamilyMemberType.Kid);
    };

    // Function to cancel editing a family member
    const cancelEditingMember = () => {
        setEditingMemberId(null);
        setEditedMemberName("");
        setEditedMemberIsKid(false);
    };

    // Function to save family member changes
    const saveMemberChanges = (familyId: string, member: FamilyMember) => {
        // Only update if values have changed
        if (editedMemberName === member.name && editedMemberIsKid === member.isKid) {
            cancelEditingMember();
            return;
        }

        updateFamilyMemberMutation.mutate({
            familyId,
            memberId: member.id,
            name: editedMemberName,
            isKid: editedMemberIsKid
        }, {
            onSuccess: () => {
                // Reset editing state
                cancelEditingMember();
            },
            onError: (error) => {
                console.error("Failed to update family member:", error);
            }
        });
    };

    // Function to remove a family member
    const removeMember = (familyId: string, member: FamilyMember, isOwner: boolean) => {
        // Prevent removing yourself if you are the owner of the family
        if (member.isYou && isOwner) {
            console.error("Cannot remove yourself as the owner of the family");
            return;
        }

        removeFamilyMemberMutation.mutate({
            familyId,
            memberId: member.id
        }, {
            onError: (error) => {
                console.error("Failed to remove family member:", error);
            }
        });
    };

    // Function to revoke (unlink) a linked account of a family member
    const revokeMember = (familyId: string, member: FamilyMember, isOwner: boolean) => {
        if (!isOwner) {
            return;
        }
        if (member.isYou && isOwner) {
            console.error("Cannot revoke your own account as the owner of the family");
            return;
        }

        revokeFamilyMemberMutation.mutate({
            familyId,
            memberId: member.id
        }, {
            onSuccess: () => {
                setUnlinkOpen(false);
            },
            onError: (error) => {
                console.error("Failed to revoke family member:", error);
            }
        });
    };
    
    const formatMemberType = (memberType: FamilyMemberType): string => {
        switch (memberType) {
            case FamilyMemberType.Adult:
                return "Adult";
            case FamilyMemberType.Owner:
                return "Owner";
            default:
                return "Kid";
        }
    }

    return (
        <TableRow key={member.id}>
            <TableCell className="font-medium">
                {editingMemberId === member.id ? (
                    <Input
                        value={editedMemberName}
                        onChange={(e) => setEditedMemberName(e.target.value)}
                        className="w-full"
                        placeholder="Member name"
                    />
                ) : (
                    <div className="flex items-center gap-2">
                        <span>
                            {member.name} {member.isYou && <i>(You)</i>}
                        </span>
                        <TooltipProvider delayDuration={200}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <span
                                        className="inline-flex h-6 w-6 items-center justify-center rounded-full hover:bg-muted/50"
                                        aria-label={member.hasAccount ? "Account linked" : "No account linked"}
                                    >
                                        {member.hasAccount ? (
                                            <UserCheck className="h-3.5 w-3.5 text-emerald-600"/>
                                        ) : (
                                            <UserX className="h-3.5 w-3.5 text-muted-foreground"/>
                                        )}
                                    </span>
                                </TooltipTrigger>
                                <TooltipContent side="top">
                                    <p className="text-xs">{member.hasAccount ? "Account linked" : "No account linked"}</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    </div>
                )}
            </TableCell>
            <TableCell>
                {editingMemberId === member.id && member.type !== FamilyMemberType.Owner ? (
                    <div className="flex items-center space-x-2">
                        <Checkbox
                            checked={editedMemberIsKid}
                            onCheckedChange={(checked) => setEditedMemberIsKid(checked === true)}
                            id={`kid-checkbox-${member.id}`}
                        />
                        <label htmlFor={`kid-checkbox-${member.id}`}>Kid</label>
                    </div>
                ) : (
                    formatMemberType(member.type)
                )}
            </TableCell>
            <TableCell>
                {editingMemberId === member.id ? (
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => saveMemberChanges(familyId, member)}
                            disabled={updateFamilyMemberMutation.isPending}
                        >
                            {updateFamilyMemberMutation.isPending ? (
                                <Loader2 className="h-4 w-4 animate-spin mr-2"/>
                            ) : (
                                <CheckIcon className="h-4 w-4 mr-2"/>
                            )}
                            Save
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={cancelEditingMember}
                            disabled={updateFamilyMemberMutation.isPending}
                        >
                            <XIcon className="h-4 w-4 mr-2"/>
                            Cancel
                        </Button>
                    </div>
                ) : (
                    <>
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon" className="h-8 w-8">
                                    <MoreVertical className="h-4 w-4" />
                                    <span className="sr-only">Open actions</span>
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                                <DropdownMenuItem onClick={() => startEditingMember(member)}>
                                    Edit
                                </DropdownMenuItem>
                                {isOwner && !member.hasAccount && !member.isYou && (
                                    <DropdownMenuItem onClick={() => setInviteOpen(true)}>
                                        Invite
                                    </DropdownMenuItem>
                                )}
                                {isOwner && member.hasAccount && !member.isYou && (
                                    <DropdownMenuItem
                                        variant="destructive"
                                        onClick={() => setUnlinkOpen(true)}
                                        disabled={revokeFamilyMemberMutation.isPending}
                                    >
                                        <UserX className="h-4 w-4" />
                                        Unlink account
                                    </DropdownMenuItem>
                                )}
                                {/* Don't show remove button for yourself if you're the owner */}
                                {!(member.isYou && isOwner) && (
                                    <>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem
                                            variant="destructive"
                                            onClick={() => removeMember(familyId, member, isOwner)}
                                            disabled={removeFamilyMemberMutation.isPending}
                                        >
                                            {removeFamilyMemberMutation.isPending ? (
                                                <>
                                                    <Loader2 className="h-4 w-4 animate-spin" />
                                                    Removing...
                                                </>
                                            ) : (
                                                <>
                                                    <XIcon className="h-4 w-4" />
                                                    Remove
                                                </>
                                            )}
                                        </DropdownMenuItem>
                                    </>
                                )}
                            </DropdownMenuContent>
                        </DropdownMenu>

                        <Dialog open={inviteOpen} onOpenChange={(o) => {
                            setInviteOpen(o);
                            if (!o) { setInviteCode(null); form.reset(); }
                        }}>
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
                                    <Button variant="ghost" onClick={() => setInviteOpen(false)}>Close</Button>
                                </DialogFooter>
                            </DialogContent>
                        </Dialog>

                        <AlertDialog open={unlinkOpen} onOpenChange={setUnlinkOpen}>
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
                                        onClick={() => revokeMember(familyId, member, isOwner)}
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
                    </>
                )}
            </TableCell>
        </TableRow>
    );
};