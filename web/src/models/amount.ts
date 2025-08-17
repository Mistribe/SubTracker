import type {AmountModel} from "@/api/models";

export interface Amount {
    value: number;
    currency: string;
    source?: Amount;
}

export function fromModel(model: AmountModel): Amount {
    return {
        value: model.value || 0,
        currency: model.currency || 'USD',
        source: model.source ? fromModel(model.source) : undefined,
    }

}

export const zeroAmount: Amount = {value: 0, currency: 'USD'};