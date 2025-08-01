// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './patch_family_member_model_type.dart';

/// auto generated
/// Model for updating family member details
class PatchFamilyMemberModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Optional member ID. If not provided, new member will be created
    String? id;
    ///  member's name
    String? name;
    ///  Indicates if the member is a kid
    PatchFamilyMemberModelType? type_;
    ///  Optional timestamp of the last update
    DateTime? updatedAt;
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
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getEnumValue<PatchFamilyMemberModelType>((stringValue) => PatchFamilyMemberModelType.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('id', id);
        writer.writeStringValue('name', name);
        writer.writeEnumValue<PatchFamilyMemberModelType>('type', type_, (e) => e?.value);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
