import 'dart:convert';
import '../models/family.dart';
import '../models/family_member.dart';
import 'base_service.dart';

/// Service for handling family-related API requests
class FamilyService extends BaseService {
  FamilyService({
    required super.baseUrl,
    required super.authenticationService,
    super.httpClient,
  });

  /// Get all families
  /// GET /families
  Future<List<Family>> getFamilies() async {
    try {
      final response = await httpClient.get(
        Uri.parse('$baseUrl/families'),
        headers: await getHeaders(),
      );

      if (response.statusCode == 200) {
        final List<dynamic> data = json.decode(response.body);
        final families = data.map((json) => Family.fromJson(json)).toList();
        return families;
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

  Future<Family> patchFamily(Family family) async {
    try {
      final payload = family.toJson();
      final response = await httpClient.patch(
        Uri.parse('$baseUrl/families'),
        headers: await getHeaders(),
        body: json.encode(payload),
      );

      if (response.statusCode == 200) {
        final j = json.decode(response.body);
        return Family.fromJson(j);
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
      final response = await httpClient.get(
        Uri.parse('$baseUrl/families/members/$id'),
        headers: await getHeaders(),
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

  /// Delete family by ID
  /// DELETE /families/{id}
  Future<void> deleteFamily(String id) async {
    // todo not implemented
  }
}
