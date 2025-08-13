import {type FamilyMemberModel} from "@/api/models";
import {FamilyMemberType, fromHttpApi} from "@/models/familyMemberType.ts";

export default class FamilyMember {
    private readonly _id: string;
    private readonly _familyId: string;
    private readonly _createdAt: Date;
    private readonly _etag: string;
    private readonly _isYou: boolean;
    private readonly _hasAccount: boolean;
    private readonly _type: FamilyMemberType;

    constructor(id: string,
                name: string,
                memberType: FamilyMemberType,
                familyId: string,
                isYou: boolean,
                hasAccount: boolean,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._type = memberType;
        this._isYou = isYou;
        this._familyId = familyId;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
        this._hasAccount = hasAccount;
    }

    get isYou(): boolean {
        return this._isYou;
    }

    get hasAccount(): boolean {
        return this._hasAccount;
    }

    private _name: string;

    get name(): string {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    get type(): FamilyMemberType {
        return this._type;
    }


    private _updatedAt: Date;

    get updatedAt(): Date {
        return this._updatedAt;
    }

    set updatedAt(value: Date) {
        this._updatedAt = value;
    }

    get id(): string {
        return this._id;
    }

    get familyId(): string {
        return this._familyId;
    }

    get createdAt(): Date {
        return this._createdAt;
    }

    get isKid(): boolean {
        return this._type === FamilyMemberType.Kid;
    }

    get etag(): string {
        return this._etag;
    }

    static fromModel(model: FamilyMemberModel): FamilyMember {
        return new FamilyMember(
            model.id || '',
            model.name || '',
            fromHttpApi(model.type),
            model.familyId || '',
            model.isYou || false,
            model.hasAccount || false,
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        );
    }

}
