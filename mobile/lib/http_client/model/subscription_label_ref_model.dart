//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionLabelRefModel {
  /// Returns a new [SubscriptionLabelRefModel] instance.
  SubscriptionLabelRefModel({
    required this.labelId,
    required this.source_,
  });

  String labelId;

  SubscriptionLabelRefModelSource_Enum source_;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionLabelRefModel &&
    other.labelId == labelId &&
    other.source_ == source_;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (labelId.hashCode) +
    (source_.hashCode);

  @override
  String toString() => 'SubscriptionLabelRefModel[labelId=$labelId, source_=$source_]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'label_id'] = this.labelId;
      json[r'source'] = this.source_;
    return json;
  }

  /// Returns a new [SubscriptionLabelRefModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionLabelRefModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionLabelRefModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionLabelRefModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionLabelRefModel(
        labelId: mapValueOfType<String>(json, r'label_id')!,
        source_: SubscriptionLabelRefModelSource_Enum.fromJson(json[r'source'])!,
      );
    }
    return null;
  }

  static List<SubscriptionLabelRefModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionLabelRefModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionLabelRefModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionLabelRefModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionLabelRefModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionLabelRefModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionLabelRefModel-objects as value to a dart map
  static Map<String, List<SubscriptionLabelRefModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionLabelRefModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionLabelRefModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'label_id',
    'source',
  };
}


class SubscriptionLabelRefModelSource_Enum {
  /// Instantiate a new enum with the provided [value].
  const SubscriptionLabelRefModelSource_Enum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const subscription = SubscriptionLabelRefModelSource_Enum._(r'subscription');
  static const provider = SubscriptionLabelRefModelSource_Enum._(r'provider');

  /// List of all possible values in this [enum][SubscriptionLabelRefModelSource_Enum].
  static const values = <SubscriptionLabelRefModelSource_Enum>[
    subscription,
    provider,
  ];

  static SubscriptionLabelRefModelSource_Enum? fromJson(dynamic value) => SubscriptionLabelRefModelSource_EnumTypeTransformer().decode(value);

  static List<SubscriptionLabelRefModelSource_Enum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionLabelRefModelSource_Enum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionLabelRefModelSource_Enum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [SubscriptionLabelRefModelSource_Enum] to String,
/// and [decode] dynamic data back to [SubscriptionLabelRefModelSource_Enum].
class SubscriptionLabelRefModelSource_EnumTypeTransformer {
  factory SubscriptionLabelRefModelSource_EnumTypeTransformer() => _instance ??= const SubscriptionLabelRefModelSource_EnumTypeTransformer._();

  const SubscriptionLabelRefModelSource_EnumTypeTransformer._();

  String encode(SubscriptionLabelRefModelSource_Enum data) => data.value;

  /// Decodes a [dynamic value][data] to a SubscriptionLabelRefModelSource_Enum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  SubscriptionLabelRefModelSource_Enum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'subscription': return SubscriptionLabelRefModelSource_Enum.subscription;
        case r'provider': return SubscriptionLabelRefModelSource_Enum.provider;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [SubscriptionLabelRefModelSource_EnumTypeTransformer] instance.
  static SubscriptionLabelRefModelSource_EnumTypeTransformer? _instance;
}


