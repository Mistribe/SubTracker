//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyCreateFamilyMemberModel {
  /// Returns a new [FamilyCreateFamilyMemberModel] instance.
  FamilyCreateFamilyMemberModel({
    this.createdAt,
    this.id,
    required this.name,
    required this.type,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? createdAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? id;

  String name;

  FamilyCreateFamilyMemberModelTypeEnum type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyCreateFamilyMemberModel &&
    other.createdAt == createdAt &&
    other.id == id &&
    other.name == name &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt == null ? 0 : createdAt!.hashCode) +
    (id == null ? 0 : id!.hashCode) +
    (name.hashCode) +
    (type.hashCode);

  @override
  String toString() => 'FamilyCreateFamilyMemberModel[createdAt=$createdAt, id=$id, name=$name, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.createdAt != null) {
      json[r'created_at'] = this.createdAt!.toUtc().toIso8601String();
    } else {
      json[r'created_at'] = null;
    }
    if (this.id != null) {
      json[r'id'] = this.id;
    } else {
      json[r'id'] = null;
    }
      json[r'name'] = this.name;
      json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [FamilyCreateFamilyMemberModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyCreateFamilyMemberModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyCreateFamilyMemberModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyCreateFamilyMemberModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyCreateFamilyMemberModel(
        createdAt: mapDateTime(json, r'created_at', r''),
        id: mapValueOfType<String>(json, r'id'),
        name: mapValueOfType<String>(json, r'name')!,
        type: FamilyCreateFamilyMemberModelTypeEnum.fromJson(json[r'type'])!,
      );
    }
    return null;
  }

  static List<FamilyCreateFamilyMemberModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyCreateFamilyMemberModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyCreateFamilyMemberModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyCreateFamilyMemberModel> mapFromJson(dynamic json) {
    final map = <String, FamilyCreateFamilyMemberModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyCreateFamilyMemberModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyCreateFamilyMemberModel-objects as value to a dart map
  static Map<String, List<FamilyCreateFamilyMemberModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyCreateFamilyMemberModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyCreateFamilyMemberModel.listFromJson(entry.value, growable: growable,);
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


class FamilyCreateFamilyMemberModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const FamilyCreateFamilyMemberModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const owner = FamilyCreateFamilyMemberModelTypeEnum._(r'owner');
  static const adult = FamilyCreateFamilyMemberModelTypeEnum._(r'adult');
  static const kid = FamilyCreateFamilyMemberModelTypeEnum._(r'kid');

  /// List of all possible values in this [enum][FamilyCreateFamilyMemberModelTypeEnum].
  static const values = <FamilyCreateFamilyMemberModelTypeEnum>[
    owner,
    adult,
    kid,
  ];

  static FamilyCreateFamilyMemberModelTypeEnum? fromJson(dynamic value) => FamilyCreateFamilyMemberModelTypeEnumTypeTransformer().decode(value);

  static List<FamilyCreateFamilyMemberModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyCreateFamilyMemberModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyCreateFamilyMemberModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [FamilyCreateFamilyMemberModelTypeEnum] to String,
/// and [decode] dynamic data back to [FamilyCreateFamilyMemberModelTypeEnum].
class FamilyCreateFamilyMemberModelTypeEnumTypeTransformer {
  factory FamilyCreateFamilyMemberModelTypeEnumTypeTransformer() => _instance ??= const FamilyCreateFamilyMemberModelTypeEnumTypeTransformer._();

  const FamilyCreateFamilyMemberModelTypeEnumTypeTransformer._();

  String encode(FamilyCreateFamilyMemberModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a FamilyCreateFamilyMemberModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  FamilyCreateFamilyMemberModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'owner': return FamilyCreateFamilyMemberModelTypeEnum.owner;
        case r'adult': return FamilyCreateFamilyMemberModelTypeEnum.adult;
        case r'kid': return FamilyCreateFamilyMemberModelTypeEnum.kid;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [FamilyCreateFamilyMemberModelTypeEnumTypeTransformer] instance.
  static FamilyCreateFamilyMemberModelTypeEnumTypeTransformer? _instance;
}


