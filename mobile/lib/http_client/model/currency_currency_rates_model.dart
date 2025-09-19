//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class CurrencyCurrencyRatesModel {
  /// Returns a new [CurrencyCurrencyRatesModel] instance.
  CurrencyCurrencyRatesModel({
    this.rates = const [],
    required this.timestamp,
  });

  List<CurrencyCurrencyRateModel> rates;

  DateTime timestamp;

  @override
  bool operator ==(Object other) => identical(this, other) || other is CurrencyCurrencyRatesModel &&
    _deepEquality.equals(other.rates, rates) &&
    other.timestamp == timestamp;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (rates.hashCode) +
    (timestamp.hashCode);

  @override
  String toString() => 'CurrencyCurrencyRatesModel[rates=$rates, timestamp=$timestamp]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'rates'] = this.rates;
      json[r'timestamp'] = this.timestamp.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [CurrencyCurrencyRatesModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static CurrencyCurrencyRatesModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "CurrencyCurrencyRatesModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "CurrencyCurrencyRatesModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return CurrencyCurrencyRatesModel(
        rates: CurrencyCurrencyRateModel.listFromJson(json[r'rates']),
        timestamp: mapDateTime(json, r'timestamp', r'')!,
      );
    }
    return null;
  }

  static List<CurrencyCurrencyRatesModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <CurrencyCurrencyRatesModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = CurrencyCurrencyRatesModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, CurrencyCurrencyRatesModel> mapFromJson(dynamic json) {
    final map = <String, CurrencyCurrencyRatesModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = CurrencyCurrencyRatesModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of CurrencyCurrencyRatesModel-objects as value to a dart map
  static Map<String, List<CurrencyCurrencyRatesModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<CurrencyCurrencyRatesModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = CurrencyCurrencyRatesModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'rates',
    'timestamp',
  };
}

