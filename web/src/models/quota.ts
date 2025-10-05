import type {DtoQuotaUsageModel} from "@/api";
import {FeatureId, FeatureType, isValidFeatureId, isValidFeatureType} from "@/models/billing.ts";

export default class Quota {
    private readonly _feature: FeatureId;
    private readonly _type: FeatureType;
    private readonly _enabled: boolean;
    private readonly _limit?: number;
    private readonly _used?: number;
    private readonly _remaining?: number;

    constructor(feature: FeatureId,
                type: FeatureType,
                enabled: boolean,
                limit?: number,
                used?: number,
                remaining?: number) {
        this._feature = feature;
        this._type = type;
        this._enabled = enabled;
        this._limit = limit;
        this._used = used;
        this._remaining = remaining;
    }

    static fromModel(model: DtoQuotaUsageModel): Quota {
        let featureId: FeatureId = FeatureId.Unknown;
        if (model.feature && isValidFeatureId(model.feature)) {
            featureId = model.feature as FeatureId;
        }
        let featureType: FeatureType = FeatureType.Unknown;
        if (model.type && isValidFeatureType(model.type)) {
            featureType = model.type as FeatureType;
        }
        return new Quota(
            featureId,
            featureType,
            model.enabled || false,
            model.limit || undefined,
            model.used || undefined,
            model.remaining || undefined,
        );
    }

    get feature(): FeatureId {
        return this._feature;
    }

    get type(): FeatureType {
        return this._type;
    }

    get enabled(): boolean {
        return this._enabled;
    }

    get limit(): number | undefined {
        return this._limit;
    }

    get used(): number | undefined {
        return this._used;
    }

    get remaining(): number | undefined {
        return this._remaining;
    }
}
