// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../models/dto/paginated_response_model_label_label_model.dart';
import '../models/ginx/http_error_response.dart';
import '../models/label/create_label_model.dart';
import '../models/label/label_model.dart';
import './default_/default_request_builder.dart';
import './item/labels_item_request_builder.dart';
import './labels_request_builder_get_query_parameters.dart';

/// auto generated
/// Builds and executes requests for operations under \labels
class LabelsRequestBuilder extends BaseRequestBuilder<LabelsRequestBuilder> {
    ///  The default property
    DefaultRequestBuilder get default_ {
        return DefaultRequestBuilder(pathParameters, requestAdapter);
    }
    /// Gets an item from the ApiSdk.labels.item collection
    ///  [id] Label ID (UUID format)
    LabelsItemRequestBuilder byId(String id) {
        var urlTplParams = Map.of(pathParameters);
        urlTplParams.putIfAbsent('id', () => id);
        return LabelsItemRequestBuilder(urlTplParams, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    LabelsRequestBuilder clone() {
        return LabelsRequestBuilder(pathParameters, requestAdapter);
    }
    /// Instantiates a new [LabelsRequestBuilder] and sets the default values.
    ///  [pathParameters] Path parameters for the request
    ///  [requestAdapter] The request adapter to use to execute the requests.
    LabelsRequestBuilder(Map<String, dynamic> pathParameters, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/labels{?limit*,offset*,search*}", pathParameters) ;
    /// Instantiates a new [LabelsRequestBuilder] and sets the default values.
    ///  [rawUrl] The raw URL to use for the request builder.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    LabelsRequestBuilder.withUrl(String rawUrl, RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}/labels{?limit*,offset*,search*}", {RequestInformation.rawUrlKey : rawUrl}) ;
    /// Retrieve a paginated list of labels with optional filtering by owner type and search text
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<PaginatedResponseModelLabelLabelModel?> getAsync([void Function(RequestConfiguration<LabelsRequestBuilderGetQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toGetRequestInformation(requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<PaginatedResponseModelLabelLabelModel>(requestInfo, PaginatedResponseModelLabelLabelModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Create a new label with specified name, color, and owner information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    Future<LabelModel?> postAsync(CreateLabelModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) async {
        var requestInfo = toPostRequestInformation(body, requestConfiguration);
        final errorMapping = <String, ParsableFactory<Parsable>>{
            '400' :  HttpErrorResponse.createFromDiscriminatorValue,
            '401' :  HttpErrorResponse.createFromDiscriminatorValue,
            '500' :  HttpErrorResponse.createFromDiscriminatorValue,
        };
        return await requestAdapter.send<LabelModel>(requestInfo, LabelModel.createFromDiscriminatorValue, errorMapping);
    }
    /// Retrieve a paginated list of labels with optional filtering by owner type and search text
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toGetRequestInformation([void Function(RequestConfiguration<LabelsRequestBuilderGetQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.get, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<LabelsRequestBuilderGetQueryParameters>(requestConfiguration, () => LabelsRequestBuilderGetQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        return requestInfo;
    }
    /// Create a new label with specified name, color, and owner information
    ///  [body] The request body
    ///  [requestConfiguration] Configuration for the request such as headers, query parameters, and middleware options.
    RequestInformation toPostRequestInformation(CreateLabelModel body, [void Function(RequestConfiguration<DefaultQueryParameters>)? requestConfiguration]) {
        var requestInfo = RequestInformation(httpMethod : HttpMethod.post, urlTemplate : urlTemplate, pathParameters :  pathParameters);
        requestInfo.configure<DefaultQueryParameters>(requestConfiguration, () => DefaultQueryParameters());
        requestInfo.headers.put('Accept', 'application/json');
        requestInfo.setContentFromParsable(requestAdapter, 'application/json', body);
        return requestInfo;
    }
}
