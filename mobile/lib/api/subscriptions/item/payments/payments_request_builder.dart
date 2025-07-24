// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/create_subscription_payment_model.dart';
import '../../../models/http_error.dart';
import '../../../models/subscription_model.dart';
import './item/with_payment_item_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \subscriptions\{-id}\payments
class PaymentsRequestBuilder extends BaseRequestBuilder<PaymentsRequestBuilder> {
    /// Gets an item from the ApiSdk.subscriptions.item.payments.item collection
    ///  [paymentId] Payment ID
    WithPaymentItemRequestBuilder byPaymentId(String paymentId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('paymentId', () => paymentId);
        return WithPaymentItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    PaymentsRequestBuilder clone() {
        return PaymentsRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [PaymentsRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PaymentsRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/{%2Did}/payments", pathParameters) ;
    /// Instantiates a new [PaymentsRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PaymentsRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/{%2Did}/payments", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Create a new subscription payment
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionModel?> postAsync(CreateSubscriptionPaymentModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionModel>(requestInfo, SubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new subscription payment
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreateSubscriptionPaymentModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
