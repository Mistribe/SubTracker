// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../models/create_family_model.dart';
import '../models/family_model.dart';
import '../models/http_error.dart';
import '../models/paginated_response_model_endpoints_family_model.dart';
import '../models/patch_family_model.dart';
import './families_request_builder_get_query_parameters.dart';
import './item/with_family_item_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \families
class FamiliesRequestBuilder extends BaseRequestBuilder<FamiliesRequestBuilder> {
    /// Gets an item from the ApiSdk.families.item collection
    ///  [familyId] Family member ID
    WithFamilyItemRequestBuilder byFamilyId(String familyId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('familyId', () => familyId);
        return WithFamilyItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    FamiliesRequestBuilder clone() {
        return FamiliesRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [FamiliesRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    FamiliesRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families{?page*,size*}", pathParameters) ;
    /// Instantiates a new [FamiliesRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    FamiliesRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families{?page*,size*}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Get all family members
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PaginatedResponseModelEndpointsFamilyModel?> getAsync([void Function(RequestConfiguration<FamiliesRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PaginatedResponseModelEndpointsFamilyModel>(requestInfo, PaginatedResponseModelEndpointsFamilyModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Patch family with members
    ///  [body] Model for updating family details
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<FamilyModel?> patchAsync(PatchFamilyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPatchRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '401' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<FamilyModel>(requestInfo, FamilyModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<FamilyModel?> postAsync(CreateFamilyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<FamilyModel>(requestInfo, FamilyModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Get all family members
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<FamiliesRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<FamiliesRequestBuilderGetQueryParameters>(requestConfiguration, () => FamiliesRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Patch family with members
    ///  [body] Model for updating family details
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPatchRequestInformation(PatchFamilyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.patch, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
    /// Create a new family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreateFamilyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
