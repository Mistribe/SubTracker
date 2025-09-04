import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Link } from "react-router-dom";
import { CreditCard } from "lucide-react";

const EmptySubscriptionsState = () => {
  return (
    <div className="mx-auto max-w-3xl">
      <Card>
        <CardHeader className="text-center">
          <div className="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary">
            <CreditCard className="h-5 w-5" />
          </div>
          <CardTitle className="text-2xl">No subscriptions yet</CardTitle>
          <CardDescription className="text-base">
            Start tracking your subscriptions to see spending, upcoming renewals, and insights on your dashboard.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ul className="grid gap-3 text-left sm:grid-cols-2">
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary">•</div>
              <div>
                <p className="font-medium">Stay ahead of renewals</p>
                <p className="text-muted-foreground">Get a clear view of what’s renewing soon.</p>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <div className="mt-1 text-primary">•</div>
              <div>
                <p className="font-medium">Understand your spend</p>
                <p className="text-muted-foreground">See monthly and yearly totals at a glance.</p>
              </div>
            </li>
          </ul>
        </CardContent>
        <CardFooter className="justify-center">
          <Button asChild>
            <Link to="/subscriptions/create">Add your first subscription</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default EmptySubscriptionsState;
