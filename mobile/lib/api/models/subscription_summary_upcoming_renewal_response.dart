// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class SubscriptionSummaryUpcomingRenewalResponse implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The at property
    DateTime? at;
    ///  The provider_id property
    String? providerId;
    ///  The total property
    double? total;
    /// Instantiates a new [SubscriptionSummaryUpcomingRenewalResponse] and sets the default values.
    SubscriptionSummaryUpcomingRenewalResponse() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionSummaryUpcomingRenewalResponse createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionSummaryUpcomingRenewalResponse();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['at'] = (node) => at = node.getDateTimeValue();
        deserializerMap['provider_id'] = (node) => providerId = node.getStringValue();
        deserializerMap['total'] = (node) => total = node.getDoubleValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('at', at);
        writer.writeStringValue('provider_id', providerId);
        writer.writeDoubleValue('total', total);
        writer.writeAdditionalData(additionalData);
    }
}
