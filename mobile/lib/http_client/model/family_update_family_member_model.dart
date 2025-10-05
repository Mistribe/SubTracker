//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyUpdateFamilyMemberModel {
  /// Returns a new [FamilyUpdateFamilyMemberModel] instance.
  FamilyUpdateFamilyMemberModel({
    required this.name,
    required this.type,
    this.updatedAt,
  });

  String name;

  FamilyUpdateFamilyMemberModelTypeEnum type;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyUpdateFamilyMemberModel &&
    other.name == name &&
    other.type == type &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (name.hashCode) +
    (type.hashCode) +
    (updatedAt == null ? 0 : updatedAt!.hashCode);

  @override
  String toString() => 'FamilyUpdateFamilyMemberModel[name=$name, type=$type, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'name'] = this.name;
      json[r'type'] = this.type;
    if (this.updatedAt != null) {
      json[r'updated_at'] = this.updatedAt!.toUtc().toIso8601String();
    } else {
      json[r'updated_at'] = null;
    }
    return json;
  }

  /// Returns a new [FamilyUpdateFamilyMemberModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyUpdateFamilyMemberModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyUpdateFamilyMemberModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyUpdateFamilyMemberModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyUpdateFamilyMemberModel(
        name: mapValueOfType<String>(json, r'name')!,
        type: FamilyUpdateFamilyMemberModelTypeEnum.fromJson(json[r'type'])!,
        updatedAt: mapDateTime(json, r'updated_at', r''),
      );
    }
    return null;
  }

  static List<FamilyUpdateFamilyMemberModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyUpdateFamilyMemberModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyUpdateFamilyMemberModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyUpdateFamilyMemberModel> mapFromJson(dynamic json) {
    final map = <String, FamilyUpdateFamilyMemberModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyUpdateFamilyMemberModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyUpdateFamilyMemberModel-objects as value to a dart map
  static Map<String, List<FamilyUpdateFamilyMemberModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyUpdateFamilyMemberModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyUpdateFamilyMemberModel.listFromJson(entry.value, growable: growable,);
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


class FamilyUpdateFamilyMemberModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const FamilyUpdateFamilyMemberModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const owner = FamilyUpdateFamilyMemberModelTypeEnum._(r'owner');
  static const adult = FamilyUpdateFamilyMemberModelTypeEnum._(r'adult');
  static const kid = FamilyUpdateFamilyMemberModelTypeEnum._(r'kid');

  /// List of all possible values in this [enum][FamilyUpdateFamilyMemberModelTypeEnum].
  static const values = <FamilyUpdateFamilyMemberModelTypeEnum>[
    owner,
    adult,
    kid,
  ];

  static FamilyUpdateFamilyMemberModelTypeEnum? fromJson(dynamic value) => FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer().decode(value);

  static List<FamilyUpdateFamilyMemberModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyUpdateFamilyMemberModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyUpdateFamilyMemberModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [FamilyUpdateFamilyMemberModelTypeEnum] to String,
/// and [decode] dynamic data back to [FamilyUpdateFamilyMemberModelTypeEnum].
class FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer {
  factory FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer() => _instance ??= const FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer._();

  const FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer._();

  String encode(FamilyUpdateFamilyMemberModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a FamilyUpdateFamilyMemberModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  FamilyUpdateFamilyMemberModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'owner': return FamilyUpdateFamilyMemberModelTypeEnum.owner;
        case r'adult': return FamilyUpdateFamilyMemberModelTypeEnum.adult;
        case r'kid': return FamilyUpdateFamilyMemberModelTypeEnum.kid;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer] instance.
  static FamilyUpdateFamilyMemberModelTypeEnumTypeTransformer? _instance;
}


