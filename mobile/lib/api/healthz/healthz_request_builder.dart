// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './live/live_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \healthz
class HealthzRequestBuilder extends BaseRequestBuilder<HealthzRequestBuilder> {
    ///  The live property
    LiveRequestBuilder get live {
        return LiveRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    HealthzRequestBuilder clone() {
        return HealthzRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [HealthzRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    HealthzRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/healthz", pathParameters) ;
    /// Instantiates a new [HealthzRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    HealthzRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/healthz", {RequestInformation.rawUrlKey : rawUrl}) ;
}
