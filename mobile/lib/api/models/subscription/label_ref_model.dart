// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './label_ref_model_source.dart';

/// auto generated
class LabelRefModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The label_id property
    String? labelId;
    ///  The source property
    LabelRefModelSource? source;
    /// Instantiates a new [LabelRefModel] and sets the default values.
    LabelRefModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static LabelRefModel createFromDiscriminatorValue(ParseNode parseNode) {
        return LabelRefModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['label_id'] = (node) => labelId = node.getStringValue();
        deserializerMap['source'] = (node) => source = node.getEnumValue<LabelRefModelSource>((stringValue) => LabelRefModelSource.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('label_id', labelId);
        writer.writeEnumValue<LabelRefModelSource>('source', source, (e) => e?.value);
        writer.writeAdditionalData(additionalData);
    }
}
