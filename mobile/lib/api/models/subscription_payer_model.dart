// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './subscription_payer_model_type.dart';

/// auto generated
/// @Description Information about who pays for this subscription within the family
class SubscriptionPayerModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
    String? etag;
    ///  @Description ID of the family associated with this payer
    String? familyId;
    ///  @Description ID of the specific family member who pays (required when type is family_member)
    String? memberId;
    ///  @Description Type of payer (family or family member)
    SubscriptionPayerModelType? type_;
    /// Instantiates a new [SubscriptionPayerModel] and sets the default values.
    SubscriptionPayerModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionPayerModel createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionPayerModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['memberId'] = (node) => memberId = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getEnumValue<SubscriptionPayerModelType>((stringValue) => SubscriptionPayerModelType.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('memberId', memberId);
        writer.writeEnumValue<SubscriptionPayerModelType>('type', type_, (e) => e?.value);
        writer.writeAdditionalData(additionalData);
    }
}
