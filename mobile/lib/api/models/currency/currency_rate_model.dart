// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class CurrencyRateModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The currency property
    String? currency;
    ///  The rate property
    double? rate;
    /// Instantiates a new [CurrencyRateModel] and sets the default values.
    CurrencyRateModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static CurrencyRateModel createFromDiscriminatorValue(ParseNode parseNode) {
        return CurrencyRateModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['currency'] = (node) => currency = node.getStringValue();
        deserializerMap['rate'] = (node) => rate = node.getDoubleValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('currency', currency);
        writer.writeDoubleValue('rate', rate);
        writer.writeAdditionalData(additionalData);
    }
}
