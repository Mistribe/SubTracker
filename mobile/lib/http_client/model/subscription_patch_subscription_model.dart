//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;

class SubscriptionPatchSubscriptionModel {
  /// Returns a new [SubscriptionPatchSubscriptionModel] instance.
  SubscriptionPatchSubscriptionModel({
    this.customPrice,
    this.customRecurrency,
    this.endDate,
    this.freeTrial,
    this.friendlyName,
    this.id,
    this.labels = const [],
    required this.owner,
    this.payer,
    this.planId,
    this.priceId,
    required this.providerId,
    required this.recurrency,
    this.serviceUsers = const [],
    required this.startDate,
    this.updatedAt,
  });

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DtoAmountModel? customPrice;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  int? customRecurrency;

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
  SubscriptionSubscriptionFreeTrialModel? freeTrial;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? friendlyName;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  String? id;

  List<String> labels;

  DtoEditableOwnerModel owner;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  SubscriptionEditableSubscriptionPayerModel? payer;

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
  String? priceId;

  String providerId;

  String recurrency;

  List<String> serviceUsers;

  DateTime startDate;

  ///
  /// Please note: This property should have been non-nullable! Since the specification file
  /// does not include a default value (using the "default:" property), however, the generated
  /// source code must fall back to having a nullable type.
  /// Consider adding a "default:" property in the specification file to hide this note.
  ///
  DateTime? updatedAt;

  @override
  bool operator ==(Object other) => identical(this, other) || other is SubscriptionPatchSubscriptionModel &&
    other.customPrice == customPrice &&
    other.customRecurrency == customRecurrency &&
    other.endDate == endDate &&
    other.freeTrial == freeTrial &&
    other.friendlyName == friendlyName &&
    other.id == id &&
    _deepEquality.equals(other.labels, labels) &&
    other.owner == owner &&
    other.payer == payer &&
    other.planId == planId &&
    other.priceId == priceId &&
    other.providerId == providerId &&
    other.recurrency == recurrency &&
    _deepEquality.equals(other.serviceUsers, serviceUsers) &&
    other.startDate == startDate &&
    other.updatedAt == updatedAt;

  @override
  int get hashCode =>
    // ignore: unnecessary_parenthesis
    (customPrice == null ? 0 : customPrice!.hashCode) +
    (customRecurrency == null ? 0 : customRecurrency!.hashCode) +
    (endDate == null ? 0 : endDate!.hashCode) +
    (freeTrial == null ? 0 : freeTrial!.hashCode) +
    (friendlyName == null ? 0 : friendlyName!.hashCode) +
    (id == null ? 0 : id!.hashCode) +
    (labels.hashCode) +
    (owner.hashCode) +
    (payer == null ? 0 : payer!.hashCode) +
    (planId == null ? 0 : planId!.hashCode) +
    (priceId == null ? 0 : priceId!.hashCode) +
    (providerId.hashCode) +
    (recurrency.hashCode) +
    (serviceUsers.hashCode) +
    (startDate.hashCode) +
    (updatedAt == null ? 0 : updatedAt!.hashCode);

