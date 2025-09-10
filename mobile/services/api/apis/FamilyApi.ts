// TODO: better import syntax?
import {BaseAPIRequestFactory, RequiredError, COLLECTION_FORMATS} from './baseapi';
import {Configuration} from '../configuration';
import {RequestContext, HttpMethod, ResponseContext, HttpFile, HttpInfo} from '../http/http';
import {ObjectSerializer} from '../models/ObjectSerializer';
import {ApiException} from './exception';
import {canConsumeForm, isCodeInRange} from '../util';
import {SecurityAuthentication} from '../auth/auth';


import { FamilyCreateFamilyMemberModel } from '../models/FamilyCreateFamilyMemberModel';
import { FamilyCreateFamilyModel } from '../models/FamilyCreateFamilyModel';
import { FamilyFamilyAcceptInvitationRequest } from '../models/FamilyFamilyAcceptInvitationRequest';
import { FamilyFamilyDeclineInvitationRequest } from '../models/FamilyFamilyDeclineInvitationRequest';
import { FamilyFamilyInviteRequest } from '../models/FamilyFamilyInviteRequest';
import { FamilyFamilyInviteResponse } from '../models/FamilyFamilyInviteResponse';
import { FamilyFamilyModel } from '../models/FamilyFamilyModel';
import { FamilyFamilySeeInvitationResponse } from '../models/FamilyFamilySeeInvitationResponse';
import { FamilyPatchFamilyModel } from '../models/FamilyPatchFamilyModel';
import { FamilyUpdateFamilyMemberModel } from '../models/FamilyUpdateFamilyMemberModel';
import { FamilyUpdateFamilyModel } from '../models/FamilyUpdateFamilyModel';
import { GinxHttpErrorResponse } from '../models/GinxHttpErrorResponse';

/**
 * no description
 */
