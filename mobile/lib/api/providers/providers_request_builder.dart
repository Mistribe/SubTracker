// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../models/dto/paginated_response_model_provider_provider_model.dart';
import '../models/ginx/http_error_response.dart';
import '../models/provider/create_provider_model.dart';
import '../models/provider/provider_model.dart';
import './item/with_provider_item_request_builder.dart';
import './providers_request_builder_get_query_parameters.dart';

/// auto generated
/// Builds and executes requests for operations under \providers
class ProvidersRequestBuilder extends BaseRequestBuilder<ProvidersRequestBuilder> {
    /// Gets an item from the ApiSdk.providers.item collection
    ///  [providerId] Provider ID (UUID format)
    WithProviderItemRequestBuilder byProviderId(String providerId) {
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
    ProvidersRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers{?limit*,offset*,search*}", pathParameters) ;
    /// Instantiates a new [ProvidersRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ProvidersRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers{?limit*,offset*,search*}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Retrieve a paginated list of all providers with their plans and prices
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PaginatedResponseModelProviderProviderModel?> getAsync([void Function(RequestConfiguration<ProvidersRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PaginatedResponseModelProviderProviderModel>(requestInfo, PaginatedResponseModelProviderProviderModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new service provider with labels and owner information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<ProviderModel?> postAsync(CreateProviderModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<ProviderModel>(requestInfo, ProviderModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Retrieve a paginated list of all providers with their plans and prices
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<ProvidersRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<ProvidersRequestBuilderGetQueryParameters>(requestConfiguration, () => ProvidersRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Create a new service provider with labels and owner information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreateProviderModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
