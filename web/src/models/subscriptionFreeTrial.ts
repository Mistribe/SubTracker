import type {SubscriptionFreeTrialModel} from "@/api/models";

export class SubscriptionFreeTrial {
    private readonly _endDate: Date | null;
    private readonly _startDate: Date | null;

    constructor(endDate: Date | null, startDate: Date | null) {
        this._endDate = endDate;
        this._startDate = startDate;
    }

    get endDate(): Date | null {
        return this._endDate;
    }

    get startDate(): Date | null {
        return this._startDate;
    }

    static fromModel(model: SubscriptionFreeTrialModel): SubscriptionFreeTrial {
        return new SubscriptionFreeTrial(
            model.endDate ? new Date(model.endDate) : null,
            model.startDate ? new Date(model.startDate) : null
        );
    }
}