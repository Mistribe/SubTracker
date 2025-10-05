import {isValidOwnerType, OwnerType} from "@/models/ownerType.ts";
import type { DtoOwnerModel as OwnerModel } from "@/api/models/DtoOwnerModel";

export default class Owner {
    constructor(type: OwnerType,
                etag: string,
                familyId: string | null,
                userId: string | null) {
        this._type = type;
        this._etag = etag;
        this._familyId = familyId;
        this._userId = userId;
    }

    _familyId: string | null;

    public get familyId(): string | null {
        return this._familyId;
    }

    _userId: string | null;

    public get userId(): string | null {
        return this._userId;
    }

    public get isPersonal(): boolean {
        return this._type === OwnerType.Personal;
    }

    public get isFamily(): boolean {
        return this._type === OwnerType.Family;
    }

    public get isSystem(): boolean {
        return this._type === OwnerType.System;
    }

    _type: OwnerType;

    public get type(): OwnerType {
        return this._type;
    }

    _etag: string;

    public get etag(): string {
        return this._etag;
    }

    public static fromModel(model: OwnerModel): Owner {
        let ownerType: OwnerType = OwnerType.System; // Default
        if (model.type && isValidOwnerType(model.type)) {
            ownerType = model.type as OwnerType;
        }

        return new Owner(
            ownerType,
            model.etag || '',
            model.familyId ?? null,
            model.userId ?? null);
    }
}