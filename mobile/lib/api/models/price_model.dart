// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Price object defining the cost of a subscription plan with currency and time validity
class PriceModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description Price amount in the specified currency (supports decimal values)
    double? amount;
    ///  @Description ISO 8601 timestamp when the price was originally created
    DateTime? createdAt;
    ///  @Description ISO 4217 currency code for the price
    String? currency;
    ///  @Description ISO 8601 timestamp when this price expires (null means indefinite)
    DateTime? endDate;
    ///  @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
    String? etag;
    ///  @Description Unique identifier for the price (UUID format)
    String? id;
    ///  @Description ISO 8601 timestamp when this price becomes effective
    DateTime? startDate;
    ///  @Description ISO 8601 timestamp when the price was last modified
    DateTime? updatedAt;
    /// Instantiates a new [PriceModel] and sets the default values.
    PriceModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PriceModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PriceModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['amount'] = (node) => amount = node.getDoubleValue();
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['currency'] = (node) => currency = node.getStringValue();
        deserializerMap['end_date'] = (node) => endDate = node.getDateTimeValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['start_date'] = (node) => startDate = node.getDateTimeValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDoubleValue('amount', amount);
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('currency', currency);
        writer.writeDateTimeValue('end_date', endDate);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('id', id);
        writer.writeDateTimeValue('start_date', startDate);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
