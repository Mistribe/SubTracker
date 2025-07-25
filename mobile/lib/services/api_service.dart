import 'package:microsoft_kiota_bundle/microsoft_kiota_bundle.dart';
import 'package:subscription_tracker/api/api_client.dart';
import 'package:subscription_tracker/api/models/create_label_model.dart';
import 'package:subscription_tracker/api/models/label_model.dart';
import 'package:subscription_tracker/api/models/patch_subscription_model.dart';
import 'package:subscription_tracker/api/models/update_label_model.dart';
import 'package:subscription_tracker/providers/kinde_auth_provider.dart';
import 'package:subscription_tracker/api/models/family_member_model.dart';
import 'package:subscription_tracker/api/models/family_model.dart';
import 'package:subscription_tracker/api/models/patch_family_member_model.dart';
import 'package:subscription_tracker/api/models/subscription_model.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';
import 'package:subscription_tracker/services/authentication_service.dart';
import '../api/models/patch_family_model.dart';
import '../api/models/patch_subscription_payment_model.dart';
import '../api/models/subscription_payment_model.dart';
import '../models/paginated.dart';
import '../models/family.dart';
import '../models/subscription.dart';
import '../models/family_member.dart';
import '../models/label.dart';

class ApiService {
  late final ApiClient _client;

  ApiService(String baseUrl, AuthenticationService authenticationService)
    : _client = _createApiClient(baseUrl, authenticationService);

  static ApiClient _createApiClient(
    String baseUrl,
    AuthenticationService authenticationService,
  ) {
    var authenticationProvider = KindeAuthenticationProvider(
      authenticationService,
    );
    var requestAdapter = DefaultRequestAdapter(
      authProvider: authenticationProvider,
    );
    requestAdapter.baseUrl = baseUrl;
    return ApiClient(requestAdapter);
  }

  SubscriptionPayment _fromSubscriptionPaymentModel(
    SubscriptionPaymentModel model,
  ) {
    return SubscriptionPayment(
      id: model.id!,
      price: model.price!,
      startDate: model.startDate!,
      endDate: model.endDate,
      months: model.months!,
      currency: model.currency!,
      createdAt: model.createdAt!,
      updatedAt: model.updatedAt!,
      eTag: model.etag!,
    );
  }

  Subscription _fromSubscriptionModel(SubscriptionModel model) {
    return Subscription(
      id: model.id!,
      name: model.name!,
      subscriptionPayments:
          model.payments?.map(_fromSubscriptionPaymentModel).toList() ?? [],
      labelIds: model.labels?.toList(),
      userFamilyMemberIds: model.familyMembers?.toList(),
      payerId: model.payerIdId,
      payedByJointAccount: model.payedByJointAccount ?? false,
      eTag: model.etag ?? '',
      familyId: model.familyId,
      createdAt: model.createdAt,
      updatedAt: model.updatedAt,
    );
  }

  FamilyMember _fromFamilyMemberModel(FamilyMemberModel model) {
    return FamilyMember(
      id: model.id!,
      name: model.name!,
      familyId: model.familyId!,
      email: '',
      isKid: model.isKid!,
      createdAt: model.createdAt!,
      updatedAt: model.updatedAt!,
      eTag: model.etag!,
    );
  }

  Family _fromFamilyModel(FamilyModel model) {
    return Family(
      id: model.id!,
      name: model.name!,
      isOwner: model.isOwner!,
      haveJointAccount: model.haveJointAccount!,
      familyMembers: model.members?.map(_fromFamilyMemberModel).toList() ?? [],
      createdAt: model.createdAt!,
      updatedAt: model.updatedAt!,
      eTag: model.etag!,
    );
  }

  Label _fromLabelModel(LabelModel model) {
    return Label(
      id: model.id!,
      color: model.color!,
      name: model.name!,
      isDefault: model.isDefault!,
      createdAt: model.createdAt!,
      updatedAt: model.updatedAt!,
      eTag: model.etag!,
    );
  }

  Future<Paginated<Subscription>> getSubscriptions({
    int size = 10,
    int page = 1,
  }) async {
    final result = await _client.subscriptions.getAsync((params) {
      params.queryParameters.page = page;
      params.queryParameters.size = size;
    });
    if (result == null) {
      return Paginated.empty<Subscription>();
    }
    return Paginated(
      result.data?.map(_fromSubscriptionModel).toList(),
      result.length,
      result.total,
    );
  }

