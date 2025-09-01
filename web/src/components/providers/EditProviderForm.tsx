import {useState} from "react";
import {useProvidersMutations} from "@/hooks/providers/useProvidersMutations";
import {useFamilyQuery} from "@/hooks/families/useFamilyQuery.ts";
import {Button} from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import {Alert, AlertDescription} from "@/components/ui/alert";
import {AlertCircle} from "lucide-react";
import Provider from "@/models/provider";
import {ProviderDetailsForm, type ProviderFormValues} from "./ProviderDetailsForm";
import {OwnerType} from "@/models/ownerType.ts";

interface EditProviderFormProps {
    isOpen: boolean;
    onClose: () => void;
    provider: Provider;
}

export function EditProviderForm({isOpen, onClose, provider}: EditProviderFormProps) {
    const [error, setError] = useState<string | null>(null);
    const {updateProviderMutation} = useProvidersMutations();
    const {data: familyData} = useFamilyQuery();
    const families = familyData ? [familyData] : [];
    const [isSubmitting, setIsSubmitting] = useState(false);

    // Handle provider form submission
    const handleProviderSubmit = async (data: ProviderFormValues) => {
        try {
            setError(null);
            setIsSubmitting(true);
            await updateProviderMutation.mutateAsync({
                id: provider.id,
                etag: provider.etag,
                name: data.name,
                description: data.description || undefined,
                url: data.url || undefined,
                iconUrl: data.iconUrl || undefined,
                pricingPageUrl: data.pricingPageUrl || undefined,
                labels: provider.labels, // Preserve existing labels
                ownerType: provider.owner.type, // Always use the original owner type
                familyId: provider.owner.type === OwnerType.Family ? data.familyId : undefined,
            });
            onClose();
        } catch (err) {
            setError("Failed to update provider. Please try again.");
            console.error(err);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setError(null);
        onClose();
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleClose}>
            <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Edit Provider</DialogTitle>
                    <DialogDescription>
                        Update provider details.
                    </DialogDescription>
                </DialogHeader>

                {error && (
                    <Alert variant="destructive">
                        <AlertCircle className="h-4 w-4"/>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}

                <div className="py-4">
                    <ProviderDetailsForm
                        provider={provider}
                        families={families}
                        onSubmit={handleProviderSubmit}
                        isSubmitting={isSubmitting}
                    />
                </div>

                <DialogFooter>
                    <Button type="button" onClick={handleClose}>
                        Close
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}