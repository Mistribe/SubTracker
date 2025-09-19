//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyFamilySeeInvitationResponse {
  /// Returns a new [FamilyFamilySeeInvitationResponse] instance.
  FamilyFamilySeeInvitationResponse({
    this.family,
    this.invitedInasmuchAs,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  FamilyFamilyModel? family;

  /// Role of the invited member
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? invitedInasmuchAs;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyFamilySeeInvitationResponse &&
    other.family == family &&
    other.invitedInasmuchAs == invitedInasmuchAs;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (family == null ? 0 : family!.hashCode) +
    (invitedInasmuchAs == null ? 0 : invitedInasmuchAs!.hashCode);

  @override
  String toString() => 'FamilyFamilySeeInvitationResponse[family=$family, invitedInasmuchAs=$invitedInasmuchAs]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.family != null) {
      json[r'family'] = this.family;
    } else {
      json[r'family'] = null;
    }
    if (this.invitedInasmuchAs != null) {
      json[r'invited_inasmuch_as'] = this.invitedInasmuchAs;
    } else {
      json[r'invited_inasmuch_as'] = null;
    }
    return json;
  }

  /// Returns a new [FamilyFamilySeeInvitationResponse] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyFamilySeeInvitationResponse? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyFamilySeeInvitationResponse[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyFamilySeeInvitationResponse[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyFamilySeeInvitationResponse(
        family: FamilyFamilyModel.fromJson(json[r'family']),
        invitedInasmuchAs: mapValueOfType<String>(json, r'invited_inasmuch_as'),
      );
    }
    return null;
  }

  static List<FamilyFamilySeeInvitationResponse> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilySeeInvitationResponse>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilySeeInvitationResponse.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyFamilySeeInvitationResponse> mapFromJson(dynamic json) {
    final map = <String, FamilyFamilySeeInvitationResponse>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyFamilySeeInvitationResponse.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyFamilySeeInvitationResponse-objects as value to a dart map
  static Map<String, List<FamilyFamilySeeInvitationResponse>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyFamilySeeInvitationResponse>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyFamilySeeInvitationResponse.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
  };
}

