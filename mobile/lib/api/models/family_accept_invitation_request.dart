// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class FamilyAcceptInvitationRequest implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  ID of the family member accepting the invitation
    String? familyMemberId;
    ///  Code received in the invitation
    String? invitationCode;
    /// Instantiates a new [FamilyAcceptInvitationRequest] and sets the default values.
    FamilyAcceptInvitationRequest() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static FamilyAcceptInvitationRequest createFromDiscriminatorValue(ParseNode parseNode) {
        return FamilyAcceptInvitationRequest();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['family_member_id'] = (node) => familyMemberId = node.getStringValue();
        deserializerMap['invitation_code'] = (node) => invitationCode = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('family_member_id', familyMemberId);
        writer.writeStringValue('invitation_code', invitationCode);
        writer.writeAdditionalData(additionalData);
    }
}
