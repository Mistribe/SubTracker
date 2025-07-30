import type {FamilyModel} from "@/api/models";
import FamilyMember from "@/models/familyMember.ts";

export default class Family {
    private readonly _id: string;
    private readonly _isOwner: boolean;
    private readonly _createdAt: Date;
    private readonly _etag: string;

    constructor(id: string,
                name: string,
                isOwner: boolean,
                members: FamilyMember[],
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._isOwner = isOwner;
        this._members = members;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
    }

    private _name: string;

    get name(): string {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    private _members: FamilyMember[];

    get members(): FamilyMember[] {
        return this._members;
    }

    set members(value: FamilyMember[]) {
        this._members = value;
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

    get isOwner(): boolean {
        return this._isOwner;
    }

    get createdAt(): Date {
        return this._createdAt;
    }

    get etag(): string {
        return this._etag;
    }

    static fromModel(model: FamilyModel): Family {
        return new Family(
            model.id || '',
            model.name || '',
            model.isOwner || false,
            model.members?.map(member => FamilyMember.fromModel(member)) || [],
            model.createdAt ? model.createdAt : new Date(),
            model.updatedAt ? model.updatedAt : new Date(),
            model.etag || ''
        )
    }


}