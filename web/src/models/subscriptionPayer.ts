import type {SubscriptionPayerModel} from "@/api/models";

export class SubscriptionPayer {
    private readonly _familyId: string | undefined;
    private readonly _memberId: string | undefined;
    private readonly _type: 'family' | 'family_member';
    private readonly _etag: string | undefined;

    constructor(
        familyId: string | undefined,
        memberId: string | undefined,
        type: 'family' | 'family_member',
        etag: string | undefined
    ) {
        this._familyId = familyId;
        this._memberId = memberId;
        this._type = type;
        this._etag = etag;
    }

    get familyId(): string | undefined {
        return this._familyId;
    }

    get memberId(): string | undefined {
        return this._memberId;
    }

    get type(): 'family' | 'family_member' {
        return this._type;
    }

    get etag(): string | undefined {
        return this._etag;
    }

    static fromModel(model: SubscriptionPayerModel): SubscriptionPayer {
        return new SubscriptionPayer(
            model.familyId || undefined,
            model.memberId || undefined,
            (model.type as 'family' | 'family_member') || 'family',
            model.etag || undefined
        );
    }
}