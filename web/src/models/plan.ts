import Price from "@/models/price.ts";
import type {PlanModel} from "@/api/models";

export default class Plan {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _name: string;
    private readonly _description: string | null;

    constructor(id: string,
                name: string,
                description: string | null,
                prices: Price[],
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
        this._name = name;
        this._description = description;
        this._prices = prices;
    }

    private _prices: Price[];

    get prices(): Price[] {
        return this._prices;
    }

    set prices(value: Price[]) {
        this._prices = value;
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

    get name(): string {
        return this._name;
    }

    get description(): string | null {
        return this._description;
    }

    static fromModel(model: PlanModel): Plan {
        return new Plan(
            model.id || '',
            model.name || '',
            model.description || null,
            model.prices?.map(price => Price.fromModel(price)) || [],
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        )
    }
}
