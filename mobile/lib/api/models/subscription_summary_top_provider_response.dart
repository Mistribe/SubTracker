// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class SubscriptionSummaryTopProviderResponse implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The provider_id property
    String? providerId;
    ///  The total property
    double? total;
    /// Instantiates a new [SubscriptionSummaryTopProviderResponse] and sets the default values.
    SubscriptionSummaryTopProviderResponse() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionSummaryTopProviderResponse createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionSummaryTopProviderResponse();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['provider_id'] = (node) => providerId = node.getStringValue();
        deserializerMap['total'] = (node) => total = node.getDoubleValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('provider_id', providerId);
        writer.writeDoubleValue('total', total);
        writer.writeAdditionalData(additionalData);
    }
}
