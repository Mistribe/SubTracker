// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/family_invite_request.dart';
import '../../../models/family_invite_response.dart';
import '../../../models/http_error_response.dart';

/// auto generated
/// Builds and executes requests for operations under \families\{familyId}\invite
class InviteRequestBuilder extends BaseRequestBuilder<InviteRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    InviteRequestBuilder clone() {
        return InviteRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [InviteRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    InviteRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/invite", pathParameters) ;
    /// Instantiates a new [InviteRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    InviteRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/invite", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Creates an invitation for a new member to join the family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<FamilyInviteResponse?> postAsync(FamilyInviteRequest body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<FamilyInviteResponse>(requestInfo, FamilyInviteResponse.createFromDiscriminatorValue, errorMapping);
    }
    /// Creates an invitation for a new member to join the family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(FamilyInviteRequest body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
