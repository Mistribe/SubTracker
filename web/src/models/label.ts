import type {LabelModel} from "@/api/models";
import Owner from "@/models/owner.ts";

export default class Label {
    private readonly _id: string;
    private readonly _createdAt: Date;
    private readonly _etag: string;
    private readonly _owner: Owner;

    constructor(id: string,
                name: string,
                color: string,
                owner: Owner,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._color = color;
        this._owner = owner;
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

    private _color: string;

    get color(): string {
        return this._color;
    }

    set color(value: string) {
        this._color = value;
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

    get owner(): Owner {
        return this._owner;
    }

    get createdAt(): Date {
        return this._createdAt;
    }

    get etag(): string {
        return this._etag;
    }

    static fromModel(labelModel: LabelModel): Label {
        return new Label(
            labelModel.id || '',
            labelModel.name || '',
            labelModel.color || '',
            Owner.fromModel(labelModel.owner || {}),
            labelModel.createdAt ? labelModel.createdAt : new Date(),
            labelModel.updatedAt ? labelModel.updatedAt : new Date(),
            labelModel.etag || ''
        );
    }
}