//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionSubscriptionModel {
  /// Returns a new [SubscriptionSubscriptionModel] instance.
  SubscriptionSubscriptionModel({
    required this.createdAt,
    this.customPrice,
    this.customRecurrency,
    this.endDate,
    required this.etag,
    this.freeTrial,
    this.friendlyName,
    required this.id,
    required this.isActive,
    this.labelRefs = const [],
    required this.owner,
    this.payer,
    this.planId,
    this.price,
    this.priceId,
    required this.providerId,
    required this.recurrency,
    this.serviceUsers = const [],
    required this.startDate,
    required this.updatedAt,
  });

  /// @Description ISO 8601 timestamp when the subscription was originally created
  DateTime createdAt;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? customPrice;

  /// @Description CustomRecurrency recurrency interval in days (required when recurrency is custom)
  ///
  /// Minimum value: 1
  /// Maximum value: 3650
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  int? customRecurrency;

  /// @Description ISO 8601 timestamp when the subscription expires (null for ongoing subscriptions)
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? endDate;

  /// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
  String etag;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  SubscriptionSubscriptionFreeTrialModel? freeTrial;

  /// @Description Optional custom name for easy identification of the subscription
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? friendlyName;

  /// @Description Unique identifier for the subscription (UUID format)
  String id;

  /// @Description Indicates whether the subscription is currently active or not
  bool isActive;

  /// @Description List of labels associated with this subscription
  List<SubscriptionLabelRefModel> labelRefs;

  DtoOwnerModel owner;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  SubscriptionSubscriptionPayerModel? payer;

  /// @Description ID of the specific plan being subscribed to
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? planId;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  SubscriptionSubscriptionPriceModel? price;

  /// @Description ID of the pricing tier for this subscription
  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? priceId;

  /// @Description ID of the service provider offering this subscription
  String providerId;

  /// @Description Billing recurrency pattern (monthly, yearly, custom, etc.)
  SubscriptionSubscriptionModelRecurrencyEnum recurrency;

  /// @Description List of family member IDs who use this service (for shared subscriptions)
  List<String> serviceUsers;

  /// @Description ISO 8601 timestamp when the subscription becomes active
  DateTime startDate;

  /// @Description ISO 8601 timestamp when the subscription was last modified
  DateTime updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionSubscriptionModel &&
    other.createdAt == createdAt &&
    other.customPrice == customPrice &&
    other.customRecurrency == customRecurrency &&
    other.endDate == endDate &&
    other.etag == etag &&
    other.freeTrial == freeTrial &&
    other.friendlyName == friendlyName &&
    other.id == id &&
    other.isActive == isActive &&
    _deepEquality.equals(other.labelRefs, labelRefs) &&
    other.owner == owner &&
    other.payer == payer &&
    other.planId == planId &&
    other.price == price &&
    other.priceId == priceId &&
    other.providerId == providerId &&
    other.recurrency == recurrency &&
    _deepEquality.equals(other.serviceUsers, serviceUsers) &&
    other.startDate == startDate &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (createdAt.hashCode) +
    (customPrice == null ? 0 : customPrice!.hashCode) +
    (customRecurrency == null ? 0 : customRecurrency!.hashCode) +
    (endDate == null ? 0 : endDate!.hashCode) +
    (etag.hashCode) +
    (freeTrial == null ? 0 : freeTrial!.hashCode) +
    (friendlyName == null ? 0 : friendlyName!.hashCode) +
    (id.hashCode) +
    (isActive.hashCode) +
    (labelRefs.hashCode) +
    (owner.hashCode) +
    (payer == null ? 0 : payer!.hashCode) +
    (planId == null ? 0 : planId!.hashCode) +
    (price == null ? 0 : price!.hashCode) +
    (priceId == null ? 0 : priceId!.hashCode) +
    (providerId.hashCode) +
    (recurrency.hashCode) +
    (serviceUsers.hashCode) +
    (startDate.hashCode) +
    (updatedAt.hashCode);

  @override
  String toString() => 'SubscriptionSubscriptionModel[createdAt=$createdAt, customPrice=$customPrice, customRecurrency=$customRecurrency, endDate=$endDate, etag=$etag, freeTrial=$freeTrial, friendlyName=$friendlyName, id=$id, isActive=$isActive, labelRefs=$labelRefs, owner=$owner, payer=$payer, planId=$planId, price=$price, priceId=$priceId, providerId=$providerId, recurrency=$recurrency, serviceUsers=$serviceUsers, startDate=$startDate, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
      json[r'created_at'] = this.createdAt.toUtc().toIso8601String();
    if (this.customPrice != null) {
      json[r'custom_price'] = this.customPrice;
    } else {
      json[r'custom_price'] = null;
    }
    if (this.customRecurrency != null) {
      json[r'custom_recurrency'] = this.customRecurrency;
    } else {
      json[r'custom_recurrency'] = null;
    }
    if (this.endDate != null) {
      json[r'end_date'] = this.endDate!.toUtc().toIso8601String();
    } else {
      json[r'end_date'] = null;
    }
      json[r'etag'] = this.etag;
    if (this.freeTrial != null) {
      json[r'free_trial'] = this.freeTrial;
    } else {
      json[r'free_trial'] = null;
    }
    if (this.friendlyName != null) {
      json[r'friendly_name'] = this.friendlyName;
    } else {
      json[r'friendly_name'] = null;
    }
      json[r'id'] = this.id;
      json[r'is_active'] = this.isActive;
      json[r'label_refs'] = this.labelRefs;
      json[r'owner'] = this.owner;
    if (this.payer != null) {
      json[r'payer'] = this.payer;
    } else {
      json[r'payer'] = null;
    }
    if (this.planId != null) {
      json[r'plan_id'] = this.planId;
    } else {
      json[r'plan_id'] = null;
    }
    if (this.price != null) {
      json[r'price'] = this.price;
    } else {
      json[r'price'] = null;
    }
    if (this.priceId != null) {
      json[r'price_id'] = this.priceId;
    } else {
      json[r'price_id'] = null;
    }
      json[r'provider_id'] = this.providerId;
      json[r'recurrency'] = this.recurrency;
      json[r'service_users'] = this.serviceUsers;
      json[r'start_date'] = this.startDate.toUtc().toIso8601String();
      json[r'updated_at'] = this.updatedAt.toUtc().toIso8601String();
    return json;
  }

  /// Returns a new [SubscriptionSubscriptionModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionSubscriptionModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionSubscriptionModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionSubscriptionModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionSubscriptionModel(
        createdAt: mapDateTime(json, r'created_at', r'')!,
        customPrice: DtoAmountModel.fromJson(json[r'custom_price']),
        customRecurrency: mapValueOfType<int>(json, r'custom_recurrency'),
        endDate: mapDateTime(json, r'end_date', r''),
        etag: mapValueOfType<String>(json, r'etag')!,
        freeTrial: SubscriptionSubscriptionFreeTrialModel.fromJson(json[r'free_trial']),
        friendlyName: mapValueOfType<String>(json, r'friendly_name'),
        id: mapValueOfType<String>(json, r'id')!,
        isActive: mapValueOfType<bool>(json, r'is_active')!,
        labelRefs: SubscriptionLabelRefModel.listFromJson(json[r'label_refs']),
        owner: DtoOwnerModel.fromJson(json[r'owner'])!,
        payer: SubscriptionSubscriptionPayerModel.fromJson(json[r'payer']),
        planId: mapValueOfType<String>(json, r'plan_id'),
        price: SubscriptionSubscriptionPriceModel.fromJson(json[r'price']),
        priceId: mapValueOfType<String>(json, r'price_id'),
        providerId: mapValueOfType<String>(json, r'provider_id')!,
        recurrency: SubscriptionSubscriptionModelRecurrencyEnum.fromJson(json[r'recurrency'])!,
        serviceUsers: json[r'service_users'] is Iterable
            ? (json[r'service_users'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        startDate: mapDateTime(json, r'start_date', r'')!,
        updatedAt: mapDateTime(json, r'updated_at', r'')!,
      );
    }
    return null;
  }

  static List<SubscriptionSubscriptionModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionSubscriptionModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionSubscriptionModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionSubscriptionModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionSubscriptionModel-objects as value to a dart map
  static Map<String, List<SubscriptionSubscriptionModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionSubscriptionModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionSubscriptionModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'created_at',
    'etag',
    'id',
    'is_active',
    'owner',
    'provider_id',
    'recurrency',
    'start_date',
    'updated_at',
  };
}

