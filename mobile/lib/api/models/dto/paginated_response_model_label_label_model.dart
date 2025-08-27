// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../label/label_model.dart';

/// auto generated
class PaginatedResponseModelLabelLabelModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Data contains the list of items for the current page
    Iterable<LabelModel>? data;
    ///  Length represents the number of items in the current page
    int? length;
    ///  Total represents the total number of items available
    int? total;
    /// Instantiates a new [PaginatedResponseModelLabelLabelModel] and sets the default values.
    PaginatedResponseModelLabelLabelModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PaginatedResponseModelLabelLabelModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PaginatedResponseModelLabelLabelModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['data'] = (node) => data = node.getCollectionOfObjectValues<LabelModel>(LabelModel.createFromDiscriminatorValue);
        deserializerMap['length'] = (node) => length = node.getIntValue();
        deserializerMap['total'] = (node) => total = node.getIntValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeCollectionOfObjectValues<LabelModel>('data', data);
        writer.writeIntValue('length', length);
        writer.writeIntValue('total', total);
        writer.writeAdditionalData(additionalData);
    }
}
