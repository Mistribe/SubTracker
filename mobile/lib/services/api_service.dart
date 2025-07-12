import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/subscription.dart';
import '../models/family_member.dart';
import '../models/label.dart';
import '../models/subscription_payment.dart';

/// Service for handling API requests to the backend
/// API implementation based on the SubTracker swagger specification
class ApiService {
  final String baseUrl;
  final http.Client _httpClient;

  ApiService({required this.baseUrl, http.Client? httpClient})
    : _httpClient = httpClient ?? http.Client();

  /// Get all subscriptions from the backend
  /// GET /subscriptions
  Future<List<Subscription>> getSubscriptions() async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/subscriptions'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(response.body);
        return data.map((json) => Subscription.fromJson(json)).toList();
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to load subscriptions: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Get subscription by ID
  /// GET /subscriptions/{id}
  Future<Subscription> getSubscription(String id) async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/subscriptions/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        return Subscription.fromJson(json.decode(response.body));
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Subscription not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to get subscription: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Create a new subscription on the backend
  /// POST /subscriptions
  Future<Subscription> createSubscription(Subscription subscription) async {
    try {
      final response = await _httpClient.post(
        Uri.parse('$baseUrl/subscriptions'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'id': subscription.id,
          'name': subscription.name,
          'payments': subscription.subscriptionPayments
              .map(
                (p) => {
                  'id': p.id,
                  'price': p.price,
                  'start_date': p.startDate.toUtc().toIso8601String(),
                  'end_date': p.endDate?.toUtc().toIso8601String(),
                  'months': p.months,
                  'currency': p.currency,
                  'created_at': p.createdAt.toUtc().toIso8601String(),
                },
              )
              .toList(),
          'labels': subscription.labelIds.toList(),
          'family_members': subscription.userFamilyMemberIds.toList(),
          'payer': subscription.payerFamilyMemberId,
          'created_at': subscription.createdAt.toUtc().toIso8601String(),
        }),
      );

      if (response.statusCode == 201) {
        final responseBody = json.decode(response.body);
        return Subscription.fromJson(responseBody);
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to create subscription: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Update an existing subscription on the backend
  /// PUT /subscriptions/{id}
  Future<Subscription> updateSubscription(Subscription subscription) async {
    try {
      // Prepare the update payload according to swagger updateSubscriptionModel
      final updatePayload = {
        'name': subscription.name,
        'payer': subscription.payerFamilyMemberId,
        'family_members': subscription.userFamilyMemberIds.toList(),
        'labels': subscription.labelIds.toList(),
        'updated_at': subscription.updatedAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.put(
        Uri.parse('$baseUrl/subscriptions/${subscription.id}'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(updatePayload),
      );

      if (response.statusCode == 200) {
        final responseBody = json.decode(response.body);
        return Subscription.fromJson(responseBody);
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Subscription not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to update subscription: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Delete a subscription from the backend
  /// DELETE /subscriptions/{id}
  Future<void> deleteSubscription(String id) async {
    try {
      final response = await _httpClient.delete(
        Uri.parse('$baseUrl/subscriptions/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 204) {
        return; // Success with no content
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Subscription not found');
      } else if (response.statusCode != 200) {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to delete subscription: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Delete a payment from a subscription
  /// DELETE /subscriptions/{id}/payments/{paymentId}
  Future<void> deleteSubscriptionPayment(
    String subscriptionId,
    String paymentId,
  ) async {
    try {
      final response = await _httpClient.delete(
        Uri.parse('$baseUrl/subscriptions/$subscriptionId/payments/$paymentId'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 204) {
        return; // Success with no content
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to delete subscription payment: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Add a new payment to a subscription
  /// POST /subscriptions/{subscription_id}/payments
  Future<SubscriptionPayment> createSubscriptionPayment(
    String subscriptionId,
    SubscriptionPayment payment,
  ) async {
    try {
      // Create payment payload according to createSubscriptionPaymentModel in swagger
      final paymentPayload = {
        'price': payment.price,
        'currency': payment.currency,
        'months': payment.months,
        'start_date': payment.startDate.toUtc().toIso8601String(),
        'end_date': payment.endDate?.toUtc().toIso8601String(),
        'created_at': payment.createdAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.post(
        Uri.parse('$baseUrl/subscriptions/$subscriptionId/payments'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(paymentPayload),
      );

      if (response.statusCode == 201) {
        // The API returns the updated subscription, but we need to extract the payment
        final subscription = Subscription.fromJson(json.decode(response.body));
        // Find the payment with the same ID or the most recently added payment
        final createdPayment = subscription.subscriptionPayments.firstWhere(
          (p) => p.id == payment.id,
          orElse: () => subscription.getLastPaymentDetail(),
        );
        return createdPayment;
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to create subscription payment: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Update a subscription payment
  /// This method is used for synchronization purposes
  Future<SubscriptionPayment> updateSubscriptionPayment(
    String subscriptionId,
    SubscriptionPayment payment,
  ) async {
    try {
      // Update payment payload
      final paymentPayload = {
        'price': payment.price,
        'currency': payment.currency,
        'months': payment.months,
        'start_date': payment.startDate.toUtc().toIso8601String(),
        'end_date': payment.endDate?.toUtc().toIso8601String(),
        'updated_at': payment.updatedAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.put(
        Uri.parse(
          '$baseUrl/subscriptions/$subscriptionId/payments/${payment.id}',
        ),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(paymentPayload),
      );

      if (response.statusCode == 200) {
        // The API returns the updated subscription, but we need to extract the payment
        final subscription = Subscription.fromJson(json.decode(response.body));
        // Find the payment with the same ID
        final updatedPayment = subscription.subscriptionPayments.firstWhere(
          (p) => p.id == payment.id,
        );
        return updatedPayment;
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to update subscription payment: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Get all family members
  /// GET /families/members
  Future<List<FamilyMember>> getFamilyMembers() async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/families/members'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(response.body);
        return data.map((json) => FamilyMember.fromJson(json)).toList();
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to load family members: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Create a new family member
  /// POST /families/members
  Future<FamilyMember> createFamilyMember(FamilyMember member) async {
    try {
      // Create family member payload according to createFamilyMemberModel in swagger
      final memberPayload = {
        'id': member.id,
        'name': member.name,
        'is_kid': member.isKid,
        'created_at': member.createdAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.post(
        Uri.parse('$baseUrl/families/members'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(memberPayload),
      );

      if (response.statusCode == 201) {
        return FamilyMember.fromJson(json.decode(response.body));
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to create family member: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Get family member by ID
  /// GET /families/members/{id}
  Future<FamilyMember> getFamilyMember(String id) async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/families/members/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        return FamilyMember.fromJson(json.decode(response.body));
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Family member not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to get family member: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Update family member by ID
  /// PUT /families/members/{id}
  Future<FamilyMember> updateFamilyMember(FamilyMember member) async {
    try {
      // Update family member payload according to updateFamilyMemberModel in swagger
      final updatePayload = {
        'name': member.name,
        'id_kid': member.isKid, // Note: field is id_kid in the API, not is_kid
        'updated_at': member.updatedAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.put(
        Uri.parse('$baseUrl/families/members/${member.id}'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(updatePayload),
      );

      if (response.statusCode == 200) {
        return FamilyMember.fromJson(json.decode(response.body));
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Family member not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to update family member: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Delete family member by ID
  /// DELETE /families/members/{id}
  Future<void> deleteFamilyMember(String id) async {
    try {
      final response = await _httpClient.delete(
        Uri.parse('$baseUrl/families/members/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 204) {
        return; // Success with no content
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Family member not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to delete family member: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Get all labels
  /// GET /labels
  Future<List<Label>> getLabels({bool withDefault = true}) async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/labels?with_default=$withDefault'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(response.body);
        return data.map((json) => Label.fromJson(json)).toList();
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to load labels: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  Future<List<Label>> getDefaultLabels() async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/labels/default'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(response.body);
        return data.map((json) => Label.fromJson(json)).toList();
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to load labels: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Get label by ID
  /// GET /labels/{id}
  Future<Label> getLabel(String id) async {
    try {
      final response = await _httpClient.get(
        Uri.parse('$baseUrl/labels/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        return Label.fromJson(json.decode(response.body));
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Label not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ?? 'Failed to get label: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Update label by ID
  /// PUT /labels/{id}
  Future<Label> updateLabel(Label label) async {
    try {
      // Update label payload according to updateLabelModel in swagger
      final updatePayload = {
        'name': label.name,
        'color': label.color,
        'updated_at': label.updatedAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.put(
        Uri.parse('$baseUrl/labels/${label.id}'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(updatePayload),
      );

      if (response.statusCode == 200) {
        return Label.fromJson(json.decode(response.body));
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Label not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to update label: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Delete label by ID
  /// DELETE /labels/{id}
  Future<void> deleteLabel(String id) async {
    try {
      final response = await _httpClient.delete(
        Uri.parse('$baseUrl/labels/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 204) {
        return; // Success with no content
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Label not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to delete label: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Create a new label
  /// POST /labels
  Future<Label> createLabel(Label label) async {
    try {
      // Create label payload according to createLabelModel in swagger
      final labelPayload = {
        'id': label.id,
        'name': label.name,
        'color': label.color,
        'created_at': label.createdAt.toUtc().toIso8601String(),
      };

      final response = await _httpClient.post(
        Uri.parse('$baseUrl/labels'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(labelPayload),
      );

      if (response.statusCode == 201) {
        return Label.fromJson(json.decode(response.body));
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to create label: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Close the HTTP client when done
  void dispose() {
    _httpClient.close();
  }
}
