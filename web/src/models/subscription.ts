import Owner from "@/models/owner.ts";
import {SubscriptionPayer} from "@/models/subscriptionPayer.ts";
import {SubscriptionFreeTrial} from "@/models/subscriptionFreeTrial.ts";
import {fromHttpApi, SubscriptionRecurrency} from "@/models/subscriptionRecurrency.ts";
import {addMonths, addYears} from "date-fns";
import {daysBetween, monthsBetween} from "@/utils/date.ts";
import {type Amount, fromModel} from "@/models/amount.ts";
import type {DtoSubscriptionModel as SubscriptionModel} from "@/api/models/DtoSubscriptionModel";
import {SubscriptionStatus, subscriptionStatusFromHttpApi} from "@/models/subscriptionStatus.ts";

export default class Subscription {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _updatedAt: Date;
    private readonly _etag: string;
    private readonly _friendlyName: string | undefined;
    private readonly _providerId: string;
    private readonly _recurrency: SubscriptionRecurrency;
    private readonly _customRecurrency: number | undefined;
    private readonly _startDate: Date;
    private readonly _endDate: Date | undefined;
    private readonly _owner: Owner;
    private readonly _payer: SubscriptionPayer | undefined;
    private readonly _familyUsers: string[];
    private readonly _price: Amount | undefined;
    private readonly _freeTrial: SubscriptionFreeTrial | undefined;
    private readonly _status: SubscriptionStatus;

    constructor(
        id: string,
        createdAt: Date,
        updatedAt: Date,
        etag: string,
        friendlyName: string | undefined,
        providerId: string,
        recurrency: SubscriptionRecurrency,
        customRecurrency: number | undefined,
        startDate: Date,
        endDate: Date | undefined,
        owner: Owner,
        payer: SubscriptionPayer | undefined,
        familyUsers: string[],
        customPrice: Amount | undefined,
        freeTrial: SubscriptionFreeTrial | undefined,
        status: SubscriptionStatus) {
        this._id = id;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
        this._friendlyName = friendlyName;
        this._providerId = providerId;
        this._recurrency = recurrency;
        this._customRecurrency = customRecurrency;
        this._startDate = startDate;
        this._endDate = endDate;
        this._owner = owner;
        this._payer = payer;
        this._familyUsers = familyUsers;
        this._price = customPrice;
        this._freeTrial = freeTrial;
        this._status = status;
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
        return this._status === SubscriptionStatus.Active;
    }

    get status(): SubscriptionStatus {
        return this._status;
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

    get familyUsers(): string[] {
        return this._familyUsers;
    }

    get price(): Amount | undefined {
        return this._price;
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
            fromHttpApi(model.recurrency),
            model.customRecurrency || undefined,
            model.startDate ? new Date(model.startDate) : new Date(),
            model.endDate ? new Date(model.endDate) : undefined,
            Owner.fromModel(model.owner || {}),
            model.payer ? SubscriptionPayer.fromModel(model.payer) : undefined,
            model.familyUsers || [],
            model.price ? fromModel(model.price) : undefined,
            model.freeTrial ? SubscriptionFreeTrial.fromModel(model.freeTrial) : undefined,
            subscriptionStatusFromHttpApi(model.status),
        );
    }

    getDailyPrice(): number {
        if (!this._price?.value) {
            return 0;
        }
        return this.getMonthlyPrice() / 30;
    }

    getMonthlyPrice(): number {
        if (!this._price?.value) {
            return 0;
        }

        switch (this._recurrency) {
            case SubscriptionRecurrency.Monthly:
                return this._price.value;
            case SubscriptionRecurrency.Quarterly:
                return this._price.value / 3;
            case SubscriptionRecurrency.HalfYearly:
                return this._price.value / 6;
            case SubscriptionRecurrency.Yearly:
                return this._price.value / 12;
            case SubscriptionRecurrency.Custom:
                if (this._customRecurrency) {
                    return this._price.value / this._customRecurrency;
                }
                return 0;
            case SubscriptionRecurrency.OneTime: {
                if (!this._endDate) {
                    return 0;
                }
                const numberOfMonths = monthsBetween(this._startDate, this._endDate, {includePartial: true});
                if (numberOfMonths <= 0) {
                    return 0;
                }
                return this._price.value / numberOfMonths;
            }
            default:
                return 0;
        }
    }

    getYearlyPrice(): number {
        if (!this._price?.value) {
            return 0;
        }

        switch (this._recurrency) {
            case SubscriptionRecurrency.Monthly:
                return this._price.value * 12;
            case SubscriptionRecurrency.Quarterly:
                return this._price.value * 4;
            case SubscriptionRecurrency.HalfYearly:
                return this._price.value * 2;
            case SubscriptionRecurrency.Yearly:
                return this._price.value;
            case SubscriptionRecurrency.Custom:
                if (this._customRecurrency) {
                    return (this._price.value * 12) / this._customRecurrency;
                }
                return 0;
            case SubscriptionRecurrency.OneTime: {
                if (!this._endDate) {
                    return 0;
                }
                const numberOfMonths = monthsBetween(this._startDate, this._endDate, {includePartial: true});
                if (numberOfMonths <= 0) {
                    return 0;
                }

                return (this._price.value * 12) / numberOfMonths;
            }
            default:
                return 0;
        }
    }

    public getTotalAmount(): number {
        const startDate = this._startDate;
        let endDate = this._endDate;
        if (!endDate) {
            endDate = new Date();
        }
        const numberOfMonths = daysBetween(startDate, endDate, {includePartial: true});
        if (numberOfMonths <= 0) {
            return 0;
        }

        const dailyPrice = this.getDailyPrice();
        if (dailyPrice <= 0) {
            return 0;
        }
        return numberOfMonths * dailyPrice;
    }

    public getAmount(): number {
        if (this._price) {
            return this._price.value;
        }
        // todo price plan are not implemented yet
        return 0;
    }

    public getCurrency(): string {
        if (this._price) {
            return this._price.currency;
        }

        // todo price plan are not implemented yet
        return 'USD'
    }

    public getNextRenewalDate(): Date | undefined {
        const dates = this.getNextRenewalDates(1);
        if (dates.length > 0) {
            return dates[0];
        }
        return undefined;
    }

    public getNextRenewalDates(n: number): Date[] {
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
                const intervalMonths = this._customRecurrency || 0;
                if (intervalMonths <= 0) {
                    return [];
                }
                while (current < today) {
                    current = addMonths(current, intervalMonths);
                }
                for (let i = 0; i < n && withinEndDate(current); i++) {
                    dates.push(new Date(current));
                    current = addMonths(current, intervalMonths);
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

