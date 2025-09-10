# .SubscriptionsApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**subscriptionsGet**](SubscriptionsApi.md#subscriptionsGet) | **GET** /subscriptions | Get all subscriptions
[**subscriptionsPatch**](SubscriptionsApi.md#subscriptionsPatch) | **PATCH** /subscriptions | Patch subscription
[**subscriptionsPost**](SubscriptionsApi.md#subscriptionsPost) | **POST** /subscriptions | Create a new subscription
[**subscriptionsSubscriptionIdDelete**](SubscriptionsApi.md#subscriptionsSubscriptionIdDelete) | **DELETE** /subscriptions/{subscriptionId} | Delete subscription by ID
[**subscriptionsSubscriptionIdGet**](SubscriptionsApi.md#subscriptionsSubscriptionIdGet) | **GET** /subscriptions/{subscriptionId} | Get subscription by ID
[**subscriptionsSubscriptionIdPut**](SubscriptionsApi.md#subscriptionsSubscriptionIdPut) | **PUT** /subscriptions/{subscriptionId} | Update subscription by ID
[**subscriptionsSummaryGet**](SubscriptionsApi.md#subscriptionsSummaryGet) | **GET** /subscriptions/summary | Get subscription summary


# **subscriptionsGet**
> DtoPaginatedResponseModelSubscriptionSubscriptionModel subscriptionsGet()

Retrieve a paginated list of all subscriptions for the authenticated user

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsGetRequest = {
    // Search text (optional)
  search: "search_example",
    // Filter by recurrency types (optional)
  recurrencies: [
    "recurrencies_example",
  ],
    // Filter by start date (RFC3339) (optional)
  fromDate: "from_date_example",
    // Filter by end date (RFC3339) (optional)
  toDate: "to_date_example",
    // Filter by user IDs (optional)
  users: [
    "users_example",
  ],
    // Include inactive subscriptions (optional)
  withInactive: true,
    // Filter by provider IDs (optional)
  providers: [
    "providers_example",
  ],
    // Number of items per page (default: 10) (optional)
  limit: 1,
    // Page number (default: 0) (optional)
  offset: 1,
};

const data = await apiInstance.subscriptionsGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **search** | [**string**] | Search text | (optional) defaults to undefined
 **recurrencies** | **Array&lt;string&gt;** | Filter by recurrency types | (optional) defaults to undefined
 **fromDate** | [**string**] | Filter by start date (RFC3339) | (optional) defaults to undefined
 **toDate** | [**string**] | Filter by end date (RFC3339) | (optional) defaults to undefined
 **users** | **Array&lt;string&gt;** | Filter by user IDs | (optional) defaults to undefined
 **withInactive** | [**boolean**] | Include inactive subscriptions | (optional) defaults to undefined
 **providers** | **Array&lt;string&gt;** | Filter by provider IDs | (optional) defaults to undefined
 **limit** | [**number**] | Number of items per page (default: 10) | (optional) defaults to undefined
 **offset** | [**number**] | Page number (default: 0) | (optional) defaults to undefined


### Return type

**DtoPaginatedResponseModelSubscriptionSubscriptionModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Paginated list of subscriptions |  -  |
**400** | Bad Request - Invalid query parameters |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **subscriptionsPatch**
> SubscriptionSubscriptionModel subscriptionsPatch(subscriptionPatchSubscriptionModel)

Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsPatchRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsPatchRequest = {
    // Complete subscription data
  subscriptionPatchSubscriptionModel: {
    customPrice: {
      currency: "USD",
      source: ,
      value: 100,
    },
    customRecurrency: 1,
    endDate: new Date('1970-01-01T00:00:00.00Z'),
    freeTrial: {
      endDate: new Date('1970-01-01T00:00:00.00Z'),
      startDate: new Date('1970-01-01T00:00:00.00Z'),
    },
    friendlyName: "friendlyName_example",
    id: "id_example",
    labels: [
      "labels_example",
    ],
    owner: {
      familyId: "familyId_example",
      type: "personal",
    },
    payer: {
      familyId: "123e4567-e89b-12d3-a456-426614174000",
      memberId: "123e4567-e89b-12d3-a456-426614174001",
      type: "family_member",
    },
    planId: "planId_example",
    priceId: "priceId_example",
    providerId: "providerId_example",
    recurrency: "recurrency_example",
    serviceUsers: [
      "serviceUsers_example",
    ],
    startDate: new Date('1970-01-01T00:00:00.00Z'),
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.subscriptionsPatch(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionPatchSubscriptionModel** | **SubscriptionPatchSubscriptionModel**| Complete subscription data |


### Return type

**SubscriptionSubscriptionModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated subscription |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Subscription not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **subscriptionsPost**
> SubscriptionSubscriptionModel subscriptionsPost(subscriptionCreateSubscriptionModel)

Create a new subscription with provider, plan, pricing, and payment information

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsPostRequest = {
    // Subscription creation data
  subscriptionCreateSubscriptionModel: {
    createdAt: "createdAt_example",
    customPrice: {
      currency: "USD",
      source: ,
      value: 100,
    },
    customRecurrency: 1,
    endDate: new Date('1970-01-01T00:00:00.00Z'),
    freeTrial: {
      endDate: new Date('1970-01-01T00:00:00.00Z'),
      startDate: new Date('1970-01-01T00:00:00.00Z'),
    },
    friendlyName: "friendlyName_example",
    id: "id_example",
    labels: [
      "labels_example",
    ],
    owner: {
      familyId: "familyId_example",
      type: "personal",
    },
    payer: {
      familyId: "123e4567-e89b-12d3-a456-426614174000",
      memberId: "123e4567-e89b-12d3-a456-426614174001",
      type: "family_member",
    },
    planId: "planId_example",
    priceId: "priceId_example",
    providerId: "providerId_example",
    recurrency: "recurrency_example",
    serviceUsers: [
      "serviceUsers_example",
    ],
    startDate: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.subscriptionsPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionCreateSubscriptionModel** | **SubscriptionCreateSubscriptionModel**| Subscription creation data |


### Return type

**SubscriptionSubscriptionModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully created subscription |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **subscriptionsSubscriptionIdDelete**
> void subscriptionsSubscriptionIdDelete()

Permanently delete an existing subscription

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsSubscriptionIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsSubscriptionIdDeleteRequest = {
    // Subscription ID (UUID format)
  subscriptionId: "subscriptionId_example",
};

const data = await apiInstance.subscriptionsSubscriptionIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | [**string**] | Subscription ID (UUID format) | defaults to undefined


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
**204** | No Content - Subscription successfully deleted |  -  |
**400** | Bad Request - Invalid subscription ID format |  -  |
**404** | Subscription not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **subscriptionsSubscriptionIdGet**
> SubscriptionSubscriptionModel subscriptionsSubscriptionIdGet()

Retrieve a single subscription with all its details including provider, plan, and pricing information

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsSubscriptionIdGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsSubscriptionIdGetRequest = {
    // Subscription ID (UUID format)
  subscriptionId: "subscriptionId_example",
};

const data = await apiInstance.subscriptionsSubscriptionIdGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionId** | [**string**] | Subscription ID (UUID format) | defaults to undefined


### Return type

**SubscriptionSubscriptionModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully retrieved subscription |  -  |
**400** | Bad Request - Invalid subscription ID format |  -  |
**404** | Subscription not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **subscriptionsSubscriptionIdPut**
> SubscriptionSubscriptionModel subscriptionsSubscriptionIdPut(subscriptionUpdateSubscriptionModel)

Update an existing subscription\'s details including provider, plan, pricing, and payment information

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsSubscriptionIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsSubscriptionIdPutRequest = {
    // Subscription ID (UUID format)
  subscriptionId: "subscriptionId_example",
    // Updated subscription data
  subscriptionUpdateSubscriptionModel: {
    customPrice: {
      currency: "USD",
      source: ,
      value: 100,
    },
    customRecurrency: 1,
    endDate: new Date('1970-01-01T00:00:00.00Z'),
    freeTrial: {
      endDate: new Date('1970-01-01T00:00:00.00Z'),
      startDate: new Date('1970-01-01T00:00:00.00Z'),
    },
    friendlyName: "friendlyName_example",
    labels: [
      "labels_example",
    ],
    owner: {
      familyId: "familyId_example",
      type: "personal",
    },
    payer: {
      familyId: "123e4567-e89b-12d3-a456-426614174000",
      memberId: "123e4567-e89b-12d3-a456-426614174001",
      type: "family_member",
    },
    planId: "planId_example",
    priceId: "priceId_example",
    providerId: "providerId_example",
    recurrency: "recurrency_example",
    serviceUsers: [
      "serviceUsers_example",
    ],
    startDate: new Date('1970-01-01T00:00:00.00Z'),
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.subscriptionsSubscriptionIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subscriptionUpdateSubscriptionModel** | **SubscriptionUpdateSubscriptionModel**| Updated subscription data |
 **subscriptionId** | [**string**] | Subscription ID (UUID format) | defaults to undefined


### Return type

**SubscriptionSubscriptionModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated subscription |  -  |
**400** | Bad Request - Invalid input data or subscription ID |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Subscription not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **subscriptionsSummaryGet**
> SubscriptionSubscriptionSummaryResponse subscriptionsSummaryGet()

Returns summary information about subscriptions including total costs and upcoming renewals

### Example


```typescript
import { createConfiguration, SubscriptionsApi } from '';
import type { SubscriptionsApiSubscriptionsSummaryGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new SubscriptionsApi(configuration);

const request: SubscriptionsApiSubscriptionsSummaryGetRequest = {
    // Number of top providers to return
  topProviders: 1,
    // Number of top labels to return
  topLabels: 1,
    // Number of upcoming renewals to return
  upcomingRenewals: 1,
    // Include monthly total costs
  totalMonthly: true,
    // Include yearly total costs
  totalYearly: true,
};

const data = await apiInstance.subscriptionsSummaryGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **topProviders** | [**number**] | Number of top providers to return | defaults to undefined
 **topLabels** | [**number**] | Number of top labels to return | defaults to undefined
 **upcomingRenewals** | [**number**] | Number of upcoming renewals to return | defaults to undefined
 **totalMonthly** | [**boolean**] | Include monthly total costs | defaults to undefined
 **totalYearly** | [**boolean**] | Include yearly total costs | defaults to undefined


### Return type

**SubscriptionSubscriptionSummaryResponse**

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

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


