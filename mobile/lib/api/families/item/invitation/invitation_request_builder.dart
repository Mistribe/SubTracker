// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../models/family/family_see_invitation_response.dart';
import '../../../models/ginx/http_error_response.dart';
import './invitation_request_builder_get_query_parameters.dart';

/// auto generated
/// Builds and executes requests for operations under \families\{familyId}\invitation
class InvitationRequestBuilder extends BaseRequestBuilder<InvitationRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    InvitationRequestBuilder clone() {
        return InvitationRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [InvitationRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    InvitationRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/invitation?code={code}&family_member_id={family_member_id}", pathParameters) ;
    /// Instantiates a new [InvitationRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    InvitationRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/invitation?code={code}&family_member_id={family_member_id}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Get information about a family invitation using invitation code
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<FamilySeeInvitationResponse?> getAsync([void Function(RequestConfiguration<InvitationRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<FamilySeeInvitationResponse>(requestInfo, FamilySeeInvitationResponse.createFromDiscriminatorValue, errorMapping);
    }
    /// Get information about a family invitation using invitation code
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<InvitationRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<InvitationRequestBuilderGetQueryParameters>(requestConfiguration, () => InvitationRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
