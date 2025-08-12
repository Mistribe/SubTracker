// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/http_error_response.dart';
import '../../../models/update_preferred_currency_model.dart';
import '../../../models/user_preferred_currency_model.dart';

/// auto generated
/// Builds and executes requests for operations under \users\preferred\currency
class CurrencyRequestBuilder extends BaseRequestBuilder<CurrencyRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    CurrencyRequestBuilder clone() {
        return CurrencyRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [CurrencyRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    CurrencyRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users/preferred/currency", pathParameters) ;
    /// Instantiates a new [CurrencyRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    CurrencyRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users/preferred/currency", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Returns the preferred currency for the authenticated user
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<UserPreferredCurrencyModel?> getAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<UserPreferredCurrencyModel>(requestInfo, UserPreferredCurrencyModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Updates the preferred currency for the authenticated user
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<void> putAsync(UpdatePreferredCurrencyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.sendNoContent(requestInfo, errorMapping);
    }
    /// Returns the preferred currency for the authenticated user
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Updates the preferred currency for the authenticated user
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPutRequestInformation(UpdatePreferredCurrencyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.put, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
