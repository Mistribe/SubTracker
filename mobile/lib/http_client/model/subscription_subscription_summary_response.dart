//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionSubscriptionSummaryResponse {
  /// Returns a new [SubscriptionSubscriptionSummaryResponse] instance.
  SubscriptionSubscriptionSummaryResponse({
    this.active,
    this.topLabels = const [],
    this.topProviders = const [],
    this.totalLastMonth,
    this.totalLastYear,
    this.totalMonthly,
    this.totalYearly,
    this.upcomingRenewals = const [],
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  int? active;

  List<SubscriptionSubscriptionSummaryTopLabelResponse> topLabels;

  List<SubscriptionSubscriptionSummaryTopProviderResponse> topProviders;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? totalLastMonth;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? totalLastYear;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? totalMonthly;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? totalYearly;

  List<SubscriptionSubscriptionSummaryUpcomingRenewalResponse> upcomingRenewals;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionSubscriptionSummaryResponse &&
    other.active == active &&
    _deepEquality.equals(other.topLabels, topLabels) &&
    _deepEquality.equals(other.topProviders, topProviders) &&
    other.totalLastMonth == totalLastMonth &&
    other.totalLastYear == totalLastYear &&
    other.totalMonthly == totalMonthly &&
    other.totalYearly == totalYearly &&
    _deepEquality.equals(other.upcomingRenewals, upcomingRenewals);

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (active == null ? 0 : active!.hashCode) +
    (topLabels.hashCode) +
    (topProviders.hashCode) +
    (totalLastMonth == null ? 0 : totalLastMonth!.hashCode) +
    (totalLastYear == null ? 0 : totalLastYear!.hashCode) +
    (totalMonthly == null ? 0 : totalMonthly!.hashCode) +
    (totalYearly == null ? 0 : totalYearly!.hashCode) +
    (upcomingRenewals.hashCode);

  @override
  String toString() => 'SubscriptionSubscriptionSummaryResponse[active=$active, topLabels=$topLabels, topProviders=$topProviders, totalLastMonth=$totalLastMonth, totalLastYear=$totalLastYear, totalMonthly=$totalMonthly, totalYearly=$totalYearly, upcomingRenewals=$upcomingRenewals]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.active != null) {
      json[r'active'] = this.active;
    } else {
      json[r'active'] = null;
    }
      json[r'top_labels'] = this.topLabels;
      json[r'top_providers'] = this.topProviders;
    if (this.totalLastMonth != null) {
      json[r'total_last_month'] = this.totalLastMonth;
    } else {
      json[r'total_last_month'] = null;
    }
    if (this.totalLastYear != null) {
      json[r'total_last_year'] = this.totalLastYear;
    } else {
      json[r'total_last_year'] = null;
    }
    if (this.totalMonthly != null) {
      json[r'total_monthly'] = this.totalMonthly;
    } else {
      json[r'total_monthly'] = null;
    }
    if (this.totalYearly != null) {
      json[r'total_yearly'] = this.totalYearly;
    } else {
      json[r'total_yearly'] = null;
    }
      json[r'upcoming_renewals'] = this.upcomingRenewals;
    return json;
  }

  /// Returns a new [SubscriptionSubscriptionSummaryResponse] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionSubscriptionSummaryResponse? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionSubscriptionSummaryResponse[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionSubscriptionSummaryResponse[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionSubscriptionSummaryResponse(
        active: mapValueOfType<int>(json, r'active'),
        topLabels: SubscriptionSubscriptionSummaryTopLabelResponse.listFromJson(json[r'top_labels']),
        topProviders: SubscriptionSubscriptionSummaryTopProviderResponse.listFromJson(json[r'top_providers']),
        totalLastMonth: DtoAmountModel.fromJson(json[r'total_last_month']),
        totalLastYear: DtoAmountModel.fromJson(json[r'total_last_year']),
        totalMonthly: DtoAmountModel.fromJson(json[r'total_monthly']),
        totalYearly: DtoAmountModel.fromJson(json[r'total_yearly']),
        upcomingRenewals: SubscriptionSubscriptionSummaryUpcomingRenewalResponse.listFromJson(json[r'upcoming_renewals']),
      );
    }
    return null;
  }

  static List<SubscriptionSubscriptionSummaryResponse> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionSummaryResponse>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionSummaryResponse.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionSubscriptionSummaryResponse> mapFromJson(dynamic json) {
    final map = <String, SubscriptionSubscriptionSummaryResponse>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionSubscriptionSummaryResponse.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionSubscriptionSummaryResponse-objects as value to a dart map
  static Map<String, List<SubscriptionSubscriptionSummaryResponse>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionSubscriptionSummaryResponse>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionSubscriptionSummaryResponse.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
  };
}

