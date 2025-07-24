// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../models/label_model.dart';

/// auto generated
/// Builds and executes requests for operations under \default
class DefaultRequestBuilder extends BaseRequestBuilder<DefaultRequestBuilder> {
    /// Clones the requestbuilder.
    @override
    DefaultRequestBuilder clone() {
        return DefaultRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [DefaultRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    DefaultRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/default", pathParameters) ;
    /// Instantiates a new [DefaultRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    DefaultRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/default", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Retrieves a list of default labels
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<Iterable<LabelModel>?> getAsync([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        return await requestAdapter.sendCollection<LabelModel>(requestInfo, LabelModel.createFromDiscriminatorValue, {});
    }
    /// Retrieves a list of default labels
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
}
