// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './owner_model.dart';
import './plan_model.dart';

/// auto generated
/// Provider object containing information about a subscription service provider and their available plans
class ProviderModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description ISO 8601 timestamp when the provider was originally created
    DateTime? createdAt;
    ///  @Description Optional detailed description of the provider and their services
    String? description;
    ///  @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
    String? etag;
    ///  @Description Optional URL to the provider's icon or logo image
    String? iconUrl;
    ///  @Description Unique identifier for the provider (UUID format)
    String? id;
    ///  The key property
    String? key;
    ///  @Description List of label IDs associated with this provider for categorization
    Iterable<String>? labels;
    ///  @Description Display name of the service provider
    String? name;
    ///  @Description Ownership information specifying whether this subscription belongs to a user or family
    OwnerModel? owner;
    ///  @Description List of subscription plans offered by this provider
    Iterable<PlanModel>? plans;
    ///  @Description Optional URL to the provider's pricing information page
    String? pricingPageUrl;
    ///  @Description ISO 8601 timestamp when the provider was last modified
    DateTime? updatedAt;
    ///  @Description Optional URL to the provider's main website
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
        deserializerMap['key'] = (node) => key = node.getStringValue();
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
        writer.writeStringValue('key', key);
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
