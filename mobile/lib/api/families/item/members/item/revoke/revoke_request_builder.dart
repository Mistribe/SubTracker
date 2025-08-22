// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../../../../models/http_error_response.dart';
import './revoke_post_request_body.dart';

/// auto generated
/// Builds and executes requests for operations under \families\{familyId}\members\{familyMember-id}\revoke
class RevokeRequestBuilder extends BaseRequestBuilder<RevokeRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    RevokeRequestBuilder clone() {
        return RevokeRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [RevokeRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    RevokeRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/members/{familyMember%2Did}/revoke", pathParameters) ;
    /// Instantiates a new [RevokeRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    RevokeRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/families/{familyId}/members/{familyMember%2Did}/revoke", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Revokes a member from the family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<void> postAsync(RevokePostRequestBody body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '404' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.sendNoContent(requestInfo, errorMapping);
    }
    /// Revokes a member from the family
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(RevokePostRequestBody body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
