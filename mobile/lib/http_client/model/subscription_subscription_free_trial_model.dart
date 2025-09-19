//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionSubscriptionFreeTrialModel {
  /// Returns a new [SubscriptionSubscriptionFreeTrialModel] instance.
  SubscriptionSubscriptionFreeTrialModel({
    this.endDate,
    this.startDate,
  });

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
  DateTime? startDate;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionSubscriptionFreeTrialModel &&
    other.endDate == endDate &&
    other.startDate == startDate;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (endDate == null ? 0 : endDate!.hashCode) +
    (startDate == null ? 0 : startDate!.hashCode);

  @override
  String toString() => 'SubscriptionSubscriptionFreeTrialModel[endDate=$endDate, startDate=$startDate]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.endDate != null) {
      json[r'end_date'] = this.endDate!.toUtc().toIso8601String();
    } else {
      json[r'end_date'] = null;
    }
    if (this.startDate != null) {
      json[r'start_date'] = this.startDate!.toUtc().toIso8601String();
    } else {
      json[r'start_date'] = null;
    }
    return json;
  }

  /// Returns a new [SubscriptionSubscriptionFreeTrialModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionSubscriptionFreeTrialModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionSubscriptionFreeTrialModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionSubscriptionFreeTrialModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionSubscriptionFreeTrialModel(
        endDate: mapDateTime(json, r'end_date', r''),
        startDate: mapDateTime(json, r'start_date', r''),
      );
    }
    return null;
  }

  static List<SubscriptionSubscriptionFreeTrialModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionFreeTrialModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionFreeTrialModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionSubscriptionFreeTrialModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionSubscriptionFreeTrialModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionSubscriptionFreeTrialModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionSubscriptionFreeTrialModel-objects as value to a dart map
  static Map<String, List<SubscriptionSubscriptionFreeTrialModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionSubscriptionFreeTrialModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionSubscriptionFreeTrialModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
  };
}

