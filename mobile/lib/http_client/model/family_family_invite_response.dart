//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyFamilyInviteResponse {
  /// Returns a new [FamilyFamilyInviteResponse] instance.
  FamilyFamilyInviteResponse({
    required this.code,
    required this.familyId,
    required this.familyMemberId,
  });

  String code;

  String familyId;

  String familyMemberId;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyFamilyInviteResponse &&
    other.code == code &&
    other.familyId == familyId &&
    other.familyMemberId == familyMemberId;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (code.hashCode) +
    (familyId.hashCode) +
    (familyMemberId.hashCode);

  @override
  String toString() => 'FamilyFamilyInviteResponse[code=$code, familyId=$familyId, familyMemberId=$familyMemberId]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'code'] = this.code;
      json[r'family_id'] = this.familyId;
      json[r'family_member_id'] = this.familyMemberId;
    return json;
  }

  /// Returns a new [FamilyFamilyInviteResponse] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyFamilyInviteResponse? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyFamilyInviteResponse[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyFamilyInviteResponse[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyFamilyInviteResponse(
        code: mapValueOfType<String>(json, r'code')!,
        familyId: mapValueOfType<String>(json, r'family_id')!,
        familyMemberId: mapValueOfType<String>(json, r'family_member_id')!,
      );
    }
    return null;
  }

  static List<FamilyFamilyInviteResponse> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyInviteResponse>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyInviteResponse.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyFamilyInviteResponse> mapFromJson(dynamic json) {
    final map = <String, FamilyFamilyInviteResponse>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyFamilyInviteResponse.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyFamilyInviteResponse-objects as value to a dart map
  static Map<String, List<FamilyFamilyInviteResponse>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyFamilyInviteResponse>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyFamilyInviteResponse.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'code',
    'family_id',
    'family_member_id',
  };
}

