import type { DtoSubscriptionPayerModel as SubscriptionPayerModel } from "@/api/models/DtoSubscriptionPayerModel";
import { DtoSubscriptionPayerModelTypeEnum } from "@/api/models/DtoSubscriptionPayerModel";
import {PayerType} from "@/models/payerType.ts";

export class SubscriptionPayer {
    private readonly _memberId: string | undefined;
    private readonly _type: PayerType;
    private readonly _etag: string | undefined;

    constructor(
        memberId: string | undefined,
        type: PayerType,
        etag: string | undefined
    ) {
        this._memberId = memberId;
        this._type = type;
        this._etag = etag;
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
            case DtoSubscriptionPayerModelTypeEnum.Family:
                payerType = PayerType.Family;
                break;
            case DtoSubscriptionPayerModelTypeEnum.FamilyMember:
                payerType = PayerType.FamilyMember;
                break;
            default:
                payerType = PayerType.Family;
        }
        return new SubscriptionPayer(
            model.memberId || undefined,
            payerType,
            model.etag || undefined
        );
    }
}