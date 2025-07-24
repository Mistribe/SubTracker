// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './patch_family_member_model.dart';

/// auto generated
/// Model for updating family details
class PatchFamilyModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  Indicates if family has joint account
    bool? haveJointAccount;
    ///  Family ID
    String? id;
    ///  List of family members
    Iterable<PatchFamilyMemberModel>? members;
    ///  Family name
    String? name;
    ///  Optional timestamp of the last update
    DateTime? updatedAt;
    /// Instantiates a new [PatchFamilyModel] and sets the default values.
    PatchFamilyModel() :  
        additionalData = {};
    /// Creates a new instance of the appropriate class based on discriminator value
    ///  [parseNode] The parse node to use to read the discriminator value and create the object
    static PatchFamilyModel createFromDiscriminatorValue(ParseNode parseNode) {
        return PatchFamilyModel();
    }
    /// The deserialization information for the current model
    @override
    Map<String, void Function(ParseNode)> getFieldDeserializers() {
        var deserializerMap = <String, void Function(ParseNode)>{};
        deserializerMap['have_joint_account'] = (node) => haveJointAccount = node.getBoolValue();
        deserializerMap['id'] = (node) => id = node.getStringValue();
        deserializerMap['members'] = (node) => members = node.getCollectionOfObjectValues<PatchFamilyMemberModel>(PatchFamilyMemberModel.createFromDiscriminatorValue);
        deserializerMap['name'] = (node) => name = node.getStringValue();
        deserializerMap['updated_at'] = (node) => updatedAt = node.getDateTimeValue();
        return deserializerMap;
    }
    /// Serializes information the current object
    ///  [writer] Serialization writer to use to serialize this model
    @override
    void serialize(SerializationWriter writer) {
        writer.writeBoolValue('have_joint_account', value:haveJointAccount);
        writer.writeStringValue('id', id);
        writer.writeCollectionOfObjectValues<PatchFamilyMemberModel>('members', members);
        writer.writeStringValue('name', name);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
