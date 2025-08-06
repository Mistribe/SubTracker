import {useMemo, useState} from "react";
import {useAllProvidersQuery} from "@/hooks/providers/useAllProvidersQuery";
import {useAllLabelsQuery} from "@/hooks/labels/useAllLabelsQuery";
import {AddProviderForm} from "@/components/providers/AddProviderForm";
import {EditProviderForm} from "@/components/providers/EditProviderForm";
import {ProviderCard} from "@/components/providers/ui/ProviderCard";
import {ProviderCardSkeletonGrid} from "@/components/providers/ui/ProviderCardSkeleton";
import {ErrorState} from "@/components/providers/ui/ErrorState";
import {NoProviders} from "@/components/providers/ui/EmptyStates";
import {PageHeader} from "@/components/ui/page-header";
import Provider from "@/models/provider";

const ProvidersPage = () => {
    const [isAddingProvider, setIsAddingProvider] = useState(false);
    const [editingProvider, setEditingProvider] = useState<Provider | null>(null);
    const [searchText, setSearchText] = useState("");
    // Query all providers using the dedicated hook
    const {
        data,
        isLoading,
        isError
    } = useAllProvidersQuery();

    // Query all labels
    const {data: labelsData} = useAllLabelsQuery();

    // Get all available labels
    const availableLabels = useMemo(() => {
        if (!labelsData?.pages) return [];
        return labelsData.pages.flatMap(page => page.labels);
    }, [labelsData]);

    // Flatten all providers from all pages
    const allProviders = data?.pages.flatMap(page => page.providers) || [];

    // Filter providers based on search text
    const filteredProviders = useMemo(() => {
        if (searchText === "") return allProviders;

        const searchLower = searchText.toLowerCase();

        return allProviders.filter(provider => {
            // Filter by provider name or description
            if (provider.name.toLowerCase().includes(searchLower) ||
                (provider.description && provider.description.toLowerCase().includes(searchLower))) {
                return true;
            }

            // Filter by label name
            if (provider.labels.length > 0) {
                // Find matching labels
                const matchingLabels = provider.labels.some(labelId => {
                    const label = availableLabels.find(l => l.id === labelId);
                    return label && label.name.toLowerCase().includes(searchLower);
                });

                if (matchingLabels) return true;
            }

            return false;
        });
    }, [allProviders, searchText, availableLabels]);


    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Providers"
                description="Manage your providers"
                onAddProvider={() => setIsAddingProvider(true)}
                searchText={searchText}
                onSearchChange={setSearchText}
            />

            {isLoading ? (
                <ProviderCardSkeletonGrid/>
            ) : isError ? (
                <ErrorState/>
            ) : (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                        {filteredProviders.map((provider) => (
                            <ProviderCard
                                key={provider.id}
                                provider={provider}
                                onEdit={setEditingProvider}
                            />
                        ))}
                    </div>

                    {filteredProviders.length === 0 && (
                        <div className="text-center mt-8">
                            <p className="text-muted-foreground">No providers match your search criteria.</p>
                        </div>
                    )}

                    <NoProviders
                        show={allProviders.length === 0}
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