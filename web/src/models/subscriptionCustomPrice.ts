import type {SubscriptionCustomPriceModel} from "@/api/models";

export class SubscriptionCustomPrice {
    private readonly _amount: number;
    private readonly _currency: string;

    constructor(amount: number, currency: string) {
        this._amount = amount;
        this._currency = currency;
    }

    get amount(): number {
        return this._amount;
    }

    get currency(): string {
        return this._currency;
    }

    static fromModel(model: SubscriptionCustomPriceModel): SubscriptionCustomPrice {
        return new SubscriptionCustomPrice(
            model.amount || 0,
            model.currency || 'USD'
        );
    }
}