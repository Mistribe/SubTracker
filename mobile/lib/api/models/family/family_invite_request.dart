// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './family_invite_request_type.dart';

/// auto generated
class FamilyInviteRequest implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Email of the invited member
    String? email;
    ///  ID of the family member to be invited
    String? familyMemberId;
    ///  Name of the invited member
    String? name;
    ///  Type of the member (adult or kid)
    FamilyInviteRequestType? type_;
    /// Instantiates a new [FamilyInviteRequest] and sets the default values.
    FamilyInviteRequest() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static FamilyInviteRequest createFromDiscriminatorValue(ParseNode parseNode) {
        return FamilyInviteRequest();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['email'] = (node) => email = node.getStringValue();
        deserializerMap['family_member_id'] = (node) => familyMemberId = node.getStringValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getEnumValue<FamilyInviteRequestType>((stringValue) => FamilyInviteRequestType.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('email', email);
        writer.writeStringValue('family_member_id', familyMemberId);
        writer.writeStringValue('name', name);
        writer.writeEnumValue<FamilyInviteRequestType>('type', type_, (e) => e?.value);
        writer.writeAdditionalData(additionalData);
    }
}
