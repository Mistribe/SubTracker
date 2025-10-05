//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class CurrenciesApi {
  CurrenciesApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Get Currency Rates
  ///
  /// Get exchange rates for all currencies at a specific date
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] date:
  ///   Conversion date in RFC3339 format (default: current time)
  Future<Response> currenciesRatesGetWithHttpInfo({ String? date, }) async {
    // ignore: prefer_const_declarations
    final path = r'/currencies/rates';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    if (date != null) {
      queryParams.addAll(_queryParams('', 'date', date));
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

  /// Get Currency Rates
  ///
  /// Get exchange rates for all currencies at a specific date
  ///
  /// Parameters:
  ///
  /// * [String] date:
  ///   Conversion date in RFC3339 format (default: current time)
  Future<CurrencyCurrencyRatesModel?> currenciesRatesGet({ String? date, }) async {
    final response = await currenciesRatesGetWithHttpInfo( date: date, );
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'CurrencyCurrencyRatesModel',) as CurrencyCurrencyRatesModel;
    
    }
    return null;
  }

  /// Get Supported Currencies
  ///
  /// get details of all supported currencies
  ///
  /// Note: This method returns the HTTP [Response].
  Future<Response> currenciesSupportedGetWithHttpInfo() async {
    // ignore: prefer_const_declarations
    final path = r'/currencies/supported';

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

  /// Get Supported Currencies
  ///
  /// get details of all supported currencies
  Future<List<String>?> currenciesSupportedGet() async {
    final response = await currenciesSupportedGetWithHttpInfo();
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      final responseBody = await _decodeBodyBytes(response);
      return (await apiClient.deserializeAsync(responseBody, 'List<String>') as List)
        .cast<String>()
        .toList(growable: false);

    }
    return null;
  }
}
