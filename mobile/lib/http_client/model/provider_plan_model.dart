//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ProviderPlanModel {
  /// Returns a new [ProviderPlanModel] instance.
  ProviderPlanModel({
    required this.createdAt,
    this.description,
    required this.etag,
    required this.id,
    required this.name,
    this.prices = const [],
    required this.updatedAt,
  });

  /// @Description ISO 8601 timestamp when the plan was originally created
  DateTime createdAt;

  /// @Description Optional detailed description of the plan features and benefits
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? description;

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

  /// @Description Unique identifier for the plan (UUID format)
  String id;

  /// @Description Display name of the subscription plan
  String name;

  /// @Description List of pricing options available for this plan (different currencies, time periods, etc.)
  List<ProviderPriceModel> prices;

  /// @Description ISO 8601 timestamp when the plan was last modified
  DateTime updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ProviderPlanModel &&
    other.createdAt == createdAt &&
    other.description == description &&
    other.etag == etag &&
    other.id == id &&
    other.name == name &&
    _deepEquality.equals(other.prices, prices) &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt.hashCode) +
    (description == null ? 0 : description!.hashCode) +
    (etag.hashCode) +
    (id.hashCode) +
    (name.hashCode) +
    (prices.hashCode) +
    (updatedAt.hashCode);

  @override
  String toString() => 'ProviderPlanModel[createdAt=$createdAt, description=$description, etag=$etag, id=$id, name=$name, prices=$prices, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    if (this.description != null) {
      json[r'description'] = this.description;
    } else {
      json[r'description'] = null;
    }
      json[r'etag'] = this.etag;
      json[r'id'] = this.id;
      json[r'name'] = this.name;
      json[r'prices'] = this.prices;
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [ProviderPlanModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ProviderPlanModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ProviderPlanModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ProviderPlanModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ProviderPlanModel(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        description: mapValueOfType<String>(json, r'description'),
        etag: mapValueOfType<String>(json, r'etag')!,
        id: mapValueOfType<String>(json, r'id')!,
        name: mapValueOfType<String>(json, r'name')!,
        prices: ProviderPriceModel.listFromJson(json[r'prices']),
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<ProviderPlanModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ProviderPlanModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ProviderPlanModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ProviderPlanModel> mapFromJson(dynamic json) {
    final map = <String, ProviderPlanModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ProviderPlanModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ProviderPlanModel-objects as value to a dart map
  static Map<String, List<ProviderPlanModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ProviderPlanModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ProviderPlanModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'created_at',
    'etag',
    'id',
    'name',
    'prices',
    'updated_at',
  };
}

