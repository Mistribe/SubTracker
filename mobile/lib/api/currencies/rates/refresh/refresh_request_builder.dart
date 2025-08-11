// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/currency_rates_model.dart';
import '../../../models/http_error_response.dart';

/// auto generated
/// Builds and executes requests for operations under \currencies\rates\refresh
class RefreshRequestBuilder extends BaseRequestBuilder<RefreshRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    RefreshRequestBuilder clone() {
        return RefreshRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [RefreshRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    RefreshRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies/rates/refresh", pathParameters) ;
    /// Instantiates a new [RefreshRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    RefreshRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/currencies/rates/refresh", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// refresh exchange rates for all currencies
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<CurrencyRatesModel?> postAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<CurrencyRatesModel>(requestInfo, CurrencyRatesModel.createFromDiscriminatorValue, errorMapping);
    }
    /// refresh exchange rates for all currencies
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
