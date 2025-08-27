// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../subscription/subscription_model.dart';

/// auto generated
class PaginatedResponseModelSubscriptionSubscriptionModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Data contains the list of items for the current page
    Iterable<SubscriptionModel>? data;
    ///  Length represents the number of items in the current page
    int? length;
    ///  Total represents the total number of items available
    int? total;
    /// Instantiates a new [PaginatedResponseModelSubscriptionSubscriptionModel] and sets the default values.
    PaginatedResponseModelSubscriptionSubscriptionModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PaginatedResponseModelSubscriptionSubscriptionModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PaginatedResponseModelSubscriptionSubscriptionModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['data'] = (node) => data = node.getCollectionOfObjectValues<SubscriptionModel>(SubscriptionModel.createFromDiscriminatorValue);
        deserializerMap['length'] = (node) => length = node.getIntValue();
        deserializerMap['total'] = (node) => total = node.getIntValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeCollectionOfObjectValues<SubscriptionModel>('data', data);
        writer.writeIntValue('length', length);
        writer.writeIntValue('total', total);
        writer.writeAdditionalData(additionalData);
    }
}
