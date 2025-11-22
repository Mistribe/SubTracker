import type {Amount} from "@/models/amount.ts";

export default interface UpcomingRenewal {
    at: Date;
    providerId: string;
    subscriptionId: string;
    total: Amount;
    source?: Amount;
}