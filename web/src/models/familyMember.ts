import type {FamilyMemberModel} from "@/api/models";

export default class FamilyMember {
    private readonly _id: string;
    private _name: string;
    private _isKid: boolean;
    private _email: string | null;
    private readonly _familyId: string;
    private readonly _createdAt: Date;
    private _updatedAt: Date;
    private readonly _etag: string;


    constructor(id: string,
                name: string,
                isKid: boolean,
                email: string | null,
                familyId: string,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._isKid = isKid;
        this._email = email;
        this._familyId = familyId;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
    }

    static fromModel(model: FamilyMemberModel): FamilyMember {
        return new FamilyMember(
            model.id || '',
            model.name || '',
            model.isKid || false,
            null,
            model.familyId || '',
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        );
    }

    set name(value: string) {
        this._name = value;
    }

    set isKid(value: boolean) {
        this._isKid = value;
    }

    set email(value: string | null) {
        this._email = value;
    }

    set updatedAt(value: Date) {
        this._updatedAt = value;
    }

    get id(): string {
        return this._id;
    }

    get name(): string {
        return this._name;
    }

    get isKid(): boolean {
        return this._isKid;
    }

    get email(): string | null {
        return this._email;
    }

    get familyId(): string {
        return this._familyId;
    }

    get createdAt(): Date {
        return this._createdAt;
    }

    get updatedAt(): Date {
        return this._updatedAt;
    }

    get etag(): string {
        return this._etag;
    }

}