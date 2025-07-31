import {Button} from "@/components/ui/button";
import {Skeleton} from "@/components/ui/skeleton";
import type {RefObject, LegacyRef} from "react";

interface InfiniteLoadingProps {
    loadMoreRef?: RefObject<HTMLDivElement> | LegacyRef<HTMLDivElement>;
    isFetchingNextPage: boolean;
    hasNextPage: boolean;
    onLoadMore: () => void;
}

export const InfiniteLoading = ({
                                    loadMoreRef,
                                    isFetchingNextPage,
                                    hasNextPage,
                                    onLoadMore
                                }: InfiniteLoadingProps) => {
    if (!hasNextPage && !isFetchingNextPage) {
        return null;
    }

    return (
        <div ref={loadMoreRef} className="flex justify-center mt-8 pb-8">
            {isFetchingNextPage ? (
                <div className="flex flex-col items-center">
                    <Skeleton className="h-8 w-8 rounded-full"/>
                    <p className="text-sm text-muted-foreground mt-2">Loading more providers...</p>
                </div>
            ) : (
                <Button variant="ghost" onClick={onLoadMore}>
                    Load More
                </Button>
            )}
        </div>
    );
};