import 'package:http/http.dart' as http;
import 'package:microsoft_kiota_bundle/microsoft_kiota_bundle.dart';
import 'package:subscription_tracker/api/api_client.dart';
import 'package:subscription_tracker/api/kinde_auth_provider.dart';
import 'package:subscription_tracker/api/models/subscription_model.dart';
import 'package:subscription_tracker/services/authentication_service.dart';
import '../models/Paginated.dart';
import '../models/family.dart';
import '../models/subscription.dart';
import '../models/family_member.dart';
import '../models/label.dart';
import '../models/subscription_payment.dart';
import 'subscription_service.dart';
import 'family_service.dart';
import 'label_service.dart';

/// Service for handling API requests to the backend
/// API implementation based on the SubTracker swagger specification
/// This class delegates API calls to specialized services for each resource group
class ApiService {
  late final ApiClient client;

  ApiService(String baseUrl, AuthenticationService authenticationService) {
    var authenticationProvider = KindeAuthenticationProvider(
      authenticationService,
    );
    var requestAdapter = DefaultRequestAdapter(
      authProvider: authenticationProvider,
    );
    requestAdapter.baseUrl = baseUrl;
    client = ApiClient(requestAdapter);
  }

  Subscription fromSubscriptionModel(SubscriptionModel model) {
    return Subscription(
        id: model.id!,
        name: model.name!,
        subscriptionPayments: null,
        labelIds: model.labels?.toList(),
        userFamilyMemberIds: model.familyMembers?.toList(),
        payerFamilyMemberId: model.payerIdId,
        payedByJointAccount: model.payedByJointAccount ?? false,
        eTag: model.etag ?? '',
        familyId: model.familyId,
        createdAt: model.createdAt,
        updatedAt: model.updatedAt
    );
  }

  /// Get all subscriptions from the backend
  /// GET /subscriptions
  Future<Paginated<Subscription>> getSubscriptions({
    int size = 10,
    int page = 1,
  }) async {
    final result = await client.subscriptions.getAsync((params) {
      params.queryParameters.page = page;
      params.queryParameters.size = size;
    });
    return Paginated(result.data!, result.length!, result.total!)
  }

  /// Get subscription by ID
  /// GET /subscriptions/{id}
  Future<Subscription> getSubscription(String id) async {
    return _subscriptionService.getSubscription(id);
  }

  /// Delete a subscription from the backend
  /// DELETE /subscriptions/{id}
  Future<void> deleteSubscription(String id) async {
    return _subscriptionService.deleteSubscription(id);
  }

  /// Get all families
  /// GET /families
  Future<List<Family>> getFamilies() async {
    return _familyService.getFamilies();
  }

  Future<Family> patchFamily(Family family) async {
    return _familyService.patchFamily(family);
  }

  /// Get family member by ID
  /// GET /families/members/{id}
  Future<FamilyMember> getFamilyMember(String id) async {
    return _familyService.getFamilyMember(id);
  }

  /// Delete family by ID
  /// DELETE /families/{id}
  Future<void> deleteFamily(String id) async {
    return _familyService.deleteFamily(id);
  }

  /// Get all labels
  /// GET /labels
  Future<List<Label>> getLabels({bool withDefault = true}) async {
    return _labelService.getLabels(withDefault: withDefault);
  }

  /// Get default labels
  /// GET /labels/default
  Future<List<Label>> getDefaultLabels() async {
    return _labelService.getDefaultLabels();
  }

  /// Get label by ID
  /// GET /labels/{id}
  Future<Label> getLabel(String id) async {
    return _labelService.getLabel(id);
  }

  /// Update label by ID
  /// PUT /labels/{id}
  Future<Label> updateLabel(Label label) async {
    return _labelService.updateLabel(label);
  }

  /// Delete label by ID
  /// DELETE /labels/{id}
  Future<void> deleteLabel(String id) async {
    return _labelService.deleteLabel(id);
  }

  /// Create a new label
  /// POST /labels
  Future<Label> createLabel(Label label) async {
    return _labelService.createLabel(label);
  }

  Future<Subscription> patchSubscription(Subscription subscription) async {
    return _subscriptionService.patchSubscription(subscription);
  }

  /// Close the HTTP client when done
  void dispose() {
    _subscriptionService.dispose();
    _familyService.dispose();
    _labelService.dispose();
  }
}
