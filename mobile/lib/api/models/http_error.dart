// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class HttpError extends ApiException implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    /// Instantiates a new [HttpError] and sets the default values.
    HttpError( {
        super.message,
        super.statusCode,
        super.responseHeaders,
        super.innerExceptions,
        required this.additionalData,
    });
    /// Creates a copy of the object.
    @override
    HttpError copyWith({int? statusCode, String? message, Map<String, List<String>>? responseHeaders, Iterable<Object?>? innerExceptions, Map<String, Object?>? additionalData }){
        return HttpError(
        message : message ?? this.message, 
        statusCode : statusCode ?? this.statusCode, 
        responseHeaders : responseHeaders ?? this.responseHeaders, 
        innerExceptions : innerExceptions ?? this.innerExceptions, 
        additionalData : additionalData ?? this.additionalData, 
        );
    }
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static HttpError createFromDiscriminatorValue(ParseNode parseNode) {
        return HttpError(additionalData: {});
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeAdditionalData(additionalData);
    }
}
