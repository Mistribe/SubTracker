import type {FamilyModel} from "@/api/models";
import FamilyMember from "@/models/familyMember.ts";

export default class Family {
    private readonly _id: string;
    private _name: string;
    private readonly _isOwner: boolean;
    private _haveJointAccount: boolean;
    private readonly _createdAt: Date;
    private _members: FamilyMember[];
    private _updatedAt: Date;
    private readonly _etag: string;


    constructor(id: string,
                name: string,
                isOwner: boolean,
                haveJointAccount: boolean,
                members: FamilyMember[],
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._isOwner = isOwner;
        this._haveJointAccount = haveJointAccount;
        this._members = members;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
    }

    static fromModel(model: FamilyModel): Family {
        return new Family(
            model.id || '',
            model.name || '',
            model.isOwner || false,
            model.haveJointAccount || false,
            model.members?.map(member => FamilyMember.fromModel(member)) || [],
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        )
    }

    get id(): string {
        return this._id;
    }

    get members(): FamilyMember[] {
        return this._members;
    }

    get name(): string {
        return this._name;
    }

    get isOwner(): boolean {
        return this._isOwner;
    }

    set members(value: FamilyMember[]) {
        this._members = value;
    }


    get haveJointAccount(): boolean {
        return this._haveJointAccount;
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

    set name(value: string) {
        this._name = value;
    }

    set haveJointAccount(value: boolean) {
        this._haveJointAccount = value;
    }

    set updatedAt(value: Date) {
        this._updatedAt = value;
    }


}