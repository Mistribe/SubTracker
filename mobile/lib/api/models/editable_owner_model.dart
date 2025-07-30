// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class EditableOwnerModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The family_id property
    String? familyId;
    ///  The type property
    String? type_;
    /// Instantiates a new [EditableOwnerModel] and sets the default values.
    EditableOwnerModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static EditableOwnerModel createFromDiscriminatorValue(ParseNode parseNode) {
        return EditableOwnerModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('type', type_);
        writer.writeAdditionalData(additionalData);
    }
}
