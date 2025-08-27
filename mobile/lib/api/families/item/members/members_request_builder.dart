// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/family/create_family_member_model.dart';
import '../../../models/family/family_model.dart';
import '../../../models/ginx/http_error_response.dart';
import './item/family_member_item_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \families\{familyId}\members
class MembersRequestBuilder extends BaseRequestBuilder<MembersRequestBuilder> {
    /// Gets an item from the ApiSdk.families.item.members.item collection
    ///  [familyMemberId] Family member ID (UUID format)
    FamilyMemberItemRequestBuilder byFamilyMemberId(String familyMemberId) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('familyMember%2Did', () => familyMemberId);
        return FamilyMemberItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    MembersRequestBuilder clone() {
        return MembersRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [MembersRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    MembersRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/members", pathParameters) ;
    /// Instantiates a new [MembersRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    MembersRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/members", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Add a new member to an existing family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<FamilyModel?> postAsync(CreateFamilyMemberModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<FamilyModel>(requestInfo, FamilyModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Add a new member to an existing family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreateFamilyMemberModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
