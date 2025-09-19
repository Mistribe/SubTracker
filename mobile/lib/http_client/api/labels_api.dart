//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class LabelsApi {
  LabelsApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Get default labels
  ///
  /// Retrieves a list of default system labels available to all users
  ///
  /// Note: This method returns the HTTP [Response].
  Future<Response> labelsDefaultGetWithHttpInfo() async {
    // ignore: prefer_const_declarations
    final path = r'/labels/default';

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

  /// Get default labels
  ///
  /// Retrieves a list of default system labels available to all users
  Future<List<LabelLabelModel>?> labelsDefaultGet() async {
    final response = await labelsDefaultGetWithHttpInfo();
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      final responseBody = await _decodeBodyBytes(response);
      return (await apiClient.deserializeAsync(responseBody, 'List<LabelLabelModel>') as List)
        .cast<LabelLabelModel>()
        .toList(growable: false);

    }
    return null;
  }

  /// Get all labels
  ///
  /// Retrieve a paginated list of labels with optional filtering by owner type and search text
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] search:
  ///   Search text to filter labels by name
  ///
  /// * [int] limit:
  ///   Maximum number of items to return (default: 10)
  ///
  /// * [int] offset:
  ///   Number of items to skip for pagination (default: 0)
  Future<Response> labelsGetWithHttpInfo({ String? search, int? limit, int? offset, }) async {
    // ignore: prefer_const_declarations
    final path = r'/labels';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    if (search != null) {
      queryParams.addAll(_queryParams('', 'search', search));
    }
    if (limit != null) {
      queryParams.addAll(_queryParams('', 'limit', limit));
    }
    if (offset != null) {
      queryParams.addAll(_queryParams('', 'offset', offset));
    }

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

  /// Get all labels
  ///
  /// Retrieve a paginated list of labels with optional filtering by owner type and search text
  ///
  /// Parameters:
  ///
  /// * [String] search:
  ///   Search text to filter labels by name
  ///
  /// * [int] limit:
  ///   Maximum number of items to return (default: 10)
  ///
  /// * [int] offset:
  ///   Number of items to skip for pagination (default: 0)
  Future<DtoPaginatedResponseModelLabelLabelModel?> labelsGet({ String? search, int? limit, int? offset, }) async {
    final response = await labelsGetWithHttpInfo( search: search, limit: limit, offset: offset, );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'DtoPaginatedResponseModelLabelLabelModel',) as DtoPaginatedResponseModelLabelLabelModel;
    
    }
    return null;
  }

  /// Delete label by ID
  ///
  /// Permanently delete a label by its unique identifier
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Label ID (UUID format)
  Future<Response> labelsIdDeleteWithHttpInfo(String id,) async {
    // ignore: prefer_const_declarations
    final path = r'/labels/{id}'
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

  /// Delete label by ID
  ///
  /// Permanently delete a label by its unique identifier
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Label ID (UUID format)
  Future<void> labelsIdDelete(String id,) async {
    final response = await labelsIdDeleteWithHttpInfo(id,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Get label by ID
  ///
  /// Retrieve a single label by its unique identifier
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Label ID (UUID format)
  Future<Response> labelsIdGetWithHttpInfo(String id,) async {
    // ignore: prefer_const_declarations
    final path = r'/labels/{id}'
      .replaceAll('{id}', id);

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

  /// Get label by ID
  ///
  /// Retrieve a single label by its unique identifier
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Label ID (UUID format)
  Future<LabelLabelModel?> labelsIdGet(String id,) async {
    final response = await labelsIdGetWithHttpInfo(id,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'LabelLabelModel',) as LabelLabelModel;
    
    }
    return null;
  }

  /// Update label by ID
  ///
  /// Update an existing label's name and color by its unique identifier
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Label ID (UUID format)
  ///
  /// * [LabelUpdateLabelModel] labelUpdateLabelModel (required):
  ///   Updated label data
  Future<Response> labelsIdPutWithHttpInfo(String id, LabelUpdateLabelModel labelUpdateLabelModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/labels/{id}'
      .replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody = labelUpdateLabelModel;

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

  /// Update label by ID
  ///
  /// Update an existing label's name and color by its unique identifier
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Label ID (UUID format)
  ///
  /// * [LabelUpdateLabelModel] labelUpdateLabelModel (required):
  ///   Updated label data
  Future<LabelLabelModel?> labelsIdPut(String id, LabelUpdateLabelModel labelUpdateLabelModel,) async {
    final response = await labelsIdPutWithHttpInfo(id, labelUpdateLabelModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'LabelLabelModel',) as LabelLabelModel;
    
    }
    return null;
  }

  /// Create a new label
  ///
  /// Create a new label with specified name, color, and owner information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [LabelCreateLabelModel] labelCreateLabelModel (required):
  ///   Label creation data
  Future<Response> labelsPostWithHttpInfo(LabelCreateLabelModel labelCreateLabelModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/labels';

    // ignore: prefer_final_locals
    Object? postBody = labelCreateLabelModel;

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

  /// Create a new label
  ///
  /// Create a new label with specified name, color, and owner information
  ///
  /// Parameters:
  ///
  /// * [LabelCreateLabelModel] labelCreateLabelModel (required):
  ///   Label creation data
  Future<LabelLabelModel?> labelsPost(LabelCreateLabelModel labelCreateLabelModel,) async {
    final response = await labelsPostWithHttpInfo(labelCreateLabelModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'LabelLabelModel',) as LabelLabelModel;
    
    }
    return null;
  }
}
