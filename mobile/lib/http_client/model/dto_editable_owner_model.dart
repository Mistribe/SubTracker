//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class DtoEditableOwnerModel {
  /// Returns a new [DtoEditableOwnerModel] instance.
  DtoEditableOwnerModel({
    this.familyId,
    required this.type,
  });

  /// @Description Family ID when an ownership type is family (required for family ownership)
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? familyId;

  /// @Description Type of ownership (personal, family or system)
  DtoEditableOwnerModelTypeEnum type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is DtoEditableOwnerModel &&
    other.familyId == familyId &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (familyId == null ? 0 : familyId!.hashCode) +
    (type.hashCode);

  @override
  String toString() => 'DtoEditableOwnerModel[familyId=$familyId, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.familyId != null) {
      json[r'family_id'] = this.familyId;
    } else {
      json[r'family_id'] = null;
    }
      json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [DtoEditableOwnerModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static DtoEditableOwnerModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "DtoEditableOwnerModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "DtoEditableOwnerModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return DtoEditableOwnerModel(
        familyId: mapValueOfType<String>(json, r'family_id'),
        type: DtoEditableOwnerModelTypeEnum.fromJson(json[r'type'])!,
      );
    }
    return null;
  }

  static List<DtoEditableOwnerModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DtoEditableOwnerModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DtoEditableOwnerModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, DtoEditableOwnerModel> mapFromJson(dynamic json) {
    final map = <String, DtoEditableOwnerModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = DtoEditableOwnerModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of DtoEditableOwnerModel-objects as value to a dart map
  static Map<String, List<DtoEditableOwnerModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<DtoEditableOwnerModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = DtoEditableOwnerModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'type',
  };
}

/// @Description Type of ownership (personal, family or system)
class DtoEditableOwnerModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const DtoEditableOwnerModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const personal = DtoEditableOwnerModelTypeEnum._(r'personal');
  static const family = DtoEditableOwnerModelTypeEnum._(r'family');
  static const system = DtoEditableOwnerModelTypeEnum._(r'system');

  /// List of all possible values in this [enum][DtoEditableOwnerModelTypeEnum].
  static const values = <DtoEditableOwnerModelTypeEnum>[
    personal,
    family,
    system,
  ];

  static DtoEditableOwnerModelTypeEnum? fromJson(dynamic value) => DtoEditableOwnerModelTypeEnumTypeTransformer().decode(value);

  static List<DtoEditableOwnerModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <DtoEditableOwnerModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = DtoEditableOwnerModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [DtoEditableOwnerModelTypeEnum] to String,
/// and [decode] dynamic data back to [DtoEditableOwnerModelTypeEnum].
class DtoEditableOwnerModelTypeEnumTypeTransformer {
  factory DtoEditableOwnerModelTypeEnumTypeTransformer() => _instance ??= const DtoEditableOwnerModelTypeEnumTypeTransformer._();

  const DtoEditableOwnerModelTypeEnumTypeTransformer._();

  String encode(DtoEditableOwnerModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a DtoEditableOwnerModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  DtoEditableOwnerModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'personal': return DtoEditableOwnerModelTypeEnum.personal;
        case r'family': return DtoEditableOwnerModelTypeEnum.family;
        case r'system': return DtoEditableOwnerModelTypeEnum.system;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [DtoEditableOwnerModelTypeEnumTypeTransformer] instance.
  static DtoEditableOwnerModelTypeEnumTypeTransformer? _instance;
}


