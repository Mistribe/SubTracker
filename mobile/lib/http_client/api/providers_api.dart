//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class ProvidersApi {
  ProvidersApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Get all providers
  ///
  /// Retrieve a paginated list of all providers with their plans and prices
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] search:
  ///   Search term
  ///
  /// * [int] offset:
  ///   Offset (default: 0)
  ///
  /// * [int] limit:
  ///   Limit per request (default: 10)
  Future<Response> providersGetWithHttpInfo({ String? search, int? offset, int? limit, }) async {
    // ignore: prefer_const_declarations
    final path = r'/providers';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    if (search != null) {
      queryParams.addAll(_queryParams('', 'search', search));
    }
    if (offset != null) {
      queryParams.addAll(_queryParams('', 'offset', offset));
    }
    if (limit != null) {
      queryParams.addAll(_queryParams('', 'limit', limit));
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

  /// Get all providers
  ///
  /// Retrieve a paginated list of all providers with their plans and prices
  ///
  /// Parameters:
  ///
  /// * [String] search:
  ///   Search term
  ///
  /// * [int] offset:
  ///   Offset (default: 0)
  ///
  /// * [int] limit:
  ///   Limit per request (default: 10)
  Future<DtoPaginatedResponseModelProviderProviderModel?> providersGet({ String? search, int? offset, int? limit, }) async {
    final response = await providersGetWithHttpInfo( search: search, offset: offset, limit: limit, );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'DtoPaginatedResponseModelProviderProviderModel',) as DtoPaginatedResponseModelProviderProviderModel;
    
    }
    return null;
  }

  /// Create a new provider
  ///
  /// Create a new service provider with labels and owner information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [ProviderCreateProviderModel] providerCreateProviderModel (required):
  ///   Provider creation data
  Future<Response> providersPostWithHttpInfo(ProviderCreateProviderModel providerCreateProviderModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers';

    // ignore: prefer_final_locals
    Object? postBody = providerCreateProviderModel;

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

  /// Create a new provider
  ///
  /// Create a new service provider with labels and owner information
  ///
  /// Parameters:
  ///
  /// * [ProviderCreateProviderModel] providerCreateProviderModel (required):
  ///   Provider creation data
  Future<ProviderProviderModel?> providersPost(ProviderCreateProviderModel providerCreateProviderModel,) async {
    final response = await providersPostWithHttpInfo(providerCreateProviderModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderProviderModel',) as ProviderProviderModel;
    
    }
    return null;
  }

  /// Delete provider by ID
  ///
  /// Permanently delete a provider and all its associated plans and prices
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  Future<Response> providersProviderIdDeleteWithHttpInfo(String providerId,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}'
      .replaceAll('{providerId}', providerId);

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

  /// Delete provider by ID
  ///
  /// Permanently delete a provider and all its associated plans and prices
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  Future<void> providersProviderIdDelete(String providerId,) async {
    final response = await providersProviderIdDeleteWithHttpInfo(providerId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Get provider by ID
  ///
  /// Retrieve a single provider with all its plans and prices by ID
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  Future<Response> providersProviderIdGetWithHttpInfo(String providerId,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}'
      .replaceAll('{providerId}', providerId);

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

  /// Get provider by ID
  ///
  /// Retrieve a single provider with all its plans and prices by ID
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  Future<ProviderProviderModel?> providersProviderIdGet(String providerId,) async {
    final response = await providersProviderIdGetWithHttpInfo(providerId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderProviderModel',) as ProviderProviderModel;
    
    }
    return null;
  }

  /// Delete provider plan by ID
  ///
  /// Permanently delete a provider plan and all its associated prices
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  Future<Response> providersProviderIdPlansPlanIdDeleteWithHttpInfo(String providerId, String planId,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}/plans/{planId}'
      .replaceAll('{providerId}', providerId)
      .replaceAll('{planId}', planId);

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

  /// Delete provider plan by ID
  ///
  /// Permanently delete a provider plan and all its associated prices
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  Future<void> providersProviderIdPlansPlanIdDelete(String providerId, String planId,) async {
    final response = await providersProviderIdPlansPlanIdDeleteWithHttpInfo(providerId, planId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Create a new provider price
  ///
  /// Create a new pricing option for a specific provider plan
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [ProviderCreatePriceModel] providerCreatePriceModel (required):
  ///   Price creation data
  Future<Response> providersProviderIdPlansPlanIdPricesPostWithHttpInfo(String providerId, String planId, ProviderCreatePriceModel providerCreatePriceModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}/plans/{planId}/prices'
      .replaceAll('{providerId}', providerId)
      .replaceAll('{planId}', planId);

    // ignore: prefer_final_locals
    Object? postBody = providerCreatePriceModel;

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

  /// Create a new provider price
  ///
  /// Create a new pricing option for a specific provider plan
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [ProviderCreatePriceModel] providerCreatePriceModel (required):
  ///   Price creation data
  Future<ProviderPriceModel?> providersProviderIdPlansPlanIdPricesPost(String providerId, String planId, ProviderCreatePriceModel providerCreatePriceModel,) async {
    final response = await providersProviderIdPlansPlanIdPricesPostWithHttpInfo(providerId, planId, providerCreatePriceModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderPriceModel',) as ProviderPriceModel;
    
    }
    return null;
  }

  /// Delete provider price by ID
  ///
  /// Permanently delete a specific price from a provider plan
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [String] priceId (required):
  ///   Price ID (UUID format)
  Future<Response> providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(String providerId, String planId, String priceId,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}/plans/{planId}/prices/{priceId}'
      .replaceAll('{providerId}', providerId)
      .replaceAll('{planId}', planId)
      .replaceAll('{priceId}', priceId);

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

  /// Delete provider price by ID
  ///
  /// Permanently delete a specific price from a provider plan
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [String] priceId (required):
  ///   Price ID (UUID format)
  Future<void> providersProviderIdPlansPlanIdPricesPriceIdDelete(String providerId, String planId, String priceId,) async {
    final response = await providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(providerId, planId, priceId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Update provider price by ID
  ///
  /// Update an existing price for a specific provider plan
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [String] priceId (required):
  ///   Price ID (UUID format)
  ///
  /// * [ProviderUpdatePriceModel] providerUpdatePriceModel (required):
  ///   Updated price data
  Future<Response> providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(String providerId, String planId, String priceId, ProviderUpdatePriceModel providerUpdatePriceModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}/plans/{planId}/prices/{priceId}'
      .replaceAll('{providerId}', providerId)
      .replaceAll('{planId}', planId)
      .replaceAll('{priceId}', priceId);

    // ignore: prefer_final_locals
    Object? postBody = providerUpdatePriceModel;

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

  /// Update provider price by ID
  ///
  /// Update an existing price for a specific provider plan
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [String] priceId (required):
  ///   Price ID (UUID format)
  ///
  /// * [ProviderUpdatePriceModel] providerUpdatePriceModel (required):
  ///   Updated price data
  Future<ProviderPriceModel?> providersProviderIdPlansPlanIdPricesPriceIdPut(String providerId, String planId, String priceId, ProviderUpdatePriceModel providerUpdatePriceModel,) async {
    final response = await providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(providerId, planId, priceId, providerUpdatePriceModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderPriceModel',) as ProviderPriceModel;
    
    }
    return null;
  }

  /// Update provider plan by ID
  ///
  /// Update an existing provider plan's information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [ProviderUpdatePlanModel] providerUpdatePlanModel (required):
  ///   Updated plan data
  Future<Response> providersProviderIdPlansPlanIdPutWithHttpInfo(String providerId, String planId, ProviderUpdatePlanModel providerUpdatePlanModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}/plans/{planId}'
      .replaceAll('{providerId}', providerId)
      .replaceAll('{planId}', planId);

    // ignore: prefer_final_locals
    Object? postBody = providerUpdatePlanModel;

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

  /// Update provider plan by ID
  ///
  /// Update an existing provider plan's information
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [String] planId (required):
  ///   Plan ID (UUID format)
  ///
  /// * [ProviderUpdatePlanModel] providerUpdatePlanModel (required):
  ///   Updated plan data
  Future<ProviderPlanModel?> providersProviderIdPlansPlanIdPut(String providerId, String planId, ProviderUpdatePlanModel providerUpdatePlanModel,) async {
    final response = await providersProviderIdPlansPlanIdPutWithHttpInfo(providerId, planId, providerUpdatePlanModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderPlanModel',) as ProviderPlanModel;
    
    }
    return null;
  }

  /// Create a new provider plan
  ///
  /// Create a new subscription plan for an existing provider
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [ProviderCreatePlanModel] providerCreatePlanModel (required):
  ///   Plan creation data
  Future<Response> providersProviderIdPlansPostWithHttpInfo(String providerId, ProviderCreatePlanModel providerCreatePlanModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}/plans'
      .replaceAll('{providerId}', providerId);

    // ignore: prefer_final_locals
    Object? postBody = providerCreatePlanModel;

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

  /// Create a new provider plan
  ///
  /// Create a new subscription plan for an existing provider
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [ProviderCreatePlanModel] providerCreatePlanModel (required):
  ///   Plan creation data
  Future<ProviderPlanModel?> providersProviderIdPlansPost(String providerId, ProviderCreatePlanModel providerCreatePlanModel,) async {
    final response = await providersProviderIdPlansPostWithHttpInfo(providerId, providerCreatePlanModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderPlanModel',) as ProviderPlanModel;
    
    }
    return null;
  }

  /// Update provider by ID
  ///
  /// Update an existing provider's basic information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [ProviderUpdateProviderModel] providerUpdateProviderModel (required):
  ///   Updated provider data
  Future<Response> providersProviderIdPutWithHttpInfo(String providerId, ProviderUpdateProviderModel providerUpdateProviderModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/providers/{providerId}'
      .replaceAll('{providerId}', providerId);

    // ignore: prefer_final_locals
    Object? postBody = providerUpdateProviderModel;

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

  /// Update provider by ID
  ///
  /// Update an existing provider's basic information
  ///
  /// Parameters:
  ///
  /// * [String] providerId (required):
  ///   Provider ID (UUID format)
  ///
  /// * [ProviderUpdateProviderModel] providerUpdateProviderModel (required):
  ///   Updated provider data
  Future<ProviderProviderModel?> providersProviderIdPut(String providerId, ProviderUpdateProviderModel providerUpdateProviderModel,) async {
    final response = await providersProviderIdPutWithHttpInfo(providerId, providerUpdateProviderModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ProviderProviderModel',) as ProviderProviderModel;
    
    }
    return null;
  }
}
