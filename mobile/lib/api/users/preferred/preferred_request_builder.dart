// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './currency/currency_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \users\preferred
class PreferredRequestBuilder extends BaseRequestBuilder<PreferredRequestBuilder> {
    ///  The currency property
    CurrencyRequestBuilder get currency {
        return CurrencyRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    PreferredRequestBuilder clone() {
        return PreferredRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [PreferredRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PreferredRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users/preferred", pathParameters) ;
    /// Instantiates a new [PreferredRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PreferredRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users/preferred", {RequestInformation.rawUrlKey : rawUrl}) ;
}
