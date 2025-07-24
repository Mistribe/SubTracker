// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class UpdateLabelModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The color property
    String? color;
    ///  The name property
    String? name;
    ///  The updated_at property
    String? updatedAt;
    /// Instantiates a new [UpdateLabelModel] and sets the default values.
    UpdateLabelModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static UpdateLabelModel createFromDiscriminatorValue(ParseNode parseNode) {
        return UpdateLabelModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['color'] = (node) => color = node.getStringValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('color', color);
        writer.writeStringValue('name', name);
        writer.writeStringValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