export class FamilyApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param familyId Family ID (UUID format)
     * @param familyFamilyAcceptInvitationRequest Invitation acceptance details
     */
    public async familiesFamilyIdAcceptPost(familyId: string, familyFamilyAcceptInvitationRequest: FamilyFamilyAcceptInvitationRequest, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdAcceptPost", "familyId");
        }


        // verify required parameter 'familyFamilyAcceptInvitationRequest' is not null or undefined
        if (familyFamilyAcceptInvitationRequest === null || familyFamilyAcceptInvitationRequest === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdAcceptPost", "familyFamilyAcceptInvitationRequest");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/accept'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyFamilyAcceptInvitationRequest, "FamilyFamilyAcceptInvitationRequest", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param familyId Family ID
     * @param familyFamilyDeclineInvitationRequest Decline invitation request
     */
    public async familiesFamilyIdDeclinePost(familyId: string, familyFamilyDeclineInvitationRequest: FamilyFamilyDeclineInvitationRequest, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdDeclinePost", "familyId");
        }


        // verify required parameter 'familyFamilyDeclineInvitationRequest' is not null or undefined
        if (familyFamilyDeclineInvitationRequest === null || familyFamilyDeclineInvitationRequest === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdDeclinePost", "familyFamilyDeclineInvitationRequest");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/decline'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyFamilyDeclineInvitationRequest, "FamilyFamilyDeclineInvitationRequest", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param familyId Family ID (UUID format)
     */
    public async familiesFamilyIdDelete(familyId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdDelete", "familyId");
        }


        // Path Params
        const localVarPath = '/families/{familyId}'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.DELETE);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param familyId Family ID
     * @param code Invitation code
     * @param familyMemberId Family member ID
     */
    public async familiesFamilyIdInvitationGet(familyId: string, code: string, familyMemberId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdInvitationGet", "familyId");
        }


        // verify required parameter 'code' is not null or undefined
        if (code === null || code === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdInvitationGet", "code");
        }


        // verify required parameter 'familyMemberId' is not null or undefined
        if (familyMemberId === null || familyMemberId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdInvitationGet", "familyMemberId");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/invitation'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (code !== undefined) {
            requestContext.setQueryParam("code", ObjectSerializer.serialize(code, "string", ""));
        }

        // Query Params
        if (familyMemberId !== undefined) {
            requestContext.setQueryParam("family_member_id", ObjectSerializer.serialize(familyMemberId, "string", ""));
        }


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param familyId Family ID (UUID format)
     * @param familyFamilyInviteRequest Invitation details including email, name, member ID and type (adult/kid)
     */
    public async familiesFamilyIdInvitePost(familyId: string, familyFamilyInviteRequest: FamilyFamilyInviteRequest, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdInvitePost", "familyId");
        }


        // verify required parameter 'familyFamilyInviteRequest' is not null or undefined
        if (familyFamilyInviteRequest === null || familyFamilyInviteRequest === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdInvitePost", "familyFamilyInviteRequest");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/invite'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyFamilyInviteRequest, "FamilyFamilyInviteRequest", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param familyId Family ID (UUID format)
     * @param familyMemberId Family Member ID (UUID format)
     * @param body 
     */
    public async familiesFamilyIdMembersFamilyMemberIdRevokePost(familyId: string, familyMemberId: string, body?: any, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersFamilyMemberIdRevokePost", "familyId");
        }


        // verify required parameter 'familyMemberId' is not null or undefined
        if (familyMemberId === null || familyMemberId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersFamilyMemberIdRevokePost", "familyMemberId");
        }



        // Path Params
        const localVarPath = '/families/{familyId}/members/{familyMemberId}/revoke'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)))
            .replace('{' + 'familyMemberId' + '}', encodeURIComponent(String(familyMemberId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(body, "any", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     */
    public async familiesFamilyIdMembersIdDelete(familyId: string, id: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersIdDelete", "familyId");
        }


        // verify required parameter 'id' is not null or undefined
        if (id === null || id === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersIdDelete", "id");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/members/{id}'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)))
            .replace('{' + 'id' + '}', encodeURIComponent(String(id)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.DELETE);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     * @param familyUpdateFamilyMemberModel Updated family member data
     */
    public async familiesFamilyIdMembersIdPut(familyId: string, id: string, familyUpdateFamilyMemberModel: FamilyUpdateFamilyMemberModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersIdPut", "familyId");
        }


        // verify required parameter 'id' is not null or undefined
        if (id === null || id === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersIdPut", "id");
        }


        // verify required parameter 'familyUpdateFamilyMemberModel' is not null or undefined
        if (familyUpdateFamilyMemberModel === null || familyUpdateFamilyMemberModel === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersIdPut", "familyUpdateFamilyMemberModel");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/members/{id}'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)))
            .replace('{' + 'id' + '}', encodeURIComponent(String(id)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyUpdateFamilyMemberModel, "FamilyUpdateFamilyMemberModel", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param familyId Family ID (UUID format)
     * @param familyCreateFamilyMemberModel Family member creation data
     */
    public async familiesFamilyIdMembersPost(familyId: string, familyCreateFamilyMemberModel: FamilyCreateFamilyMemberModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersPost", "familyId");
        }


        // verify required parameter 'familyCreateFamilyMemberModel' is not null or undefined
        if (familyCreateFamilyMemberModel === null || familyCreateFamilyMemberModel === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdMembersPost", "familyCreateFamilyMemberModel");
        }


        // Path Params
        const localVarPath = '/families/{familyId}/members'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyCreateFamilyMemberModel, "FamilyCreateFamilyMemberModel", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param familyId Family ID (UUID format)
     * @param familyUpdateFamilyModel Updated family data
     */
    public async familiesFamilyIdPut(familyId: string, familyUpdateFamilyModel: FamilyUpdateFamilyModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyId' is not null or undefined
        if (familyId === null || familyId === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdPut", "familyId");
        }


        // verify required parameter 'familyUpdateFamilyModel' is not null or undefined
        if (familyUpdateFamilyModel === null || familyUpdateFamilyModel === undefined) {
            throw new RequiredError("FamilyApi", "familiesFamilyIdPut", "familyUpdateFamilyModel");
        }


        // Path Params
        const localVarPath = '/families/{familyId}'
            .replace('{' + 'familyId' + '}', encodeURIComponent(String(familyId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyUpdateFamilyModel, "FamilyUpdateFamilyModel", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     */
    public async familiesMeGet(_options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // Path Params
        const localVarPath = '/families/me';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param familyPatchFamilyModel Family update data with members
     */
    public async familiesPatch(familyPatchFamilyModel: FamilyPatchFamilyModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyPatchFamilyModel' is not null or undefined
        if (familyPatchFamilyModel === null || familyPatchFamilyModel === undefined) {
            throw new RequiredError("FamilyApi", "familiesPatch", "familyPatchFamilyModel");
        }


        // Path Params
        const localVarPath = '/families';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PATCH);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyPatchFamilyModel, "FamilyPatchFamilyModel", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param familyCreateFamilyModel Family creation data
     */
    public async familiesPost(familyCreateFamilyModel: FamilyCreateFamilyModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'familyCreateFamilyModel' is not null or undefined
        if (familyCreateFamilyModel === null || familyCreateFamilyModel === undefined) {
            throw new RequiredError("FamilyApi", "familiesPost", "familyCreateFamilyModel");
        }


        // Path Params
        const localVarPath = '/families';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(familyCreateFamilyModel, "FamilyCreateFamilyModel", ""),
            contentType
        );
        requestContext.setBody(serializedBody);

        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

}

export class FamilyApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdAcceptPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdAcceptPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid or missing authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdDeclinePost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdDeclinePostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid family ID format", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdInvitationGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdInvitationGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilySeeInvitationResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: FamilyFamilySeeInvitationResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilySeeInvitationResponse", ""
            ) as FamilyFamilySeeInvitationResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilySeeInvitationResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilySeeInvitationResponse", ""
            ) as FamilyFamilySeeInvitationResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdInvitePost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdInvitePostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyInviteResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: FamilyFamilyInviteResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyInviteResponse", ""
            ) as FamilyFamilyInviteResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid or missing authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyInviteResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyInviteResponse", ""
            ) as FamilyFamilyInviteResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdMembersFamilyMemberIdRevokePost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid or missing authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family or member not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdMembersIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdMembersIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid ID format", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family or family member not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: void = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "void", ""
            ) as void;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdMembersIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdMembersIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or ID format", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid user authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family or family member not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdMembersPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdMembersPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or family ID", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid user authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesFamilyIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesFamilyIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or family ID", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid user authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesMeGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesMeGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid ID format", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid user authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesPatch
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesPatchWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid user authentication", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Family not found", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to familiesPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async familiesPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<FamilyFamilyModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data", body, response.headers);
        }
        if (isCodeInRange("401", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Unauthorized - Invalid user authentication", body, response.headers);
        }
        if (isCodeInRange("500", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Internal Server Error", body, response.headers);
        }

        // Work around for missing responses in specification, e.g. for petstore.yaml
        if (response.httpStatusCode >= 200 && response.httpStatusCode <= 299) {
            const body: FamilyFamilyModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "FamilyFamilyModel", ""
            ) as FamilyFamilyModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
