// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../../models/http_error.dart';
import '../../../../models/price_model.dart';
import '../../../../models/update_price_model.dart';
import './prices/prices_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \providers\{providerId}\plans\{planId}
class WithPlanItemRequestBuilder extends BaseRequestBuilder<WithPlanItemRequestBuilder> {
    ///  The prices property
    PricesRequestBuilder get prices {
        return PricesRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    WithPlanItemRequestBuilder clone() {
        return WithPlanItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithPlanItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithPlanItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans/{planId}", pathParameters) ;
    /// Instantiates a new [WithPlanItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithPlanItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans/{planId}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Update a new price for a provider plan
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PriceModel?> putAsync(UpdatePriceModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PriceModel>(requestInfo, PriceModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Update a new price for a provider plan
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
