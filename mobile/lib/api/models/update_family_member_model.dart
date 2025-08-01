// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './update_family_member_model_type.dart';

/// auto generated
class UpdateFamilyMemberModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The name property
    String? name;
    ///  The type property
    UpdateFamilyMemberModelType? type_;
    ///  The updated_at property
    DateTime? updatedAt;
    /// Instantiates a new [UpdateFamilyMemberModel] and sets the default values.
    UpdateFamilyMemberModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static UpdateFamilyMemberModel createFromDiscriminatorValue(ParseNode parseNode) {
        return UpdateFamilyMemberModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getEnumValue<UpdateFamilyMemberModelType>((stringValue) => UpdateFamilyMemberModelType.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('name', name);
        writer.writeEnumValue<UpdateFamilyMemberModelType>('type', type_, (e) => e?.value);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
