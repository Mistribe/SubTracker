import 'dart:convert';
import '../models/subscription.dart';
import '../models/subscription_payment.dart';
import 'base_service.dart';

/// Service for handling subscription-related API requests
class SubscriptionService extends BaseService {
  SubscriptionService({
    required super.baseUrl,
    required super.authenticationService,
    super.httpClient,
  });

  /// Get all subscriptions from the backend
  /// GET /subscriptions
  Future<List<Subscription>> getSubscriptions() async {
    try {
      final response = await httpClient.get(
        Uri.parse('$baseUrl/subscriptions'),
        headers: await getHeaders(),
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
      final response = await httpClient.get(
        Uri.parse('$baseUrl/subscriptions/$id'),
        headers: await getHeaders(),
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
      final response = await httpClient.post(
        Uri.parse('$baseUrl/subscriptions'),
        headers: await getHeaders(),
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

      final response = await httpClient.put(
        Uri.parse('$baseUrl/subscriptions/${subscription.id}'),
        headers: await getHeaders(),
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
      final response = await httpClient.delete(
        Uri.parse('$baseUrl/subscriptions/$id'),
        headers: await getHeaders(),
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
      final response = await httpClient.delete(
        Uri.parse('$baseUrl/subscriptions/$subscriptionId/payments/$paymentId'),
        headers: await getHeaders(),
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

      final response = await httpClient.post(
        Uri.parse('$baseUrl/subscriptions/$subscriptionId/payments'),
        headers: await getHeaders(),
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

      final response = await httpClient.put(
        Uri.parse(
          '$baseUrl/subscriptions/$subscriptionId/payments/${payment.id}',
        ),
        headers: await getHeaders(),
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
}