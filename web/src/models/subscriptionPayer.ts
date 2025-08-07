import {type SubscriptionPayerModel, SubscriptionPayerModel_typeObject} from "@/api/models";
import {PayerType} from "@/models/payerType.ts";

export class SubscriptionPayer {
    private readonly _familyId: string | undefined;
    private readonly _memberId: string | undefined;
    private readonly _type: PayerType;
    private readonly _etag: string | undefined;

    constructor(
        familyId: string | undefined,
        memberId: string | undefined,
        type: PayerType,
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
        let payerType: PayerType;
        switch (model.type) {
            case SubscriptionPayerModel_typeObject.Family:
                payerType = PayerType.Family;
                break;
            case SubscriptionPayerModel_typeObject.Family_member:
                payerType = PayerType.FamilyMember;
                break;
            default:
                payerType = PayerType.Family;
        }
        return new SubscriptionPayer(
            model.familyId || undefined,
            model.memberId || undefined,
            payerType,
            model.etag || undefined
        );
    }
}