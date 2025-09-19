//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyFamilyAcceptInvitationRequest {
  /// Returns a new [FamilyFamilyAcceptInvitationRequest] instance.
  FamilyFamilyAcceptInvitationRequest({
    required this.familyMemberId,
    required this.invitationCode,
  });

  /// ID of the family member accepting the invitation
  String familyMemberId;

  /// Code received in the invitation
  String invitationCode;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyFamilyAcceptInvitationRequest &&
    other.familyMemberId == familyMemberId &&
    other.invitationCode == invitationCode;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (familyMemberId.hashCode) +
    (invitationCode.hashCode);

  @override
  String toString() => 'FamilyFamilyAcceptInvitationRequest[familyMemberId=$familyMemberId, invitationCode=$invitationCode]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'family_member_id'] = this.familyMemberId;
      json[r'invitation_code'] = this.invitationCode;
    return json;
  }

  /// Returns a new [FamilyFamilyAcceptInvitationRequest] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyFamilyAcceptInvitationRequest? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyFamilyAcceptInvitationRequest[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyFamilyAcceptInvitationRequest[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyFamilyAcceptInvitationRequest(
        familyMemberId: mapValueOfType<String>(json, r'family_member_id')!,
        invitationCode: mapValueOfType<String>(json, r'invitation_code')!,
      );
    }
    return null;
  }

  static List<FamilyFamilyAcceptInvitationRequest> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyAcceptInvitationRequest>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyAcceptInvitationRequest.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyFamilyAcceptInvitationRequest> mapFromJson(dynamic json) {
    final map = <String, FamilyFamilyAcceptInvitationRequest>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyFamilyAcceptInvitationRequest.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyFamilyAcceptInvitationRequest-objects as value to a dart map
  static Map<String, List<FamilyFamilyAcceptInvitationRequest>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyFamilyAcceptInvitationRequest>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyFamilyAcceptInvitationRequest.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'family_member_id',
    'invitation_code',
  };
}

