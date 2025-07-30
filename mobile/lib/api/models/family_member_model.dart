// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Family member object containing member information
class FamilyMemberModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description Timestamp when the member was created
    DateTime? createdAt;
    ///  @Description Entity tag for optimistic concurrency control
    String? etag;
    ///  @Description ID of the family this member belongs to
    String? familyId;
    ///  @Description Unique identifier for the family member
    String? id;
    ///  @Description Whether this member is a child (affects permissions and features)
    bool? isKid;
    ///  @Description Name of the family member
    String? name;
    ///  @Description Timestamp when the member was last updated
    DateTime? updatedAt;
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
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['is_kid'] = (node) => isKid = node.getBoolValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_kid', value:isKid);
        writer.writeStringValue('name', name);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
