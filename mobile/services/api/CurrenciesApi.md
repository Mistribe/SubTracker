# .CurrenciesApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**currenciesRatesGet**](CurrenciesApi.md#currenciesRatesGet) | **GET** /currencies/rates | Get Currency Rates
[**currenciesSupportedGet**](CurrenciesApi.md#currenciesSupportedGet) | **GET** /currencies/supported | Get Supported Currencies


# **currenciesRatesGet**
> CurrencyCurrencyRatesModel currenciesRatesGet()

Get exchange rates for all currencies at a specific date

### Example


```typescript
import { createConfiguration, CurrenciesApi } from '';
import type { CurrenciesApiCurrenciesRatesGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new CurrenciesApi(configuration);

const request: CurrenciesApiCurrenciesRatesGetRequest = {
    // Conversion date in RFC3339 format (default: current time) (optional)
  date: "date_example",
};

const data = await apiInstance.currenciesRatesGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **date** | [**string**] | Conversion date in RFC3339 format (default: current time) | (optional) defaults to undefined


### Return type

**CurrencyCurrencyRatesModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **currenciesSupportedGet**
> Array<string> currenciesSupportedGet()

get details of all supported currencies

### Example


```typescript
import { createConfiguration, CurrenciesApi } from '';

const configuration = createConfiguration();
const apiInstance = new CurrenciesApi(configuration);

const request = {};

const data = await apiInstance.currenciesSupportedGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**Array<string>**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | currencies |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


