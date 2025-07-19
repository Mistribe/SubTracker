import 'package:http/http.dart' as http;
import 'package:subscription_tracker/services/authentication_service.dart';
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
  final String baseUrl;
  final AuthenticationService authenticationService;

  // Specialized services
  final SubscriptionService _subscriptionService;
  final FamilyService _familyService;
  final LabelService _labelService;

  ApiService({
    required this.baseUrl,
    required this.authenticationService,
    http.Client? httpClient,
  }) : _subscriptionService = SubscriptionService(
         baseUrl: baseUrl,
         authenticationService: authenticationService,
         httpClient: httpClient,
       ),
       _familyService = FamilyService(
         baseUrl: baseUrl,
         authenticationService: authenticationService,
         httpClient: httpClient,
       ),
       _labelService = LabelService(
         baseUrl: baseUrl,
         authenticationService: authenticationService,
         httpClient: httpClient,
       );

  /// Check if the user is authenticated
  Future<bool> isAuthenticated() async {
    return _subscriptionService.isAuthenticated();
  }

  /// Get all subscriptions from the backend
  /// GET /subscriptions
  Future<List<Subscription>> getSubscriptions() async {
    return _subscriptionService.getSubscriptions();
  }

  /// Get subscription by ID
  /// GET /subscriptions/{id}
  Future<Subscription> getSubscription(String id) async {
    return _subscriptionService.getSubscription(id);
  }

  /// Create a new subscription on the backend
  /// POST /subscriptions
  Future<Subscription> createSubscription(Subscription subscription) async {
    return _subscriptionService.createSubscription(subscription);
  }

  /// Update an existing subscription on the backend
  /// PUT /subscriptions/{id}
  Future<Subscription> updateSubscription(Subscription subscription) async {
    return _subscriptionService.updateSubscription(subscription);
  }

  /// Delete a subscription from the backend
  /// DELETE /subscriptions/{id}
  Future<void> deleteSubscription(String id) async {
    return _subscriptionService.deleteSubscription(id);
  }

  /// Delete a payment from a subscription
  /// DELETE /subscriptions/{id}/payments/{paymentId}
  Future<void> deleteSubscriptionPayment(
    String subscriptionId,
    String paymentId,
  ) async {
    return _subscriptionService.deleteSubscriptionPayment(
      subscriptionId,
      paymentId,
    );
  }

  /// Add a new payment to a subscription
  /// POST /subscriptions/{subscription_id}/payments
  Future<SubscriptionPayment> createSubscriptionPayment(
    String subscriptionId,
    SubscriptionPayment payment,
  ) async {
    return _subscriptionService.createSubscriptionPayment(
      subscriptionId,
      payment,
    );
  }

  /// Update a subscription payment
  /// This method is used for synchronization purposes
  Future<SubscriptionPayment> updateSubscriptionPayment(
    String subscriptionId,
    SubscriptionPayment payment,
  ) async {
    return _subscriptionService.updateSubscriptionPayment(
      subscriptionId,
      payment,
    );
  }

  /// Get all families
  /// GET /families
  Future<List<Family>> getFamilies() async {
    return _familyService.getFamilies();
  }

  /// Create a new family member
  /// POST /families/members
  Future<Family> createFamilyMember(FamilyMember member) async {
    return _familyService.createFamilyMember(member);
  }

  /// Get family member by ID
  /// GET /families/members/{id}
  Future<FamilyMember> getFamilyMember(String id) async {
    return _familyService.getFamilyMember(id);
  }

  /// Update family member by ID
  /// PUT /families/members/{id}
  Future<Family> updateFamilyMember(FamilyMember member) async {
    return _familyService.updateFamilyMember(member);
  }

  /// Delete family member by ID
  /// DELETE /families/members/{id}
  Future<void> deleteFamilyMember(String familyId, String memberId) async {
    return _familyService.deleteFamilyMember(familyId, memberId);
  }

  /// Delete family by ID
  /// DELETE /families/{id}
  Future<void> deleteFamily(String id) async {
    return _familyService.deleteFamily(id);
  }

  /// Create a new family
  /// POST /families
  Future<Family> createFamily(Family family) async {
    return _familyService.createFamily(family);
  }

  /// Update family by ID
  /// PUT /families/{id}
  Future<Family> updateFamily(Family family) async {
    return _familyService.updateFamily(family);
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

  /// Close the HTTP client when done
  void dispose() {
    _subscriptionService.dispose();
    _familyService.dispose();
    _labelService.dispose();
  }
}
