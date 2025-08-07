// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../../models/http_error.dart';
import '../../models/update_profile_model.dart';
import '../../models/user_profile_model.dart';

/// auto generated
/// Builds and executes requests for operations under \users\profile
class ProfileRequestBuilder extends BaseRequestBuilder<ProfileRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    ProfileRequestBuilder clone() {
        return ProfileRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [ProfileRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ProfileRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users/profile", pathParameters) ;
    /// Instantiates a new [ProfileRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ProfileRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/users/profile", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Returns the profile information for the authenticated user
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<UserProfileModel?> getAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '401' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<UserProfileModel>(requestInfo, UserProfileModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Updates the currency preference in user's profile
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<UserProfileModel?> putAsync(UpdateProfileModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPutRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpError.createFromDiscriminatorValue,
            '401' :  HttpError.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<UserProfileModel>(requestInfo, UserProfileModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Returns the profile information for the authenticated user
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Updates the currency preference in user's profile
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPutRequestInformation(UpdateProfileModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.put, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
