//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class SubscriptionsApi {
  SubscriptionsApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Get all subscriptions
  ///
  /// Retrieve a paginated list of all subscriptions for the authenticated user
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] search:
  ///   Search text
  ///
  /// * [List<String>] recurrencies:
  ///   Filter by recurrency types
  ///
  /// * [String] fromDate:
  ///   Filter by start date (RFC3339)
  ///
  /// * [String] toDate:
  ///   Filter by end date (RFC3339)
  ///
  /// * [List<String>] users:
  ///   Filter by user IDs
  ///
  /// * [bool] withInactive:
  ///   Include inactive subscriptions
  ///
  /// * [List<String>] providers:
  ///   Filter by provider IDs
  ///
  /// * [int] limit:
  ///   Number of items per page (default: 10)
  ///
  /// * [int] offset:
  ///   Page number (default: 0)
  Future<Response> subscriptionsGetWithHttpInfo({ String? search, List<String>? recurrencies, String? fromDate, String? toDate, List<String>? users, bool? withInactive, List<String>? providers, int? limit, int? offset, }) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    if (search != null) {
      queryParams.addAll(_queryParams('', 'search', search));
    }
    if (recurrencies != null) {
      queryParams.addAll(_queryParams('multi', 'recurrencies', recurrencies));
    }
    if (fromDate != null) {
      queryParams.addAll(_queryParams('', 'from_date', fromDate));
    }
    if (toDate != null) {
      queryParams.addAll(_queryParams('', 'to_date', toDate));
    }
    if (users != null) {
      queryParams.addAll(_queryParams('multi', 'users', users));
    }
    if (withInactive != null) {
      queryParams.addAll(_queryParams('', 'with_inactive', withInactive));
    }
    if (providers != null) {
      queryParams.addAll(_queryParams('multi', 'providers', providers));
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

  /// Get all subscriptions
  ///
  /// Retrieve a paginated list of all subscriptions for the authenticated user
  ///
  /// Parameters:
  ///
  /// * [String] search:
  ///   Search text
  ///
  /// * [List<String>] recurrencies:
  ///   Filter by recurrency types
  ///
  /// * [String] fromDate:
  ///   Filter by start date (RFC3339)
  ///
  /// * [String] toDate:
  ///   Filter by end date (RFC3339)
  ///
  /// * [List<String>] users:
  ///   Filter by user IDs
  ///
  /// * [bool] withInactive:
  ///   Include inactive subscriptions
  ///
  /// * [List<String>] providers:
  ///   Filter by provider IDs
  ///
  /// * [int] limit:
  ///   Number of items per page (default: 10)
  ///
  /// * [int] offset:
  ///   Page number (default: 0)
  Future<DtoPaginatedResponseModelSubscriptionSubscriptionModel?> subscriptionsGet({ String? search, List<String>? recurrencies, String? fromDate, String? toDate, List<String>? users, bool? withInactive, List<String>? providers, int? limit, int? offset, }) async {
    final response = await subscriptionsGetWithHttpInfo( search: search, recurrencies: recurrencies, fromDate: fromDate, toDate: toDate, users: users, withInactive: withInactive, providers: providers, limit: limit, offset: offset, );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'DtoPaginatedResponseModelSubscriptionSubscriptionModel',) as DtoPaginatedResponseModelSubscriptionSubscriptionModel;
    
    }
    return null;
  }

  /// Patch subscription
  ///
  /// Update or create a subscription with complete details. If subscription doesn't exist, it will be created.
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [SubscriptionPatchSubscriptionModel] subscriptionPatchSubscriptionModel (required):
  ///   Complete subscription data
  Future<Response> subscriptionsPatchWithHttpInfo(SubscriptionPatchSubscriptionModel subscriptionPatchSubscriptionModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions';

    // ignore: prefer_final_locals
    Object? postBody = subscriptionPatchSubscriptionModel;

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

  /// Patch subscription
  ///
  /// Update or create a subscription with complete details. If subscription doesn't exist, it will be created.
  ///
  /// Parameters:
  ///
  /// * [SubscriptionPatchSubscriptionModel] subscriptionPatchSubscriptionModel (required):
  ///   Complete subscription data
  Future<SubscriptionSubscriptionModel?> subscriptionsPatch(SubscriptionPatchSubscriptionModel subscriptionPatchSubscriptionModel,) async {
    final response = await subscriptionsPatchWithHttpInfo(subscriptionPatchSubscriptionModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'SubscriptionSubscriptionModel',) as SubscriptionSubscriptionModel;
    
    }
    return null;
  }

  /// Create a new subscription
  ///
  /// Create a new subscription with provider, plan, pricing, and payment information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [SubscriptionCreateSubscriptionModel] subscriptionCreateSubscriptionModel (required):
  ///   Subscription creation data
  Future<Response> subscriptionsPostWithHttpInfo(SubscriptionCreateSubscriptionModel subscriptionCreateSubscriptionModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions';

    // ignore: prefer_final_locals
    Object? postBody = subscriptionCreateSubscriptionModel;

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

  /// Create a new subscription
  ///
  /// Create a new subscription with provider, plan, pricing, and payment information
  ///
  /// Parameters:
  ///
  /// * [SubscriptionCreateSubscriptionModel] subscriptionCreateSubscriptionModel (required):
  ///   Subscription creation data
  Future<SubscriptionSubscriptionModel?> subscriptionsPost(SubscriptionCreateSubscriptionModel subscriptionCreateSubscriptionModel,) async {
    final response = await subscriptionsPostWithHttpInfo(subscriptionCreateSubscriptionModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'SubscriptionSubscriptionModel',) as SubscriptionSubscriptionModel;
    
    }
    return null;
  }

  /// Delete subscription by ID
  ///
  /// Permanently delete an existing subscription
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] subscriptionId (required):
  ///   Subscription ID (UUID format)
  Future<Response> subscriptionsSubscriptionIdDeleteWithHttpInfo(String subscriptionId,) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions/{subscriptionId}'
      .replaceAll('{subscriptionId}', subscriptionId);

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

  /// Delete subscription by ID
  ///
  /// Permanently delete an existing subscription
  ///
  /// Parameters:
  ///
  /// * [String] subscriptionId (required):
  ///   Subscription ID (UUID format)
  Future<void> subscriptionsSubscriptionIdDelete(String subscriptionId,) async {
    final response = await subscriptionsSubscriptionIdDeleteWithHttpInfo(subscriptionId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Get subscription by ID
  ///
  /// Retrieve a single subscription with all its details including provider, plan, and pricing information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] subscriptionId (required):
  ///   Subscription ID (UUID format)
  Future<Response> subscriptionsSubscriptionIdGetWithHttpInfo(String subscriptionId,) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions/{subscriptionId}'
      .replaceAll('{subscriptionId}', subscriptionId);

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

  /// Get subscription by ID
  ///
  /// Retrieve a single subscription with all its details including provider, plan, and pricing information
  ///
  /// Parameters:
  ///
  /// * [String] subscriptionId (required):
  ///   Subscription ID (UUID format)
  Future<SubscriptionSubscriptionModel?> subscriptionsSubscriptionIdGet(String subscriptionId,) async {
    final response = await subscriptionsSubscriptionIdGetWithHttpInfo(subscriptionId,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'SubscriptionSubscriptionModel',) as SubscriptionSubscriptionModel;
    
    }
    return null;
  }

  /// Update subscription by ID
  ///
  /// Update an existing subscription's details including provider, plan, pricing, and payment information
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] subscriptionId (required):
  ///   Subscription ID (UUID format)
  ///
  /// * [SubscriptionUpdateSubscriptionModel] subscriptionUpdateSubscriptionModel (required):
  ///   Updated subscription data
  Future<Response> subscriptionsSubscriptionIdPutWithHttpInfo(String subscriptionId, SubscriptionUpdateSubscriptionModel subscriptionUpdateSubscriptionModel,) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions/{subscriptionId}'
      .replaceAll('{subscriptionId}', subscriptionId);

    // ignore: prefer_final_locals
    Object? postBody = subscriptionUpdateSubscriptionModel;

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

  /// Update subscription by ID
  ///
  /// Update an existing subscription's details including provider, plan, pricing, and payment information
  ///
  /// Parameters:
  ///
  /// * [String] subscriptionId (required):
  ///   Subscription ID (UUID format)
  ///
  /// * [SubscriptionUpdateSubscriptionModel] subscriptionUpdateSubscriptionModel (required):
  ///   Updated subscription data
  Future<SubscriptionSubscriptionModel?> subscriptionsSubscriptionIdPut(String subscriptionId, SubscriptionUpdateSubscriptionModel subscriptionUpdateSubscriptionModel,) async {
    final response = await subscriptionsSubscriptionIdPutWithHttpInfo(subscriptionId, subscriptionUpdateSubscriptionModel,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'SubscriptionSubscriptionModel',) as SubscriptionSubscriptionModel;
    
    }
    return null;
  }

  /// Get subscription summary
  ///
  /// Returns summary information about subscriptions including total costs and upcoming renewals
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [int] topProviders (required):
  ///   Number of top providers to return
  ///
  /// * [int] topLabels (required):
  ///   Number of top labels to return
  ///
  /// * [int] upcomingRenewals (required):
  ///   Number of upcoming renewals to return
  ///
  /// * [bool] totalMonthly (required):
  ///   Include monthly total costs
  ///
  /// * [bool] totalYearly (required):
  ///   Include yearly total costs
  Future<Response> subscriptionsSummaryGetWithHttpInfo(int topProviders, int topLabels, int upcomingRenewals, bool totalMonthly, bool totalYearly,) async {
    // ignore: prefer_const_declarations
    final path = r'/subscriptions/summary';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

      queryParams.addAll(_queryParams('', 'top_providers', topProviders));
      queryParams.addAll(_queryParams('', 'top_labels', topLabels));
      queryParams.addAll(_queryParams('', 'upcoming_renewals', upcomingRenewals));
      queryParams.addAll(_queryParams('', 'total_monthly', totalMonthly));
      queryParams.addAll(_queryParams('', 'total_yearly', totalYearly));

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

  /// Get subscription summary
  ///
  /// Returns summary information about subscriptions including total costs and upcoming renewals
  ///
  /// Parameters:
  ///
  /// * [int] topProviders (required):
  ///   Number of top providers to return
  ///
  /// * [int] topLabels (required):
  ///   Number of top labels to return
  ///
  /// * [int] upcomingRenewals (required):
  ///   Number of upcoming renewals to return
  ///
  /// * [bool] totalMonthly (required):
  ///   Include monthly total costs
  ///
  /// * [bool] totalYearly (required):
  ///   Include yearly total costs
  Future<SubscriptionSubscriptionSummaryResponse?> subscriptionsSummaryGet(int topProviders, int topLabels, int upcomingRenewals, bool totalMonthly, bool totalYearly,) async {
    final response = await subscriptionsSummaryGetWithHttpInfo(topProviders, topLabels, upcomingRenewals, totalMonthly, totalYearly,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'SubscriptionSubscriptionSummaryResponse',) as SubscriptionSubscriptionSummaryResponse;
    
    }
    return null;
  }
}
