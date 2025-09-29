import React from "react";
import {Skeleton} from "@/components/ui/skeleton";

export const QuotaUsageSkeleton: React.FC<{ className?: string }> = ({className}) => (
  <div className={"border rounded-md p-4 bg-card flex flex-col gap-2 " + (className || "")}>
    <div className="flex items-center justify-between">
      <Skeleton className="h-4 w-28"/>
      <Skeleton className="h-3 w-16"/>
    </div>
    <Skeleton className="h-2 w-full"/>
    <Skeleton className="h-3 w-20"/>
  </div>
);

export default QuotaUsageSkeleton;

