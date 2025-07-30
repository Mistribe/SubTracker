// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class SubscriptionPayerModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The etag property
    String? etag;
    ///  The family_id property
    String? familyId;
    ///  The memberId property
    String? memberId;
    ///  The type property
    String? type_;
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
        deserializerMap['type'] = (node) => type_ = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('family_id', familyId);
        writer.writeStringValue('memberId', memberId);
        writer.writeStringValue('type', type_);
        writer.writeAdditionalData(additionalData);
    }
}
