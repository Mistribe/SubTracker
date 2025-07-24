// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../../models/http_error.dart';
import '../../../../models/subscription_model.dart';
import '../../../../models/update_payment_model.dart';
import './with_payment_delete_request_body.dart';

/// auto generated
/// Builds and executes requests for operations under \subscriptions\{-id}\payments\{paymentId}
class WithPaymentItemRequestBuilder extends BaseRequestBuilder<WithPaymentItemRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    WithPaymentItemRequestBuilder clone() {
        return WithPaymentItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithPaymentItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithPaymentItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/{%2Did}/payments/{paymentId}", pathParameters) ;
    /// Instantiates a new [WithPaymentItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithPaymentItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/{%2Did}/payments/{paymentId}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Delete a subscription payment
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<void> deleteAsync(WithPaymentDeleteRequestBody body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toDeleteRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.sendNoContent(requestInfo, errorMapping);
    }
    /// Update payment details for a specific subscription
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionModel?> putAsync(UpdatePaymentModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '404' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionModel>(requestInfo, SubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Delete a subscription payment
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toDeleteRequestInformation(WithPaymentDeleteRequestBody body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.delete, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
    /// Update payment details for a specific subscription
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPutRequestInformation(UpdatePaymentModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.put, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
