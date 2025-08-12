// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './subscription_summary_top_provider_response.dart';
import './subscription_summary_upcoming_renewal_response.dart';

/// auto generated
/// Response containing subscription summary information
class SubscriptionSummaryResponse implements AdditionalDataHolder, Parsable {
    ///  The active property
    int? active;
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The top_providers property
    Iterable<SubscriptionSummaryTopProviderResponse>? topProviders;
    ///  The total_monthly property
    double? totalMonthly;
    ///  The total_yearly property
    double? totalYearly;
    ///  The upcoming_renewals property
    Iterable<SubscriptionSummaryUpcomingRenewalResponse>? upcomingRenewals;
    /// Instantiates a new [SubscriptionSummaryResponse] and sets the default values.
    SubscriptionSummaryResponse() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionSummaryResponse createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionSummaryResponse();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['active'] = (node) => active = node.getIntValue();
        deserializerMap['top_providers'] = (node) => topProviders = node.getCollectionOfObjectValues<SubscriptionSummaryTopProviderResponse>(SubscriptionSummaryTopProviderResponse.createFromDiscriminatorValue);
        deserializerMap['total_monthly'] = (node) => totalMonthly = node.getDoubleValue();
        deserializerMap['total_yearly'] = (node) => totalYearly = node.getDoubleValue();
        deserializerMap['upcoming_renewals'] = (node) => upcomingRenewals = node.getCollectionOfObjectValues<SubscriptionSummaryUpcomingRenewalResponse>(SubscriptionSummaryUpcomingRenewalResponse.createFromDiscriminatorValue);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeIntValue('active', active);
        writer.writeCollectionOfObjectValues<SubscriptionSummaryTopProviderResponse>('top_providers', topProviders);
        writer.writeDoubleValue('total_monthly', totalMonthly);
        writer.writeDoubleValue('total_yearly', totalYearly);
        writer.writeCollectionOfObjectValues<SubscriptionSummaryUpcomingRenewalResponse>('upcoming_renewals', upcomingRenewals);
        writer.writeAdditionalData(additionalData);
    }
}
