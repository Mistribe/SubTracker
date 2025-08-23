import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Loader2, MoreVertical, UserX, XIcon } from "lucide-react";

interface FamilyMemberActionsMenuProps {
  isOwner: boolean;
  isSelf: boolean;
  hasAccount: boolean;
  onEdit: () => void;
  onInvite: () => void;
  onUnlink: () => void;
  onRemove: () => void;
  isRemoving: boolean;
}

export function FamilyMemberActionsMenu({
  isOwner,
  isSelf,
  hasAccount,
  onEdit,
  onInvite,
  onUnlink,
  onRemove,
  isRemoving,
}: FamilyMemberActionsMenuProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="h-8 w-8">
          <MoreVertical className="h-4 w-4" />
          <span className="sr-only">Open actions</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={onEdit}>Edit</DropdownMenuItem>
        {isOwner && !hasAccount && !isSelf && (
          <DropdownMenuItem onClick={onInvite}>Invite</DropdownMenuItem>
        )}
        {isOwner && hasAccount && !isSelf && (
          <DropdownMenuItem variant="destructive" onClick={onUnlink}>
            <UserX className="h-4 w-4" />
            Unlink account
          </DropdownMenuItem>
        )}
        {!isSelf && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              variant="destructive"
              onClick={onRemove}
              disabled={isRemoving}
            >
              {isRemoving ? (
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
  );
}
