// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './owner_model_type.dart';

/// auto generated
/// @Description Ownership information specifying whether this subscription belongs to a user or family
class OwnerModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description Entity tag for optimistic concurrency control
    String? etag;
    ///  @Description Family ID when an ownership type is family (required for family ownership)
    String? familyId;
    ///  @Description Type of ownership (personal, family or system)
    OwnerModelType? type_;
    ///  @Description UserProfile ID when an ownership type is personal (required for personal ownership)
    String? userId;
    /// Instantiates a new [OwnerModel] and sets the default values.
    OwnerModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static OwnerModel createFromDiscriminatorValue(ParseNode parseNode) {
        return OwnerModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getEnumValue<OwnerModelType>((stringValue) => OwnerModelType.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        deserializerMap['userId'] = (node) => userId = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('family_id', familyId);
        writer.writeEnumValue<OwnerModelType>('type', type_, (e) => e?.value);
        writer.writeStringValue('userId', userId);
        writer.writeAdditionalData(additionalData);
    }
}
