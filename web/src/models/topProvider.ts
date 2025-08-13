import type {Amount} from "@/models/amount.ts";

export default interface TopProvider {
    providerId: string;
    total: Amount;
}