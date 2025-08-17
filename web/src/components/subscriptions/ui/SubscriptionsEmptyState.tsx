import React from "react";

export function SubscriptionsEmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-12">
      <h3 className="text-xl font-semibold mb-2">No Subscriptions Found</h3>
      <p className="text-muted-foreground mb-6">
        You don't have any subscriptions yet.
      </p>
    </div>
  );
}

export default SubscriptionsEmptyState;
