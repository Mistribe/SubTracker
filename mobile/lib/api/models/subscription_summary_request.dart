// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Request parameters for subscription summary
class SubscriptionSummaryRequest implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The top_providers property
    int? topProviders;
    ///  The total_monthly property
    bool? totalMonthly;
    ///  The total_yearly property
    bool? totalYearly;
    ///  The upcoming_renewals property
    int? upcomingRenewals;
    /// Instantiates a new [SubscriptionSummaryRequest] and sets the default values.
    SubscriptionSummaryRequest() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionSummaryRequest createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionSummaryRequest();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['top_providers'] = (node) => topProviders = node.getIntValue();
        deserializerMap['total_monthly'] = (node) => totalMonthly = node.getBoolValue();
        deserializerMap['total_yearly'] = (node) => totalYearly = node.getBoolValue();
        deserializerMap['upcoming_renewals'] = (node) => upcomingRenewals = node.getIntValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeIntValue('top_providers', topProviders);
        writer.writeBoolValue('total_monthly', value:totalMonthly);
        writer.writeBoolValue('total_yearly', value:totalYearly);
        writer.writeIntValue('upcoming_renewals', upcomingRenewals);
        writer.writeAdditionalData(additionalData);
    }
}
