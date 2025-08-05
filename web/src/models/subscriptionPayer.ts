import type {SubscriptionPayerModel} from "@/api/models";

export class SubscriptionPayer {
    private readonly _familyId: string | null;
    private readonly _memberId: string | null;
    private readonly _type: 'family' | 'family_member';
    private readonly _etag: string | null;

    constructor(
        familyId: string | null,
        memberId: string | null,
        type: 'family' | 'family_member',
        etag: string | null
    ) {
        this._familyId = familyId;
        this._memberId = memberId;
        this._type = type;
        this._etag = etag;
    }

    get familyId(): string | null {
        return this._familyId;
    }

    get memberId(): string | null {
        return this._memberId;
    }

    get type(): 'family' | 'family_member' {
        return this._type;
    }

    get etag(): string | null {
        return this._etag;
    }

    static fromModel(model: SubscriptionPayerModel): SubscriptionPayer {
        return new SubscriptionPayer(
            model.familyId || null,
            model.memberId || null,
            (model.type as 'family' | 'family_member') || 'family',
            model.etag || null
        );
    }
}