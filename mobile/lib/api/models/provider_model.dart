// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './owner_model.dart';
import './plan_model.dart';

/// auto generated
class ProviderModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    DateTime? createdAt;
    ///  The description property
    String? description;
    ///  The etag property
    String? etag;
    ///  The icon_url property
    String? iconUrl;
    ///  The id property
    String? id;
    ///  The labels property
    Iterable<String>? labels;
    ///  The name property
    String? name;
    ///  The owner property
    OwnerModel? owner;
    ///  The plans property
    Iterable<PlanModel>? plans;
    ///  The pricing_page_url property
    String? pricingPageUrl;
    ///  The updated_at property
    DateTime? updatedAt;
    ///  The url property
    String? url;
    /// Instantiates a new [ProviderModel] and sets the default values.
    ProviderModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static ProviderModel createFromDiscriminatorValue(ParseNode parseNode) {
        return ProviderModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['description'] = (node) => description = node.getStringValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['icon_url'] = (node) => iconUrl = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['labels'] = (node) => labels = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['owner'] = (node) => owner = node.getObjectValue<OwnerModel>(OwnerModel.createFromDiscriminatorValue);
        deserializerMap['plans'] = (node) => plans = node.getCollectionOfObjectValues<PlanModel>(PlanModel.createFromDiscriminatorValue);
        deserializerMap['pricing_page_url'] = (node) => pricingPageUrl = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        deserializerMap['url'] = (node) => url = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('description', description);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('icon_url', iconUrl);
        writer.writeStringValue('id', id);
        writer.writeCollectionOfPrimitiveValues<String?>('labels', labels);
        writer.writeStringValue('name', name);
        writer.writeObjectValue<OwnerModel>('owner', owner);
        writer.writeCollectionOfObjectValues<PlanModel>('plans', plans);
        writer.writeStringValue('pricing_page_url', pricingPageUrl);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeStringValue('url', url);
        writer.writeAdditionalData(additionalData);
    }
}
