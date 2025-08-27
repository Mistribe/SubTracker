import { Skeleton } from "@/components/ui/skeleton";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

export function SubscriptionsTableSkeleton() {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
          <TableHead>
            <Skeleton className="h-4 w-20" />
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {Array(8)
          .fill(0)
          .map((_, index) => (
            <TableRow key={index}>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-4 w-full" />
              </TableCell>
            </TableRow>
          ))}
      </TableBody>
    </Table>
  );
}

export default SubscriptionsTableSkeleton;
