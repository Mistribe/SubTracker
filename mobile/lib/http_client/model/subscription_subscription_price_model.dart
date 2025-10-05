//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionSubscriptionPriceModel {
  /// Returns a new [SubscriptionSubscriptionPriceModel] instance.
  SubscriptionSubscriptionPriceModel({
    required this.monthly,
    required this.yearly,
  });

  DtoAmountModel monthly;

  DtoAmountModel yearly;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionSubscriptionPriceModel &&
    other.monthly == monthly &&
    other.yearly == yearly;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (monthly.hashCode) +
    (yearly.hashCode);

  @override
  String toString() => 'SubscriptionSubscriptionPriceModel[monthly=$monthly, yearly=$yearly]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'monthly'] = this.monthly;
      json[r'yearly'] = this.yearly;
    return json;
  }

  /// Returns a new [SubscriptionSubscriptionPriceModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionSubscriptionPriceModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionSubscriptionPriceModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionSubscriptionPriceModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionSubscriptionPriceModel(
        monthly: DtoAmountModel.fromJson(json[r'monthly'])!,
        yearly: DtoAmountModel.fromJson(json[r'yearly'])!,
      );
    }
    return null;
  }

  static List<SubscriptionSubscriptionPriceModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionPriceModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionPriceModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionSubscriptionPriceModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionSubscriptionPriceModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionSubscriptionPriceModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionSubscriptionPriceModel-objects as value to a dart map
  static Map<String, List<SubscriptionSubscriptionPriceModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionSubscriptionPriceModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionSubscriptionPriceModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'monthly',
    'yearly',
  };
}

