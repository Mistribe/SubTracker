import type {LabelModel} from "@/api/models";

export default class Label {
    private readonly _id: string;
    private _name: string;
    private _isDefault: boolean;
    private _color: string;
    private readonly _createdAt: Date;
    private _updatedAt: Date;
    private readonly _etag: string;


    constructor(id: string,
                name: string,
                isDefault: boolean,
                color: string,
                createdAt: Date,
                updatedAt: Date,
                etag: string) {
        this._id = id;
        this._name = name;
        this._isDefault = isDefault;
        this._color = color;
        this._createdAt = createdAt;
        this._updatedAt = updatedAt;
        this._etag = etag;
    }

    static fromModel(labelModel: LabelModel): Label {
        return new Label(
            labelModel.id || '',
            labelModel.name || '',
            labelModel.isDefault || false,
            labelModel.color || '',
            labelModel.createdAt ? labelModel.createdAt : new Date(),
            labelModel.updatedAt ? labelModel.updatedAt : new Date(),
            labelModel.etag || ''
        );
    }


    get id(): string {
        return this._id;
    }

    get name(): string {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    get isDefault(): boolean {
        return this._isDefault;
    }

    set isDefault(value: boolean) {
        this._isDefault = value;
    }

    get color(): string {
        return this._color;
    }

    set color(value: string) {
        this._color = value;
    }

    get createdAt(): Date {
        return this._createdAt;
    }

    get updatedAt(): Date {
        return this._updatedAt;
    }

    set updatedAt(value: Date) {
        this._updatedAt = value;
    }

    get etag(): string {
        return this._etag;
    }
}