//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ProviderProviderModel {
  /// Returns a new [ProviderProviderModel] instance.
  ProviderProviderModel({
    required this.createdAt,
    this.description,
    required this.etag,
    this.iconUrl,
    required this.id,
    this.key,
    this.labels = const [],
    required this.name,
    required this.owner,
    this.plans = const [],
    this.pricingPageUrl,
    required this.updatedAt,
    this.url,
  });

  /// @Description ISO 8601 timestamp when the provider was originally created
  DateTime createdAt;

  /// @Description Optional detailed description of the provider and their services
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? description;

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

  /// @Description Optional URL to the provider's icon or logo image
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? iconUrl;

  /// @Description Unique identifier for the provider (UUID format)
  String id;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? key;

  /// @Description List of label IDs associated with this provider for categorization
  List<String> labels;

  /// @Description Display name of the service provider
  String name;

  DtoOwnerModel owner;

  /// @Description List of subscription plans offered by this provider
  List<ProviderPlanModel> plans;

  /// @Description Optional URL to the provider's pricing information page
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? pricingPageUrl;

  /// @Description ISO 8601 timestamp when the provider was last modified
  DateTime updatedAt;

  /// @Description Optional URL to the provider's main website
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? url;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ProviderProviderModel &&
    other.createdAt == createdAt &&
    other.description == description &&
    other.etag == etag &&
    other.iconUrl == iconUrl &&
    other.id == id &&
    other.key == key &&
    _deepEquality.equals(other.labels, labels) &&
    other.name == name &&
    other.owner == owner &&
    _deepEquality.equals(other.plans, plans) &&
    other.pricingPageUrl == pricingPageUrl &&
    other.updatedAt == updatedAt &&
    other.url == url;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt.hashCode) +
    (description == null ? 0 : description!.hashCode) +
    (etag.hashCode) +
    (iconUrl == null ? 0 : iconUrl!.hashCode) +
    (id.hashCode) +
    (key == null ? 0 : key!.hashCode) +
    (labels.hashCode) +
    (name.hashCode) +
    (owner.hashCode) +
    (plans.hashCode) +
    (pricingPageUrl == null ? 0 : pricingPageUrl!.hashCode) +
    (updatedAt.hashCode) +
    (url == null ? 0 : url!.hashCode);

  @override
  String toString() => 'ProviderProviderModel[createdAt=$createdAt, description=$description, etag=$etag, iconUrl=$iconUrl, id=$id, key=$key, labels=$labels, name=$name, owner=$owner, plans=$plans, pricingPageUrl=$pricingPageUrl, updatedAt=$updatedAt, url=$url]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    if (this.description != null) {
      json[r'description'] = this.description;
    } else {
      json[r'description'] = null;
    }
      json[r'etag'] = this.etag;
    if (this.iconUrl != null) {
      json[r'icon_url'] = this.iconUrl;
    } else {
      json[r'icon_url'] = null;
    }
      json[r'id'] = this.id;
    if (this.key != null) {
      json[r'key'] = this.key;
    } else {
      json[r'key'] = null;
    }
      json[r'labels'] = this.labels;
      json[r'name'] = this.name;
      json[r'owner'] = this.owner;
      json[r'plans'] = this.plans;
    if (this.pricingPageUrl != null) {
      json[r'pricing_page_url'] = this.pricingPageUrl;
    } else {
      json[r'pricing_page_url'] = null;
    }
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    if (this.url != null) {
      json[r'url'] = this.url;
    } else {
      json[r'url'] = null;
    }
    return json;
  }

  /// Returns a new [ProviderProviderModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ProviderProviderModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ProviderProviderModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ProviderProviderModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ProviderProviderModel(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        description: mapValueOfType<String>(json, r'description'),
        etag: mapValueOfType<String>(json, r'etag')!,
        iconUrl: mapValueOfType<String>(json, r'icon_url'),
        id: mapValueOfType<String>(json, r'id')!,
        key: mapValueOfType<String>(json, r'key'),
        labels: json[r'labels'] is Iterable
            ? (json[r'labels'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        name: mapValueOfType<String>(json, r'name')!,
        owner: DtoOwnerModel.fromJson(json[r'owner'])!,
        plans: ProviderPlanModel.listFromJson(json[r'plans']),
        pricingPageUrl: mapValueOfType<String>(json, r'pricing_page_url'),
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
        url: mapValueOfType<String>(json, r'url'),
      );
    }
    return null;
  }

  static List<ProviderProviderModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ProviderProviderModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ProviderProviderModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ProviderProviderModel> mapFromJson(dynamic json) {
    final map = <String, ProviderProviderModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ProviderProviderModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ProviderProviderModel-objects as value to a dart map
  static Map<String, List<ProviderProviderModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ProviderProviderModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ProviderProviderModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'created_at',
    'etag',
    'id',
    'labels',
    'name',
    'owner',
    'plans',
    'updated_at',
  };
}

