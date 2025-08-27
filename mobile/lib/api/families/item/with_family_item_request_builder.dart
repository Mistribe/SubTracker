// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../models/family/family_model.dart';
import '../../models/family/update_family_model.dart';
import '../../models/ginx/http_error_response.dart';
import './accept/accept_request_builder.dart';
import './decline/decline_request_builder.dart';
import './invitation/invitation_request_builder.dart';
import './invite/invite_request_builder.dart';
import './members/members_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \families\{familyId}
class WithFamilyItemRequestBuilder extends BaseRequestBuilder<WithFamilyItemRequestBuilder> {
    ///  The accept property
    AcceptRequestBuilder get accept {
        return AcceptRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The decline property
    DeclineRequestBuilder get decline {
        return DeclineRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The invitation property
    InvitationRequestBuilder get invitation {
        return InvitationRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The invite property
    InviteRequestBuilder get invite {
        return InviteRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The members property
    MembersRequestBuilder get members {
        return MembersRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    WithFamilyItemRequestBuilder clone() {
        return WithFamilyItemRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [WithFamilyItemRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithFamilyItemRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}", pathParameters) ;
    /// Instantiates a new [WithFamilyItemRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    WithFamilyItemRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Permanently delete a family and all its members
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
    /// Update family information such as name and other details
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<FamilyModel?> putAsync(UpdateFamilyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<FamilyModel>(requestInfo, FamilyModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Permanently delete a family and all its members
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toDeleteRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.delete, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Update family information such as name and other details
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPutRequestInformation(UpdateFamilyModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.put, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
