import { useFormContext } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { FormValues } from "./SubscriptionFormSchema";
import { OwnerType } from "@/models/ownerType";
import { Family } from "@/models/family";

interface OwnershipSectionProps {
    families: Family[];
}

export const OwnershipSection = ({ families }: OwnershipSectionProps) => {
    const form = useFormContext<FormValues>();
    const selectedOwnerType = form.watch("ownerType");
    
    return (
        <div className="space-y-4">
            <h2 className="text-xl font-semibold">Ownership</h2>

            <div>
                <Label>Owner Type</Label>
                <RadioGroup
                    defaultValue={OwnerType.Personal}
                    value={selectedOwnerType}
                    onValueChange={(value) => form.setValue("ownerType", value as OwnerType)}
                    className="grid grid-cols-2 gap-4 mt-2"
                >
                    <div className="flex items-center space-x-2">
                        <RadioGroupItem value={OwnerType.Personal} id="personal"/>
                        <Label htmlFor="personal">Personal</Label>
                    </div>
                    <div className="flex items-center space-x-2">
                        <RadioGroupItem value={OwnerType.Family} id="family"/>
                        <Label htmlFor="family">Family</Label>
                    </div>
                </RadioGroup>
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