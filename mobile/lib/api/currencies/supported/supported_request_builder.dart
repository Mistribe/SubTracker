// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Builds and executes requests for operations under \currencies\supported
class SupportedRequestBuilder extends BaseRequestBuilder<SupportedRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    SupportedRequestBuilder clone() {
        return SupportedRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [SupportedRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    SupportedRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies/supported", pathParameters) ;
    /// Instantiates a new [SupportedRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    SupportedRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies/supported", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// get details of all supported currencies
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<Iterable<String>?> getAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        return await requestAdapter.sendPrimitiveCollection<String>(requestInfo, {});
    }
    /// get details of all supported currencies
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
