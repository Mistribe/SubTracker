// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './family_member_model.dart';

/// auto generated
class FamilyModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  The created_at property
    DateTime? createdAt;
    ///  The etag property
    String? etag;
    ///  The have_joint_account property
    bool? haveJointAccount;
    ///  The id property
    String? id;
    ///  The is_owner property
    bool? isOwner;
    ///  The members property
    Iterable<FamilyMemberModel>? members;
    ///  The name property
    String? name;
    ///  The updated_at property
    DateTime? updatedAt;
    /// Instantiates a new [FamilyModel] and sets the default values.
    FamilyModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static FamilyModel createFromDiscriminatorValue(ParseNode parseNode) {
        return FamilyModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['created_at'] = (node) => createdAt = node.getDateTimeValue();
        deserializerMap['etag'] = (node) => etag = node.getStringValue();
        deserializerMap['have_joint_account'] = (node) => haveJointAccount = node.getBoolValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['is_owner'] = (node) => isOwner = node.getBoolValue();
        deserializerMap['members'] = (node) => members = node.getCollectionOfObjectValues<FamilyMemberModel>(FamilyMemberModel.createFromDiscriminatorValue);
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeDateTimeValue('created_at', createdAt);
        writer.writeStringValue('etag', etag);
        writer.writeBoolValue('have_joint_account', value:haveJointAccount);
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_owner', value:isOwner);
        writer.writeCollectionOfObjectValues<FamilyMemberModel>('members', members);
        writer.writeStringValue('name', name);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
