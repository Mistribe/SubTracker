import Owner from "@/models/owner.ts";
import Plan from "@/models/plan.ts";
import type {ProviderModel} from "@/api/models/provider";

export default class Provider {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _name: string;
    private readonly _description: string | null;
    private readonly _iconUrl: string | null;
    private readonly _url: string | null;
    private readonly _pricingPageUrl: string | null;
    private readonly _labels: string[];
    private readonly _plans: Plan[];
    private readonly _owner: Owner;


    constructor(id: string,
                name: string,
                description: string | null,
                iconUrl: string | null,
                url: string | null,
                pricingPageUrl: string | null,
                labels: string[],
                plans: Plan[],
                owner: Owner,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
        this._name = name;
        this._description = description;
        this._iconUrl = iconUrl;
        this._url = url;
        this._pricingPageUrl = pricingPageUrl;
        this._labels = labels;
        this._plans = plans;
        this._owner = owner;
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

    get iconUrl(): string | null {
        return this._iconUrl;
    }

    get url(): string | null {
        return this._url;
    }

    get pricingPageUrl(): string | null {
        return this._pricingPageUrl;
    }

    get labels(): string[] {
        return this._labels;
    }

    get plans(): Plan[] {
        return this._plans;
    }

    get owner(): Owner {
        return this._owner;
    }

    static fromModel(model: ProviderModel): Provider {
        return new Provider(
            model.id || '',
            model.name || '',
            model.description || null,
            model.iconUrl || null,
            model.url || null,
            model.pricingPageUrl || null,
            model.labels || [],
            model.plans?.map(plan => Plan.fromModel(plan)) || [],
            Owner.fromModel(model.owner || {}),
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        )
    }
}