  Future<Subscription> patchSubscription(Subscription subscription) async {
    final payload = PatchSubscriptionModel();
    payload.id = subscription.id;
    payload.name = subscription.name;
    payload.labels = subscription.labelIds;
    payload.familyMembers = subscription.familyMemberIds;
    payload.payerId = subscription.payerId;
    payload.payedByJointAccount = subscription.payedByJointAccount;
    payload.familyId = subscription.familyId;
    payload.updatedAt = subscription.updatedAt;
    payload.payments = subscription.payments.map((payment) {
      final model = PatchSubscriptionPaymentModel();
      model.id = payment.id;
      model.currency = payment.currency;
      model.updatedAt = payment.updatedAt;
      model.months = payment.months;
      model.endDate = payment.endDate;
      model.startDate = payment.startDate;
      model.price = payment.price;
      return model;
    });

    final result = await _client.subscriptions.patchAsync(payload);
    if (result == null) {
      throw Exception('Failed to patch subscription');
    }
    return _fromSubscriptionModel(result);
  }

  Future<Subscription?> getSubscription(String id) async {
    final result = await _client.subscriptions.byId(id).getAsync();
    if (result == null) {
      return null;
    }
    return _fromSubscriptionModel(result);
  }

  Future<void> deleteSubscription(String id) async {
    await _client.subscriptions.byId(id).deleteAsync();
  }

  Future<Paginated<Family>> getFamilies() async {
    final result = await _client.families.getAsync();
    if (result == null) {
      return Paginated.empty<Family>();
    }
    return Paginated(
      result.data?.map(_fromFamilyModel).toList(),
      result.length,
      result.total,
    );
  }

  Future<Family> patchFamily(Family family) async {
    final payload = PatchFamilyModel();
    payload.id = family.id;
    payload.name = family.name;
    payload.haveJointAccount = family.haveJointAccount;
    payload.members = family.members.map((member) {
      final model = PatchFamilyMemberModel();
      model.id = member.id;
      model.name = member.name;
      model.isKid = member.isKid;
      model.email = member.email;
      model.updatedAt = member.updatedAt;
      return model;
    });
    payload.updatedAt = family.updatedAt;
    final result = await _client.families.patchAsync(payload);
    if (result == null) {
      throw Exception('Failed to patch family');
    }
    return _fromFamilyModel(result);
  }

  Future<void> deleteFamily(String familyId) async {
    await _client.families.byFamilyId(familyId).deleteAsync();
  }

  Future<Paginated<Label>> getLabels({bool withDefault = true}) async {
    final result = await _client.labels.getAsync((params) {
      params.queryParameters.withDefault = withDefault;
    });
    if (result == null) {
      return Paginated.empty<Label>();
    }
    return Paginated(
      result.data?.map(_fromLabelModel).toList() ?? [],
      result.length,
      result.total,
    );
  }

  Future<List<Label>> getDefaultLabels() async {
    final result = await _client.labels.default_.getAsync();
    if (result == null) {
      return [];
    }
    return result.map(_fromLabelModel).toList();
  }

  Future<Label> getLabel(String id) async {
    final result = await _client.labels.byId(id).getAsync();
    if (result == null) {
      throw Exception('Label not found');
    }
    return _fromLabelModel(result);
  }

  Future<Label> updateLabel(Label label) async {
    final payload = UpdateLabelModel();
    payload.name = label.name;
    payload.color = label.color;
    payload.updatedAt = label.updatedAt;
    final result = await _client.labels.byId(label.id).putAsync(payload);
    if (result == null) {
      throw Exception('Failed to update label');
    }
    return _fromLabelModel(result);
  }

  Future<void> deleteLabel(String id) async {
    await _client.labels.byId(id).deleteAsync();
  }

  Future<Label> createLabel(Label label) async {
    final payload = CreateLabelModel();
    payload.id = label.id;
    payload.name = label.name;
    payload.color = label.color;
    payload.createdAt = label.createdAt;

    final result = await _client.labels.postAsync(payload);
    if (result == null) {
      throw Exception('Failed to create label');
    }
    return _fromLabelModel(result);
  }
}
