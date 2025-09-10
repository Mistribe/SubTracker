# .ProvidersApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**providersGet**](ProvidersApi.md#providersGet) | **GET** /providers | Get all providers
[**providersPost**](ProvidersApi.md#providersPost) | **POST** /providers | Create a new provider
[**providersProviderIdDelete**](ProvidersApi.md#providersProviderIdDelete) | **DELETE** /providers/{providerId} | Delete provider by ID
[**providersProviderIdGet**](ProvidersApi.md#providersProviderIdGet) | **GET** /providers/{providerId} | Get provider by ID
[**providersProviderIdPlansPlanIdDelete**](ProvidersApi.md#providersProviderIdPlansPlanIdDelete) | **DELETE** /providers/{providerId}/plans/{planId} | Delete provider plan by ID
[**providersProviderIdPlansPlanIdPricesPost**](ProvidersApi.md#providersProviderIdPlansPlanIdPricesPost) | **POST** /providers/{providerId}/plans/{planId}/prices | Create a new provider price
[**providersProviderIdPlansPlanIdPricesPriceIdDelete**](ProvidersApi.md#providersProviderIdPlansPlanIdPricesPriceIdDelete) | **DELETE** /providers/{providerId}/plans/{planId}/prices/{priceId} | Delete provider price by ID
[**providersProviderIdPlansPlanIdPricesPriceIdPut**](ProvidersApi.md#providersProviderIdPlansPlanIdPricesPriceIdPut) | **PUT** /providers/{providerId}/plans/{planId}/prices/{priceId} | Update provider price by ID
[**providersProviderIdPlansPlanIdPut**](ProvidersApi.md#providersProviderIdPlansPlanIdPut) | **PUT** /providers/{providerId}/plans/{planId} | Update provider plan by ID
[**providersProviderIdPlansPost**](ProvidersApi.md#providersProviderIdPlansPost) | **POST** /providers/{providerId}/plans | Create a new provider plan
[**providersProviderIdPut**](ProvidersApi.md#providersProviderIdPut) | **PUT** /providers/{providerId} | Update provider by ID


# **providersGet**
> DtoPaginatedResponseModelProviderProviderModel providersGet()

Retrieve a paginated list of all providers with their plans and prices

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersGetRequest = {
    // Search term (optional)
  search: "search_example",
    // Offset (default: 0) (optional)
  offset: 1,
    // Limit per request (default: 10) (optional)
  limit: 1,
};

const data = await apiInstance.providersGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **search** | [**string**] | Search term | (optional) defaults to undefined
 **offset** | [**number**] | Offset (default: 0) | (optional) defaults to undefined
 **limit** | [**number**] | Limit per request (default: 10) | (optional) defaults to undefined


### Return type

**DtoPaginatedResponseModelProviderProviderModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Paginated list of providers |  -  |
**400** | Bad Request - Invalid query parameters |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersPost**
> ProviderProviderModel providersPost(providerCreateProviderModel)

Create a new service provider with labels and owner information

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersPostRequest = {
    // Provider creation data
  providerCreateProviderModel: {
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    description: "description_example",
    iconUrl: "iconUrl_example",
    id: "id_example",
    labels: [
      "labels_example",
    ],
    name: "name_example",
    owner: {
      familyId: "familyId_example",
      type: "personal",
    },
    pricingPageUrl: "pricingPageUrl_example",
    url: "url_example",
  },
};

const data = await apiInstance.providersPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerCreateProviderModel** | **ProviderCreateProviderModel**| Provider creation data |


### Return type

**ProviderProviderModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully created provider |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdDelete**
> void providersProviderIdDelete()

Permanently delete a provider and all its associated plans and prices

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdDeleteRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
};

const data = await apiInstance.providersProviderIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined


### Return type

**void**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content - Provider successfully deleted |  -  |
**400** | Bad Request - Invalid provider ID format |  -  |
**404** | Provider not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdGet**
> ProviderProviderModel providersProviderIdGet()

Retrieve a single provider with all its plans and prices by ID

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdGetRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
};

const data = await apiInstance.providersProviderIdGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined


### Return type

**ProviderProviderModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully retrieved provider |  -  |
**400** | Bad Request - Invalid provider ID format |  -  |
**404** | Provider not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPlansPlanIdDelete**
> void providersProviderIdPlansPlanIdDelete()

Permanently delete a provider plan and all its associated prices

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPlansPlanIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPlansPlanIdDeleteRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Plan ID (UUID format)
  planId: "planId_example",
};

const data = await apiInstance.providersProviderIdPlansPlanIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined
 **planId** | [**string**] | Plan ID (UUID format) | defaults to undefined


### Return type

**void**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content - Plan successfully deleted |  -  |
**400** | Bad Request - Invalid ID format |  -  |
**404** | Provider or plan not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPlansPlanIdPricesPost**
> ProviderPriceModel providersProviderIdPlansPlanIdPricesPost(providerCreatePriceModel)

Create a new pricing option for a specific provider plan

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPlansPlanIdPricesPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPlansPlanIdPricesPostRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Plan ID (UUID format)
  planId: "planId_example",
    // Price creation data
  providerCreatePriceModel: {
    amount: 3.14,
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    currency: "currency_example",
    endDate: new Date('1970-01-01T00:00:00.00Z'),
    id: "id_example",
    startDate: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.providersProviderIdPlansPlanIdPricesPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerCreatePriceModel** | **ProviderCreatePriceModel**| Price creation data |
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined
 **planId** | [**string**] | Plan ID (UUID format) | defaults to undefined


### Return type

**ProviderPriceModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully created price |  -  |
**400** | Bad Request - Invalid input data or IDs |  -  |
**404** | Provider or plan not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPlansPlanIdPricesPriceIdDelete**
> void providersProviderIdPlansPlanIdPricesPriceIdDelete()

Permanently delete a specific price from a provider plan

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdDeleteRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Plan ID (UUID format)
  planId: "planId_example",
    // Price ID (UUID format)
  priceId: "priceId_example",
};

const data = await apiInstance.providersProviderIdPlansPlanIdPricesPriceIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined
 **planId** | [**string**] | Plan ID (UUID format) | defaults to undefined
 **priceId** | [**string**] | Price ID (UUID format) | defaults to undefined


### Return type

**void**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content - Price successfully deleted |  -  |
**400** | Bad Request - Invalid ID format |  -  |
**404** | Provider, plan, or price not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPlansPlanIdPricesPriceIdPut**
> ProviderPriceModel providersProviderIdPlansPlanIdPricesPriceIdPut(providerUpdatePriceModel)

Update an existing price for a specific provider plan

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdPutRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Plan ID (UUID format)
  planId: "planId_example",
    // Price ID (UUID format)
  priceId: "priceId_example",
    // Updated price data
  providerUpdatePriceModel: {
    amount: 3.14,
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    currency: "currency_example",
    endDate: new Date('1970-01-01T00:00:00.00Z'),
    startDate: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.providersProviderIdPlansPlanIdPricesPriceIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerUpdatePriceModel** | **ProviderUpdatePriceModel**| Updated price data |
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined
 **planId** | [**string**] | Plan ID (UUID format) | defaults to undefined
 **priceId** | [**string**] | Price ID (UUID format) | defaults to undefined


### Return type

**ProviderPriceModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated price |  -  |
**400** | Bad Request - Invalid input data or IDs |  -  |
**404** | Provider, plan, or price not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPlansPlanIdPut**
> ProviderPlanModel providersProviderIdPlansPlanIdPut(providerUpdatePlanModel)

Update an existing provider plan\'s information

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPlansPlanIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPlansPlanIdPutRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Plan ID (UUID format)
  planId: "planId_example",
    // Updated plan data
  providerUpdatePlanModel: {
    description: "description_example",
    name: "name_example",
    updateAt: "updateAt_example",
  },
};

const data = await apiInstance.providersProviderIdPlansPlanIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerUpdatePlanModel** | **ProviderUpdatePlanModel**| Updated plan data |
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined
 **planId** | [**string**] | Plan ID (UUID format) | defaults to undefined


### Return type

**ProviderPlanModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated plan |  -  |
**400** | Bad Request - Invalid input data or IDs |  -  |
**404** | Provider or plan not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPlansPost**
> ProviderPlanModel providersProviderIdPlansPost(providerCreatePlanModel)

Create a new subscription plan for an existing provider

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPlansPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPlansPostRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Plan creation data
  providerCreatePlanModel: {
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    description: "description_example",
    id: "id_example",
    name: "name_example",
  },
};

const data = await apiInstance.providersProviderIdPlansPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerCreatePlanModel** | **ProviderCreatePlanModel**| Plan creation data |
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined


### Return type

**ProviderPlanModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully created plan |  -  |
**400** | Bad Request - Invalid input data or provider ID |  -  |
**404** | Provider not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **providersProviderIdPut**
> ProviderProviderModel providersProviderIdPut(providerUpdateProviderModel)

Update an existing provider\'s basic information

### Example


```typescript
import { createConfiguration, ProvidersApi } from '';
import type { ProvidersApiProvidersProviderIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new ProvidersApi(configuration);

const request: ProvidersApiProvidersProviderIdPutRequest = {
    // Provider ID (UUID format)
  providerId: "providerId_example",
    // Updated provider data
  providerUpdateProviderModel: {
    description: "description_example",
    iconUrl: "iconUrl_example",
    labels: [
      "labels_example",
    ],
    name: "name_example",
    pricingPageUrl: "pricingPageUrl_example",
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
    url: "url_example",
  },
};

const data = await apiInstance.providersProviderIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **providerUpdateProviderModel** | **ProviderUpdateProviderModel**| Updated provider data |
 **providerId** | [**string**] | Provider ID (UUID format) | defaults to undefined


### Return type

**ProviderProviderModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated provider |  -  |
**400** | Bad Request - Invalid input data or provider ID |  -  |
**404** | Provider not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


