# .HealthApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**healthzLiveGet**](HealthApi.md#healthzLiveGet) | **GET** /healthz/live | Health check endpoint


# **healthzLiveGet**
> { [key: string]: string; } healthzLiveGet()

Returns the health status of the application

### Example


```typescript
import { createConfiguration, HealthApi } from '';

const configuration = createConfiguration();
const apiInstance = new HealthApi(configuration);

const request = {};

const data = await apiInstance.healthzLiveGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**{ [key: string]: string; }**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Health status |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


