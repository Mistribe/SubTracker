import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { ColorPicker } from "@/components/ui/color-picker";
import { argbToRgba } from "@/components/ui/utils/color-utils";
import { Loader2, PlusIcon } from "lucide-react";
import { OwnerType } from "@/models/ownerType";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import Family from "@/models/family";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { useLabelsQuotaQuery } from "@/hooks/labels/useLabelsQuotaQuery";
import { useQuotaLimit, getQuotaTooltip } from "@/hooks/quotas/useFeature";
import { FeatureId } from "@/models/billing";
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip";

interface AddLabelDialogProps {
    onAddLabel: (name: string, color: string, ownerType?: OwnerType, familyId?: string) => void;
    isAdding: boolean;
    families: Family[];
    trigger?: React.ReactNode;
}

export function AddLabelDialog({
    onAddLabel,
    isAdding,
    families,
    trigger
}: AddLabelDialogProps) {
    const [open, setOpen] = useState(false);
    const [newLabel, setNewLabel] = useState("");
    const [newLabelColor, setNewLabelColor] = useState("#FF000000"); // Default ARGB color
    const [labelType, setLabelType] = useState<"personal" | "family">("personal");
    const [selectedFamilyId, setSelectedFamilyId] = useState<string>(
        families.length > 0 ? families[0].id : ""
    );

    // Check labels quota
    const { data: labelsQuota } = useLabelsQuotaQuery();
    const { enabled: labelsEnabled, canAdd: canAddLabels, used: labelsUsed, limit: labelsLimit } = useQuotaLimit(
        labelsQuota,
        FeatureId.CustomLabelsCount
    );
    const isDisabled = !labelsEnabled || !canAddLabels;
    const tooltipMessage = getQuotaTooltip(labelsEnabled, canAddLabels, "custom labels");

    const handleAddLabel = () => {
        if (newLabel.trim()) {
            if (labelType === "personal") {
                onAddLabel(newLabel, newLabelColor, OwnerType.Personal);
            } else if (labelType === "family" && selectedFamilyId) {
                onAddLabel(newLabel, newLabelColor, OwnerType.Family, selectedFamilyId);
            }

            // Reset form and close dialog
            resetForm();
            setOpen(false);
        }
    };

    const resetForm = () => {
        setNewLabel("");
        setNewLabelColor("#FF000000");
        setLabelType("personal");
        setSelectedFamilyId(families.length > 0 ? families[0].id : "");
    };

    const handleOpenChange = (newOpen: boolean) => {
        if (!newOpen) {
            resetForm();
        }
        setOpen(newOpen);
    };

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <TooltipProvider>
                <Tooltip>
                    <TooltipTrigger asChild>
                        <span>
                            <DialogTrigger asChild>
                                {trigger || (
                                    <Button disabled={isDisabled}>
                                        <PlusIcon className="h-4 w-4 mr-2" />
                                        Add Label
                                        {labelsEnabled && labelsLimit !== undefined && (
                                            <span className="ml-2 text-xs opacity-70">
                                                ({labelsUsed}/{labelsLimit})
                                            </span>
                                        )}
                                    </Button>
                                )}
                            </DialogTrigger>
                        </span>
                    </TooltipTrigger>
                    {tooltipMessage && (
                        <TooltipContent>
                            <p>{tooltipMessage}</p>
                        </TooltipContent>
                    )}
                </Tooltip>
            </TooltipProvider>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Add New Label</DialogTitle>
                    <DialogDescription>
                        Create a new label to organize your subscriptions.
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                    {/* Owner type selector */}
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Label Type</label>
                        <Select
                            value={labelType}
                            onValueChange={(value: "personal" | "family") => setLabelType(value)}
                        >
                            <SelectTrigger className="w-full">
                                <SelectValue placeholder="Label type" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="personal">Personal</SelectItem>
                                <SelectItem value="family" disabled={families.length === 0}>Family</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Family selector - only shown when family type is selected */}
                    {labelType === "family" && families.length > 0 && (
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Family</label>
                            <Select
                                value={selectedFamilyId}
                                onValueChange={setSelectedFamilyId}
                            >
                                <SelectTrigger className="w-full">
                                    <SelectValue placeholder="Select family" />
                                </SelectTrigger>
                                <SelectContent>
                                    {families.map(family => (
                                        <SelectItem key={family.id} value={family.id}>
                                            {family.name}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    )}

                    {/* Label name input */}
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Label Name</label>
                        <div className="flex items-center gap-2">
                            <div
                                className="w-4 h-4 rounded-full flex-shrink-0"
                                style={{ backgroundColor: argbToRgba(newLabelColor) }}
                            />
                            <Input
                                placeholder="Enter label name"
                                value={newLabel}
                                onChange={(e) => setNewLabel(e.target.value)}
                            />
                        </div>
                    </div>

                    {/* Color picker */}
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Label Color</label>
                        <Popover>
                            <PopoverTrigger asChild>
                                <Button variant="outline" className="w-full justify-start">
                                    <div
                                        className="w-4 h-4 rounded-md mr-2"
                                        style={{ backgroundColor: argbToRgba(newLabelColor) }}
                                    />
                                    {argbToRgba(newLabelColor)}
                                </Button>
                            </PopoverTrigger>
                            <PopoverContent className="w-auto p-0" align="start">
                                <ColorPicker
                                    color={newLabelColor}
                                    onChange={setNewLabelColor}
                                />
                            </PopoverContent>
                        </Popover>
                    </div>
                </div>

                <DialogFooter>
                    <Button
                        onClick={handleAddLabel}
                        disabled={isAdding || !newLabel.trim()}
                    >
                        {isAdding ? (
                            <>
                                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                                Adding...
                            </>
                        ) : (
                            "Add Label"
                        )}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}