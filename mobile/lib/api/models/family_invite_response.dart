// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class FamilyInviteResponse implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The code property
    String? code;
    ///  The family_id property
    String? familyId;
    ///  The family_member_id property
    String? familyMemberId;
    /// Instantiates a new [FamilyInviteResponse] and sets the default values.
    FamilyInviteResponse() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static FamilyInviteResponse createFromDiscriminatorValue(ParseNode parseNode) {
        return FamilyInviteResponse();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['code'] = (node) => code = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['family_member_id'] = (node) => familyMemberId = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('code', code);
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('family_member_id', familyMemberId);
        writer.writeAdditionalData(additionalData);
    }
}
