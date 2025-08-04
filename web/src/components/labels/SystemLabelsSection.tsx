import {useEffect} from "react";
import {LabelItem} from "./LabelItem";
import {Loader2} from "lucide-react";
import {useAllLabelsQuery} from "@/hooks/labels/useAllLabelsQuery";
import {OwnerType} from "@/models/ownerType";

interface SystemLabelsSectionProps {
    /** Number of items to request per page – backend maximum is 10. */
    limit?: number;
}

export const SystemLabelsSection = ({limit = 10}: SystemLabelsSectionProps) => {
    // Fetch *all* system labels, page after page
    const {
        data,
        isLoading,
        error,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
    } = useAllLabelsQuery({
        ownerTypes: [OwnerType.System],
        limit,
    });

    // Keep requesting next pages until every label is fetched
    useEffect(() => {
        if (hasNextPage && !isFetchingNextPage) {
            fetchNextPage().catch(() => void 0);
        }
    }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

    const labels =
        data?.pages.flatMap((page) => page.labels) ?? [];

    if (isLoading) {
        return (
            <div className="mb-8">
                <h2 className="text-xl font-semibold mb-4">System Labels</h2>
                <div className="flex flex-col items-center justify-center h-32">
                    <Loader2 className="h-6 w-6 animate-spin text-primary mb-2"/>
                    <p className="text-muted-foreground text-sm">Loading system labels...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="mb-8">
                <h2 className="text-xl font-semibold mb-4">System Labels</h2>
                <div className="p-4 border rounded-md bg-destructive/10">
                    <p className="text-destructive text-sm">Error loading system labels</p>
                </div>
            </div>
        );
    }

    return (
        <div className="mb-8">
            <h2 className="text-xl font-semibold mb-4">System Labels</h2>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
                {labels.length === 0 ? (
                    <p className="text-muted-foreground col-span-full text-center py-4">No system labels found</p>
                ) : (
                    labels.map((label) => (
                        <LabelItem
                            key={label.id}
                            label={label}
                            isReadOnly={true}
                        />
                    ))
                )}
            </div>
        </div>
    );
};