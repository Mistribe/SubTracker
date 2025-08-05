import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

export const ProviderCardSkeleton = () => {
  return (
    <Card className="overflow-hidden">
      <Skeleton className="w-full h-28" />
      <CardHeader className="py-2 px-3">
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-3 w-1/2" />
      </CardHeader>
      <CardContent className="py-2 px-3">
        <Skeleton className="h-16 w-full" />
      </CardContent>
      <CardFooter className="py-2 px-3">
        <Skeleton className="h-7 w-20" />
      </CardFooter>
    </Card>
  );
};

export const ProviderCardSkeletonGrid = ({ count = 8 }: { count?: number }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      {Array(count).fill(0).map((_, index) => (
        <ProviderCardSkeleton key={index} />
      ))}
    </div>
  );
};