// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';
import './family_member_model.dart';

/// auto generated
/// Family object containing family information and members
class FamilyModel implements AdditionalDataHolder, Parsable {
    ///  Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    @override
    Map<String, Object?> additionalData;
    ///  @Description ISO 8601 timestamp indicating when the family was originally created
    DateTime? createdAt;
    ///  @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
    String? etag;
    ///  @Description Unique identifier for the family (UUID format)
    String? id;
    ///  @Description Indicates whether the current authenticated user is the owner of this family
    bool? isOwner;
    ///  @Description Complete list of all members belonging to this family
    Iterable<FamilyMemberModel>? members;
    ///  @Description Display name of the family
    String? name;
    ///  @Description ISO 8601 timestamp indicating when the family information was last modified
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
        writer.writeStringValue('id', id);
        writer.writeBoolValue('is_owner', value:isOwner);
        writer.writeCollectionOfObjectValues<FamilyMemberModel>('members', members);
        writer.writeStringValue('name', name);
        writer.writeDateTimeValue('updated_at', updatedAt);
        writer.writeAdditionalData(additionalData);
    }
}
