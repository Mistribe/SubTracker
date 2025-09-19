//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class FamilyFamilyModel {
  /// Returns a new [FamilyFamilyModel] instance.
  FamilyFamilyModel({
    required this.createdAt,
    required this.etag,
    required this.id,
    required this.isOwner,
    this.members = const [],
    required this.name,
    required this.updatedAt,
  });

  /// @Description ISO 8601 timestamp indicating when the family was originally created
  DateTime createdAt;

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

  /// @Description Unique identifier for the family (UUID format)
  String id;

  /// @Description Indicates whether the current authenticated user is the owner of this family
  bool isOwner;

  /// @Description Complete list of all members belonging to this family
  List<FamilyFamilyMemberModel> members;

  /// @Description Display name of the family
  String name;

  /// @Description ISO 8601 timestamp indicating when the family information was last modified
  DateTime updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is FamilyFamilyModel &&
    other.createdAt == createdAt &&
    other.etag == etag &&
    other.id == id &&
    other.isOwner == isOwner &&
    _deepEquality.equals(other.members, members) &&
    other.name == name &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt.hashCode) +
    (etag.hashCode) +
    (id.hashCode) +
    (isOwner.hashCode) +
    (members.hashCode) +
    (name.hashCode) +
    (updatedAt.hashCode);

  @override
  String toString() => 'FamilyFamilyModel[createdAt=$createdAt, etag=$etag, id=$id, isOwner=$isOwner, members=$members, name=$name, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
      json[r'etag'] = this.etag;
      json[r'id'] = this.id;
      json[r'is_owner'] = this.isOwner;
      json[r'members'] = this.members;
      json[r'name'] = this.name;
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [FamilyFamilyModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static FamilyFamilyModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "FamilyFamilyModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "FamilyFamilyModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return FamilyFamilyModel(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        etag: mapValueOfType<String>(json, r'etag')!,
        id: mapValueOfType<String>(json, r'id')!,
        isOwner: mapValueOfType<bool>(json, r'is_owner')!,
        members: FamilyFamilyMemberModel.listFromJson(json[r'members']),
        name: mapValueOfType<String>(json, r'name')!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<FamilyFamilyModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <FamilyFamilyModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = FamilyFamilyModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, FamilyFamilyModel> mapFromJson(dynamic json) {
    final map = <String, FamilyFamilyModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = FamilyFamilyModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of FamilyFamilyModel-objects as value to a dart map
  static Map<String, List<FamilyFamilyModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<FamilyFamilyModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = FamilyFamilyModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'created_at',
    'etag',
    'id',
    'is_owner',
    'members',
    'name',
    'updated_at',
  };
}

