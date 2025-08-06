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
                        <ToggleGroupItem value={OwnerType.Family} className="flex-1 justify-center py-6">
                            Family
                        </ToggleGroupItem>
                    </ToggleGroup>
                </div>
            </div>

            {selectedOwnerType === OwnerType.Family && (
                <div className="max-w-md mx-auto mt-6">
                    <Label htmlFor="familyId" className="text-lg mb-2 block">Which family is this for?</Label>
                    <Select
                        onValueChange={(value) => form.setValue("familyId", value)}
                        value={form.watch("familyId") || ""}
                    >
                        <SelectTrigger className="h-12">
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