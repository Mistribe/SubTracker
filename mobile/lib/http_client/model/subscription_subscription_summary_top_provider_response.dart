//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionSubscriptionSummaryTopProviderResponse {
  /// Returns a new [SubscriptionSubscriptionSummaryTopProviderResponse] instance.
  SubscriptionSubscriptionSummaryTopProviderResponse({
    this.duration,
    required this.providerId,
    this.total,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? duration;

  String providerId;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? total;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionSubscriptionSummaryTopProviderResponse &&
    other.duration == duration &&
    other.providerId == providerId &&
    other.total == total;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (duration == null ? 0 : duration!.hashCode) +
    (providerId.hashCode) +
    (total == null ? 0 : total!.hashCode);

  @override
  String toString() => 'SubscriptionSubscriptionSummaryTopProviderResponse[duration=$duration, providerId=$providerId, total=$total]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (this.duration != null) {
      json[r'duration'] = this.duration;
    } else {
      json[r'duration'] = null;
    }
      json[r'provider_id'] = this.providerId;
    if (this.total != null) {
      json[r'total'] = this.total;
    } else {
      json[r'total'] = null;
    }
    return json;
  }

  /// Returns a new [SubscriptionSubscriptionSummaryTopProviderResponse] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionSubscriptionSummaryTopProviderResponse? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionSubscriptionSummaryTopProviderResponse[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionSubscriptionSummaryTopProviderResponse[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionSubscriptionSummaryTopProviderResponse(
        duration: mapValueOfType<String>(json, r'duration'),
        providerId: mapValueOfType<String>(json, r'provider_id')!,
        total: DtoAmountModel.fromJson(json[r'total']),
      );
    }
    return null;
  }

  static List<SubscriptionSubscriptionSummaryTopProviderResponse> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionSummaryTopProviderResponse>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionSummaryTopProviderResponse.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionSubscriptionSummaryTopProviderResponse> mapFromJson(dynamic json) {
    final map = <String, SubscriptionSubscriptionSummaryTopProviderResponse>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionSubscriptionSummaryTopProviderResponse.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionSubscriptionSummaryTopProviderResponse-objects as value to a dart map
  static Map<String, List<SubscriptionSubscriptionSummaryTopProviderResponse>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionSubscriptionSummaryTopProviderResponse>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionSubscriptionSummaryTopProviderResponse.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'provider_id',
  };
}

