import { useCallback, useEffect, useState } from "react";
import { useProvidersQuery } from "@/hooks/providers/useProvidersQuery";
import { useInView } from "react-intersection-observer";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { OwnerType } from "@/models/ownerType";
import { AddProviderForm } from "@/components/providers/AddProviderForm";

const ProvidersPage = () => {
  const { ref: loadMoreRef, inView } = useInView();
  const [isAddingProvider, setIsAddingProvider] = useState(false);
  // Query providers with infinite loading
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError
  } = useProvidersQuery({ limit: 12 });

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

  // Get badge color based on owner type
  const getBadgeVariant = (ownerType: OwnerType) => {
    switch (ownerType) {
      case OwnerType.System:
        return "secondary";
      case OwnerType.Family:
        return "outline";
      case OwnerType.Personal:
        return "default";
      default:
        return "default";
    }
  };

  // Get badge text based on owner type
  const getBadgeText = (ownerType: OwnerType) => {
    switch (ownerType) {
      case OwnerType.System:
        return "System";
      case OwnerType.Family:
        return "Family";
      case OwnerType.Personal:
        return "Personal";
      default:
        return "Unknown";
    }
  };

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Providers</h1>
        <Button onClick={() => setIsAddingProvider(true)}>Add Provider</Button>
      </div>

      {isLoading ? (
        // Loading skeleton
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array(6).fill(0).map((_, index) => (
            <Card key={index} className="overflow-hidden">
              <CardHeader className="pb-2">
                <Skeleton className="h-6 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-20 w-full" />
              </CardContent>
              <CardFooter>
                <Skeleton className="h-8 w-20" />
              </CardFooter>
            </Card>
          ))}
        </div>
      ) : isError ? (
        // Error state
        <div className="text-center py-10">
          <p className="text-red-500">Failed to load providers. Please try again later.</p>
          <Button variant="outline" className="mt-4" onClick={() => window.location.reload()}>
            Retry
          </Button>
        </div>
      ) : (
        // Provider cards grid
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {providers.map((provider) => (
              <Card key={provider.id} className="overflow-hidden">
                <CardHeader className="pb-2">
                  <div className="flex justify-between items-start">
                    <CardTitle>{provider.name}</CardTitle>
                    <Badge variant={getBadgeVariant(provider.owner.type)}>
                      {getBadgeText(provider.owner.type)}
                    </Badge>
                  </div>
                  {provider.url && (
                    <CardDescription>
                      <a href={provider.url} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                        {provider.url}
                      </a>
                    </CardDescription>
                  )}
                </CardHeader>
                <CardContent>
                  {provider.description && <p>{provider.description}</p>}
                  {provider.labels.length > 0 && (
                    <div className="flex flex-wrap gap-2 mt-2">
                      {provider.labels.map((label, index) => (
                        <Badge key={index} variant="outline">{label}</Badge>
                      ))}
                    </div>
                  )}
                </CardContent>
                <CardFooter className="flex justify-between">
                  {provider.pricingPageUrl && (
                    <a href={provider.pricingPageUrl} target="_blank" rel="noopener noreferrer">
                      <Button variant="outline" size="sm">View Pricing</Button>
                    </a>
                  )}
                  {provider.plans.length > 0 && (
                    <Badge variant="secondary">{provider.plans.length} Plans</Badge>
                  )}
                </CardFooter>
              </Card>
            ))}
          </div>

          {/* Infinite loading trigger */}
          {(hasNextPage || isFetchingNextPage) && (
            <div ref={loadMoreRef} className="flex justify-center mt-8 pb-8">
              {isFetchingNextPage ? (
                <div className="flex flex-col items-center">
                  <Skeleton className="h-8 w-8 rounded-full" />
                  <p className="text-sm text-muted-foreground mt-2">Loading more providers...</p>
                </div>
              ) : (
                <Button variant="ghost" onClick={() => fetchNextPage()}>
                  Load More
                </Button>
              )}
            </div>
          )}

          {!hasNextPage && providers.length > 0 && (
            <p className="text-center text-muted-foreground mt-8 pb-8">
              No more providers to load
            </p>
          )}

          {providers.length === 0 && (
            <div className="text-center py-10">
              <p className="text-muted-foreground">No providers found</p>
              <Button variant="outline" className="mt-4" onClick={() => setIsAddingProvider(true)}>
                Add Your First Provider
              </Button>
            </div>
          )}
        </>
      )}

      {/* Provider form modal */}
      <AddProviderForm 
        isOpen={isAddingProvider}
        onClose={() => setIsAddingProvider(false)}
      />
    </div>
  );
};

export default ProvidersPage;