// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './v1/v1_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \api
class ApiRequestBuilder extends BaseRequestBuilder<ApiRequestBuilder> {
    ///  The v1 property
    V1RequestBuilder get v1 {
        return V1RequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    ApiRequestBuilder clone() {
        return ApiRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [ApiRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ApiRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api", pathParameters) ;
    /// Instantiates a new [ApiRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ApiRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api", {RequestInformation.rawUrlKey : rawUrl}) ;
}
