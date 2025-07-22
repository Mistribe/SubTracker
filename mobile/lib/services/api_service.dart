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
