// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class UpdateProfileModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The family_name property
    String? familyName;
    ///  The given_name property
    String? givenName;
    /// Instantiates a new [UpdateProfileModel] and sets the default values.
    UpdateProfileModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static UpdateProfileModel createFromDiscriminatorValue(ParseNode parseNode) {
        return UpdateProfileModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['family_name'] = (node) => familyName = node.getStringValue();
        deserializerMap['given_name'] = (node) => givenName = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('family_name', familyName);
        writer.writeStringValue('given_name', givenName);
        writer.writeAdditionalData(additionalData);
    }
}
