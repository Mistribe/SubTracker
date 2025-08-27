// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import '../dto/amount_model.dart';
import '../dto/owner_model.dart';
import './label_ref_model.dart';
import './subscription_free_trial_model.dart';
import './subscription_model_recurrency.dart';
import './subscription_payer_model.dart';

/// auto generated
/// Subscription object containing all information about an active subscription including billing and usage details
class SubscriptionModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description ISO 8601 timestamp when the subscription was originally created
    DateTime? createdAt;
    ///  The custom_price property
    AmountModel? customPrice;
    ///  @Description CustomRecurrency recurrency interval in days (required when recurrency is custom)
    int? customRecurrency;
    ///  @Description ISO 8601 timestamp when the subscription expires (null for ongoing subscriptions)
    DateTime? endDate;
    ///  @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
    String? etag;
    ///  @Description Number of free trial days remaining (null if no trial or trial expired)
    SubscriptionFreeTrialModel? freeTrial;
    ///  @Description Optional custom name for easy identification of the subscription
    String? friendlyName;
    ///  @Description Unique identifier for the subscription (UUID format)
    String? id;
    ///  @Description Indicates whether the subscription is currently active or not
    bool? isActive;
    ///  @Description List of labels associated with this subscription
    Iterable<LabelRefModel>? labelRefs;
    ///  @Description Ownership information specifying whether this subscription belongs to a user or family
    OwnerModel? owner;
    ///  @Description Information about who pays for this subscription within the family
    SubscriptionPayerModel? payer;
    ///  @Description ID of the specific plan being subscribed to
    String? planId;
    ///  @Description ID of the pricing tier for this subscription
    String? priceId;
    ///  @Description ID of the service provider offering this subscription
    String? providerId;
    ///  @Description Billing recurrency pattern (monthly, yearly, custom, etc.)
    SubscriptionModelRecurrency? recurrency;
    ///  @Description List of family member IDs who use this service (for shared subscriptions)
    Iterable<String>? serviceUsers;
    ///  @Description ISO 8601 timestamp when the subscription becomes active
    DateTime? startDate;
    ///  @Description ISO 8601 timestamp when the subscription was last modified
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
        deserializerMap['custom_price'] = (node) => customPrice = node.getObjectValue<AmountModel>(AmountModel.createFromDiscriminatorValue);
        deserializerMap['custom_recurrency'] = (node) => customRecurrency = node.getIntValue();
        deserializerMap['end_date'] = (node) => endDate = node.getDateTimeValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['free_trial'] = (node) => freeTrial = node.getObjectValue<SubscriptionFreeTrialModel>(SubscriptionFreeTrialModel.createFromDiscriminatorValue);
        deserializerMap['friendly_name'] = (node) => friendlyName = node.getStringValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['is_active'] = (node) => isActive = node.getBoolValue();
        deserializerMap['label_refs'] = (node) => labelRefs = node.getCollectionOfObjectValues<LabelRefModel>(LabelRefModel.createFromDiscriminatorValue);
        deserializerMap['owner'] = (node) => owner = node.getObjectValue<OwnerModel>(OwnerModel.createFromDiscriminatorValue);
        deserializerMap['payer'] = (node) => payer = node.getObjectValue<SubscriptionPayerModel>(SubscriptionPayerModel.createFromDiscriminatorValue);
        deserializerMap['plan_id'] = (node) => planId = node.getStringValue();
        deserializerMap['price_id'] = (node) => priceId = node.getStringValue();
        deserializerMap['provider_id'] = (node) => providerId = node.getStringValue();
        deserializerMap['recurrency'] = (node) => recurrency = node.getEnumValue<SubscriptionModelRecurrency>((stringValue) => SubscriptionModelRecurrency.values.where((enumVal) => enumVal.value == stringValue).firstOrNull);
        deserializerMap['service_users'] = (node) => serviceUsers = node.getCollectionOfPrimitiveValues<String>();
        deserializerMap['start_date'] = (node) => startDate = node.getDateTimeValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeObjectValue<AmountModel>('custom_price', customPrice);
        writer.writeIntValue('custom_recurrency', customRecurrency);
        writer.writeDateTimeValue('end_date', endDate);
        writer.writeStringValue('etag', etag);
        writer.writeObjectValue<SubscriptionFreeTrialModel>('free_trial', freeTrial);
        writer.writeStringValue('friendly_name', friendlyName);
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_active', value:isActive);
        writer.writeCollectionOfObjectValues<LabelRefModel>('label_refs', labelRefs);
        writer.writeObjectValue<OwnerModel>('owner', owner);
        writer.writeObjectValue<SubscriptionPayerModel>('payer', payer);
        writer.writeStringValue('plan_id', planId);
        writer.writeStringValue('price_id', priceId);
        writer.writeStringValue('provider_id', providerId);
        writer.writeEnumValue<SubscriptionModelRecurrency>('recurrency', recurrency, (e) => e?.value);
        writer.writeCollectionOfPrimitiveValues<String?>('service_users', serviceUsers);
        writer.writeDateTimeValue('start_date', startDate);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
