import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:kinde_flutter_sdk/kinde_flutter_sdk.dart';
import 'package:subscription_tracker/services/authentication_service.dart';

/// Base service class for API services
/// Contains common functionality for all API services
class BaseService {
  final String baseUrl;
  final AuthenticationService authenticationService;
  final http.Client _httpClient;

  BaseService({
    required this.baseUrl,
    required this.authenticationService,
    http.Client? httpClient,
  }) : _httpClient = httpClient ?? http.Client();

  /// Get authentication headers
  Future<Map<String, String>> getHeaders({bool requiresAuth = true}) async {
    final headers = {'Content-Type': 'application/json'};

    if (requiresAuth) {
      try {
        final token = await authenticationService.getToken();
        if (token != null && token.isNotEmpty) {
          headers['Authorization'] = 'Bearer $token';
        }
      } catch (e) {
        print('Error getting authentication token: $e');
      }
    }

    return headers;
  }

  /// Check if the user is authenticated
  Future<bool> isAuthenticated() async {
    try {
      return await KindeFlutterSDK.instance.isAuthenticated();
    } catch (e) {
      print('Error checking authentication status: $e');
      return false;
    }
  }

  /// Close the HTTP client when done
  void dispose() {
    _httpClient.close();
  }

  /// Get HTTP client
  http.Client get httpClient => _httpClient;
}