import Subscription from "@/models/subscription";

export interface SubscriptionWithNextRenewal {
    subscription: Subscription;
    nextRenewalDate: Date;
}

