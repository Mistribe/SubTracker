import type {Amount} from "@/models/amount.ts";

export default class SubscriptionPrice {
    private readonly _monthly: Amount | undefined;
    private readonly _yearly: Amount | undefined;

    constructor(monthly: Amount | null | undefined, yearly: Amount | null |undefined) {
        this._monthly = monthly ?? undefined;
        this._yearly = yearly ?? undefined;
    }


    get monthly(): Amount | undefined {
        return this._monthly;
    }

    get yearly(): Amount | undefined {
        return this._yearly;
    }
}