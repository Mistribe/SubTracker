// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../dto/amount_model.dart';

/// auto generated
/// @Description Price details for this subscription
class SubscriptionPriceModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description Custom price for this subscription
    AmountModel? monthly;
    ///  @Description Custom price for this subscription
    AmountModel? yearly;
    /// Instantiates a new [SubscriptionPriceModel] and sets the default values.
    SubscriptionPriceModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionPriceModel createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionPriceModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['monthly'] = (node) => monthly = node.getObjectValue<AmountModel>(AmountModel.createFromDiscriminatorValue);
        deserializerMap['yearly'] = (node) => yearly = node.getObjectValue<AmountModel>(AmountModel.createFromDiscriminatorValue);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeObjectValue<AmountModel>('monthly', monthly);
        writer.writeObjectValue<AmountModel>('yearly', yearly);
        writer.writeAdditionalData(additionalData);
    }
}
