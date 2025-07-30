// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './live_get_response.dart';

/// auto generated
/// Builds and executes requests for operations under \healthz\live
class LiveRequestBuilder extends BaseRequestBuilder<LiveRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    LiveRequestBuilder clone() {
        return LiveRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [LiveRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    LiveRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/healthz/live", pathParameters) ;
    /// Instantiates a new [LiveRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    LiveRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/healthz/live", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Returns the health status of the application
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<LiveGetResponse?> getAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        return await requestAdapter.send<LiveGetResponse>(requestInfo, LiveGetResponse.createFromDiscriminatorValue, {});
    }
    /// Returns the health status of the application
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
