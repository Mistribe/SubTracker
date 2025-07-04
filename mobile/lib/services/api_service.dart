import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/subscription.dart';

/// Service for handling API requests to the backend
class ApiService {
  final String baseUrl;
  final http.Client _httpClient;

  ApiService({
    required this.baseUrl,
    http.Client? httpClient,
  }) : _httpClient = httpClient ?? http.Client();

  /// Get all subscriptions from the backend
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
        throw Exception('Failed to load subscriptions: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Create a new subscription on the backend
  Future<Subscription> createSubscription(Subscription subscription) async {
    try {
      final response = await _httpClient.post(
        Uri.parse('$baseUrl/subscriptions'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(subscription.toJson()),
      );

      if (response.statusCode == 201) {
        return Subscription.fromJson(json.decode(response.body));
      } else {
        throw Exception('Failed to create subscription: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Update an existing subscription on the backend
  Future<Subscription> updateSubscription(Subscription subscription) async {
    try {
      final response = await _httpClient.put(
        Uri.parse('$baseUrl/subscriptions/${subscription.id}'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode(subscription.toJson()),
      );

      if (response.statusCode == 200) {
        return Subscription.fromJson(json.decode(response.body));
      } else {
        throw Exception('Failed to update subscription: ${response.statusCode}');
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Delete a subscription from the backend
  Future<void> deleteSubscription(String id) async {
    try {
      final response = await _httpClient.delete(
        Uri.parse('$baseUrl/subscriptions/$id'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode != 204 && response.statusCode != 200) {
        throw Exception('Failed to delete subscription: ${response.statusCode}');
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