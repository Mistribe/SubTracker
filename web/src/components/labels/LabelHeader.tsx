import { Input } from "@/components/ui/input";
import { AddLabelForm } from "./AddLabelForm";
import { OwnerType } from "@/models/ownerType";
import Family from "@/models/family";

interface LabelHeaderProps {
    searchText: string;
    onSearchChange: (value: string) => void;
    onAddLabel: (name: string, color: string, ownerType?: OwnerType, familyId?: string) => void;
    isAdding: boolean;
    families: Family[];
}

export const LabelHeader = ({ 
    searchText, 
    onSearchChange, 
    onAddLabel, 
    isAdding,
    families 
}: LabelHeaderProps) => {

    return (
        <div className="flex flex-col gap-4 mb-6">
            <div className="flex justify-between items-center">
                <h1 className="text-3xl font-bold">Labels</h1>
                <div className="flex items-center gap-2">
                    <AddLabelForm
                        onAddLabel={onAddLabel}
                        isAdding={isAdding}
                        families={families}
                    />
                </div>
            </div>
            <div className="flex justify-center">
                <Input
                    placeholder="Search labels..."
                    value={searchText}
                    onChange={(e) => onSearchChange(e.target.value)}
                    className="max-w-md"
                />
            </div>
        </div>
    );
};