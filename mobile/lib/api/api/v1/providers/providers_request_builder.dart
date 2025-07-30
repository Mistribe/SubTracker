// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import 'package:uuid/uuid.dart';
import '../../../models/paginated_response_model_endpoints_provider_model.dart';
import './item/with_provider_item_request_builder.dart';
import './providers_get_request_body.dart';
import './providers_request_builder_get_query_parameters.dart';

/// auto generated
/// Builds and executes requests for operations under \api\v1\providers
class ProvidersRequestBuilder extends BaseRequestBuilder<ProvidersRequestBuilder> {
    /// Gets an item from the ApiSdk.api.v1.providers.item collection
    ///  [providerId] Provider ID
    WithProviderItemRequestBuilder byProviderId(UuidValue providerId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('providerId', () => providerId);
        return WithProviderItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    ProvidersRequestBuilder clone() {
        return ProvidersRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [ProvidersRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ProvidersRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api/v1/providers{?page*,size*}", pathParameters) ;
    /// Instantiates a new [ProvidersRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ProvidersRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/api/v1/providers{?page*,size*}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Returns a paginated list of providers
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PaginatedResponseModelEndpointsProviderModel?> getAsync(ProvidersGetRequestBody body, [void Function(RequestConfiguration<ProvidersRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(body, requestConfiguration);
        return await requestAdapter.send<PaginatedResponseModelEndpointsProviderModel>(requestInfo, PaginatedResponseModelEndpointsProviderModel.createFromDiscriminatorValue, {});
    }
    /// Returns a paginated list of providers
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation(ProvidersGetRequestBody body, [void Function(RequestConfiguration<ProvidersRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<ProvidersRequestBuilderGetQueryParameters>(requestConfiguration, () => ProvidersRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