/// @Description Billing recurrency pattern (monthly, yearly, custom, etc.)
class SubscriptionSubscriptionModelRecurrencyEnum {
  /// Instantiate a new enum with the provided [value].
  const SubscriptionSubscriptionModelRecurrencyEnum._(this.value);

  /// The underlying value of this enum member.
  final String value;

  @override
  String toString() => value;

  String toJson() => value;

  static const unknown = SubscriptionSubscriptionModelRecurrencyEnum._(r'unknown');
  static const oneTime = SubscriptionSubscriptionModelRecurrencyEnum._(r'one_time');
  static const monthly = SubscriptionSubscriptionModelRecurrencyEnum._(r'monthly');
  static const quarterly = SubscriptionSubscriptionModelRecurrencyEnum._(r'quarterly');
  static const halfYearly = SubscriptionSubscriptionModelRecurrencyEnum._(r'half_yearly');
  static const yearly = SubscriptionSubscriptionModelRecurrencyEnum._(r'yearly');
  static const custom = SubscriptionSubscriptionModelRecurrencyEnum._(r'custom');

  /// List of all possible values in this [enum][SubscriptionSubscriptionModelRecurrencyEnum].
  static const values = <SubscriptionSubscriptionModelRecurrencyEnum>[
    unknown,
    oneTime,
    monthly,
    quarterly,
    halfYearly,
    yearly,
    custom,
  ];

