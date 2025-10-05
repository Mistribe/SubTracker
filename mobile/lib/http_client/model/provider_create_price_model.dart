//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class ProviderCreatePriceModel {
  /// Returns a new [ProviderCreatePriceModel] instance.
  ProviderCreatePriceModel({
    required this.amount,
    this.createdAt,
    required this.currency,
    this.endDate,
    this.id,
    required this.startDate,
  });

  num amount;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? createdAt;

  String currency;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? endDate;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? id;

  DateTime startDate;

  @override
  bool operator ==(Object other) => identical(this, other) || other is ProviderCreatePriceModel &&
    other.amount == amount &&
    other.createdAt == createdAt &&
    other.currency == currency &&
    other.endDate == endDate &&
    other.id == id &&
    other.startDate == startDate;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (amount.hashCode) +
    (createdAt == null ? 0 : createdAt!.hashCode) +
    (currency.hashCode) +
    (endDate == null ? 0 : endDate!.hashCode) +
    (id == null ? 0 : id!.hashCode) +
    (startDate.hashCode);

  @override
  String toString() => 'ProviderCreatePriceModel[amount=$amount, createdAt=$createdAt, currency=$currency, endDate=$endDate, id=$id, startDate=$startDate]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'amount'] = this.amount;
    if (this.createdAt != null) {
      json[r'created_at'] = this.createdAt!.toUtc().toIso8601String();
    } else {
      json[r'created_at'] = null;
    }
      json[r'currency'] = this.currency;
    if (this.endDate != null) {
      json[r'end_date'] = this.endDate!.toUtc().toIso8601String();
    } else {
      json[r'end_date'] = null;
    }
    if (this.id != null) {
      json[r'id'] = this.id;
    } else {
      json[r'id'] = null;
    }
      json[r'start_date'] = this.startDate.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [ProviderCreatePriceModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static ProviderCreatePriceModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "ProviderCreatePriceModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "ProviderCreatePriceModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return ProviderCreatePriceModel(
        amount: num.parse('${json[r'amount']}'),
        createdAt: mapDateTime(json, r'created_at', r''),
        currency: mapValueOfType<String>(json, r'currency')!,
        endDate: mapDateTime(json, r'end_date', r''),
        id: mapValueOfType<String>(json, r'id'),
        startDate: mapDateTime(json, r'start_date', r'')!,
      );
    }
    return null;
  }

  static List<ProviderCreatePriceModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <ProviderCreatePriceModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = ProviderCreatePriceModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, ProviderCreatePriceModel> mapFromJson(dynamic json) {
    final map = <String, ProviderCreatePriceModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = ProviderCreatePriceModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of ProviderCreatePriceModel-objects as value to a dart map
  static Map<String, List<ProviderCreatePriceModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<ProviderCreatePriceModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = ProviderCreatePriceModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'amount',
    'currency',
    'start_date',
  };
}

