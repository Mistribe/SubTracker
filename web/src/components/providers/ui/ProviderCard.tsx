import {useState} from "react";
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card";
import {Badge} from "@/components/ui/badge";
import {Button} from "@/components/ui/button";
import Provider from "@/models/provider";
import Plan from "@/models/plan";
import {getBadgeText, getBadgeVariant} from "../utils/badgeUtils";
import {Edit, MoreVertical, Plus, Trash2} from "lucide-react";
import {useProvidersMutations} from "@/hooks/providers/useProvidersMutations";
import {AddPlanDialog} from "../AddPlanDialog";
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
import {PlanDetailsDialog} from "@/components/providers/plan-details";

interface ProviderCardProps {
    provider: Provider;
    onEdit: (provider: Provider) => void;
}

export const ProviderCard = ({provider, onEdit}: ProviderCardProps) => {
    const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null);
    const [isAddingPlan, setIsAddingPlan] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const {canModifyProvider, canDeleteProvider, deleteProviderMutation} = useProvidersMutations();
    const isEditable = canModifyProvider(provider);
    const isDeletable = canDeleteProvider(provider);

    const handleDeleteProvider = async () => {
        try {
            await deleteProviderMutation.mutateAsync(provider.id);
            setIsDeleteDialogOpen(false);
        } catch (error) {
            console.error("Failed to delete provider:", error);
        }
    };

    return (
        <Card key={provider.id} className="overflow-hidden">
            <CardHeader className="pb-2">
                <div className="flex justify-between items-start">
                    <CardTitle>{provider.name}</CardTitle>
                    <div className="flex items-center space-x-2">
                        <Badge variant={getBadgeVariant(provider.owner.type)}>
                            {getBadgeText(provider.owner.type)}
                        </Badge>
                        {isEditable && (
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        className="h-8 w-8 p-0"
                                    >
                                        <MoreVertical className="h-4 w-4"/>
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem onClick={() => onEdit(provider)}>
                                        <Edit className="h-4 w-4 mr-2"/>
                                        Edit Provider
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => setIsAddingPlan(true)}>
                                        <Plus className="h-4 w-4 mr-2"/>
                                        Add Plan
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
                           className="text-blue-500 hover:underline">
                            {provider.url}
                        </a>
                    </CardDescription>
                )}
            </CardHeader>
            <CardContent>
                {provider.description && <p>{provider.description}</p>}
                {provider.labels.length > 0 && (
                    <div className="flex flex-wrap gap-2 mt-2">
                        {provider.labels.map((label, index) => (
                            <Badge key={index} variant="outline">{label}</Badge>
                        ))}
                    </div>
                )}
            </CardContent>
            <CardFooter className="flex flex-col space-y-2 w-full">
                <div className="flex justify-between w-full">
                    {provider.pricingPageUrl && (
                        <a href={provider.pricingPageUrl} target="_blank" rel="noopener noreferrer">
                            <Button variant="outline" size="sm">View Pricing</Button>
                        </a>
                    )}
                </div>

                {provider.plans.length > 0 && (
                    <div className="flex flex-wrap gap-2 w-full">
                        {provider.plans.map((plan) => {
                            // Count active prices
                            const activePricesCount = plan.prices.filter(price => price.isActive).length;

                            return (
                                <Badge
                                    key={plan.id}
                                    variant="secondary"
                                    className="cursor-pointer hover:bg-secondary/80 text-sm py-1.5 px-3"
                                    onClick={() => setSelectedPlan(plan)}
                                >
                                    {plan.name} {activePricesCount > 0 &&
                                    <span className="ml-1 font-bold">{activePricesCount}</span>}
                                </Badge>
                            );
                        })}
                    </div>
                )}
            </CardFooter>

            {/* Plan Details Dialog */}
            {selectedPlan && (
                <PlanDetailsDialog
                    isOpen={!!selectedPlan}
                    onClose={() => setSelectedPlan(null)}
                    plan={selectedPlan}
                    providerId={provider.id}
                />
            )}

            {/* Add Plan Dialog */}
            {isAddingPlan && (
                <AddPlanDialog
                    isOpen={isAddingPlan}
                    onClose={() => setIsAddingPlan(false)}
                    providerId={provider.id}
                />
            )}

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