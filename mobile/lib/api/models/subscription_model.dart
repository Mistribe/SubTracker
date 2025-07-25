// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './subscription_payment_model.dart';

/// auto generated
class SubscriptionModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    DateTime? createdAt;
    ///  The etag property
    String? etag;
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
    ///  The payer_id_id property
    String? payerIdId;
    ///  The payments property
    Iterable<SubscriptionPaymentModel>? payments;
    ///  The updated_at property
    DateTime? updatedAt;
    /// Instantiates a new [SubscriptionModel] and sets the default values.
    SubscriptionModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static SubscriptionModel createFromDiscriminatorValue(ParseNode parseNode) {
        return SubscriptionModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['family_id'] = (node) => familyId = node.getStringValue();
        deserializerMap['family_members'] = (node) => familyMembers = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['labels'] = (node) => labels = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['payed_by_joint_account'] = (node) => payedByJointAccount = node.getBoolValue();
        deserializerMap['payer_id_id'] = (node) => payerIdId = node.getStringValue();
        deserializerMap['payments'] = (node) => payments = node.getCollectionOfObjectValues<SubscriptionPaymentModel>(SubscriptionPaymentModel.createFromDiscriminatorValue);
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('etag', etag);
        writer.writeStringValue('family_id', familyId);
        writer.writeCollectionOfPrimitiveValues<String?>('family_members', familyMembers);
        writer.writeStringValue('id', id);
        writer.writeCollectionOfPrimitiveValues<String?>('labels', labels);
        writer.writeStringValue('name', name);
        writer.writeBoolValue('payed_by_joint_account', value:payedByJointAccount);
        writer.writeStringValue('payer_id_id', payerIdId);
        writer.writeCollectionOfObjectValues<SubscriptionPaymentModel>('payments', payments);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
