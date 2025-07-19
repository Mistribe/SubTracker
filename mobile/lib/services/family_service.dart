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

  /// Create a new family member
  /// POST /families/members
  Future<Family> createFamilyMember(FamilyMember member) async {
    try {
      // Create family member payload according to createFamilyMemberModel in swagger
      final memberPayload = {
        'id': member.id,
        'name': member.name,
        'is_kid': member.isKid,
        'created_at': member.createdAt.toUtc().toIso8601String(),
      };

      final familyId = member.familyId;
      final response = await httpClient.post(
        Uri.parse('$baseUrl/families/$familyId/members'),
        headers: await getHeaders(),
        body: json.encode(memberPayload),
      );

      if (response.statusCode == 201) {
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

  /// Update family member by ID
  /// PUT /families/members/{id}
  Future<Family> updateFamilyMember(FamilyMember member) async {
    try {
      // Update family member payload according to updateFamilyMemberModel in swagger
      final updatePayload = {
        'name': member.name,
        'id_kid': member.isKid, // Note: field is id_kid in the API, not is_kid
        'updated_at': member.updatedAt.toUtc().toIso8601String(),
      };

      final response = await httpClient.put(
        Uri.parse('$baseUrl/families/${member.familyId}/members/${member.id}'),
        headers: await getHeaders(),
        body: json.encode(updatePayload),
      );

      if (response.statusCode == 200) {
        return Family.fromJson(json.decode(response.body));
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
  Future<void> deleteFamilyMember(String familyId, String memberId) async {
    try {
      final response = await httpClient.delete(
        Uri.parse('$baseUrl/families/$familyId/members/$memberId'),
        headers: await getHeaders(),
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

  /// Delete family by ID
  /// DELETE /families/{id}
  Future<void> deleteFamily(String id) async {
    // todo not implemented
  }

  /// Create a new family
  /// POST /families
  Future<Family> createFamily(Family family) async {
    try {
      // Create family payload according to createFamilyModel in swagger
      final familyPayload = {
        'id': family.id,
        'name': family.name,
        'have_joint_account': family.haveJointAccount,
        'created_at': family.createdAt.toUtc().toIso8601String(),
      };

      final response = await httpClient.post(
        Uri.parse('$baseUrl/families'),
        headers: await getHeaders(),
        body: json.encode(familyPayload),
      );

      if (response.statusCode == 201) {
        final j = json.decode(response.body);
        return Family.fromJson(j);
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to create family: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }

  /// Update family by ID
  /// PUT /families/{id}
  Future<Family> updateFamily(Family family) async {
    try {
      // Update family payload according to updateFamilyModel in swagger
      final updatePayload = {
        'name': family.name,
        'have_joint_account': family.haveJointAccount,
        'updated_at': family.updatedAt.toUtc().toIso8601String(),
      };

      final response = await httpClient.put(
        Uri.parse('$baseUrl/families/${family.id}'),
        headers: await getHeaders(),
        body: json.encode(updatePayload),
      );

      if (response.statusCode == 200) {
        return Family.fromJson(json.decode(response.body));
      } else if (response.statusCode == 404) {
        final errorData = json.decode(response.body);
        throw Exception(errorData['message'] ?? 'Family not found');
      } else {
        final errorData = json.decode(response.body);
        throw Exception(
          errorData['message'] ??
              'Failed to update family: ${response.statusCode}',
        );
      }
    } catch (e) {
      throw Exception('Failed to connect to the server: $e');
    }
  }
}
