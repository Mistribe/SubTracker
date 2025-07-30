// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../../models/http_error.dart';
import '../../../../models/provider_model.dart';
import './with_provider_get_request_body.dart';

/// auto generated
/// Builds and executes requests for operations under \api\v1\providers\{providerId}
class WithProviderItemRequestBuilder extends BaseRequestBuilder<WithProviderItemRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    WithProviderItemRequestBuilder clone() {
        return WithProviderItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithProviderItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithProviderItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api/v1/providers/{providerId}", pathParameters) ;
    /// Instantiates a new [WithProviderItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithProviderItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api/v1/providers/{providerId}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Returns a single provider by its ID
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<ProviderModel?> getAsync(WithProviderGetRequestBody body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<ProviderModel>(requestInfo, ProviderModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Returns a single provider by its ID
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation(WithProviderGetRequestBody body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
