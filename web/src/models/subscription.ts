import type {SubscriptionModel} from "@/api/models";
import Owner from "@/models/owner.ts";
import {SubscriptionPayer} from "@/models/subscriptionPayer.ts";
import {SubscriptionCustomPrice} from "@/models/subscriptionCustomPrice.ts";
import {SubscriptionFreeTrial} from "@/models/subscriptionFreeTrial.ts";
import {fromHttpApi, SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";
import {addDays, addMonths, addYears} from "date-fns";

export default class Subscription {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _friendlyName: string | undefined;
    private readonly _providerId: string;
    private readonly _planId: string;
    private readonly _priceId: string;
    private readonly _recurrency: SubscriptionRecurrency;
    private readonly _customRecurrency: number | undefined;
    private readonly _startDate: Date;
    private readonly _endDate: Date | undefined;
    private readonly _owner: Owner;
    private readonly _payer: SubscriptionPayer | undefined;
    private readonly _serviceUsers: string[];
    private readonly _customPrice: SubscriptionCustomPrice | undefined;
    private readonly _freeTrial: SubscriptionFreeTrial | undefined;

    constructor(
        id: string,
        createdAt: Date,
        updatedAt: Date,
        etag: string,
        friendlyName: string | undefined,
        providerId: string,
        planId: string,
        priceId: string,
        recurrency: SubscriptionRecurrency,
        customRecurrency: number | undefined,
        startDate: Date,
        endDate: Date | undefined,
        owner: Owner,
        payer: SubscriptionPayer | undefined,
        serviceUsers: string[],
        customPrice: SubscriptionCustomPrice | undefined,
        freeTrial: SubscriptionFreeTrial | undefined
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
        return this._endDate === undefined || this._endDate > new Date();
    }

    get etag(): string {
        return this._etag;
    }

    get friendlyName(): string | undefined {
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

    get customRecurrency(): number | undefined {
        return this._customRecurrency;
    }

    get startDate(): Date {
        return this._startDate;
    }

    get endDate(): Date | undefined {
        return this._endDate;
    }

    get owner(): Owner {
        return this._owner;
    }

    get payer(): SubscriptionPayer | undefined {
        return this._payer;
    }

    get serviceUsers(): string[] {
        return this._serviceUsers;
    }

    get customPrice(): SubscriptionCustomPrice | undefined {
        return this._customPrice;
    }

    get freeTrial(): SubscriptionFreeTrial | undefined {
        return this._freeTrial;
    }

    static fromModel(model: SubscriptionModel): Subscription {
        return new Subscription(
            model.id || '',
            model.createdAt ? new Date(model.createdAt) : new Date(),
            model.updatedAt ? new Date(model.updatedAt) : new Date(),
            model.etag || '',
            model.friendlyName || undefined,
            model.providerId || '',
            model.planId || '',
            model.priceId || '',
            fromHttpApi(model.recurrency),
            model.customRecurrency || undefined,
            model.startDate ? new Date(model.startDate) : new Date(),
            model.endDate ? new Date(model.endDate) : undefined,
            Owner.fromModel(model.owner || {}),
            model.payer ? SubscriptionPayer.fromModel(model.payer) : undefined,
            model.serviceUsers || [],
            model.customPrice ? SubscriptionCustomPrice.fromModel(model.customPrice) : undefined,
            model.freeTrial ? SubscriptionFreeTrial.fromModel(model.freeTrial) : undefined
        );
    }

    nextRenewDates(n: number): Date[] {
        // No renewals for unknown or one-time subscriptions
        if (this._recurrency === SubscriptionRecurrency.Unknown ||
            this._recurrency === SubscriptionRecurrency.OneTime) {
            return [];
        }

        const dates: Date[] = [];
        const today = new Date();
        const end = this._endDate;

        // Helper to stop if endDate exists and current is after endDate
        const withinEndDate = (d: Date) => !end || d <= end;

        let current = new Date(this._startDate);

        switch (this._recurrency) {
            case SubscriptionRecurrency.Monthly: {
                while (current < today) {
                    current = addMonths(current, 1);
                }
                for (let i = 0; i < n && withinEndDate(current); i++) {
                    dates.push(new Date(current));
                    current = addMonths(current, 1);
                }
                break;
            }
            case SubscriptionRecurrency.Quarterly: {
                while (current < today) {
                    current = addMonths(current, 3);
                }
                for (let i = 0; i < n && withinEndDate(current); i++) {
                    dates.push(new Date(current));
                    current = addMonths(current, 3);
                }
                break;
            }
            case SubscriptionRecurrency.HalfYearly: {
                while (current < today) {
                    current = addMonths(current, 6);
                }
                for (let i = 0; i < n && withinEndDate(current); i++) {
                    dates.push(new Date(current));
                    current = addMonths(current, 6);
                }
                break;
            }
            case SubscriptionRecurrency.Yearly: {
                while (current < today) {
                    current = addYears(current, 1);
                }
                for (let i = 0; i < n && withinEndDate(current); i++) {
                    dates.push(new Date(current));
                    current = addYears(current, 1);
                }
                break;
            }
            case SubscriptionRecurrency.Custom: {
                const intervalDays = this._customRecurrency || 0;
                if (intervalDays <= 0) {
                    return [];
                }
                while (current < today) {
                    current = addDays(current, intervalDays);
                }
                for (let i = 0; i < n && withinEndDate(current); i++) {
                    dates.push(new Date(current));
                    current = addDays(current, intervalDays);
                }
                break;
            }
            default:
                // Safety: If a new recurrency type is added and not handled, return empty.
                return [];
        }

        return dates;
    }
}

