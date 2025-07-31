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
import Price from "@/models/price";

interface DeletePriceDialogProps {
    price: Price | null;
    isOpen: boolean;
    isSubmitting: boolean;
    onClose: () => void;
    onConfirm: () => void;
}

export function DeletePriceDialog({
    price,
    isOpen,
    isSubmitting,
    onClose,
    onConfirm
}: DeletePriceDialogProps) {
    return (
        <AlertDialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Delete Price</AlertDialogTitle>
                    <AlertDialogDescription>
                        {price ? (
                            <>
                                Are you sure you want to delete the price of {price.amount} {price.currency}
                                {price.endDate 
                                    ? ` (${price.startDate.toLocaleDateString()} - ${price.endDate.toLocaleDateString()})`
                                    : ` (from ${price.startDate.toLocaleDateString()})`}? 
                                This action cannot be undone.
                            </>
                        ) : (
                            <>Are you sure you want to delete this price? This action cannot be undone.</>
                        )}
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