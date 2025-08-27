// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../../../models/ginx/http_error_response.dart';
import '../../../../../models/provider/create_price_model.dart';
import '../../../../../models/provider/price_model.dart';
import './item/with_price_item_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \providers\{providerId}\plans\{planId}\prices
class PricesRequestBuilder extends BaseRequestBuilder<PricesRequestBuilder> {
    /// Gets an item from the ApiSdk.providers.item.plans.item.prices.item collection
    ///  [priceId] Price ID (UUID format)
    WithPriceItemRequestBuilder byPriceId(String priceId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('priceId', () => priceId);
        return WithPriceItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    PricesRequestBuilder clone() {
        return PricesRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [PricesRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PricesRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans/{planId}/prices", pathParameters) ;
    /// Instantiates a new [PricesRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PricesRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans/{planId}/prices", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Create a new pricing option for a specific provider plan
    ///  [body] Price information for a plan
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PriceModel?> postAsync(CreatePriceModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PriceModel>(requestInfo, PriceModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new pricing option for a specific provider plan
    ///  [body] Price information for a plan
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreatePriceModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
