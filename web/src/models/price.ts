import type {PriceModel} from "@/api/models";

export default class Price {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _amount: number;
    private readonly _currency: string;
    private readonly _startDate: Date;
    private readonly _endDate: Date | null;

    constructor(id: string,
                amount: number,
                currency: string,
                startDate: Date,
                endDate: Date | null,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
        this._amount = amount;
        this._currency = currency;
        this._startDate = startDate;
        this._endDate = endDate;
    }

    get id(): string {
        return this._id;
    }

    get createdAt(): Date {
        return this._createdAt;
    }

    get updatedAt(): Date {
        return this._updatedAt;
    }

    get etag(): string {
        return this._etag;
    }

    get amount(): number {
        return this._amount;
    }

    get currency(): string {
        return this._currency;
    }

    get startDate(): Date {
        return this._startDate;
    }

    get endDate(): Date | null {
        return this._endDate;
    }

    static fromModel(model: PriceModel): Price {
        return new Price(
            model.id || '',
            model.amount || 0,
            model.currency || '',
            model.startDate || new Date(),
            model.endDate || null,
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        )
    }
}