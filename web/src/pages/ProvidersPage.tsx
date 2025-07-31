import {useCallback, useEffect, useState} from "react";
import {useProvidersQuery} from "@/hooks/providers/useProvidersQuery";
import {useInView} from "react-intersection-observer";
import {AddProviderForm} from "@/components/providers/AddProviderForm";
import {EditProviderForm} from "@/components/providers/EditProviderForm";
import {ProviderCard} from "@/components/providers/ui/ProviderCard";
import {ProviderCardSkeletonGrid} from "@/components/providers/ui/ProviderCardSkeleton";
import {ErrorState} from "@/components/providers/ui/ErrorState";
import {InfiniteLoading} from "@/components/providers/ui/InfiniteLoading";
import {NoMoreProviders, NoProviders} from "@/components/providers/ui/EmptyStates";
import {PageHeader} from "@/components/providers/ui/PageHeader";
import Provider from "@/models/provider";

const ProvidersPage = () => {
    const {ref, inView} = useInView();
    const [loadMoreRef, setLoadMoreRef] = useState<HTMLDivElement | null>(null);

    // This effect connects the loadMoreRef element to the InView observer
    useEffect(() => {
        if (loadMoreRef) {
            ref(loadMoreRef);
        }
    }, [loadMoreRef, ref]);
    const [isAddingProvider, setIsAddingProvider] = useState(false);
    const [editingProvider, setEditingProvider] = useState<Provider | null>(null);
    // Query providers with infinite loading
    const {
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        isLoading,
        isError
    } = useProvidersQuery({limit: 12});

    // Load more providers when the user scrolls to the bottom
    const handleLoadMore = useCallback(() => {
        if (inView && hasNextPage && !isFetchingNextPage) {
            fetchNextPage();
        }
    }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

    // Call handleLoadMore when inView changes
    useEffect(() => {
        handleLoadMore();
    }, [inView, handleLoadMore]);

    // Flatten all providers from all pages
    const providers = data?.pages.flatMap(page => page.providers) || [];

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Providers"
                onAddProvider={() => setIsAddingProvider(true)}
            />

            {isLoading ? (
                <ProviderCardSkeletonGrid/>
            ) : isError ? (
                <ErrorState/>
            ) : (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {providers.map((provider) => (
                            <ProviderCard 
                                key={provider.id} 
                                provider={provider}
                                onEdit={setEditingProvider}
                            />
                        ))}
                    </div>

                    <InfiniteLoading
                        loadMoreRef={setLoadMoreRef}
                        isFetchingNextPage={isFetchingNextPage}
                        hasNextPage={hasNextPage || false}
                        onLoadMore={() => fetchNextPage()}
                    />

                    <NoMoreProviders show={!hasNextPage && providers.length > 0}/>

                    <NoProviders
                        show={providers.length === 0}
                        onAddProvider={() => setIsAddingProvider(true)}
                    />
                </>
            )}

            {/* Provider form modals */}
            <AddProviderForm
                isOpen={isAddingProvider}
                onClose={() => setIsAddingProvider(false)}
            />
            
            {editingProvider && (
                <EditProviderForm
                    isOpen={!!editingProvider}
                    onClose={() => setEditingProvider(null)}
                    provider={editingProvider}
                />
            )}
        </div>
    );
};

export default ProvidersPage;