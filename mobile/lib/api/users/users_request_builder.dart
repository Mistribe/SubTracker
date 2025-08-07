// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './profile/profile_request_builder.dart';

/// auto generated
/// Builds and executes requests for operations under \users
class UsersRequestBuilder extends BaseRequestBuilder<UsersRequestBuilder> {
    ///  The profile property
    ProfileRequestBuilder get profile {
        return ProfileRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    UsersRequestBuilder clone() {
        return UsersRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [UsersRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    UsersRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users", pathParameters) ;
    /// Instantiates a new [UsersRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    UsersRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users", {RequestInformation.rawUrlKey : rawUrl}) ;
}