  @override
  String toString() => 'SubscriptionPatchSubscriptionModel[customPrice=$customPrice, customRecurrency=$customRecurrency, endDate=$endDate, freeTrial=$freeTrial, friendlyName=$friendlyName, id=$id, labels=$labels, owner=$owner, payer=$payer, planId=$planId, priceId=$priceId, providerId=$providerId, recurrency=$recurrency, serviceUsers=$serviceUsers, startDate=$startDate, updatedAt=$updatedAt]';

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
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
    if (this.id != null) {
      json[r'id'] = this.id;
    } else {
      json[r'id'] = null;
    }
      json[r'labels'] = this.labels;
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
    if (this.priceId != null) {
      json[r'price_id'] = this.priceId;
    } else {
      json[r'price_id'] = null;
    }
      json[r'provider_id'] = this.providerId;
      json[r'recurrency'] = this.recurrency;
      json[r'service_users'] = this.serviceUsers;
      json[r'start_date'] = this.startDate.toUtc().toIso8601String();
    if (this.updatedAt != null) {
      json[r'updated_at'] = this.updatedAt!.toUtc().toIso8601String();
    } else {
      json[r'updated_at'] = null;
    }
    return json;
  }

  /// Returns a new [SubscriptionPatchSubscriptionModel] instance and imports its values from
  /// [value] if it's a [Map], null otherwise.
  // ignore: prefer_constructors_over_static_methods
  static SubscriptionPatchSubscriptionModel? fromJson(dynamic value) {
    if (value is Map) {
      final json = value.cast<String, dynamic>();

      // Ensure that the map contains the required keys.
      // Note 1: the values aren't checked for validity beyond being non-null.
      // Note 2: this code is stripped in release mode!
      assert(() {
        requiredKeys.forEach((key) {
          assert(json.containsKey(key), 'Required key "SubscriptionPatchSubscriptionModel[$key]" is missing from JSON.');
          assert(json[key] != null, 'Required key "SubscriptionPatchSubscriptionModel[$key]" has a null value in JSON.');
        });
        return true;
      }());

      return SubscriptionPatchSubscriptionModel(
        customPrice: DtoAmountModel.fromJson(json[r'custom_price']),
        customRecurrency: mapValueOfType<int>(json, r'custom_recurrency'),
        endDate: mapDateTime(json, r'end_date', r''),
        freeTrial: SubscriptionSubscriptionFreeTrialModel.fromJson(json[r'free_trial']),
        friendlyName: mapValueOfType<String>(json, r'friendly_name'),
        id: mapValueOfType<String>(json, r'id'),
        labels: json[r'labels'] is Iterable
            ? (json[r'labels'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        owner: DtoEditableOwnerModel.fromJson(json[r'owner'])!,
        payer: SubscriptionEditableSubscriptionPayerModel.fromJson(json[r'payer']),
        planId: mapValueOfType<String>(json, r'plan_id'),
        priceId: mapValueOfType<String>(json, r'price_id'),
        providerId: mapValueOfType<String>(json, r'provider_id')!,
        recurrency: mapValueOfType<String>(json, r'recurrency')!,
        serviceUsers: json[r'service_users'] is Iterable
            ? (json[r'service_users'] as Iterable).cast<String>().toList(growable: false)
            : const [],
        startDate: mapDateTime(json, r'start_date', r'')!,
        updatedAt: mapDateTime(json, r'updated_at', r''),
      );
    }
    return null;
  }

  static List<SubscriptionPatchSubscriptionModel> listFromJson(dynamic json, {bool growable = false,}) {
    final result = <SubscriptionPatchSubscriptionModel>[];
    if (json is List && json.isNotEmpty) {
      for (final row in json) {
        final value = SubscriptionPatchSubscriptionModel.fromJson(row);
        if (value != null) {
          result.add(value);
        }
      }
    }
    return result.toList(growable: growable);
  }

  static Map<String, SubscriptionPatchSubscriptionModel> mapFromJson(dynamic json) {
    final map = <String, SubscriptionPatchSubscriptionModel>{};
    if (json is Map && json.isNotEmpty) {
      json = json.cast<String, dynamic>(); // ignore: parameter_assignments
      for (final entry in json.entries) {
        final value = SubscriptionPatchSubscriptionModel.fromJson(entry.value);
        if (value != null) {
          map[entry.key] = value;
        }
      }
    }
    return map;
  }

  // maps a json object with a list of SubscriptionPatchSubscriptionModel-objects as value to a dart map
  static Map<String, List<SubscriptionPatchSubscriptionModel>> mapListFromJson(dynamic json, {bool growable = false,}) {
    final map = <String, List<SubscriptionPatchSubscriptionModel>>{};
    if (json is Map && json.isNotEmpty) {
      // ignore: parameter_assignments
      json = json.cast<String, dynamic>();
      for (final entry in json.entries) {
        map[entry.key] = SubscriptionPatchSubscriptionModel.listFromJson(entry.value, growable: growable,);
      }
    }
    return map;
  }

  /// The list of required keys that must be present in a JSON.
  static const requiredKeys = <String>{
    'owner',
    'provider_id',
    'recurrency',
    'start_date',
  };
}

