import 'dart:convert';
import '../models/subscription.dart';
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
        final responseBody = response.body;
        final List<dynamic> data = json.decode(responseBody);
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

  Future<Subscription> patchSubscription(Subscription subscription) async {
    try {
      final payload = subscription.toJson();
      final response = await httpClient.patch(
        Uri.parse('$baseUrl/subscriptions'),
        headers: await getHeaders(),
        body: json.encode(payload),
      );

      if (response.statusCode == 200) {
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
}
