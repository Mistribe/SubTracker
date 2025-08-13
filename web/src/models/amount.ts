export interface Amount {
    value: number;
    currency: string;
}

export const zeroAmount: Amount = {value: 0, currency: 'USD'};