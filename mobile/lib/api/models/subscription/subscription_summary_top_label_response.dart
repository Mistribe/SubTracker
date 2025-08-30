// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../dto/amount_model.dart';

/// auto generated
class SubscriptionSummaryTopLabelResponse implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The label_id property
    String? labelId;
    ///  @Description Custom price for this subscription
    AmountModel? total;
    /// Instantiates a new [SubscriptionSummaryTopLabelResponse] and sets the default values.
    SubscriptionSummaryTopLabelResponse() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionSummaryTopLabelResponse createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionSummaryTopLabelResponse();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['label_id'] = (node) => labelId = node.getStringValue();
        deserializerMap['total'] = (node) => total = node.getObjectValue<AmountModel>(AmountModel.createFromDiscriminatorValue);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('label_id', labelId);
        writer.writeObjectValue<AmountModel>('total', total);
        writer.writeAdditionalData(additionalData);
    }
}
