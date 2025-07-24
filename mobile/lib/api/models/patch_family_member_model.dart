// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Model for updating family member details
class PatchFamilyMemberModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Optional email address
    String? email;
    ///  Optional member ID. If not provided, new member will be created
    String? id;
    ///  Indicates if the member is a kid
    bool? isKid;
    ///  Member's name
    String? name;
    ///  Optional timestamp of the last update
    String? updatedAt;
    /// Instantiates a new [PatchFamilyMemberModel] and sets the default values.
    PatchFamilyMemberModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PatchFamilyMemberModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PatchFamilyMemberModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['email'] = (node) => email = node.getStringValue();
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
        writer.writeStringValue('email', email);
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_kid', value:isKid);
        writer.writeStringValue('name', name);
        writer.writeStringValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
