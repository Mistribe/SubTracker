import type {SubscriptionModel} from "@/api/models";
import Owner from "@/models/owner.ts";
import {SubscriptionPayer} from "@/models/subscriptionPayer.ts";
import {SubscriptionCustomPrice} from "@/models/subscriptionCustomPrice.ts";
import {SubscriptionFreeTrial} from "@/models/subscriptionFreeTrial.ts";
import {fromHttpApi, type SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";

export default class Subscription {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _friendlyName: string | null;
    private readonly _providerId: string;
    private readonly _planId: string;
    private readonly _priceId: string;
    private readonly _recurrency: SubscriptionRecurrency;
    private readonly _customRecurrency: number | null;
    private readonly _startDate: Date;
    private readonly _endDate: Date | null;
    private readonly _owner: Owner;
    private readonly _payer: SubscriptionPayer | null;
    private readonly _serviceUsers: string[];
    private readonly _customPrice: SubscriptionCustomPrice | null;
    private readonly _freeTrial: SubscriptionFreeTrial | null;

    constructor(
        id: string,
        createdAt: Date,
        updatedAt: Date,
        etag: string,
        friendlyName: string | null,
        providerId: string,
        planId: string,
        priceId: string,
        recurrency: SubscriptionRecurrency,
        customRecurrency: number | null,
        startDate: Date,
        endDate: Date | null,
        owner: Owner,
        payer: SubscriptionPayer | null,
        serviceUsers: string[],
        customPrice: SubscriptionCustomPrice | null,
        freeTrial: SubscriptionFreeTrial | null
    ) {
        this._id = id;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
        this._friendlyName = friendlyName;
        this._providerId = providerId;
        this._planId = planId;
        this._priceId = priceId;
        this._recurrency = recurrency;
        this._customRecurrency = customRecurrency;
        this._startDate = startDate;
        this._endDate = endDate;
        this._owner = owner;
        this._payer = payer;
        this._serviceUsers = serviceUsers;
        this._customPrice = customPrice;
        this._freeTrial = freeTrial;
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

    get isActive(): boolean {
        return this._endDate === null || this._endDate > new Date();
    }

    get etag(): string {
        return this._etag;
    }

    get friendlyName(): string | null {
        return this._friendlyName;
    }

    get providerId(): string {
        return this._providerId;
    }

    get planId(): string {
        return this._planId;
    }

    get priceId(): string {
        return this._priceId;
    }

    get recurrency(): SubscriptionRecurrency {
        return this._recurrency;
    }

    get customRecurrency(): number | null {
        return this._customRecurrency;
    }

    get startDate(): Date {
        return this._startDate;
    }

    get endDate(): Date | null {
        return this._endDate;
    }

    get owner(): Owner {
        return this._owner;
    }

    get payer(): SubscriptionPayer | null {
        return this._payer;
    }

    get serviceUsers(): string[] {
        return this._serviceUsers;
    }

    get customPrice(): SubscriptionCustomPrice | null {
        return this._customPrice;
    }

    get freeTrial(): SubscriptionFreeTrial | null {
        return this._freeTrial;
    }

    static fromModel(model: SubscriptionModel): Subscription {
        return new Subscription(
            model.id || '',
            model.createdAt ? new Date(model.createdAt) : new Date(),
            model.updatedAt ? new Date(model.updatedAt) : new Date(),
            model.etag || '',
            model.friendlyName || null,
            model.providerId || '',
            model.planId || '',
            model.priceId || '',
            fromHttpApi(model.recurrency),
            model.customRecurrency || null,
            model.startDate ? new Date(model.startDate) : new Date(),
            model.endDate ? new Date(model.endDate) : null,
            Owner.fromModel(model.owner || {}),
            model.payer ? SubscriptionPayer.fromModel(model.payer) : null,
            model.serviceUsers || [],
            model.customPrice ? SubscriptionCustomPrice.fromModel(model.customPrice) : null,
            model.freeTrial ? SubscriptionFreeTrial.fromModel(model.freeTrial) : null
        );
    }
}

