//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class LabelLabelModel {
  /// Returns a new [LabelLabelModel] instance.
  LabelLabelModel({
    required this.color,
    required this.createdAt,
    required this.etag,
    required this.id,
    this.key,
    required this.name,
    required this.owner,
    required this.updatedAt,
  });

  /// @Description Hexadecimal color code for visual representation of the label
  String color;

  /// @Description ISO 8601 timestamp indicating when the label was originally created
  DateTime createdAt;

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

  /// @Description Unique identifier for the label (UUID format)
  String id;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? key;

  /// @Description Display name of the label
  String name;

  DtoOwnerModel owner;

  /// @Description ISO 8601 timestamp indicating when the label was last modified
  DateTime updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is LabelLabelModel &&
    other.color == color &&
    other.createdAt == createdAt &&
    other.etag == etag &&
    other.id == id &&
    other.key == key &&
    other.name == name &&
    other.owner == owner &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (color.hashCode) +
    (createdAt.hashCode) +
    (etag.hashCode) +
    (id.hashCode) +
    (key == null ? 0 : key!.hashCode) +
    (name.hashCode) +
    (owner.hashCode) +
    (updatedAt.hashCode);

  @override
  String toString() => 'LabelLabelModel[color=$color, createdAt=$createdAt, etag=$etag, id=$id, key=$key, name=$name, owner=$owner, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'color'] = this.color;
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
      json[r'etag'] = this.etag;
      json[r'id'] = this.id;
    if (this.key != null) {
      json[r'key'] = this.key;
    } else {
      json[r'key'] = null;
    }
      json[r'name'] = this.name;
      json[r'owner'] = this.owner;
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [LabelLabelModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static LabelLabelModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "LabelLabelModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "LabelLabelModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return LabelLabelModel(
        color: mapValueOfType<String>(json, r'color')!,
        createdAt: mapDateTime(json, r'created_at', r'')!,
        etag: mapValueOfType<String>(json, r'etag')!,
        id: mapValueOfType<String>(json, r'id')!,
        key: mapValueOfType<String>(json, r'key'),
        name: mapValueOfType<String>(json, r'name')!,
        owner: DtoOwnerModel.fromJson(json[r'owner'])!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<LabelLabelModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <LabelLabelModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = LabelLabelModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, LabelLabelModel> mapFromJson(dynamic json) {
    final map = <String, LabelLabelModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = LabelLabelModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of LabelLabelModel-objects as value to a dart map
  static Map<String, List<LabelLabelModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<LabelLabelModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = LabelLabelModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'color',
    'created_at',
    'etag',
    'id',
    'name',
    'owner',
    'updated_at',
  };
}

