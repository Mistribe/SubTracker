// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../models/create_subscription_model.dart';
import '../models/http_error.dart';
import '../models/paginated_response_model_endpoints_subscription_model.dart';
import '../models/patch_subscription_model.dart';
import '../models/subscription_model.dart';
import './item/with_subscription_item_request_builder.dart';
import './subscriptions_request_builder_get_query_parameters.dart';

/// auto generated
/// Builds and executes requests for operations under \subscriptions
class SubscriptionsRequestBuilder extends BaseRequestBuilder<SubscriptionsRequestBuilder> {
    /// Gets an item from the ApiSdk.subscriptions.item collection
    ///  [subscriptionId] Subscription ID (UUID format)
    WithSubscriptionItemRequestBuilder bySubscriptionId(String subscriptionId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('subscriptionId', () => subscriptionId);
        return WithSubscriptionItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    SubscriptionsRequestBuilder clone() {
        return SubscriptionsRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [SubscriptionsRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    SubscriptionsRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions{?limit*,offset*}", pathParameters) ;
    /// Instantiates a new [SubscriptionsRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    SubscriptionsRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions{?limit*,offset*}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Retrieve a paginated list of all subscriptions for the authenticated user
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PaginatedResponseModelEndpointsSubscriptionModel?> getAsync([void Function(RequestConfiguration<SubscriptionsRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '500' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PaginatedResponseModelEndpointsSubscriptionModel>(requestInfo, PaginatedResponseModelEndpointsSubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Update or create a subscription with complete details. If subscription doesn't exist, it will be created.
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionModel?> patchAsync(PatchSubscriptionModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPatchRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '401' :  HttpError.createFromDiscriminatorValue,
            '404' :  HttpError.createFromDiscriminatorValue,
            '500' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionModel>(requestInfo, SubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new subscription with provider, plan, pricing, and payment information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionModel?> postAsync(CreateSubscriptionModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '401' :  HttpError.createFromDiscriminatorValue,
            '500' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionModel>(requestInfo, SubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Retrieve a paginated list of all subscriptions for the authenticated user
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<SubscriptionsRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<SubscriptionsRequestBuilderGetQueryParameters>(requestConfiguration, () => SubscriptionsRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Update or create a subscription with complete details. If subscription doesn't exist, it will be created.
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPatchRequestInformation(PatchSubscriptionModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.patch, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
    /// Create a new subscription with provider, plan, pricing, and payment information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreateSubscriptionModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
