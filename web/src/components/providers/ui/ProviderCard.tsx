import {useState} from "react";
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card";
import {Badge} from "@/components/ui/badge";
import {Button} from "@/components/ui/button";
import Provider from "@/models/provider";
import {getBadgeText, getBadgeVariant} from "../utils/badgeUtils";
import {Edit, MoreVertical, Trash2} from "lucide-react";
import {useProvidersMutations} from "@/hooks/providers/useProvidersMutations";
import { useNavigate } from "react-router-dom";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
import {argbToRgba} from "@/components/ui/utils/color-utils.ts";
import { useLabelQuery } from "@/hooks/labels/useLabelQuery";

interface ProviderCardProps {
    provider: Provider;
    onEdit: (provider: Provider) => void;
}

interface LabelPillProps { labelId: string; }
const LabelPill = ({ labelId }: LabelPillProps) => {
    const { data: label, isLoading } = useLabelQuery(labelId);
    return (
        <Badge
            variant="outline"
            className="text-xs py-0.5"
            style={{ backgroundColor: label?.color ? argbToRgba(label.color) : undefined }}
        >
            {label ? label.name : (isLoading ? "..." : labelId)}
        </Badge>
    );
};

export const ProviderCard = ({provider, onEdit}: ProviderCardProps) => {
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const {canModifyProvider, canDeleteProvider, deleteProviderMutation} = useProvidersMutations();
    const isEditable = canModifyProvider(provider);
    const isDeletable = canDeleteProvider(provider);
    const navigate = useNavigate();

    const handleDeleteProvider = async () => {
        try {
            await deleteProviderMutation.mutateAsync(provider.id);
            setIsDeleteDialogOpen(false);
        } catch (error) {
            console.error("Failed to delete provider:", error);
        }
    };

    return (
        <Card key={provider.id} className="overflow-hidden cursor-pointer hover:bg-accent/50 transition-colors" onClick={() => navigate(`/providers/${provider.id}`)}>
            {provider.iconUrl && (
                <div className="w-full h-20 overflow-hidden -mt-0.5 -mx-0.5">
                    <img
                        src={provider.iconUrl}
                        alt={`${provider.name} logo`}
                        className="w-full h-full object-contain p-1"
                    />
                </div>
            )}
            <CardHeader className="py-1 px-3">
                <div className="flex justify-between items-start">
                    <CardTitle className="text-lg">{provider.name}</CardTitle>
                    <div className="flex items-center space-x-2">
                        <Badge variant={getBadgeVariant(provider.owner.type)} className="text-xs">
                            {getBadgeText(provider.owner.type)}
                        </Badge>
                        {isEditable && (
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-8 w-8 p-0"
                                        onClick={(e) => e.stopPropagation()}
                                    >
                                        <MoreVertical className="h-4 w-4"/>
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem onClick={() => onEdit(provider)}>
                                        <Edit className="h-4 w-4 mr-2"/>
                                        Edit Provider
                                    </DropdownMenuItem>
                                    {isDeletable && (
                                        <>
                                            <DropdownMenuSeparator/>
                                            <DropdownMenuItem
                                                onClick={() => setIsDeleteDialogOpen(true)}
                                                className="text-destructive focus:text-destructive"
                                            >
                                                <Trash2 className="h-4 w-4 mr-2"/>
                                                Remove Provider
                                            </DropdownMenuItem>
                                        </>
                                    )}
                                </DropdownMenuContent>
                            </DropdownMenu>
                        )}
                    </div>
                </div>
                {provider.url && (
                    <CardDescription>
                        <a href={provider.url} target="_blank" rel="noopener noreferrer"
                           className="text-blue-500 hover:underline"
                           onClick={(e) => e.stopPropagation()}>
                            {provider.url}
                        </a>
                    </CardDescription>
                )}
            </CardHeader>
            <CardContent className="py-1 px-3">
                {provider.description && <p className="text-sm">{provider.description}</p>}
                {provider.labels.length > 0 && (
                    <div className="flex flex-wrap gap-1 my-1">
                        {provider.labels.map((labelId) => (
                            <LabelPill key={labelId} labelId={labelId} />
                        ))}
                    </div>
                )}
            </CardContent>
            <CardFooter className="flex flex-col space-y-1 w-full py-1 px-3">
                <div className="flex justify-between w-full">
                    {provider.pricingPageUrl && (
                        <a href={provider.pricingPageUrl} target="_blank" rel="noopener noreferrer">
                            <Button variant="outline" size="sm" className="h-7 text-xs">View Pricing</Button>
                        </a>
                    )}
                </div>

            </CardFooter>

            {/* Delete Provider Confirmation Dialog */}
            <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                        <AlertDialogDescription>
                            This will permanently delete the provider "{provider.name}" and all its associated plans and
                            prices.
                            This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                            onClick={handleDeleteProvider}
                            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                        >
                            Delete
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </Card>
    );
};