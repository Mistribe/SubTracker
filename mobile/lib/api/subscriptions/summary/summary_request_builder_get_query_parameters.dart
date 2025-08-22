// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Returns summary information about subscriptions including total costs and upcoming renewals
class SummaryRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Number of top labels to return
    /// @QueryParameter('top_labels')
    int? topLabels;
    ///  Number of top providers to return
    /// @QueryParameter('top_providers')
    int? topProviders;
    ///  Include monthly total costs
    /// @QueryParameter('total_monthly')
    bool? totalMonthly;
    ///  Include yearly total costs
    /// @QueryParameter('total_yearly')
    bool? totalYearly;
    ///  Number of upcoming renewals to return
    /// @QueryParameter('upcoming_renewals')
    int? upcomingRenewals;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'top_labels' : topLabels,
            'top_providers' : topProviders,
            'total_monthly' : totalMonthly,
            'total_yearly' : totalYearly,
            'upcoming_renewals' : upcomingRenewals,
        };
    }
}
