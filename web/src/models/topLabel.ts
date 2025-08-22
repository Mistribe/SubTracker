import type {Amount} from "@/models/amount.ts";

export default interface TopLabel {
    labelId: string;
    total: Amount;
}