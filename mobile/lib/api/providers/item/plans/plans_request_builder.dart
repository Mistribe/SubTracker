// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/create_plan_model.dart';
import '../../../models/http_error.dart';
import '../../../models/plan_model.dart';
import './item/with_plan_item_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \providers\{providerId}\plans
class PlansRequestBuilder extends BaseRequestBuilder<PlansRequestBuilder> {
    /// Gets an item from the ApiSdk.providers.item.plans.item collection
    ///  [planId] Plan ID (UUID format)
    WithPlanItemRequestBuilder byPlanId(String planId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('planId', () => planId);
        return WithPlanItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    PlansRequestBuilder clone() {
        return PlansRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [PlansRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PlansRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans", pathParameters) ;
    /// Instantiates a new [PlansRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    PlansRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/providers/{providerId}/plans", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Create a new subscription plan for an existing provider
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PlanModel?> postAsync(CreatePlanModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '404' :  HttpError.createFromDiscriminatorValue,
            '500' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PlanModel>(requestInfo, PlanModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new subscription plan for an existing provider
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreatePlanModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
