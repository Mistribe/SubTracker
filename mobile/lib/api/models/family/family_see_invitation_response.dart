// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './family_model.dart';

/// auto generated
class FamilySeeInvitationResponse implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Family details
    FamilyModel? family;
    ///  Role of the invited member
    String? invitedInasmuchAs;
    /// Instantiates a new [FamilySeeInvitationResponse] and sets the default values.
    FamilySeeInvitationResponse() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static FamilySeeInvitationResponse createFromDiscriminatorValue(ParseNode parseNode) {
        return FamilySeeInvitationResponse();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['family'] = (node) => family = node.getObjectValue<FamilyModel>(FamilyModel.createFromDiscriminatorValue);
        deserializerMap['invited_inasmuch_as'] = (node) => invitedInasmuchAs = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeObjectValue<FamilyModel>('family', family);
        writer.writeStringValue('invited_inasmuch_as', invitedInasmuchAs);
        writer.writeAdditionalData(additionalData);
    }
}
