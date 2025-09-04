import {useState} from "react";
import {TableCell, TableRow} from "@/components/ui/table";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Checkbox} from "@/components/ui/checkbox";
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip";
import {CheckIcon, Loader2, UserCheck, UserX, XIcon} from "lucide-react";
import FamilyMember from "@/models/familyMember.ts";
import {useFamiliesMutations} from "@/hooks/families/useFamiliesMutations";
import {FamilyMemberType} from "@/models/familyMemberType.ts";
import {FamilyMemberActionsMenu} from "./FamilyMemberActionsMenu";
import {FamilyMemberInviteDialog} from "./FamilyMemberInviteDialog";
import {FamilyMemberUnlinkAlert} from "./FamilyMemberUnlinkAlert";
import {twJoin} from "tailwind-merge";

interface FamilyMemberRowProps {
    member: FamilyMember,
    familyId: string,
    isOwner: boolean,
    className?: string
}

export const FamilyMemberRow = ({member, familyId, isOwner, className}: FamilyMemberRowProps) => {
    const {updateFamilyMemberMutation, removeFamilyMemberMutation} = useFamiliesMutations();

    const [editingMemberId, setEditingMemberId] = useState<string | null>(null);
    const [editedMemberName, setEditedMemberName] = useState<string>("");
    const [editedMemberIsKid, setEditedMemberIsKid] = useState<boolean>(false);

    // Invite dialog state
    const [inviteOpen, setInviteOpen] = useState(false);
    const [unlinkOpen, setUnlinkOpen] = useState(false);

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
        <TableRow key={member.id} className={twJoin(className)}>
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
                        <FamilyMemberActionsMenu
                            isOwner={isOwner}
                            isSelf={member.isYou && isOwner}
                            hasAccount={member.hasAccount}
                            onEdit={() => startEditingMember(member)}
                            onInvite={() => setInviteOpen(true)}
                            onUnlink={() => setUnlinkOpen(true)}
                            onRemove={() => removeMember(familyId, member, isOwner)}
                            isRemoving={removeFamilyMemberMutation.isPending}
                        />

                        <FamilyMemberInviteDialog
                            open={inviteOpen}
                            onOpenChange={setInviteOpen}
                            familyId={familyId}
                            member={member}
                        />

                        <FamilyMemberUnlinkAlert
                            open={unlinkOpen}
                            onOpenChange={setUnlinkOpen}
                            familyId={familyId}
                            member={member}
                            isOwner={isOwner}
                        />
                    </>
                )}
            </TableCell>
        </TableRow>
    );
};