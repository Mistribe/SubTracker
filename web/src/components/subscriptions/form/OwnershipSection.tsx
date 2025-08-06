import {useFormContext} from "react-hook-form";
import {Label} from "@/components/ui/label";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {ToggleGroup, ToggleGroupItem} from "@/components/ui/toggle-group";
import type {FormValues} from "./SubscriptionFormSchema";
import {OwnerType} from "@/models/ownerType";
import Family from "@/models/family";

interface OwnershipSectionProps {
    families: Family[];
}

export const OwnershipSection = ({families}: OwnershipSectionProps) => {
    const form = useFormContext<FormValues>();
    const selectedOwnerType = form.watch("ownerType");

    return (
        <div className="space-y-4">
            <h2 className="text-xl font-semibold">Ownership</h2>

            <div>
                <Label>Owner Type</Label>
                <div className="mt-2">
                    <ToggleGroup
                        type="single"
                        defaultValue={OwnerType.Personal}
                        value={selectedOwnerType}
                        onValueChange={(value) => value && form.setValue("ownerType", value as OwnerType)}
                        variant="outline"
                        className="w-full"
                    >
                        <ToggleGroupItem value={OwnerType.Personal} className="flex-1 justify-center">
                            Personal
                        </ToggleGroupItem>
                        <ToggleGroupItem value={OwnerType.Family} className="flex-1 justify-center">
                            Family
                        </ToggleGroupItem>
                    </ToggleGroup>
                </div>
            </div>

            {selectedOwnerType === OwnerType.Family && (
                <div>
                    <Label htmlFor="familyId">Family</Label>
                    <Select
                        onValueChange={(value) => form.setValue("familyId", value)}
                        value={form.watch("familyId") || ""}
                    >
                        <SelectTrigger>
                            <SelectValue placeholder="Select a family"/>
                        </SelectTrigger>
                        <SelectContent>
                            {families.map((family) => (
                                <SelectItem key={family.id} value={family.id}>
                                    {family.name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                    {form.formState.errors.familyId && (
                        <p className="text-sm text-red-500 mt-1">{form.formState.errors.familyId.message}</p>
                    )}
                </div>
            )}
        </div>
    );
};