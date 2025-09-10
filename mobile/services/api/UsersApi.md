# .UsersApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**usersDelete**](UsersApi.md#usersDelete) | **DELETE** /users | Delete user
[**usersPreferredCurrencyGet**](UsersApi.md#usersPreferredCurrencyGet) | **GET** /users/preferred/currency | Get user preferred currency
[**usersPreferredCurrencyPut**](UsersApi.md#usersPreferredCurrencyPut) | **PUT** /users/preferred/currency | Update user preferred currency


# **usersDelete**
> void usersDelete()

Deletes the authenticated user\'s account

### Example


```typescript
import { createConfiguration, UsersApi } from '';

const configuration = createConfiguration();
const apiInstance = new UsersApi(configuration);

const request = {};

const data = await apiInstance.usersDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


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
**204** | No Content |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **usersPreferredCurrencyGet**
> UserUserPreferredCurrencyModel usersPreferredCurrencyGet()

Returns the preferred currency for the authenticated user

### Example


```typescript
import { createConfiguration, UsersApi } from '';

const configuration = createConfiguration();
const apiInstance = new UsersApi(configuration);

const request = {};

const data = await apiInstance.usersPreferredCurrencyGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**UserUserPreferredCurrencyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **usersPreferredCurrencyPut**
> void usersPreferredCurrencyPut(userUpdatePreferredCurrencyModel)

Updates the preferred currency for the authenticated user

### Example


```typescript
import { createConfiguration, UsersApi } from '';
import type { UsersApiUsersPreferredCurrencyPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new UsersApi(configuration);

const request: UsersApiUsersPreferredCurrencyPutRequest = {
    // Bearer token
  authorization: "Authorization_example",
    // Profile update parameters
  userUpdatePreferredCurrencyModel: {
    currency: "currency_example",
  },
};

const data = await apiInstance.usersPreferredCurrencyPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userUpdatePreferredCurrencyModel** | **UserUpdatePreferredCurrencyModel**| Profile update parameters |
 **authorization** | [**string**] | Bearer token | defaults to undefined


### Return type

**void**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