  static SubscriptionSubscriptionModelRecurrencyEnum? fromJson(dynamic value) => SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer().decode(value);

  static List<SubscriptionSubscriptionModelRecurrencyEnum> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionSubscriptionModelRecurrencyEnum>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionSubscriptionModelRecurrencyEnum.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }
}

/// Transformation class that can [encode] an instance of [SubscriptionSubscriptionModelRecurrencyEnum] to String,
/// and [decode] dynamic data back to [SubscriptionSubscriptionModelRecurrencyEnum].
class SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer {
  factory SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer() => _instance ??= const SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer._();

  const SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer._();

  String encode(SubscriptionSubscriptionModelRecurrencyEnum data) => data.value;

  /// Decodes a [dynamic value][data] to a SubscriptionSubscriptionModelRecurrencyEnum.
  ///
  /// If [allowNull] is true and the [dynamic value][data] cannot be decoded successfully,
  /// then null is returned. However, if [allowNull] is false and the [dynamic value][data]
  /// cannot be decoded successfully, then an [UnimplementedError] is thrown.
  ///
  /// The [allowNull] is very handy when an API changes and a new enum value is added or removed,
  /// and users are still using an old app with the old code.
  SubscriptionSubscriptionModelRecurrencyEnum? decode(dynamic data, {bool allowNull = true}) {
    if (data != null) {
      switch (data) {
        case r'unknown': return SubscriptionSubscriptionModelRecurrencyEnum.unknown;
        case r'one_time': return SubscriptionSubscriptionModelRecurrencyEnum.oneTime;
        case r'monthly': return SubscriptionSubscriptionModelRecurrencyEnum.monthly;
        case r'quarterly': return SubscriptionSubscriptionModelRecurrencyEnum.quarterly;
        case r'half_yearly': return SubscriptionSubscriptionModelRecurrencyEnum.halfYearly;
        case r'yearly': return SubscriptionSubscriptionModelRecurrencyEnum.yearly;
        case r'custom': return SubscriptionSubscriptionModelRecurrencyEnum.custom;
        default:
          if (!allowNull) {
            throw ArgumentError('Unknown enum value to decode: $data');
          }
      }
    }
    return null;
  }

  /// Singleton [SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer] instance.
  static SubscriptionSubscriptionModelRecurrencyEnumTypeTransformer? _instance;
}


