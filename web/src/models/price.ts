import type { DtoPriceModel as PriceModel } from "@/api/models/DtoPriceModel";
import type {Amount} from "@/models/amount.ts";

export default class Price {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _amount: Amount;
    private readonly _startDate: Date;
    private readonly _endDate: Date | null;

    constructor(id: string,
                amount: Amount,
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
        this._startDate = startDate;
        this._endDate = endDate;
    }

    get isActive(): boolean {
        return this._endDate === null || this._endDate > new Date();
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

    get amount(): Amount {
        return this._amount;
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
            {
                value: model.amount || 0,
                currency: model.currency || '',
            },
            model.startDate || new Date(),
            model.endDate || null,
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        )
    }
}