// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class PaymentModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    String? createdAt;
    ///  The currency property
    String? currency;
    ///  The end_date property
    String? endDate;
    ///  The etag property
    String? etag;
    ///  The id property
    String? id;
    ///  The months property
    int? months;
    ///  The price property
    double? price;
    ///  The start_date property
    String? startDate;
    ///  The updated_at property
    String? updatedAt;
    /// Instantiates a new [PaymentModel] and sets the default values.
    PaymentModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PaymentModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PaymentModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getStringValue();
        deserializerMap['currency'] = (node) => currency = node.getStringValue();
        deserializerMap['end_date'] = (node) => endDate = node.getStringValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['months'] = (node) => months = node.getIntValue();
        deserializerMap['price'] = (node) => price = node.getDoubleValue();
        deserializerMap['start_date'] = (node) => startDate = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getStringValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('created_at', createdAt);
        writer.writeStringValue('currency', currency);
        writer.writeStringValue('end_date', endDate);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('id', id);
        writer.writeIntValue('months', months);
        writer.writeDoubleValue('price', price);
        writer.writeStringValue('start_date', startDate);
        writer.writeStringValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
