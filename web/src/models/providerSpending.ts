import type {Amount} from "@/models/amount.ts";

export interface ProviderSpending {
    id: string;
    providerName: string;
    amount: Amount;
}