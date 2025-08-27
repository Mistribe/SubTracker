// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../../../../models/ginx/http_error_response.dart';
import '../../../../../../models/provider/price_model.dart';
import '../../../../../../models/provider/update_price_model.dart';

/// auto generated
/// Builds and executes requests for operations under \providers\{providerId}\plans\{planId}\prices\{priceId}
class WithPriceItemRequestBuilder extends BaseRequestBuilder<WithPriceItemRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    WithPriceItemRequestBuilder clone() {
        return WithPriceItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithPriceItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithPriceItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans/{planId}/prices/{priceId}", pathParameters) ;
    /// Instantiates a new [WithPriceItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithPriceItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans/{planId}/prices/{priceId}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Permanently delete a specific price from a provider plan
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<void> deleteAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toDeleteRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.sendNoContent(requestInfo, errorMapping);
    }
    /// Update an existing price for a specific provider plan
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PriceModel?> putAsync(UpdatePriceModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PriceModel>(requestInfo, PriceModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Permanently delete a specific price from a provider plan
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toDeleteRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.delete, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Update an existing price for a specific provider plan
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPutRequestInformation(UpdatePriceModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.put, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
