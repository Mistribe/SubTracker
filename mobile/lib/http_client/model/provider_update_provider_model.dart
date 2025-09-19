//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ProviderUpdateProviderModel {
  /// Returns a new [ProviderUpdateProviderModel] instance.
  ProviderUpdateProviderModel({
    this.description,
    this.iconUrl,
    this.labels = const [],
    required this.name,
    this.pricingPageUrl,
    this.updatedAt,
    this.url,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? description;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? iconUrl;

  List<String> labels;

  String name;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? pricingPageUrl;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? updatedAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? url;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ProviderUpdateProviderModel &&
    other.description == description &&
    other.iconUrl == iconUrl &&
    _deepEquality.equals(other.labels, labels) &&
    other.name == name &&
    other.pricingPageUrl == pricingPageUrl &&
    other.updatedAt == updatedAt &&
    other.url == url;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (description == null ? 0 : description!.hashCode) +
    (iconUrl == null ? 0 : iconUrl!.hashCode) +
    (labels.hashCode) +
    (name.hashCode) +
    (pricingPageUrl == null ? 0 : pricingPageUrl!.hashCode) +
    (updatedAt == null ? 0 : updatedAt!.hashCode) +
    (url == null ? 0 : url!.hashCode);

  @override
  String toString() => 'ProviderUpdateProviderModel[description=$description, iconUrl=$iconUrl, labels=$labels, name=$name, pricingPageUrl=$pricingPageUrl, updatedAt=$updatedAt, url=$url]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.description != null) {
      json[r'description'] = this.description;
    } else {
      json[r'description'] = null;
    }
    if (this.iconUrl != null) {
      json[r'icon_url'] = this.iconUrl;
    } else {
      json[r'icon_url'] = null;
    }
      json[r'labels'] = this.labels;
      json[r'name'] = this.name;
    if (this.pricingPageUrl != null) {
      json[r'pricing_page_url'] = this.pricingPageUrl;
    } else {
      json[r'pricing_page_url'] = null;
    }
    if (this.updatedAt != null) {
      json[r'updated_at'] = this.updatedAt!.toUtc().toIso8601String();
    } else {
      json[r'updated_at'] = null;
    }
    if (this.url != null) {
      json[r'url'] = this.url;
    } else {
      json[r'url'] = null;
    }
    return json;
  }

  /// Returns a new [ProviderUpdateProviderModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ProviderUpdateProviderModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ProviderUpdateProviderModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ProviderUpdateProviderModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ProviderUpdateProviderModel(
        description: mapValueOfType<String>(json, r'description'),
        iconUrl: mapValueOfType<String>(json, r'icon_url'),
        labels: json[r'labels'] is Iterable
            ? (json[r'labels'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        name: mapValueOfType<String>(json, r'name')!,
        pricingPageUrl: mapValueOfType<String>(json, r'pricing_page_url'),
        updatedAt: mapDateTime(json, r'updated_at', r''),
        url: mapValueOfType<String>(json, r'url'),
      );
    }
    return null;
  }

  static List<ProviderUpdateProviderModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ProviderUpdateProviderModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ProviderUpdateProviderModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ProviderUpdateProviderModel> mapFromJson(dynamic json) {
    final map = <String, ProviderUpdateProviderModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ProviderUpdateProviderModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ProviderUpdateProviderModel-objects as value to a dart map
  static Map<String, List<ProviderUpdateProviderModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ProviderUpdateProviderModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ProviderUpdateProviderModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'labels',
    'name',
  };
}

