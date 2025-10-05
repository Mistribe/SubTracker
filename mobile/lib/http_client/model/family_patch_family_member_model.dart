//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyPatchFamilyMemberModel {
  /// Returns a new [FamilyPatchFamilyMemberModel] instance.
  FamilyPatchFamilyMemberModel({
    this.id,
    required this.name,
    required this.type,
    this.updatedAt,
  });

  /// Optional member ID. If not provided, new member will be created
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? id;

  /// member's name
  String name;

  /// Indicates if the member is a kid
  FamilyPatchFamilyMemberModelTypeEnum type;

  /// Optional timestamp of the last update
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyPatchFamilyMemberModel &&
    other.id == id &&
    other.name == name &&
    other.type == type &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (id == null ? 0 : id!.hashCode) +
    (name.hashCode) +
    (type.hashCode) +
    (updatedAt == null ? 0 : updatedAt!.hashCode);

  @override
  String toString() => 'FamilyPatchFamilyMemberModel[id=$id, name=$name, type=$type, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.id != null) {
      json[r'id'] = this.id;
    } else {
      json[r'id'] = null;
    }
      json[r'name'] = this.name;
      json[r'type'] = this.type;
    if (this.updatedAt != null) {
      json[r'updated_at'] = this.updatedAt!.toUtc().toIso8601String();
    } else {
      json[r'updated_at'] = null;
    }
    return json;
  }

  /// Returns a new [FamilyPatchFamilyMemberModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyPatchFamilyMemberModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyPatchFamilyMemberModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyPatchFamilyMemberModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyPatchFamilyMemberModel(
        id: mapValueOfType<String>(json, r'id'),
        name: mapValueOfType<String>(json, r'name')!,
        type: FamilyPatchFamilyMemberModelTypeEnum.fromJson(json[r'type'])!,
        updatedAt: mapDateTime(json, r'updated_at', r''),
      );
    }
    return null;
  }

  static List<FamilyPatchFamilyMemberModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyPatchFamilyMemberModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyPatchFamilyMemberModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyPatchFamilyMemberModel> mapFromJson(dynamic json) {
    final map = <String, FamilyPatchFamilyMemberModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyPatchFamilyMemberModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyPatchFamilyMemberModel-objects as value to a dart map
  static Map<String, List<FamilyPatchFamilyMemberModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyPatchFamilyMemberModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyPatchFamilyMemberModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'name',
    'type',
  };
}

/// Indicates if the member is a kid
class FamilyPatchFamilyMemberModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const FamilyPatchFamilyMemberModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const owner = FamilyPatchFamilyMemberModelTypeEnum._(r'owner');
  static const adult = FamilyPatchFamilyMemberModelTypeEnum._(r'adult');
  static const kid = FamilyPatchFamilyMemberModelTypeEnum._(r'kid');

  /// List of all possible values in this [enum][FamilyPatchFamilyMemberModelTypeEnum].
  static const values = <FamilyPatchFamilyMemberModelTypeEnum>[
    owner,
    adult,
    kid,
  ];

  static FamilyPatchFamilyMemberModelTypeEnum? fromJson(dynamic value) => FamilyPatchFamilyMemberModelTypeEnumTypeTransformer().decode(value);

  static List<FamilyPatchFamilyMemberModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyPatchFamilyMemberModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyPatchFamilyMemberModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [FamilyPatchFamilyMemberModelTypeEnum] to String,
/// and [decode] dynamic data back to [FamilyPatchFamilyMemberModelTypeEnum].
class FamilyPatchFamilyMemberModelTypeEnumTypeTransformer {
  factory FamilyPatchFamilyMemberModelTypeEnumTypeTransformer() => _instance ??= const FamilyPatchFamilyMemberModelTypeEnumTypeTransformer._();

  const FamilyPatchFamilyMemberModelTypeEnumTypeTransformer._();

  String encode(FamilyPatchFamilyMemberModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a FamilyPatchFamilyMemberModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  FamilyPatchFamilyMemberModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'owner': return FamilyPatchFamilyMemberModelTypeEnum.owner;
        case r'adult': return FamilyPatchFamilyMemberModelTypeEnum.adult;
        case r'kid': return FamilyPatchFamilyMemberModelTypeEnum.kid;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [FamilyPatchFamilyMemberModelTypeEnumTypeTransformer] instance.
  static FamilyPatchFamilyMemberModelTypeEnumTypeTransformer? _instance;
}


