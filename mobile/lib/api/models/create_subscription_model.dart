// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './create_payment_model.dart';

/// auto generated
class CreateSubscriptionModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    String? createdAt;
    ///  The family_id property
    String? familyId;
    ///  The family_members property
    Iterable<String>? familyMembers;
    ///  The id property
    String? id;
    ///  The labels property
    Iterable<String>? labels;
    ///  The name property
    String? name;
    ///  The payed_by_joint_account property
    bool? payedByJointAccount;
    ///  The payer_id property
    String? payerId;
    ///  The payments property
    Iterable<CreatePaymentModel>? payments;
    /// Instantiates a new [CreateSubscriptionModel] and sets the default values.
    CreateSubscriptionModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static CreateSubscriptionModel createFromDiscriminatorValue(ParseNode parseNode) {
        return CreateSubscriptionModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['family_members'] = (node) => familyMembers = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['labels'] = (node) => labels = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['payed_by_joint_account'] = (node) => payedByJointAccount = node.getBoolValue();
        deserializerMap['payer_id'] = (node) => payerId = node.getStringValue();
        deserializerMap['payments'] = (node) => payments = node.getCollectionOfObjectValues<CreatePaymentModel>(CreatePaymentModel.createFromDiscriminatorValue);
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeStringValue('created_at', createdAt);
        writer.writeStringValue('family_id', familyId);
        writer.writeCollectionOfPrimitiveValues<String?>('family_members', familyMembers);
        writer.writeStringValue('id', id);
        writer.writeCollectionOfPrimitiveValues<String?>('labels', labels);
        writer.writeStringValue('name', name);
        writer.writeBoolValue('payed_by_joint_account', value:payedByJointAccount);
        writer.writeStringValue('payer_id', payerId);
        writer.writeCollectionOfObjectValues<CreatePaymentModel>('payments', payments);
        writer.writeAdditionalData(additionalData);
    }
}
