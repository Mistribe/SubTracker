# .FamilyApi

All URIs are relative to *https://api.subtracker.mistribe.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**familiesFamilyIdAcceptPost**](FamilyApi.md#familiesFamilyIdAcceptPost) | **POST** /families/{familyId}/accept | Accept a family invitation
[**familiesFamilyIdDeclinePost**](FamilyApi.md#familiesFamilyIdDeclinePost) | **POST** /families/{familyId}/decline | Decline family invitation
[**familiesFamilyIdDelete**](FamilyApi.md#familiesFamilyIdDelete) | **DELETE** /families/{familyId} | Delete family by ID
[**familiesFamilyIdInvitationGet**](FamilyApi.md#familiesFamilyIdInvitationGet) | **GET** /families/{familyId}/invitation | View family invitation details
[**familiesFamilyIdInvitePost**](FamilyApi.md#familiesFamilyIdInvitePost) | **POST** /families/{familyId}/invite | Invite a new member to the family
[**familiesFamilyIdMembersFamilyMemberIdRevokePost**](FamilyApi.md#familiesFamilyIdMembersFamilyMemberIdRevokePost) | **POST** /families/{familyId}/members/{familyMemberId}/revoke | Revoke family member
[**familiesFamilyIdMembersIdDelete**](FamilyApi.md#familiesFamilyIdMembersIdDelete) | **DELETE** /families/{familyId}/members/{id} | Delete family member by ID
[**familiesFamilyIdMembersIdPut**](FamilyApi.md#familiesFamilyIdMembersIdPut) | **PUT** /families/{familyId}/members/{id} | Update family member by ID
[**familiesFamilyIdMembersPost**](FamilyApi.md#familiesFamilyIdMembersPost) | **POST** /families/{familyId}/members | Add a new family member
[**familiesFamilyIdPut**](FamilyApi.md#familiesFamilyIdPut) | **PUT** /families/{familyId} | Update a family
[**familiesMeGet**](FamilyApi.md#familiesMeGet) | **GET** /families/me | Get user\&#39;s family
[**familiesPatch**](FamilyApi.md#familiesPatch) | **PATCH** /families | Patch family with members
[**familiesPost**](FamilyApi.md#familiesPost) | **POST** /families | Create a new family


# **familiesFamilyIdAcceptPost**
> void familiesFamilyIdAcceptPost(familyFamilyAcceptInvitationRequest)

Accepts an invitation to join a family using the provided invitation code

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdAcceptPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdAcceptPostRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Invitation acceptance details
  familyFamilyAcceptInvitationRequest: {
    familyMemberId: "123e4567-e89b-12d3-a456-426614174000",
    invitationCode: "123456",
  },
};

const data = await apiInstance.familiesFamilyIdAcceptPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyFamilyAcceptInvitationRequest** | **FamilyFamilyAcceptInvitationRequest**| Invitation acceptance details |
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined


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
**204** | Successfully accepted invitation |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid or missing authentication |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdDeclinePost**
> void familiesFamilyIdDeclinePost(familyFamilyDeclineInvitationRequest)

Endpoint to decline an invitation to join a family

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdDeclinePostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdDeclinePostRequest = {
    // Family ID
  familyId: "familyId_example",
    // Decline invitation request
  familyFamilyDeclineInvitationRequest: {
    familyMemberId: "123e4567-e89b-12d3-a456-426614174000",
    invitationCode: "123456",
  },
};

const data = await apiInstance.familiesFamilyIdDeclinePost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyFamilyDeclineInvitationRequest** | **FamilyFamilyDeclineInvitationRequest**| Decline invitation request |
 **familyId** | [**string**] | Family ID | defaults to undefined


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

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdDelete**
> void familiesFamilyIdDelete()

Permanently delete a family and all its members

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdDeleteRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
};

const data = await apiInstance.familiesFamilyIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined


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
**204** | No Content - Family successfully deleted |  -  |
**400** | Bad Request - Invalid family ID format |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdInvitationGet**
> FamilyFamilySeeInvitationResponse familiesFamilyIdInvitationGet()

Get information about a family invitation using invitation code

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdInvitationGetRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdInvitationGetRequest = {
    // Family ID
  familyId: "familyId_example",
    // Invitation code
  code: "code_example",
    // Family member ID
  familyMemberId: "family_member_id_example",
};

const data = await apiInstance.familiesFamilyIdInvitationGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyId** | [**string**] | Family ID | defaults to undefined
 **code** | [**string**] | Invitation code | defaults to undefined
 **familyMemberId** | [**string**] | Family member ID | defaults to undefined


### Return type

**FamilyFamilySeeInvitationResponse**

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

# **familiesFamilyIdInvitePost**
> FamilyFamilyInviteResponse familiesFamilyIdInvitePost(familyFamilyInviteRequest)

Creates an invitation for a new member to join the family

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdInvitePostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdInvitePostRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Invitation details including email, name, member ID and type (adult/kid)
  familyFamilyInviteRequest: {
    email: "email_example",
    familyMemberId: "familyMemberId_example",
    name: "name_example",
    type: "adult",
  },
};

const data = await apiInstance.familiesFamilyIdInvitePost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyFamilyInviteRequest** | **FamilyFamilyInviteRequest**| Invitation details including email, name, member ID and type (adult/kid) |
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined


### Return type

**FamilyFamilyInviteResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully created invitation with code and IDs |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid or missing authentication |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdMembersFamilyMemberIdRevokePost**
> void familiesFamilyIdMembersFamilyMemberIdRevokePost()

Revokes a member from the family

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdMembersFamilyMemberIdRevokePostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdMembersFamilyMemberIdRevokePostRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Family Member ID (UUID format)
  familyMemberId: "familyMemberId_example",
  
  body: {},
};

const data = await apiInstance.familiesFamilyIdMembersFamilyMemberIdRevokePost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **any**|  |
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined
 **familyMemberId** | [**string**] | Family Member ID (UUID format) | defaults to undefined


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
**204** | Successfully revoked member |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid or missing authentication |  -  |
**404** | Family or member not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdMembersIdDelete**
> void familiesFamilyIdMembersIdDelete()

Permanently delete a family member from a family

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdMembersIdDeleteRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdMembersIdDeleteRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Family member ID (UUID format)
  id: "id_example",
};

const data = await apiInstance.familiesFamilyIdMembersIdDelete(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined
 **id** | [**string**] | Family member ID (UUID format) | defaults to undefined


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
**204** | No Content - Family member successfully deleted |  -  |
**400** | Bad Request - Invalid ID format |  -  |
**404** | Family or family member not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdMembersIdPut**
> FamilyFamilyModel familiesFamilyIdMembersIdPut(familyUpdateFamilyMemberModel)

Update an existing family member\'s information such as name and kid status

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdMembersIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdMembersIdPutRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Family member ID (UUID format)
  id: "id_example",
    // Updated family member data
  familyUpdateFamilyMemberModel: {
    name: "name_example",
    type: "owner",
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.familiesFamilyIdMembersIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyUpdateFamilyMemberModel** | **FamilyUpdateFamilyMemberModel**| Updated family member data |
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined
 **id** | [**string**] | Family member ID (UUID format) | defaults to undefined


### Return type

**FamilyFamilyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated family member |  -  |
**400** | Bad Request - Invalid input data or ID format |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Family or family member not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdMembersPost**
> FamilyFamilyModel familiesFamilyIdMembersPost(familyCreateFamilyMemberModel)

Add a new member to an existing family

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdMembersPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdMembersPostRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Family member creation data
  familyCreateFamilyMemberModel: {
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    id: "id_example",
    name: "name_example",
    type: "owner",
  },
};

const data = await apiInstance.familiesFamilyIdMembersPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyCreateFamilyMemberModel** | **FamilyCreateFamilyMemberModel**| Family member creation data |
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined


### Return type

**FamilyFamilyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully added family member |  -  |
**400** | Bad Request - Invalid input data or family ID |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesFamilyIdPut**
> FamilyFamilyModel familiesFamilyIdPut(familyUpdateFamilyModel)

Update family information such as name and other details

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesFamilyIdPutRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesFamilyIdPutRequest = {
    // Family ID (UUID format)
  familyId: "familyId_example",
    // Updated family data
  familyUpdateFamilyModel: {
    name: "name_example",
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.familiesFamilyIdPut(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyUpdateFamilyModel** | **FamilyUpdateFamilyModel**| Updated family data |
 **familyId** | [**string**] | Family ID (UUID format) | defaults to undefined


### Return type

**FamilyFamilyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated family |  -  |
**400** | Bad Request - Invalid input data or family ID |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesMeGet**
> FamilyFamilyModel familiesMeGet()

Retrieve the user\'s family

### Example


```typescript
import { createConfiguration, FamilyApi } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request = {};

const data = await apiInstance.familiesMeGet(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters
This endpoint does not need any parameter.


### Return type

**FamilyFamilyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully retrieved family |  -  |
**400** | Bad Request - Invalid ID format |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesPatch**
> FamilyFamilyModel familiesPatch(familyPatchFamilyModel)

Update or create a family with specified members. If family doesn\'t exist, it will be created.

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesPatchRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesPatchRequest = {
    // Family update data with members
  familyPatchFamilyModel: {
    id: "id_example",
    members: [
      {
        id: "id_example",
        name: "name_example",
        type: "owner",
        updatedAt: new Date('1970-01-01T00:00:00.00Z'),
      },
    ],
    name: "name_example",
    updatedAt: new Date('1970-01-01T00:00:00.00Z'),
  },
};

const data = await apiInstance.familiesPatch(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyPatchFamilyModel** | **FamilyPatchFamilyModel**| Family update data with members |


### Return type

**FamilyFamilyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | Successfully updated family |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**404** | Family not found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)

# **familiesPost**
> FamilyFamilyModel familiesPost(familyCreateFamilyModel)

Create a new family with the authenticated user as the owner and initial member

### Example


```typescript
import { createConfiguration, FamilyApi } from '';
import type { FamilyApiFamiliesPostRequest } from '';

const configuration = createConfiguration();
const apiInstance = new FamilyApi(configuration);

const request: FamilyApiFamiliesPostRequest = {
    // Family creation data
  familyCreateFamilyModel: {
    createdAt: new Date('1970-01-01T00:00:00.00Z'),
    creatorName: "creatorName_example",
    id: "id_example",
    name: "name_example",
  },
};

const data = await apiInstance.familiesPost(request);
console.log('API called successfully. Returned data:', data);
```


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **familyCreateFamilyModel** | **FamilyCreateFamilyModel**| Family creation data |


### Return type

**FamilyFamilyModel**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
**201** | Successfully created family |  -  |
**400** | Bad Request - Invalid input data |  -  |
**401** | Unauthorized - Invalid user authentication |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](README.md#documentation-for-api-endpoints) [[Back to Model list]](README.md#documentation-for-models) [[Back to README]](README.md)


