// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './editable_subscription_payer_model_type.dart';

/// auto generated
/// Subscription payer object used for updating who pays for a subscription
class EditableSubscriptionPayerModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description ID of the family associated with this payer
    String? familyId;
    ///  @Description ID of the specific family member who pays (required when type is family_member)
    String? memberId;
    ///  @Description Type of payer (family or family member)
    EditableSubscriptionPayerModelType? type_;
    /// Instantiates a new [EditableSubscriptionPayerModel] and sets the default values.
    EditableSubscriptionPayerModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static EditableSubscriptionPayerModel createFromDiscriminatorValue(ParseNode parseNode) {
        return EditableSubscriptionPayerModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['memberId'] = (node) => memberId = node.getStringValue();
        deserializerMap['type'] = (node) => type_ = node.getEnumValue<EditableSubscriptionPayerModelType>((stringValue) => EditableSubscriptionPayerModelType.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('memberId', memberId);
        writer.writeEnumValue<EditableSubscriptionPayerModelType>('type', type_, (e) => e?.value);
        writer.writeAdditionalData(additionalData);
    }
}
