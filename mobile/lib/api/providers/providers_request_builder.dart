// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../models/create_provider_model.dart';
import '../models/http_error.dart';
import '../models/provider_model.dart';
import './item/with_provider_item_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \providers
class ProvidersRequestBuilder extends BaseRequestBuilder<ProvidersRequestBuilder> {
    /// Gets an item from the ApiSdk.providers.item collection
    ///  [providerId] Unique identifier of the item
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
    ProvidersRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers", pathParameters) ;
    /// Instantiates a new [ProvidersRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ProvidersRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Create a new provider with plans and prices
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<ProviderModel?> postAsync(CreateProviderModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<ProviderModel>(requestInfo, ProviderModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new provider with plans and prices
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
