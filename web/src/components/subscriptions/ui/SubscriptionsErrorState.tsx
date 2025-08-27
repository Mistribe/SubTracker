export function SubscriptionsErrorState() {
  return (
    <div className="flex flex-col items-center justify-center py-12">
      <h3 className="text-xl font-semibold mb-2">Error Loading Subscriptions</h3>
      <p className="text-muted-foreground mb-6">
        There was a problem loading your subscriptions. Please try again later.
      </p>
    </div>
  );
}

export default SubscriptionsErrorState;
