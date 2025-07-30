import type {FamilyMemberModel} from "@/api/models";

export default class FamilyMember {
    private readonly _id: string;
    private readonly _familyId: string;
    private readonly _createdAt: Date;
    private readonly _etag: string;
    private readonly _isYou: boolean;

    constructor(id: string,
                name: string,
                isKid: boolean,
                familyId: string,
                isYou: boolean,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._isKid = isKid;
        this._isYou = isYou;
        this._familyId = familyId;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
    }

    get isYou(): boolean {
        return this._isYou;
    }

    private _name: string;

    get name(): string {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    private _isKid: boolean;

    get isKid(): boolean {
        return this._isKid;
    }

    set isKid(value: boolean) {
        this._isKid = value;
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

    get etag(): string {
        return this._etag;
    }

    static fromModel(model: FamilyMemberModel): FamilyMember {
        return new FamilyMember(
            model.id || '',
            model.name || '',
            model.isKid || false,
            model.familyId || '',
            model.isYou || false,
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        );
    }

}