//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyFamilyInviteRequest {
  /// Returns a new [FamilyFamilyInviteRequest] instance.
  FamilyFamilyInviteRequest({
    this.email,
    required this.familyMemberId,
    this.name,
    this.type,
  });

  /// Email of the invited member
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? email;

  /// ID of the family member to be invited
  String familyMemberId;

  /// Name of the invited member
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? name;

  /// Type of the member (adult or kid)
  FamilyFamilyInviteRequestTypeEnum? type;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyFamilyInviteRequest &&
    other.email == email &&
    other.familyMemberId == familyMemberId &&
    other.name == name &&
    other.type == type;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (email == null ? 0 : email!.hashCode) +
    (familyMemberId.hashCode) +
    (name == null ? 0 : name!.hashCode) +
    (type == null ? 0 : type!.hashCode);

  @override
  String toString() => 'FamilyFamilyInviteRequest[email=$email, familyMemberId=$familyMemberId, name=$name, type=$type]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.email != null) {
      json[r'email'] = this.email;
    } else {
      json[r'email'] = null;
    }
      json[r'family_member_id'] = this.familyMemberId;
    if (this.name != null) {
      json[r'name'] = this.name;
    } else {
      json[r'name'] = null;
    }
    if (this.type != null) {
      json[r'type'] = this.type;
    } else {
      json[r'type'] = null;
    }
    return json;
  }

  /// Returns a new [FamilyFamilyInviteRequest] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyFamilyInviteRequest? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyFamilyInviteRequest[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyFamilyInviteRequest[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyFamilyInviteRequest(
        email: mapValueOfType<String>(json, r'email'),
        familyMemberId: mapValueOfType<String>(json, r'family_member_id')!,
        name: mapValueOfType<String>(json, r'name'),
        type: FamilyFamilyInviteRequestTypeEnum.fromJson(json[r'type']),
      );
    }
    return null;
  }

  static List<FamilyFamilyInviteRequest> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyInviteRequest>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyInviteRequest.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyFamilyInviteRequest> mapFromJson(dynamic json) {
    final map = <String, FamilyFamilyInviteRequest>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyFamilyInviteRequest.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyFamilyInviteRequest-objects as value to a dart map
  static Map<String, List<FamilyFamilyInviteRequest>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyFamilyInviteRequest>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyFamilyInviteRequest.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'family_member_id',
  };
}

/// Type of the member (adult or kid)
class FamilyFamilyInviteRequestTypeEnum {
  /// Instantiate a new enum with the provided [value].
  const FamilyFamilyInviteRequestTypeEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const adult = FamilyFamilyInviteRequestTypeEnum._(r'adult');
  static const kid = FamilyFamilyInviteRequestTypeEnum._(r'kid');

  /// List of all possible values in this [enum][FamilyFamilyInviteRequestTypeEnum].
  static const values = <FamilyFamilyInviteRequestTypeEnum>[
    adult,
    kid,
  ];

  static FamilyFamilyInviteRequestTypeEnum? fromJson(dynamic value) => FamilyFamilyInviteRequestTypeEnumTypeTransformer().decode(value);

  static List<FamilyFamilyInviteRequestTypeEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyInviteRequestTypeEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyInviteRequestTypeEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [FamilyFamilyInviteRequestTypeEnum] to String,
/// and [decode] dynamic data back to [FamilyFamilyInviteRequestTypeEnum].
class FamilyFamilyInviteRequestTypeEnumTypeTransformer {
  factory FamilyFamilyInviteRequestTypeEnumTypeTransformer() => _instance ??= const FamilyFamilyInviteRequestTypeEnumTypeTransformer._();

  const FamilyFamilyInviteRequestTypeEnumTypeTransformer._();

  String encode(FamilyFamilyInviteRequestTypeEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a FamilyFamilyInviteRequestTypeEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  FamilyFamilyInviteRequestTypeEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'adult': return FamilyFamilyInviteRequestTypeEnum.adult;
        case r'kid': return FamilyFamilyInviteRequestTypeEnum.kid;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [FamilyFamilyInviteRequestTypeEnumTypeTransformer] instance.
  static FamilyFamilyInviteRequestTypeEnumTypeTransformer? _instance;
}


