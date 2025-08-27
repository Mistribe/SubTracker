// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class UpdateProviderModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The description property
    String? description;
    ///  The icon_url property
    String? iconUrl;
    ///  The labels property
    Iterable<String>? labels;
    ///  The name property
    String? name;
    ///  The pricing_page_url property
    String? pricingPageUrl;
    ///  The updated_at property
    DateTime? updatedAt;
    ///  The url property
    String? url;
    /// Instantiates a new [UpdateProviderModel] and sets the default values.
    UpdateProviderModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static UpdateProviderModel createFromDiscriminatorValue(ParseNode parseNode) {
        return UpdateProviderModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['description'] = (node) => description = node.getStringValue();
        deserializerMap['icon_url'] = (node) => iconUrl = node.getStringValue();
        deserializerMap['labels'] = (node) => labels = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['pricing_page_url'] = (node) => pricingPageUrl = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        deserializerMap['url'] = (node) => url = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('description', description);
        writer.writeStringValue('icon_url', iconUrl);
        writer.writeCollectionOfPrimitiveValues<String?>('labels', labels);
        writer.writeStringValue('name', name);
        writer.writeStringValue('pricing_page_url', pricingPageUrl);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeStringValue('url', url);
        writer.writeAdditionalData(additionalData);
    }
}
