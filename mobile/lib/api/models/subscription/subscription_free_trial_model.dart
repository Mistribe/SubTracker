// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// @Description Number of free trial days remaining (null if no trial or trial expired)
class SubscriptionFreeTrialModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The end_date property
    DateTime? endDate;
    ///  The start_date property
    DateTime? startDate;
    /// Instantiates a new [SubscriptionFreeTrialModel] and sets the default values.
    SubscriptionFreeTrialModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionFreeTrialModel createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionFreeTrialModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['end_date'] = (node) => endDate = node.getDateTimeValue();
        deserializerMap['start_date'] = (node) => startDate = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('end_date', endDate);
        writer.writeDateTimeValue('start_date', startDate);
        writer.writeAdditionalData(additionalData);
    }
}
