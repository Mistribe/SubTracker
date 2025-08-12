// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../models/http_error_response.dart';
import '../../models/subscription_model.dart';
import '../../models/update_subscription_model.dart';

/// auto generated
/// Builds and executes requests for operations under \subscriptions\{subscriptionId}
class WithSubscriptionItemRequestBuilder extends BaseRequestBuilder<WithSubscriptionItemRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    WithSubscriptionItemRequestBuilder clone() {
        return WithSubscriptionItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithSubscriptionItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithSubscriptionItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/{subscriptionId}", pathParameters) ;
    /// Instantiates a new [WithSubscriptionItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithSubscriptionItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/subscriptions/{subscriptionId}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Permanently delete an existing subscription
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
    /// Retrieve a single subscription with all its details including provider, plan, and pricing information
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionModel?> getAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionModel>(requestInfo, SubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Update an existing subscription's details including provider, plan, pricing, and payment information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<SubscriptionModel?> putAsync(UpdateSubscriptionModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<SubscriptionModel>(requestInfo, SubscriptionModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Permanently delete an existing subscription
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toDeleteRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.delete, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Retrieve a single subscription with all its details including provider, plan, and pricing information
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Update an existing subscription's details including provider, plan, pricing, and payment information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPutRequestInformation(UpdateSubscriptionModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.put, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
