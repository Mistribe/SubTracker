// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './price_model.dart';

/// auto generated
/// Plan object defining a specific subscription tier with associated pricing options
class PlanModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description ISO 8601 timestamp when the plan was originally created
    DateTime? createdAt;
    ///  @Description Optional detailed description of the plan features and benefits
    String? description;
    ///  @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
    String? etag;
    ///  @Description Unique identifier for the plan (UUID format)
    String? id;
    ///  @Description Display name of the subscription plan
    String? name;
    ///  @Description List of pricing options available for this plan (different currencies, time periods, etc.)
    Iterable<PriceModel>? prices;
    ///  @Description ISO 8601 timestamp when the plan was last modified
    DateTime? updatedAt;
    /// Instantiates a new [PlanModel] and sets the default values.
    PlanModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PlanModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PlanModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['description'] = (node) => description = node.getStringValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['prices'] = (node) => prices = node.getCollectionOfObjectValues<PriceModel>(PriceModel.createFromDiscriminatorValue);
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('description', description);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('id', id);
        writer.writeStringValue('name', name);
        writer.writeCollectionOfObjectValues<PriceModel>('prices', prices);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
