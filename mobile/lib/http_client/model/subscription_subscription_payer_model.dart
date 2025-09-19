//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionSubscriptionPayerModel {
  /// Returns a new [SubscriptionSubscriptionPayerModel] instance.
  SubscriptionSubscriptionPayerModel({
    required this.etag,
    required this.familyId,
    this.memberId,
    required this.type,
  });

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

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
  SubscriptionSubscriptionPayerModelTypeEnum type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionSubscriptionPayerModel &&
    other.etag == etag &&
    other.familyId == familyId &&
    other.memberId == memberId &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (etag.hashCode) +
    (familyId.hashCode) +
    (memberId == null ? 0 : memberId!.hashCode) +
    (type.hashCode);

  @override
  String toString() => 'SubscriptionSubscriptionPayerModel[etag=$etag, familyId=$familyId, memberId=$memberId, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'etag'] = this.etag;
      json[r'family_id'] = this.familyId;
    if (this.memberId != null) {
      json[r'memberId'] = this.memberId;
    } else {
      json[r'memberId'] = null;
    }
      json[r'type'] = this.type;
    return json;
  }

  /// Returns a new [SubscriptionSubscriptionPayerModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionSubscriptionPayerModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionSubscriptionPayerModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionSubscriptionPayerModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionSubscriptionPayerModel(
        etag: mapValueOfType<String>(json, r'etag')!,
        familyId: mapValueOfType<String>(json, r'family_id')!,
        memberId: mapValueOfType<String>(json, r'memberId'),
        type: SubscriptionSubscriptionPayerModelTypeEnum.fromJson(json[r'type'])!,
      );
    }
    return null;
  }

  static List<SubscriptionSubscriptionPayerModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionPayerModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionPayerModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionSubscriptionPayerModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionSubscriptionPayerModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionSubscriptionPayerModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionSubscriptionPayerModel-objects as value to a dart map
  static Map<String, List<SubscriptionSubscriptionPayerModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionSubscriptionPayerModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionSubscriptionPayerModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'etag',
    'family_id',
    'type',
  };
}

/// @Description Type of payer (family or family member)
class SubscriptionSubscriptionPayerModelTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const SubscriptionSubscriptionPayerModelTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const family = SubscriptionSubscriptionPayerModelTypeEnum._(r'family');
  static const familyMember = SubscriptionSubscriptionPayerModelTypeEnum._(r'family_member');

  /// List of all possible values in this [enum][SubscriptionSubscriptionPayerModelTypeEnum].
  static const values = <SubscriptionSubscriptionPayerModelTypeEnum>[
    family,
    familyMember,
  ];

  static SubscriptionSubscriptionPayerModelTypeEnum? fromJson(dynamic value) => SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer().decode(value);

  static List<SubscriptionSubscriptionPayerModelTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionPayerModelTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionPayerModelTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [SubscriptionSubscriptionPayerModelTypeEnum] to String,
/// and [decode] dynamic data back to [SubscriptionSubscriptionPayerModelTypeEnum].
class SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer {
  factory SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer() => _instance ??= const SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer._();

  const SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer._();

  String encode(SubscriptionSubscriptionPayerModelTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a SubscriptionSubscriptionPayerModelTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  SubscriptionSubscriptionPayerModelTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'family': return SubscriptionSubscriptionPayerModelTypeEnum.family;
        case r'family_member': return SubscriptionSubscriptionPayerModelTypeEnum.familyMember;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer] instance.
  static SubscriptionSubscriptionPayerModelTypeEnumTypeTransformer? _instance;
}


