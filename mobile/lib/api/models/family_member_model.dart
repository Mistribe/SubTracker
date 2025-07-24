// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class FamilyMemberModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    String? createdAt;
    ///  The etag property
    String? etag;
    ///  The family_id property
    String? familyId;
    ///  The id property
    String? id;
    ///  The is_kid property
    bool? isKid;
    ///  The name property
    String? name;
    ///  The updated_at property
    String? updatedAt;
    /// Instantiates a new [FamilyMemberModel] and sets the default values.
    FamilyMemberModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static FamilyMemberModel createFromDiscriminatorValue(ParseNode parseNode) {
        return FamilyMemberModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getStringValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['is_kid'] = (node) => isKid = node.getBoolValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('created_at', createdAt);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_kid', value:isKid);
        writer.writeStringValue('name', name);
        writer.writeStringValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
