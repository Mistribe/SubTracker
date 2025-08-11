// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../models/currency_rates_model.dart';
import '../../models/http_error_response.dart';
import './rates_request_builder_get_query_parameters.dart';
import './refresh/refresh_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \currencies\rates
class RatesRequestBuilder extends BaseRequestBuilder<RatesRequestBuilder> {
    ///  The refresh property
    RefreshRequestBuilder get refresh {
        return RefreshRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    RatesRequestBuilder clone() {
        return RatesRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [RatesRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    RatesRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies/rates{?date*}", pathParameters) ;
    /// Instantiates a new [RatesRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    RatesRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies/rates{?date*}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Get exchange rates for all currencies at a specific date
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<CurrencyRatesModel?> getAsync([void Function(RequestConfiguration<RatesRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<CurrencyRatesModel>(requestInfo, CurrencyRatesModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Get exchange rates for all currencies at a specific date
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<RatesRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<RatesRequestBuilderGetQueryParameters>(requestConfiguration, () => RatesRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
