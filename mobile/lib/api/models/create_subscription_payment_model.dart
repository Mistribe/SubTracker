// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class CreateSubscriptionPaymentModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    DateTime? createdAt;
    ///  The currency property
    String? currency;
    ///  The end_date property
    DateTime? endDate;
    ///  The id property
    String? id;
    ///  The months property
    int? months;
    ///  The price property
    double? price;
    ///  The start_date property
    DateTime? startDate;
    /// Instantiates a new [CreateSubscriptionPaymentModel] and sets the default values.
    CreateSubscriptionPaymentModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static CreateSubscriptionPaymentModel createFromDiscriminatorValue(ParseNode parseNode) {
        return CreateSubscriptionPaymentModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['currency'] = (node) => currency = node.getStringValue();
        deserializerMap['end_date'] = (node) => endDate = node.getDateTimeValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['months'] = (node) => months = node.getIntValue();
        deserializerMap['price'] = (node) => price = node.getDoubleValue();
        deserializerMap['start_date'] = (node) => startDate = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('currency', currency);
        writer.writeDateTimeValue('end_date', endDate);
        writer.writeStringValue('id', id);
        writer.writeIntValue('months', months);
        writer.writeDoubleValue('price', price);
        writer.writeDateTimeValue('start_date', startDate);
        writer.writeAdditionalData(additionalData);
    }
}
