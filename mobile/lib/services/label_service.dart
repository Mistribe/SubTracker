import 'dart:convert';
import '../models/label.dart';
import 'base_service.dart';

/// Service for handling label-related API requests
class LabelService extends BaseService {
  LabelService({
    required super.baseUrl,
    required super.authenticationService,
    super.httpClient,
  });

  /// Get all labels
  /// GET /labels
  Future<List<Label>> getLabels({bool withDefault = true}) async {
    try {
      final isAuthenticated = await authenticationService.isAuthenticated();
      if (!isAuthenticated) {
        return await getDefaultLabels();
      }
      final response = await httpClient.get(
        Uri.parse('$baseUrl/labels?with_default=$withDefault'),
        headers: await getHeaders(),
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

  /// Get default labels
  /// GET /labels/default
  Future<List<Label>> getDefaultLabels() async {
    try {
      final response = await httpClient.get(
        Uri.parse('$baseUrl/labels/default'),
        headers: await getHeaders(),
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
      final response = await httpClient.get(
        Uri.parse('$baseUrl/labels/$id'),
        headers: await getHeaders(),
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

      final response = await httpClient.put(
        Uri.parse('$baseUrl/labels/${label.id}'),
        headers: await getHeaders(),
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
      final response = await httpClient.delete(
        Uri.parse('$baseUrl/labels/$id'),
        headers: await getHeaders(),
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

      final response = await httpClient.post(
        Uri.parse('$baseUrl/labels'),
        headers: await getHeaders(),
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
}