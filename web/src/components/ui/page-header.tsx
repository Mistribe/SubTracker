import type {ReactNode} from "react";
import {useEffect, useRef, useState} from "react";
import {Input} from "@/components/ui/input";
import {Button} from "@/components/ui/button.tsx";
import {SlidersHorizontalIcon} from "lucide-react";

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
     * A callback function that is invoked when a filter action is triggered.
     * Typically used to handle the application of filters in a user interface or
     * other context-dependent logic.
     */
    onFilter?: () => void;
    /**
     * Custom action button
     */
    actionButton?: ReactNode;

    /**
     * Search placeholder text
     */
    searchPlaceholder?: string;

    /**
     * Debounce delay (ms) applied before calling onSearchChange
     * Defaults to 300ms for a balance between responsiveness and fewer calls
     */
    searchDebounceMs?: number;
}

export const PageHeader = ({
                               title,
                               description,
                               searchText = "",
                               onSearchChange,
                               actionButton,
                               onFilter,
                               searchPlaceholder = "Search...",
                               searchDebounceMs = 300,
                           }: PageHeaderProps) => {
    // Local input state for instant UI feedback
    const [localValue, setLocalValue] = useState<string>(searchText ?? "");

    // Keep local value in sync when the parent-controlled value changes
    useEffect(() => {
        setLocalValue(searchText ?? "");
    }, [searchText]);

    // Track the last value that has been emitted to avoid redundant calls
    const lastEmittedRef = useRef<string | undefined>(undefined);

    // Debounce changes before notifying the parent
    useEffect(() => {
        if (!onSearchChange) return;

        // If delay is 0 or negative, emit immediately for maximum responsiveness
        if ((searchDebounceMs ?? 0) <= 0) {
            if (lastEmittedRef.current !== localValue) {
                lastEmittedRef.current = localValue;
                onSearchChange(localValue);
            }
            return;
        }

        const handle = window.setTimeout(() => {
            if (lastEmittedRef.current !== localValue) {
                lastEmittedRef.current = localValue;
                onSearchChange(localValue);
            }
        }, searchDebounceMs);

        return () => {
            window.clearTimeout(handle);
        };
    }, [localValue, onSearchChange, searchDebounceMs]);

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
                            value={localValue}
                            onChange={(e) => setLocalValue(e.target.value)}
                            className="max-w-md"
                        />
                        {onFilter && (
                            <Button size="icon" onClick={onFilter}>
                                <SlidersHorizontalIcon />
                            </Button>
                        )}
                    </div>
                )}

                {/* Action button */}
                {actionButton && actionButton}
            </div>

        </div>
    );
};