// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class CreatePlanModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    DateTime? createdAt;
    ///  The description property
    String? description;
    ///  The id property
    String? id;
    ///  The name property
    String? name;
    /// Instantiates a new [CreatePlanModel] and sets the default values.
    CreatePlanModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static CreatePlanModel createFromDiscriminatorValue(ParseNode parseNode) {
        return CreatePlanModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['description'] = (node) => description = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('description', description);
        writer.writeStringValue('id', id);
        writer.writeStringValue('name', name);
        writer.writeAdditionalData(additionalData);
    }
}
