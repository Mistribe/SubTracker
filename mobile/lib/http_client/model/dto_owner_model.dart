//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class DtoOwnerModel {
  /// Returns a new [DtoOwnerModel] instance.
  DtoOwnerModel({
    required this.etag,
    this.familyId,
    required this.type,
    this.userId,
  });

  /// @Description Entity tag for optimistic concurrency control
  String etag;

  /// @Description Family ID when an ownership type is family (required for family ownership)
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? familyId;

  /// @Description Type of ownership (personal, family or system)
  DtoOwnerModelTypeEnum type;

  /// @Description UserProfile ID when an ownership type is personal (required for personal ownership)
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? userId;

  @override
  bool operator ==(Object other) => identical(this, other) || other is DtoOwnerModel &&
    other.etag == etag &&
    other.familyId == familyId &&
    other.type == type &&
    other.userId == userId;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (etag.hashCode) +
    (familyId == null ? 0 : familyId!.hashCode) +
    (type.hashCode) +
    (userId == null ? 0 : userId!.hashCode);

  @override
  String toString() => 'DtoOwnerModel[etag=$etag, familyId=$familyId, type=$type, userId=$userId]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'etag'] = this.etag;
    if (this.familyId != null) {
      json[r'family_id'] = this.familyId;
    } else {
      json[r'family_id'] = null;
    }
      json[r'type'] = this.type;
    if (this.userId != null) {
      json[r'userId'] = this.userId;
    } else {
      json[r'userId'] = null;
    }
    return json;
  }

  /// Returns a new [DtoOwnerModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static DtoOwnerModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "DtoOwnerModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "DtoOwnerModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return DtoOwnerModel(
        etag: mapValueOfType<String>(json, r'etag')!,
        familyId: mapValueOfType<String>(json, r'family_id'),
        type: DtoOwnerModelTypeEnum.fromJson(json[r'type'])!,
        userId: mapValueOfType<String>(json, r'userId'),
      );
    }
    return null;
  }

  static List<DtoOwnerModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DtoOwnerModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DtoOwnerModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, DtoOwnerModel> mapFromJson(dynamic json) {
    final map = <String, DtoOwnerModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = DtoOwnerModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of DtoOwnerModel-objects as value to a dart map
  static Map<String, List<DtoOwnerModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<DtoOwnerModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = DtoOwnerModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'etag',
    'type',
  };
}

/// @Description Type of ownership (personal, family or system)
class DtoOwnerModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const DtoOwnerModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const personal = DtoOwnerModelTypeEnum._(r'personal');
  static const family = DtoOwnerModelTypeEnum._(r'family');
  static const system = DtoOwnerModelTypeEnum._(r'system');

  /// List of all possible values in this [enum][DtoOwnerModelTypeEnum].
  static const values = <DtoOwnerModelTypeEnum>[
    personal,
    family,
    system,
  ];

  static DtoOwnerModelTypeEnum? fromJson(dynamic value) => DtoOwnerModelTypeEnumTypeTransformer().decode(value);

  static List<DtoOwnerModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DtoOwnerModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DtoOwnerModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [DtoOwnerModelTypeEnum] to String,
/// and [decode] dynamic data back to [DtoOwnerModelTypeEnum].
class DtoOwnerModelTypeEnumTypeTransformer {
  factory DtoOwnerModelTypeEnumTypeTransformer() => _instance ??= const DtoOwnerModelTypeEnumTypeTransformer._();

  const DtoOwnerModelTypeEnumTypeTransformer._();

  String encode(DtoOwnerModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a DtoOwnerModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  DtoOwnerModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'personal': return DtoOwnerModelTypeEnum.personal;
        case r'family': return DtoOwnerModelTypeEnum.family;
        case r'system': return DtoOwnerModelTypeEnum.system;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [DtoOwnerModelTypeEnumTypeTransformer] instance.
  static DtoOwnerModelTypeEnumTypeTransformer? _instance;
}


