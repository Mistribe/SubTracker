import type { DtoSubscriptionFreeTrialModel as SubscriptionFreeTrialModel } from "@/api/models/DtoSubscriptionFreeTrialModel";

export class SubscriptionFreeTrial {
    private readonly _endDate: Date | undefined;
    private readonly _startDate: Date | undefined;

    constructor(endDate: Date | undefined, startDate: Date | undefined) {
        this._endDate = endDate;
        this._startDate = startDate;
    }

    get endDate(): Date | undefined {
        return this._endDate;
    }

    get startDate(): Date | undefined {
        return this._startDate;
    }

    static fromModel(model: SubscriptionFreeTrialModel): SubscriptionFreeTrial {
        return new SubscriptionFreeTrial(
            model.endDate ? new Date(model.endDate) : undefined,
            model.startDate ? new Date(model.startDate) : undefined
        );
    }
}