// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../models/http_error_response.dart';
import '../../models/subscription_summary_response.dart';
import './summary_request_builder_get_query_parameters.dart';

/// auto generated
/// Builds and executes requests for operations under \subscriptions\summary
class SummaryRequestBuilder extends BaseRequestBuilder<SummaryRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    SummaryRequestBuilder clone() {
        return SummaryRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [SummaryRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    SummaryRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/summary?top_labels={top_labels}&top_providers={top_providers}&total_monthly={total_monthly}&total_yearly={total_yearly}&upcoming_renewals={upcoming_renewals}", pathParameters) ;
    /// Instantiates a new [SummaryRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    SummaryRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/summary?top_labels={top_labels}&top_providers={top_providers}&total_monthly={total_monthly}&total_yearly={total_yearly}&upcoming_renewals={upcoming_renewals}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Returns summary information about subscriptions including total costs and upcoming renewals
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionSummaryResponse?> getAsync([void Function(RequestConfiguration<SummaryRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionSummaryResponse>(requestInfo, SubscriptionSummaryResponse.createFromDiscriminatorValue, errorMapping);
    }
    /// Returns summary information about subscriptions including total costs and upcoming renewals
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<SummaryRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<SummaryRequestBuilderGetQueryParameters>(requestConfiguration, () => SummaryRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
