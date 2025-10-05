//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ProviderCreateProviderModel {
  /// Returns a new [ProviderCreateProviderModel] instance.
  ProviderCreateProviderModel({
    this.createdAt,
    this.description,
    this.iconUrl,
    this.id,
    this.labels = const [],
    required this.name,
    required this.owner,
    this.pricingPageUrl,
    this.url,
  });

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
  String? description;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? iconUrl;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? id;

  List<String> labels;

  String name;

  DtoEditableOwnerModel owner;

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
  String? url;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ProviderCreateProviderModel &&
    other.createdAt == createdAt &&
    other.description == description &&
    other.iconUrl == iconUrl &&
    other.id == id &&
    _deepEquality.equals(other.labels, labels) &&
    other.name == name &&
    other.owner == owner &&
    other.pricingPageUrl == pricingPageUrl &&
    other.url == url;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt == null ? 0 : createdAt!.hashCode) +
    (description == null ? 0 : description!.hashCode) +
    (iconUrl == null ? 0 : iconUrl!.hashCode) +
    (id == null ? 0 : id!.hashCode) +
    (labels.hashCode) +
    (name.hashCode) +
    (owner.hashCode) +
    (pricingPageUrl == null ? 0 : pricingPageUrl!.hashCode) +
    (url == null ? 0 : url!.hashCode);

  @override
  String toString() => 'ProviderCreateProviderModel[createdAt=$createdAt, description=$description, iconUrl=$iconUrl, id=$id, labels=$labels, name=$name, owner=$owner, pricingPageUrl=$pricingPageUrl, url=$url]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.createdAt != null) {
      json[r'created_at'] = this.createdAt!.toUtc().toIso8601String();
    } else {
      json[r'created_at'] = null;
    }
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
    if (this.id != null) {
      json[r'id'] = this.id;
    } else {
      json[r'id'] = null;
    }
      json[r'labels'] = this.labels;
      json[r'name'] = this.name;
      json[r'owner'] = this.owner;
    if (this.pricingPageUrl != null) {
      json[r'pricing_page_url'] = this.pricingPageUrl;
    } else {
      json[r'pricing_page_url'] = null;
    }
    if (this.url != null) {
      json[r'url'] = this.url;
    } else {
      json[r'url'] = null;
    }
    return json;
  }

  /// Returns a new [ProviderCreateProviderModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ProviderCreateProviderModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ProviderCreateProviderModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ProviderCreateProviderModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ProviderCreateProviderModel(
        createdAt: mapDateTime(json, r'created_at', r''),
        description: mapValueOfType<String>(json, r'description'),
        iconUrl: mapValueOfType<String>(json, r'icon_url'),
        id: mapValueOfType<String>(json, r'id'),
        labels: json[r'labels'] is Iterable
            ? (json[r'labels'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        name: mapValueOfType<String>(json, r'name')!,
        owner: DtoEditableOwnerModel.fromJson(json[r'owner'])!,
        pricingPageUrl: mapValueOfType<String>(json, r'pricing_page_url'),
        url: mapValueOfType<String>(json, r'url'),
      );
    }
    return null;
  }

  static List<ProviderCreateProviderModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ProviderCreateProviderModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ProviderCreateProviderModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ProviderCreateProviderModel> mapFromJson(dynamic json) {
    final map = <String, ProviderCreateProviderModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ProviderCreateProviderModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ProviderCreateProviderModel-objects as value to a dart map
  static Map<String, List<ProviderCreateProviderModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ProviderCreateProviderModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ProviderCreateProviderModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'labels',
    'name',
    'owner',
  };
}

