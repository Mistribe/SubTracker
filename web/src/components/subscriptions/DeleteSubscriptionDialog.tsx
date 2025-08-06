import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import Subscription from "@/models/subscription";

interface DeleteSubscriptionDialogProps {
    subscription: Subscription;
    isOpen: boolean;
    isSubmitting: boolean;
    onClose: () => void;
    onConfirm: () => void;
}

export function DeleteSubscriptionDialog({
    subscription,
    isOpen,
    isSubmitting,
    onClose,
    onConfirm
}: DeleteSubscriptionDialogProps) {
    return (
        <AlertDialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Delete Subscription</AlertDialogTitle>
                    <AlertDialogDescription>
                        Are you sure you want to delete the subscription 
                        "{subscription.friendlyName || `Subscription ${subscription.id.substring(0, 8)}`}"? 
                        This action cannot be undone.
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                        onClick={onConfirm}
                        disabled={isSubmitting}
                    >
                        Delete
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
}