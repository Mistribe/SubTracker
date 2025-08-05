import {useState} from "react";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover";
import {ColorPicker} from "@/components/ui/color-picker";
import {argbToRgba} from "@/components/ui/utils/color-utils";
import {Loader2, PlusIcon, SaveIcon, XIcon} from "lucide-react";
import {OwnerType} from "@/models/ownerType";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import Family from "@/models/family";

interface AddLabelFormProps {
    onAddLabel: (name: string, color: string, ownerType?: OwnerType, familyId?: string) => void;
    isAdding: boolean;
    families: Family[];
    title?: string;
}

export const AddLabelForm = ({
                                 onAddLabel,
                                 isAdding,
                                 families,
                             }: AddLabelFormProps) => {
    const [newLabel, setNewLabel] = useState("");
    const [newLabelColor, setNewLabelColor] = useState("#FF000000"); // Default ARGB color
    const [isFormVisible, setIsFormVisible] = useState(false);
    const [labelType, setLabelType] = useState<"personal" | "family">("personal");
    const [selectedFamilyId, setSelectedFamilyId] = useState<string>(
        families.length > 0 ? families[0].id : ""
    );

    const handleAddLabel = () => {
        if (newLabel.trim()) {
            if (labelType === "personal") {
                onAddLabel(newLabel, newLabelColor, OwnerType.Personal);
            } else if (labelType === "family" && selectedFamilyId) {
                onAddLabel(newLabel, newLabelColor, OwnerType.Family, selectedFamilyId);
            }
            setNewLabel(""); // Clear the input after adding
            setIsFormVisible(false); // Hide the form after adding
        }
    };
    
    const toggleForm = () => {
        if (!isFormVisible) {
            setIsFormVisible(true);
        } else {
            handleAddLabel();
        }
    };
    
    const handleCancel = () => {
        setIsFormVisible(false);
        setNewLabel(""); // Reset the input field
    };

    return (
        <div className="w-full">
            <div className="flex items-center gap-3">
                {/* Always visible Add/Save button */}
                {isFormVisible ? (
                    <div className="flex gap-2">
                        <Button
                            size="sm"
                            className="h-9"
                            onClick={toggleForm}
                            disabled={isAdding}
                        >
                            {isAdding ? (
                                <Loader2 className="h-4 w-4 mr-1 animate-spin"/>
                            ) : (
                                <SaveIcon className="h-4 w-4 mr-1"/>
                            )}
                            Save
                        </Button>
                        <Button
                            size="sm"
                            variant="outline"
                            className="h-9"
                            onClick={handleCancel}
                            disabled={isAdding}
                        >
                            <XIcon className="h-4 w-4 mr-1"/>
                            Cancel
                        </Button>
                    </div>
                ) : (
                    <Button
                        size="sm"
                        className="h-9"
                        onClick={toggleForm}
                        disabled={isAdding}
                    >
                        {isAdding ? (
                            <Loader2 className="h-4 w-4 mr-1 animate-spin"/>
                        ) : (
                            <PlusIcon className="h-4 w-4 mr-1"/>
                        )}
                        Add
                    </Button>
                )}
                
                {/* Form elements that appear with animation */}
                <div 
                    className={`flex items-center gap-3 overflow-hidden transition-all duration-300 ease-in-out ${
                        isFormVisible ? "max-w-[800px] opacity-100" : "max-w-0 opacity-0"
                    }`}
                >
                    {/* Owner type selector */}
                    <Select
                        value={labelType}
                        onValueChange={(value: "personal" | "family") => setLabelType(value)}
                    >
                        <SelectTrigger className="w-[130px] h-9">
                            <SelectValue placeholder="Label type" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="personal">Personal</SelectItem>
                            <SelectItem value="family" disabled={families.length === 0}>Family</SelectItem>
                        </SelectContent>
                    </Select>
                    
                    {/* Family selector - only shown when family type is selected and there are multiple families */}
                    {labelType === "family" && families.length > 1 && (
                        <Select
                            value={selectedFamilyId}
                            onValueChange={setSelectedFamilyId}
                        >
                            <SelectTrigger className="w-[150px] h-9">
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
                    )}
                    
                    {/* Label name input */}
                    <div className="flex items-center gap-2 whitespace-nowrap">
                        <div
                            className="w-4 h-4 rounded-full flex-shrink-0"
                            style={{backgroundColor: argbToRgba(newLabelColor)}}
                        />
                        <Input
                            placeholder="Enter label name"
                            value={newLabel}
                            onChange={(e) => setNewLabel(e.target.value)}
                            className="h-9 text-sm w-48 md:w-64"
                        />
                    </div>
                    
                    {/* Color picker */}
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button variant="outline" size="sm" className="h-9 whitespace-nowrap">
                                <div
                                    className="w-4 h-4 rounded-md mr-1"
                                    style={{backgroundColor: argbToRgba(newLabelColor)}}
                                />
                                Color
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
        </div>
    );
};