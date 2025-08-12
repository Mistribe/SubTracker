// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './rates/rates_request_builder.dart';
import './supported/supported_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \currencies
class CurrenciesRequestBuilder extends BaseRequestBuilder<CurrenciesRequestBuilder> {
    ///  The rates property
    RatesRequestBuilder get rates {
        return RatesRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The supported property
    SupportedRequestBuilder get supported {
        return SupportedRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    CurrenciesRequestBuilder clone() {
        return CurrenciesRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [CurrenciesRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    CurrenciesRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies", pathParameters) ;
    /// Instantiates a new [CurrenciesRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    CurrenciesRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies", {RequestInformation.rawUrlKey : rawUrl}) ;
}
