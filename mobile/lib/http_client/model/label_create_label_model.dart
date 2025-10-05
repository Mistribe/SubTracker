//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class LabelCreateLabelModel {
  /// Returns a new [LabelCreateLabelModel] instance.
  LabelCreateLabelModel({
    required this.color,
    this.createdAt,
    this.id,
    required this.name,
    required this.owner,
  });

  String color;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? createdAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? id;

  String name;

  DtoEditableOwnerModel owner;

  @override
  bool operator ==(Object other) => identical(this, other) || other is LabelCreateLabelModel &&
    other.color == color &&
    other.createdAt == createdAt &&
    other.id == id &&
    other.name == name &&
    other.owner == owner;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (color.hashCode) +
    (createdAt == null ? 0 : createdAt!.hashCode) +
    (id == null ? 0 : id!.hashCode) +
    (name.hashCode) +
    (owner.hashCode);

  @override
  String toString() => 'LabelCreateLabelModel[color=$color, createdAt=$createdAt, id=$id, name=$name, owner=$owner]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'color'] = this.color;
    if (this.createdAt != null) {
      json[r'created_at'] = this.createdAt!.toUtc().toIso8601String();
    } else {
      json[r'created_at'] = null;
    }
    if (this.id != null) {
      json[r'id'] = this.id;
    } else {
      json[r'id'] = null;
    }
      json[r'name'] = this.name;
      json[r'owner'] = this.owner;
    return json;
  }

  /// Returns a new [LabelCreateLabelModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static LabelCreateLabelModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "LabelCreateLabelModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "LabelCreateLabelModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return LabelCreateLabelModel(
        color: mapValueOfType<String>(json, r'color')!,
        createdAt: mapDateTime(json, r'created_at', r''),
        id: mapValueOfType<String>(json, r'id'),
        name: mapValueOfType<String>(json, r'name')!,
        owner: DtoEditableOwnerModel.fromJson(json[r'owner'])!,
      );
    }
    return null;
  }

  static List<LabelCreateLabelModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <LabelCreateLabelModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = LabelCreateLabelModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, LabelCreateLabelModel> mapFromJson(dynamic json) {
    final map = <String, LabelCreateLabelModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = LabelCreateLabelModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of LabelCreateLabelModel-objects as value to a dart map
  static Map<String, List<LabelCreateLabelModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<LabelCreateLabelModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = LabelCreateLabelModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'color',
    'name',
    'owner',
  };
}

