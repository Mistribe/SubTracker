// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../dto/amount_model.dart';
import '../dto/editable_owner_model.dart';
import './editable_subscription_payer_model.dart';
import './subscription_free_trial_model.dart';

/// auto generated
class PatchSubscriptionModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The custom_price property
    AmountModel? customPrice;
    ///  The custom_recurrency property
    int? customRecurrency;
    ///  The end_date property
    DateTime? endDate;
    ///  @Description Number of free trial days remaining (null if no trial or trial expired)
    SubscriptionFreeTrialModel? freeTrial;
    ///  The friendly_name property
    String? friendlyName;
    ///  The id property
    String? id;
    ///  The labels property
    Iterable<String>? labels;
    ///  The owner property
    EditableOwnerModel? owner;
    ///  Subscription payer object used for updating who pays for a subscription
    EditableSubscriptionPayerModel? payer;
    ///  The plan_id property
    String? planId;
    ///  The price_id property
    String? priceId;
    ///  The provider_id property
    String? providerId;
    ///  The recurrency property
    String? recurrency;
    ///  The service_users property
    Iterable<String>? serviceUsers;
    ///  The start_date property
    DateTime? startDate;
    ///  The updated_at property
    DateTime? updatedAt;
    /// Instantiates a new [PatchSubscriptionModel] and sets the default values.
    PatchSubscriptionModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PatchSubscriptionModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PatchSubscriptionModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['custom_price'] = (node) => customPrice = node.getObjectValue<AmountModel>(AmountModel.createFromDiscriminatorValue);
        deserializerMap['custom_recurrency'] = (node) => customRecurrency = node.getIntValue();
        deserializerMap['end_date'] = (node) => endDate = node.getDateTimeValue();
        deserializerMap['free_trial'] = (node) => freeTrial = node.getObjectValue<SubscriptionFreeTrialModel>(SubscriptionFreeTrialModel.createFromDiscriminatorValue);
        deserializerMap['friendly_name'] = (node) => friendlyName = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['labels'] = (node) => labels = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['owner'] = (node) => owner = node.getObjectValue<EditableOwnerModel>(EditableOwnerModel.createFromDiscriminatorValue);
        deserializerMap['payer'] = (node) => payer = node.getObjectValue<EditableSubscriptionPayerModel>(EditableSubscriptionPayerModel.createFromDiscriminatorValue);
        deserializerMap['plan_id'] = (node) => planId = node.getStringValue();
        deserializerMap['price_id'] = (node) => priceId = node.getStringValue();
        deserializerMap['provider_id'] = (node) => providerId = node.getStringValue();
        deserializerMap['recurrency'] = (node) => recurrency = node.getStringValue();
        deserializerMap['service_users'] = (node) => serviceUsers = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['start_date'] = (node) => startDate = node.getDateTimeValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeObjectValue<AmountModel>('custom_price', customPrice);
        writer.writeIntValue('custom_recurrency', customRecurrency);
        writer.writeDateTimeValue('end_date', endDate);
        writer.writeObjectValue<SubscriptionFreeTrialModel>('free_trial', freeTrial);
        writer.writeStringValue('friendly_name', friendlyName);
        writer.writeStringValue('id', id);
        writer.writeCollectionOfPrimitiveValues<String?>('labels', labels);
        writer.writeObjectValue<EditableOwnerModel>('owner', owner);
        writer.writeObjectValue<EditableSubscriptionPayerModel>('payer', payer);
        writer.writeStringValue('plan_id', planId);
        writer.writeStringValue('price_id', priceId);
        writer.writeStringValue('provider_id', providerId);
        writer.writeStringValue('recurrency', recurrency);
        writer.writeCollectionOfPrimitiveValues<String?>('service_users', serviceUsers);
        writer.writeDateTimeValue('start_date', startDate);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
