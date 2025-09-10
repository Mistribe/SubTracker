# .LabelsApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**labelsDefaultGet**](LabelsApi.md#labelsDefaultGet) | **GET** /labels/default | Get default labels
[**labelsGet**](LabelsApi.md#labelsGet) | **GET** /labels | Get all labels
[**labelsIdDelete**](LabelsApi.md#labelsIdDelete) | **DELETE** /labels/{id} | Delete label by ID
[**labelsIdGet**](LabelsApi.md#labelsIdGet) | **GET** /labels/{id} | Get label by ID
[**labelsIdPut**](LabelsApi.md#labelsIdPut) | **PUT** /labels/{id} | Update label by ID
[**labelsPost**](LabelsApi.md#labelsPost) | **POST** /labels | Create a new label


# **labelsDefaultGet**
> Array<LabelLabelModel> labelsDefaultGet()

Retrieves a list of default system labels available to all users

### Example


```typescript
import { createConfiguration, LabelsApi } from '';

const configuration = createConfiguration();
const apiInstance = new LabelsApi(configuration);

const request = {};

const data = await apiInstance.labelsDefaultGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**Array<LabelLabelModel>**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | List of default labels |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **labelsGet**
> DtoPaginatedResponseModelLabelLabelModel labelsGet()

Retrieve a paginated list of labels with optional filtering by owner type and search text

### Example


```typescript
import { createConfiguration, LabelsApi } from '';
import type { LabelsApiLabelsGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new LabelsApi(configuration);

const request: LabelsApiLabelsGetRequest = {
    // Search text to filter labels by name (optional)
  search: "search_example",
    // Maximum number of items to return (default: 10) (optional)
  limit: 1,
    // Number of items to skip for pagination (default: 0) (optional)
  offset: 1,
};

const data = await apiInstance.labelsGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **search** | [**string**] | Search text to filter labels by name | (optional) defaults to undefined
 **limit** | [**number**] | Maximum number of items to return (default: 10) | (optional) defaults to undefined
 **offset** | [**number**] | Number of items to skip for pagination (default: 0) | (optional) defaults to undefined


### Return type

**DtoPaginatedResponseModelLabelLabelModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Paginated list of labels |  -  |
**400** | Bad Request - Invalid query parameters |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **labelsIdDelete**
> void labelsIdDelete()

Permanently delete a label by its unique identifier

### Example


```typescript
import { createConfiguration, LabelsApi } from '';
import type { LabelsApiLabelsIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new LabelsApi(configuration);

const request: LabelsApiLabelsIdDeleteRequest = {
    // Label ID (UUID format)
  id: "id_example",
};

const data = await apiInstance.labelsIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | [**string**] | Label ID (UUID format) | defaults to undefined


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
**204** | No Content - Label successfully deleted |  -  |
**400** | Bad Request - Invalid ID format |  -  |
**404** | Label not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **labelsIdGet**
> LabelLabelModel labelsIdGet()

Retrieve a single label by its unique identifier

### Example


```typescript
import { createConfiguration, LabelsApi } from '';
import type { LabelsApiLabelsIdGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new LabelsApi(configuration);

const request: LabelsApiLabelsIdGetRequest = {
    // Label ID (UUID format)
  id: "id_example",
};

const data = await apiInstance.labelsIdGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | [**string**] | Label ID (UUID format) | defaults to undefined


### Return type

**LabelLabelModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request - Invalid ID format |  -  |
**404** | Label not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **labelsIdPut**
> LabelLabelModel labelsIdPut(labelUpdateLabelModel)

Update an existing label\'s name and color by its unique identifier

### Example


```typescript
import { createConfiguration, LabelsApi } from '';
import type { LabelsApiLabelsIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new LabelsApi(configuration);

const request: LabelsApiLabelsIdPutRequest = {
    // Label ID (UUID format)
  id: "id_example",
    // Updated label data
  labelUpdateLabelModel: {
    color: "color_example",
    name: "name_example",
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.labelsIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **labelUpdateLabelModel** | **LabelUpdateLabelModel**| Updated label data |
 **id** | [**string**] | Label ID (UUID format) | defaults to undefined


### Return type

**LabelLabelModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated label |  -  |
**400** | Bad Request - Invalid ID format or input data |  -  |
**404** | Label not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **labelsPost**
> LabelLabelModel labelsPost(labelCreateLabelModel)

Create a new label with specified name, color, and owner information

### Example


```typescript
import { createConfiguration, LabelsApi } from '';
import type { LabelsApiLabelsPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new LabelsApi(configuration);

const request: LabelsApiLabelsPostRequest = {
    // Label creation data
  labelCreateLabelModel: {
    color: "color_example",
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    id: "id_example",
    name: "name_example",
    owner: {
      familyId: "familyId_example",
      type: "personal",
    },
  },
};

const data = await apiInstance.labelsPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **labelCreateLabelModel** | **LabelCreateLabelModel**| Label creation data |


### Return type

**LabelLabelModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully created label |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


