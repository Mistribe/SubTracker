// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
class AmountModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The currency property
    String? currency;
    ///  The source property
    AmountModel? source;
    ///  The value property
    double? value;
    /// Instantiates a new [AmountModel] and sets the default values.
    AmountModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static AmountModel createFromDiscriminatorValue(ParseNode parseNode) {
        return AmountModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['currency'] = (node) => currency = node.getStringValue();
        deserializerMap['source'] = (node) => source = node.getObjectValue<AmountModel>(AmountModel.createFromDiscriminatorValue);
        deserializerMap['value'] = (node) => value = node.getDoubleValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('currency', currency);
        writer.writeObjectValue<AmountModel>('source', source);
        writer.writeDoubleValue('value', value);
        writer.writeAdditionalData(additionalData);
    }
}
