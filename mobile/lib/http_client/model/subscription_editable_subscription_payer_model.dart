//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionEditableSubscriptionPayerModel {
  /// Returns a new [SubscriptionEditableSubscriptionPayerModel] instance.
  SubscriptionEditableSubscriptionPayerModel({
    required this.familyId,
    this.memberId,
    required this.type,
  });

  /// @Description ID of the family associated with this payer
  String familyId;

  /// @Description ID of the specific family member who pays (required when type is family_member)
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? memberId;

  /// @Description Type of payer (family or family member)
  SubscriptionEditableSubscriptionPayerModelTypeEnum type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionEditableSubscriptionPayerModel &&
    other.familyId == familyId &&
    other.memberId == memberId &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (familyId.hashCode) +
    (memberId == null ? 0 : memberId!.hashCode) +
    (type.hashCode);

  @override
  String toString() => 'SubscriptionEditableSubscriptionPayerModel[familyId=$familyId, memberId=$memberId, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'family_id'] = this.familyId;
    if (this.memberId != null) {
      json[r'memberId'] = this.memberId;
    } else {
      json[r'memberId'] = null;
    }
      json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [SubscriptionEditableSubscriptionPayerModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionEditableSubscriptionPayerModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionEditableSubscriptionPayerModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionEditableSubscriptionPayerModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionEditableSubscriptionPayerModel(
        familyId: mapValueOfType<String>(json, r'family_id')!,
        memberId: mapValueOfType<String>(json, r'memberId'),
        type: SubscriptionEditableSubscriptionPayerModelTypeEnum.fromJson(json[r'type'])!,
      );
    }
    return null;
  }

  static List<SubscriptionEditableSubscriptionPayerModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionEditableSubscriptionPayerModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionEditableSubscriptionPayerModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionEditableSubscriptionPayerModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionEditableSubscriptionPayerModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionEditableSubscriptionPayerModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionEditableSubscriptionPayerModel-objects as value to a dart map
  static Map<String, List<SubscriptionEditableSubscriptionPayerModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionEditableSubscriptionPayerModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionEditableSubscriptionPayerModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'family_id',
    'type',
  };
}

/// @Description Type of payer (family or family member)
class SubscriptionEditableSubscriptionPayerModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const SubscriptionEditableSubscriptionPayerModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const family = SubscriptionEditableSubscriptionPayerModelTypeEnum._(r'family');
  static const familyMember = SubscriptionEditableSubscriptionPayerModelTypeEnum._(r'family_member');

  /// List of all possible values in this [enum][SubscriptionEditableSubscriptionPayerModelTypeEnum].
  static const values = <SubscriptionEditableSubscriptionPayerModelTypeEnum>[
    family,
    familyMember,
  ];

  static SubscriptionEditableSubscriptionPayerModelTypeEnum? fromJson(dynamic value) => SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer().decode(value);

  static List<SubscriptionEditableSubscriptionPayerModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionEditableSubscriptionPayerModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionEditableSubscriptionPayerModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [SubscriptionEditableSubscriptionPayerModelTypeEnum] to String,
/// and [decode] dynamic data back to [SubscriptionEditableSubscriptionPayerModelTypeEnum].
class SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer {
  factory SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer() => _instance ??= const SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer._();

  const SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer._();

  String encode(SubscriptionEditableSubscriptionPayerModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a SubscriptionEditableSubscriptionPayerModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  SubscriptionEditableSubscriptionPayerModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'family': return SubscriptionEditableSubscriptionPayerModelTypeEnum.family;
        case r'family_member': return SubscriptionEditableSubscriptionPayerModelTypeEnum.familyMember;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer] instance.
  static SubscriptionEditableSubscriptionPayerModelTypeEnumTypeTransformer? _instance;
}


