import type {ReactNode} from "react";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Badge} from "@/components/ui/badge";
import {X} from "lucide-react";
import Label from "@/models/label";
import {OwnerType} from "@/models/ownerType";

export interface PageHeaderProps {
    /**
     * The title of the page
     */
    title: string;

    /**
     * Optional description of the page
     */
    description?: string;

    /**
     * Current search text
     */
    searchText?: string;

    /**
     * Function to handle search text changes
     */
    onSearchChange?: (value: string) => void;

    /**
     * Function to handle adding a provider
     */
    onAddProvider?: () => void;

    /**
     * Function to handle adding a label
     */
    onAddLabel?: (name: string, color: string, ownerType?: OwnerType, familyId?: string) => void;

    /**
     * Boolean indicating if an item is being added
     */
    isAdding?: boolean;

    /**
     * Array of selected labels
     */
    selectedLabels?: Label[];

    /**
     * Array of available labels
     */
    availableLabels?: Label[];

    /**
     * Function to handle label selection
     */
    onLabelSelect?: (label: Label) => void;

    /**
     * Function to handle label removal
     */
    onLabelRemove?: (labelId: string) => void;

    /**
     * Custom action button
     */
    actionButton?: ReactNode;

    /**
     * Search placeholder text
     */
    searchPlaceholder?: string;
}

export const PageHeader = ({
                               title,
                               description,
                               searchText = "",
                               onSearchChange,
                               onAddProvider,
                               onAddLabel,
                               isAdding,
                               selectedLabels = [],
                               availableLabels = [],
                               onLabelSelect,
                               onLabelRemove,
                               actionButton,
                               searchPlaceholder = "Search..."
                           }: PageHeaderProps) => {
    return (
        <div className="flex flex-col gap-4 mb-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold">{title}</h1>
                    {description && (
                        <p className="text-muted-foreground mt-1">{description}</p>
                    )}
                </div>

                {/* Action buttons */}
                {actionButton ? (
                    actionButton
                ) : (
                    <>
                        {onAddProvider && (
                            <Button onClick={onAddProvider}>Add Provider</Button>
                        )}
                        {onAddLabel && (
                            <div className="flex items-center gap-2">
                                {/* This is a placeholder for the AddLabelForm component */}
                                {/* In a real implementation, you would import and use the actual AddLabelForm component */}
                                <Button
                                    onClick={() => onAddLabel("New Label", "#FF0000")}
                                    disabled={isAdding}
                                >
                                    Add Label
                                </Button>
                            </div>
                        )}
                    </>
                )}
            </div>

            {/* Search input */}
            {onSearchChange && (
                <div className="flex justify-center">
                    <Input
                        placeholder={searchPlaceholder}
                        value={searchText}
                        onChange={(e) => onSearchChange(e.target.value)}
                        className="max-w-md"
                    />
                </div>
            )}

            {/* Label selection */}
            {onLabelSelect && availableLabels.length > 0 && (
                <div className="flex flex-col gap-2">
                    <div className="flex justify-center">
                        <div className="flex flex-wrap gap-2 max-w-2xl justify-center">
                            {availableLabels
                                .filter(label => !selectedLabels.some(selected => selected.id === label.id))
                                .map(label => (
                                    <Badge
                                        key={label.id}
                                        variant="outline"
                                        className="cursor-pointer hover:bg-secondary/80"
                                        style={{backgroundColor: label.color || undefined}}
                                        onClick={() => onLabelSelect(label)}
                                    >
                                        {label.name}
                                    </Badge>
                                ))
                            }
                        </div>
                    </div>

                    {selectedLabels.length > 0 && (
                        <div className="flex justify-center">
                            <div className="flex flex-wrap gap-2 max-w-2xl justify-center">
                                {selectedLabels.map(label => (
                                    <Badge
                                        key={label.id}
                                        variant="secondary"
                                        className="cursor-pointer flex items-center gap-1"
                                        style={{backgroundColor: label.color || undefined}}
                                    >
                                        {label.name}
                                        <X
                                            className="h-3 w-3 hover:text-destructive"
                                            onClick={() => onLabelRemove?.(label.id)}
                                        />
                                    </Badge>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};