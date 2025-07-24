// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './family_model.dart';

/// auto generated
class PaginatedResponseModelEndpointsFamilyModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The data property
    Iterable<FamilyModel>? data;
    ///  The length property
    int? length;
    ///  The total property
    int? total;
    /// Instantiates a new [PaginatedResponseModelEndpointsFamilyModel] and sets the default values.
    PaginatedResponseModelEndpointsFamilyModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PaginatedResponseModelEndpointsFamilyModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PaginatedResponseModelEndpointsFamilyModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['data'] = (node) => data = node.getCollectionOfObjectValues<FamilyModel>(FamilyModel.createFromDiscriminatorValue);
        deserializerMap['length'] = (node) => length = node.getIntValue();
        deserializerMap['total'] = (node) => total = node.getIntValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeCollectionOfObjectValues<FamilyModel>('data', data);
        writer.writeIntValue('length', length);
        writer.writeIntValue('total', total);
        writer.writeAdditionalData(additionalData);
    }
}
