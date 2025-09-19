//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class UsersApi {
  UsersApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Delete user
  ///
  /// Deletes the authenticated user's account
  ///
  /// Note: This method returns the HTTP [Response].
  Future<Response> usersDeleteWithHttpInfo() async {
    // ignore: prefer_const_declarations
    final path = r'/users';

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

  /// Delete user
  ///
  /// Deletes the authenticated user's account
  Future<void> usersDelete() async {
    final response = await usersDeleteWithHttpInfo();
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Get user preferred currency
  ///
  /// Returns the preferred currency for the authenticated user
  ///
  /// Note: This method returns the HTTP [Response].
  Future<Response> usersPreferredCurrencyGetWithHttpInfo() async {
    // ignore: prefer_const_declarations
    final path = r'/users/preferred/currency';

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

  /// Get user preferred currency
  ///
  /// Returns the preferred currency for the authenticated user
  Future<UserUserPreferredCurrencyModel?> usersPreferredCurrencyGet() async {
    final response = await usersPreferredCurrencyGetWithHttpInfo();
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'UserUserPreferredCurrencyModel',) as UserUserPreferredCurrencyModel;
    
    }
    return null;
  }

  /// Update user preferred currency
  ///
  /// Updates the preferred currency for the authenticated user
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] authorization (required):
  ///   Bearer token
  ///
  /// * [UserUpdatePreferredCurrencyModel] userUpdatePreferredCurrencyModel (required):
  ///   Profile update parameters
  Future<Response> usersPreferredCurrencyPutWithHttpInfo(String authorization, UserUpdatePreferredCurrencyModel userUpdatePreferredCurrencyModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/users/preferred/currency';

    // ignore: prefer_final_locals
    Object? postBody = userUpdatePreferredCurrencyModel;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    headerParams[r'Authorization'] = parameterToString(authorization);

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

  /// Update user preferred currency
  ///
  /// Updates the preferred currency for the authenticated user
  ///
  /// Parameters:
  ///
  /// * [String] authorization (required):
  ///   Bearer token
  ///
  /// * [UserUpdatePreferredCurrencyModel] userUpdatePreferredCurrencyModel (required):
  ///   Profile update parameters
  Future<void> usersPreferredCurrencyPut(String authorization, UserUpdatePreferredCurrencyModel userUpdatePreferredCurrencyModel,) async {
    final response = await usersPreferredCurrencyPutWithHttpInfo(authorization, userUpdatePreferredCurrencyModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }
}
