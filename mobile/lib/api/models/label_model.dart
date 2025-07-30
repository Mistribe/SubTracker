// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './owner_model.dart';

/// auto generated
class LabelModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The color property
    String? color;
    ///  The created_at property
    DateTime? createdAt;
    ///  The etag property
    String? etag;
    ///  The id property
    String? id;
    ///  The name property
    String? name;
    ///  The owner property
    OwnerModel? owner;
    ///  The updated_at property
    DateTime? updatedAt;
    /// Instantiates a new [LabelModel] and sets the default values.
    LabelModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static LabelModel createFromDiscriminatorValue(ParseNode parseNode) {
        return LabelModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['color'] = (node) => color = node.getStringValue();
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['owner'] = (node) => owner = node.getObjectValue<OwnerModel>(OwnerModel.createFromDiscriminatorValue);
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('color', color);
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('id', id);
        writer.writeStringValue('name', name);
        writer.writeObjectValue<OwnerModel>('owner', owner);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
