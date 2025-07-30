// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './providers/providers_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \api\v1
class V1RequestBuilder extends BaseRequestBuilder<V1RequestBuilder> {
    ///  The providers property
    ProvidersRequestBuilder get providers {
        return ProvidersRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    V1RequestBuilder clone() {
        return V1RequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [V1RequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    V1RequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api/v1", pathParameters) ;
    /// Instantiates a new [V1RequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    V1RequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api/v1", {RequestInformation.rawUrlKey : rawUrl}) ;
}
