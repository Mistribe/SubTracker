//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ProviderPriceModel {
  /// Returns a new [ProviderPriceModel] instance.
  ProviderPriceModel({
    required this.amount,
    required this.createdAt,
    required this.currency,
    this.endDate,
    required this.etag,
    required this.id,
    required this.startDate,
    required this.updatedAt,
  });

  /// @Description Price amount in the specified currency (supports decimal values)
  ///
  /// Minimum value: 0
  num amount;

  /// @Description ISO 8601 timestamp when the price was originally created
  DateTime createdAt;

  /// @Description ISO 4217 currency code for the price
  String currency;

  /// @Description ISO 8601 timestamp when this price expires (null means indefinite)
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? endDate;

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

  /// @Description Unique identifier for the price (UUID format)
  String id;

  /// @Description ISO 8601 timestamp when this price becomes effective
  DateTime startDate;

  /// @Description ISO 8601 timestamp when the price was last modified
  DateTime updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ProviderPriceModel &&
    other.amount == amount &&
    other.createdAt == createdAt &&
    other.currency == currency &&
    other.endDate == endDate &&
    other.etag == etag &&
    other.id == id &&
    other.startDate == startDate &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (amount.hashCode) +
    (createdAt.hashCode) +
    (currency.hashCode) +
    (endDate == null ? 0 : endDate!.hashCode) +
    (etag.hashCode) +
    (id.hashCode) +
    (startDate.hashCode) +
    (updatedAt.hashCode);

  @override
  String toString() => 'ProviderPriceModel[amount=$amount, createdAt=$createdAt, currency=$currency, endDate=$endDate, etag=$etag, id=$id, startDate=$startDate, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'amount'] = this.amount;
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
      json[r'currency'] = this.currency;
    if (this.endDate != null) {
      json[r'end_date'] = this.endDate!.toUtc().toIso8601String();
    } else {
      json[r'end_date'] = null;
    }
      json[r'etag'] = this.etag;
      json[r'id'] = this.id;
      json[r'start_date'] = this.startDate.toUtc().toIso8601String();
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [ProviderPriceModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ProviderPriceModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ProviderPriceModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ProviderPriceModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ProviderPriceModel(
        amount: num.parse('${json[r'amount']}'),
        createdAt: mapDateTime(json, r'created_at', r'')!,
        currency: mapValueOfType<String>(json, r'currency')!,
        endDate: mapDateTime(json, r'end_date', r''),
        etag: mapValueOfType<String>(json, r'etag')!,
        id: mapValueOfType<String>(json, r'id')!,
        startDate: mapDateTime(json, r'start_date', r'')!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<ProviderPriceModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ProviderPriceModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ProviderPriceModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ProviderPriceModel> mapFromJson(dynamic json) {
    final map = <String, ProviderPriceModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ProviderPriceModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ProviderPriceModel-objects as value to a dart map
  static Map<String, List<ProviderPriceModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ProviderPriceModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ProviderPriceModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'amount',
    'created_at',
    'currency',
    'etag',
    'id',
    'start_date',
    'updated_at',
  };
}

