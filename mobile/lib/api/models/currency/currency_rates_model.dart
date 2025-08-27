// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './currency_rate_model.dart';

/// auto generated
class CurrencyRatesModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The rates property
    Iterable<CurrencyRateModel>? rates;
    ///  The timestamp property
    DateTime? timestamp;
    /// Instantiates a new [CurrencyRatesModel] and sets the default values.
    CurrencyRatesModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static CurrencyRatesModel createFromDiscriminatorValue(ParseNode parseNode) {
        return CurrencyRatesModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['rates'] = (node) => rates = node.getCollectionOfObjectValues<CurrencyRateModel>(CurrencyRateModel.createFromDiscriminatorValue);
        deserializerMap['timestamp'] = (node) => timestamp = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeCollectionOfObjectValues<CurrencyRateModel>('rates', rates);
        writer.writeDateTimeValue('timestamp', timestamp);
        writer.writeAdditionalData(additionalData);
    }
}
