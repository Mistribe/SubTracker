//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class FamilyApi {
  FamilyApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Accept a family invitation
  ///
  /// Accepts an invitation to join a family using the provided invitation code
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyFamilyAcceptInvitationRequest] familyFamilyAcceptInvitationRequest (required):
  ///   Invitation acceptance details
  Future<Response> familiesFamilyIdAcceptPostWithHttpInfo(String familyId, FamilyFamilyAcceptInvitationRequest familyFamilyAcceptInvitationRequest,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/accept'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody = familyFamilyAcceptInvitationRequest;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Accept a family invitation
  ///
  /// Accepts an invitation to join a family using the provided invitation code
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyFamilyAcceptInvitationRequest] familyFamilyAcceptInvitationRequest (required):
  ///   Invitation acceptance details
  Future<void> familiesFamilyIdAcceptPost(String familyId, FamilyFamilyAcceptInvitationRequest familyFamilyAcceptInvitationRequest,) async {
    final response = await familiesFamilyIdAcceptPostWithHttpInfo(familyId, familyFamilyAcceptInvitationRequest,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Decline family invitation
  ///
  /// Endpoint to decline an invitation to join a family
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID
  ///
  /// * [FamilyFamilyDeclineInvitationRequest] familyFamilyDeclineInvitationRequest (required):
  ///   Decline invitation request
  Future<Response> familiesFamilyIdDeclinePostWithHttpInfo(String familyId, FamilyFamilyDeclineInvitationRequest familyFamilyDeclineInvitationRequest,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/decline'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody = familyFamilyDeclineInvitationRequest;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Decline family invitation
  ///
  /// Endpoint to decline an invitation to join a family
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID
  ///
  /// * [FamilyFamilyDeclineInvitationRequest] familyFamilyDeclineInvitationRequest (required):
  ///   Decline invitation request
  Future<void> familiesFamilyIdDeclinePost(String familyId, FamilyFamilyDeclineInvitationRequest familyFamilyDeclineInvitationRequest,) async {
    final response = await familiesFamilyIdDeclinePostWithHttpInfo(familyId, familyFamilyDeclineInvitationRequest,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Delete family by ID
  ///
  /// Permanently delete a family and all its members
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  Future<Response> familiesFamilyIdDeleteWithHttpInfo(String familyId,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Delete family by ID
  ///
  /// Permanently delete a family and all its members
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  Future<void> familiesFamilyIdDelete(String familyId,) async {
    final response = await familiesFamilyIdDeleteWithHttpInfo(familyId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// View family invitation details
  ///
  /// Get information about a family invitation using invitation code
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID
  ///
  /// * [String] code (required):
  ///   Invitation code
  ///
  /// * [String] familyMemberId (required):
  ///   Family member ID
  Future<Response> familiesFamilyIdInvitationGetWithHttpInfo(String familyId, String code, String familyMemberId,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/invitation'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

      queryParams.addAll(_queryParams('', 'code', code));
      queryParams.addAll(_queryParams('', 'family_member_id', familyMemberId));

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// View family invitation details
  ///
  /// Get information about a family invitation using invitation code
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID
  ///
  /// * [String] code (required):
  ///   Invitation code
  ///
  /// * [String] familyMemberId (required):
  ///   Family member ID
  Future<FamilyFamilySeeInvitationResponse?> familiesFamilyIdInvitationGet(String familyId, String code, String familyMemberId,) async {
    final response = await familiesFamilyIdInvitationGetWithHttpInfo(familyId, code, familyMemberId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilySeeInvitationResponse',) as FamilyFamilySeeInvitationResponse;
    
    }
    return null;
  }

  /// Invite a new member to the family
  ///
  /// Creates an invitation for a new member to join the family
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyFamilyInviteRequest] familyFamilyInviteRequest (required):
  ///   Invitation details including email, name, member ID and type (adult/kid)
  Future<Response> familiesFamilyIdInvitePostWithHttpInfo(String familyId, FamilyFamilyInviteRequest familyFamilyInviteRequest,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/invite'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody = familyFamilyInviteRequest;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Invite a new member to the family
  ///
  /// Creates an invitation for a new member to join the family
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyFamilyInviteRequest] familyFamilyInviteRequest (required):
  ///   Invitation details including email, name, member ID and type (adult/kid)
  Future<FamilyFamilyInviteResponse?> familiesFamilyIdInvitePost(String familyId, FamilyFamilyInviteRequest familyFamilyInviteRequest,) async {
    final response = await familiesFamilyIdInvitePostWithHttpInfo(familyId, familyFamilyInviteRequest,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyInviteResponse',) as FamilyFamilyInviteResponse;
    
    }
    return null;
  }

  /// Revoke family member
  ///
  /// Revokes a member from the family
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [String] familyMemberId (required):
  ///   Family Member ID (UUID format)
  ///
  /// * [Object] body:
  Future<Response> familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(String familyId, String familyMemberId, { Object? body, }) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/members/{familyMemberId}/revoke'
      .replaceAll('{familyId}', familyId)
      .replaceAll('{familyMemberId}', familyMemberId);

    // ignore: prefer_final_locals
    Object? postBody = body;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Revoke family member
  ///
  /// Revokes a member from the family
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [String] familyMemberId (required):
  ///   Family Member ID (UUID format)
  ///
  /// * [Object] body:
  Future<void> familiesFamilyIdMembersFamilyMemberIdRevokePost(String familyId, String familyMemberId, { Object? body, }) async {
    final response = await familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(familyId, familyMemberId,  body: body, );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Delete family member by ID
  ///
  /// Permanently delete a family member from a family
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [String] id (required):
  ///   Family member ID (UUID format)
  Future<Response> familiesFamilyIdMembersIdDeleteWithHttpInfo(String familyId, String id,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/members/{id}'
      .replaceAll('{familyId}', familyId)
      .replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Delete family member by ID
  ///
  /// Permanently delete a family member from a family
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [String] id (required):
  ///   Family member ID (UUID format)
  Future<void> familiesFamilyIdMembersIdDelete(String familyId, String id,) async {
    final response = await familiesFamilyIdMembersIdDeleteWithHttpInfo(familyId, id,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Update family member by ID
  ///
  /// Update an existing family member's information such as name and kid status
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [String] id (required):
  ///   Family member ID (UUID format)
  ///
  /// * [FamilyUpdateFamilyMemberModel] familyUpdateFamilyMemberModel (required):
  ///   Updated family member data
  Future<Response> familiesFamilyIdMembersIdPutWithHttpInfo(String familyId, String id, FamilyUpdateFamilyMemberModel familyUpdateFamilyMemberModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/members/{id}'
      .replaceAll('{familyId}', familyId)
      .replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody = familyUpdateFamilyMemberModel;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'PUT',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Update family member by ID
  ///
  /// Update an existing family member's information such as name and kid status
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [String] id (required):
  ///   Family member ID (UUID format)
  ///
  /// * [FamilyUpdateFamilyMemberModel] familyUpdateFamilyMemberModel (required):
  ///   Updated family member data
  Future<FamilyFamilyModel?> familiesFamilyIdMembersIdPut(String familyId, String id, FamilyUpdateFamilyMemberModel familyUpdateFamilyMemberModel,) async {
    final response = await familiesFamilyIdMembersIdPutWithHttpInfo(familyId, id, familyUpdateFamilyMemberModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyModel',) as FamilyFamilyModel;
    
    }
    return null;
  }

  /// Add a new family member
  ///
  /// Add a new member to an existing family
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyCreateFamilyMemberModel] familyCreateFamilyMemberModel (required):
  ///   Family member creation data
  Future<Response> familiesFamilyIdMembersPostWithHttpInfo(String familyId, FamilyCreateFamilyMemberModel familyCreateFamilyMemberModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}/members'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody = familyCreateFamilyMemberModel;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Add a new family member
  ///
  /// Add a new member to an existing family
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyCreateFamilyMemberModel] familyCreateFamilyMemberModel (required):
  ///   Family member creation data
  Future<FamilyFamilyModel?> familiesFamilyIdMembersPost(String familyId, FamilyCreateFamilyMemberModel familyCreateFamilyMemberModel,) async {
    final response = await familiesFamilyIdMembersPostWithHttpInfo(familyId, familyCreateFamilyMemberModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyModel',) as FamilyFamilyModel;
    
    }
    return null;
  }

  /// Update a family
  ///
  /// Update family information such as name and other details
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyUpdateFamilyModel] familyUpdateFamilyModel (required):
  ///   Updated family data
  Future<Response> familiesFamilyIdPutWithHttpInfo(String familyId, FamilyUpdateFamilyModel familyUpdateFamilyModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/families/{familyId}'
      .replaceAll('{familyId}', familyId);

    // ignore: prefer_final_locals
    Object? postBody = familyUpdateFamilyModel;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'PUT',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Update a family
  ///
  /// Update family information such as name and other details
  ///
  /// Parameters:
  ///
  /// * [String] familyId (required):
  ///   Family ID (UUID format)
  ///
  /// * [FamilyUpdateFamilyModel] familyUpdateFamilyModel (required):
  ///   Updated family data
  Future<FamilyFamilyModel?> familiesFamilyIdPut(String familyId, FamilyUpdateFamilyModel familyUpdateFamilyModel,) async {
    final response = await familiesFamilyIdPutWithHttpInfo(familyId, familyUpdateFamilyModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyModel',) as FamilyFamilyModel;
    
    }
    return null;
  }

  /// Get user's family
  ///
  /// Retrieve the user's family
  ///
  /// Note: This method returns the HTTP [Response].
  Future<Response> familiesMeGetWithHttpInfo() async {
    // ignore: prefer_const_declarations
    final path = r'/families/me';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get user's family
  ///
  /// Retrieve the user's family
  Future<FamilyFamilyModel?> familiesMeGet() async {
    final response = await familiesMeGetWithHttpInfo();
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyModel',) as FamilyFamilyModel;
    
    }
    return null;
  }

  /// Patch family with members
  ///
  /// Update or create a family with specified members. If family doesn't exist, it will be created.
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [FamilyPatchFamilyModel] familyPatchFamilyModel (required):
  ///   Family update data with members
  Future<Response> familiesPatchWithHttpInfo(FamilyPatchFamilyModel familyPatchFamilyModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/families';

    // ignore: prefer_final_locals
    Object? postBody = familyPatchFamilyModel;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'PATCH',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Patch family with members
  ///
  /// Update or create a family with specified members. If family doesn't exist, it will be created.
  ///
  /// Parameters:
  ///
  /// * [FamilyPatchFamilyModel] familyPatchFamilyModel (required):
  ///   Family update data with members
  Future<FamilyFamilyModel?> familiesPatch(FamilyPatchFamilyModel familyPatchFamilyModel,) async {
    final response = await familiesPatchWithHttpInfo(familyPatchFamilyModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyModel',) as FamilyFamilyModel;
    
    }
    return null;
  }

  /// Create a new family
  ///
  /// Create a new family with the authenticated user as the owner and initial member
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [FamilyCreateFamilyModel] familyCreateFamilyModel (required):
  ///   Family creation data
  Future<Response> familiesPostWithHttpInfo(FamilyCreateFamilyModel familyCreateFamilyModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/families';

    // ignore: prefer_final_locals
    Object? postBody = familyCreateFamilyModel;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Create a new family
  ///
  /// Create a new family with the authenticated user as the owner and initial member
  ///
  /// Parameters:
  ///
  /// * [FamilyCreateFamilyModel] familyCreateFamilyModel (required):
  ///   Family creation data
  Future<FamilyFamilyModel?> familiesPost(FamilyCreateFamilyModel familyCreateFamilyModel,) async {
    final response = await familiesPostWithHttpInfo(familyCreateFamilyModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'FamilyFamilyModel',) as FamilyFamilyModel;
    
    }
    return null;
  }
}
