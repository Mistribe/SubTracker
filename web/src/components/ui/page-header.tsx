import type {ReactNode} from "react";
import {Input} from "@/components/ui/input";

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
                {/* Search input */}
                {onSearchChange && (
                    <div className="flex justify-center w-120">
                        <Input
                            placeholder={searchPlaceholder}
                            value={searchText}
                            onChange={(e) => onSearchChange(e.target.value)}
                            className="max-w-md"
                        />
                    </div>
                )}

                {/* Action button */}
                {actionButton && actionButton}
            </div>

        </div>
    );
};