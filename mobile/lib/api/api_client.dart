// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import 'package:microsoft_kiota_serialization_form/microsoft_kiota_serialization_form.dart';
import 'package:microsoft_kiota_serialization_json/microsoft_kiota_serialization_json.dart';
import 'package:microsoft_kiota_serialization_multipart/microsoft_kiota_serialization_multipart.dart';
import 'package:microsoft_kiota_serialization_text/microsoft_kiota_serialization_text.dart';
import './api/api_request_builder.dart';
import './families/families_request_builder.dart';
import './labels/labels_request_builder.dart';
import './providers/providers_request_builder.dart';
import './subscriptions/subscriptions_request_builder.dart';

/// auto generated
/// The main entry point of the SDK, exposes the configuration and the fluent API.
class ApiClient extends BaseRequestBuilder<ApiClient> {
    ///  The api property
    ApiRequestBuilder get api {
        return ApiRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The families property
    FamiliesRequestBuilder get families {
        return FamiliesRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The labels property
    LabelsRequestBuilder get labels {
        return LabelsRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The providers property
    ProvidersRequestBuilder get providers {
        return ProvidersRequestBuilder(pathParameters, requestAdapter);
    }
    ///  The subscriptions property
    SubscriptionsRequestBuilder get subscriptions {
        return SubscriptionsRequestBuilder(pathParameters, requestAdapter);
    }
    /// Clones the requestbuilder.
    @override
    ApiClient clone() {
        return ApiClient(requestAdapter);
    }
    /// Instantiates a new [ApiClient] and sets the default values.
    ///  [requestAdapter] The request adapter to use to execute the requests.
    ApiClient(RequestAdapter requestAdapter) : super(requestAdapter, "{+baseurl}", {}) {
        ApiClientBuilder.registerDefaultSerializer(JsonSerializationWriterFactory.new);
        ApiClientBuilder.registerDefaultSerializer(TextSerializationWriterFactory.new);
        ApiClientBuilder.registerDefaultSerializer(FormSerializationWriterFactory.new);
        ApiClientBuilder.registerDefaultSerializer(MultipartSerializationWriterFactory.new);
        ApiClientBuilder.registerDefaultDeserializer(JsonParseNodeFactory.new);
        ApiClientBuilder.registerDefaultDeserializer(FormParseNodeFactory.new);
        ApiClientBuilder.registerDefaultDeserializer(TextParseNodeFactory.new);
        if (requestAdapter.baseUrl == null || requestAdapter.baseUrl!.isEmpty) {
            requestAdapter.baseUrl = 'https://api.subtracker.mistribe.com';
        }
        pathParameters['baseurl'] = requestAdapter.baseUrl;
    }
}
