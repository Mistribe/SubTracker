// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './plans/plans_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \providers\{providerId}
class WithProviderItemRequestBuilder extends BaseRequestBuilder<WithProviderItemRequestBuilder> {
    ///  The plans property
    PlansRequestBuilder get plans {
        return PlansRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    WithProviderItemRequestBuilder clone() {
        return WithProviderItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithProviderItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithProviderItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}", pathParameters) ;
    /// Instantiates a new [WithProviderItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithProviderItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}", {RequestInformation.rawUrlKey : rawUrl}) ;
}
