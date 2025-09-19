//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyFamilyMemberModel {
  /// Returns a new [FamilyFamilyMemberModel] instance.
  FamilyFamilyMemberModel({
    required this.createdAt,
    required this.etag,
    required this.familyId,
    required this.hasAccount,
    required this.id,
    required this.isYou,
    required this.name,
    required this.type,
    required this.updatedAt,
  });

  /// @Description Timestamp when the member was created
  DateTime createdAt;

  /// @Description Entity tag for optimistic concurrency control
  String etag;

  /// @Description ID of the family this member belongs to
  String familyId;

  /// @Description Indicates whether this member has an account with the service provider
  bool hasAccount;

  /// @Description Unique identifier for the family member
  String id;

  /// @Description Indicates whether this member is the current authenticated user
  bool isYou;

  /// @Description Name of the family member
  String name;

  /// @Description Whether this member is a child (affects permissions and features)
  FamilyFamilyMemberModelTypeEnum type;

  /// @Description Timestamp when the member was last updated
  DateTime updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyFamilyMemberModel &&
    other.createdAt == createdAt &&
    other.etag == etag &&
    other.familyId == familyId &&
    other.hasAccount == hasAccount &&
    other.id == id &&
    other.isYou == isYou &&
    other.name == name &&
    other.type == type &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt.hashCode) +
    (etag.hashCode) +
    (familyId.hashCode) +
    (hasAccount.hashCode) +
    (id.hashCode) +
    (isYou.hashCode) +
    (name.hashCode) +
    (type.hashCode) +
    (updatedAt.hashCode);

  @override
  String toString() => 'FamilyFamilyMemberModel[createdAt=$createdAt, etag=$etag, familyId=$familyId, hasAccount=$hasAccount, id=$id, isYou=$isYou, name=$name, type=$type, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
      json[r'etag'] = this.etag;
      json[r'family_id'] = this.familyId;
      json[r'has_account'] = this.hasAccount;
      json[r'id'] = this.id;
      json[r'is_you'] = this.isYou;
      json[r'name'] = this.name;
      json[r'type'] = this.type;
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [FamilyFamilyMemberModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyFamilyMemberModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyFamilyMemberModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyFamilyMemberModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyFamilyMemberModel(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        etag: mapValueOfType<String>(json, r'etag')!,
        familyId: mapValueOfType<String>(json, r'family_id')!,
        hasAccount: mapValueOfType<bool>(json, r'has_account')!,
        id: mapValueOfType<String>(json, r'id')!,
        isYou: mapValueOfType<bool>(json, r'is_you')!,
        name: mapValueOfType<String>(json, r'name')!,
        type: FamilyFamilyMemberModelTypeEnum.fromJson(json[r'type'])!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<FamilyFamilyMemberModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyMemberModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyMemberModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyFamilyMemberModel> mapFromJson(dynamic json) {
    final map = <String, FamilyFamilyMemberModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyFamilyMemberModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyFamilyMemberModel-objects as value to a dart map
  static Map<String, List<FamilyFamilyMemberModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyFamilyMemberModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyFamilyMemberModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'created_at',
    'etag',
    'family_id',
    'has_account',
    'id',
    'is_you',
    'name',
    'type',
    'updated_at',
  };
}

/// @Description Whether this member is a child (affects permissions and features)
class FamilyFamilyMemberModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const FamilyFamilyMemberModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const owner = FamilyFamilyMemberModelTypeEnum._(r'owner');
  static const adult = FamilyFamilyMemberModelTypeEnum._(r'adult');
  static const kid = FamilyFamilyMemberModelTypeEnum._(r'kid');

  /// List of all possible values in this [enum][FamilyFamilyMemberModelTypeEnum].
  static const values = <FamilyFamilyMemberModelTypeEnum>[
    owner,
    adult,
    kid,
  ];

  static FamilyFamilyMemberModelTypeEnum? fromJson(dynamic value) => FamilyFamilyMemberModelTypeEnumTypeTransformer().decode(value);

  static List<FamilyFamilyMemberModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyMemberModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyMemberModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [FamilyFamilyMemberModelTypeEnum] to String,
/// and [decode] dynamic data back to [FamilyFamilyMemberModelTypeEnum].
class FamilyFamilyMemberModelTypeEnumTypeTransformer {
  factory FamilyFamilyMemberModelTypeEnumTypeTransformer() => _instance ??= const FamilyFamilyMemberModelTypeEnumTypeTransformer._();

  const FamilyFamilyMemberModelTypeEnumTypeTransformer._();

  String encode(FamilyFamilyMemberModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a FamilyFamilyMemberModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  FamilyFamilyMemberModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'owner': return FamilyFamilyMemberModelTypeEnum.owner;
        case r'adult': return FamilyFamilyMemberModelTypeEnum.adult;
        case r'kid': return FamilyFamilyMemberModelTypeEnum.kid;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [FamilyFamilyMemberModelTypeEnumTypeTransformer] instance.
  static FamilyFamilyMemberModelTypeEnumTypeTransformer? _instance;
}


