// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class CreateLabelModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The color property
    String? color;
    ///  The created_at property
    String? createdAt;
    ///  The id property
    String? id;
    ///  The is_default property
    bool? isDefault;
    ///  The name property
    String? name;
    /// Instantiates a new [CreateLabelModel] and sets the default values.
    CreateLabelModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static CreateLabelModel createFromDiscriminatorValue(ParseNode parseNode) {
        return CreateLabelModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['color'] = (node) => color = node.getStringValue();
        deserializerMap['created_at'] = (node) => createdAt = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['is_default'] = (node) => isDefault = node.getBoolValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('color', color);
        writer.writeStringValue('created_at', createdAt);
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_default', value:isDefault);
        writer.writeStringValue('name', name);
        writer.writeAdditionalData(additionalData);
    }
}
