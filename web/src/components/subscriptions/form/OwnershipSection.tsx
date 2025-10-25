import { useEffect } from "react";
import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {ToggleGroup, ToggleGroupItem} from "@/components/ui/toggle-group";
import type {FormValues} from "./SubscriptionFormSchema";
import {OwnerType} from "@/models/ownerType";
import Family from "@/models/family";
import {Checkbox} from "@/components/ui/checkbox";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

interface OwnershipSectionProps {
    family?: Family | null;
}

export const OwnershipSection = ({ family }: OwnershipSectionProps) => {
    const form = useFormContext<FormValues>();
    const selectedOwnerType = form.watch("ownerType");
    const hasFamily = Boolean(family);

    useEffect(() => {
        if (selectedOwnerType === OwnerType.Family && family) {
            const currentFamilyId = form.getValues("familyId");
            if (!currentFamilyId) {
                form.setValue("familyId", family.id);
                form.setValue("serviceUsers", []);
                form.setValue("payerType", "family");
                form.setValue("payerId", family.id);
            }
        }
    }, [selectedOwnerType, family, form]);

    // Safety: if user selects Family but has no family, revert to Personal
    useEffect(() => {
        if (selectedOwnerType === OwnerType.Family && !family) {
            form.setValue("ownerType", OwnerType.Personal);
        }
    }, [selectedOwnerType, family, form]);

    return (
        <div className="space-y-6">
            <div className="text-center">
                <h2 className="text-2xl font-semibold mb-2">Who is this subscription for?</h2>
                <p className="text-muted-foreground">The subscription is for me or a family member?</p>
            </div>

            <div className="max-w-md mx-auto">
                <div className="mt-4">
                    <ToggleGroup
                        type="single"
                        defaultValue={OwnerType.Personal}
                        value={selectedOwnerType}
                        onValueChange={(value) => value && form.setValue("ownerType", value as OwnerType)}
                        variant="outline"
                        className="w-full"
                    >
                        <ToggleGroupItem value={OwnerType.Personal} className="flex-1 justify-center py-6">
                            Just for me
                        </ToggleGroupItem>
                        {hasFamily ? (
                            <ToggleGroupItem value={OwnerType.Family} className="flex-1 justify-center py-6">
                                Family
                            </ToggleGroupItem>
                        ) : (
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <span className="flex-1">
                                            <ToggleGroupItem
                                                value={OwnerType.Family}
                                                className="flex-1 justify-center py-6 cursor-not-allowed opacity-60"
                                                disabled
                                            >
                                                Family
                                            </ToggleGroupItem>
                                        </span>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        <p>You don’t have a family yet. Create a family to use this feature.</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                        )}
                    </ToggleGroup>
                </div>
            </div>

            {selectedOwnerType === OwnerType.Family && (
                <div className="max-w-md mx-auto mt-6 space-y-6">
                    <div>
                        <Label htmlFor="familyId" className="text-lg mb-2 block">Family</Label>
                        <div className="h-12 px-3 flex items-center rounded-md border bg-background text-foreground">
                            {family?.name}
                        </div>
                        {form.formState.errors.familyId && (
                            <p className="text-sm text-red-500 mt-1">{form.formState.errors.familyId.message}</p>
                        )}
                    </div>
                    
                    {form.watch("familyId") && (
                        <>
                            <div>
                                <Label className="text-lg mb-2 block">Who is using this subscription?</Label>
                                <div className="space-y-2 mt-2">
                                    {family?.members.map((member) => (
                                        <div key={member.id} className="flex items-center space-x-2">
                                            <Checkbox 
                                                id={`member-${member.id}`}
                                                checked={form.watch("serviceUsers")?.includes(member.id) || false}
                                                onCheckedChange={(checked) => {
                                                    const currentUsers = form.watch("serviceUsers") || [];
                                                    if (checked) {
                                                        form.setValue("serviceUsers", [...currentUsers, member.id]);
                                                    } else {
                                                        form.setValue("serviceUsers", currentUsers.filter(id => id !== member.id));
                                                    }
                                                }}
                                            />
                                            <Label 
                                                htmlFor={`member-${member.id}`}
                                                className="text-sm font-normal cursor-pointer"
                                            >
                                                {member.name} {member.isYou ? "(You)" : ""}
                                            </Label>
                                        </div>
                                    ))}
                                </div>
                            </div>
                            
                            <div>
                                <Label className="text-lg mb-2 block">Who is paying for this subscription?</Label>
                                <Select
                                    onValueChange={(value) => {
                                        const [type, id] = value.split(':');
                                        form.setValue("payerType", type as "family" | "family_member");
                                        form.setValue("payerId", id);
                                    }}
                                    value={`${form.watch("payerType") || "family"}:${form.watch("payerId") || form.watch("familyId")}`}
                                >
                                    <SelectTrigger className="h-12">
                                        <SelectValue placeholder="Select who is paying"/>
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value={`family:${form.watch("familyId")}`}>
                                            Family: {family?.name}
                                        </SelectItem>
                                        {family?.members.map((member) => (
                                            <SelectItem key={member.id} value={`family_member:${member.id}`}>
                                                {member.name} {member.isYou ? "(You)" : ""}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>
                        </>
                    )}
                </div>
            )}
        </div>
    );
};