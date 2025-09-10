// TODO: better import syntax?
import {BaseAPIRequestFactory, RequiredError, COLLECTION_FORMATS} from './baseapi';
import {Configuration} from '../configuration';
import {RequestContext, HttpMethod, ResponseContext, HttpFile, HttpInfo} from '../http/http';
import {ObjectSerializer} from '../models/ObjectSerializer';
import {ApiException} from './exception';
import {canConsumeForm, isCodeInRange} from '../util';
import {SecurityAuthentication} from '../auth/auth';


import { DtoPaginatedResponseModelLabelLabelModel } from '../models/DtoPaginatedResponseModelLabelLabelModel';
import { GinxHttpErrorResponse } from '../models/GinxHttpErrorResponse';
import { LabelCreateLabelModel } from '../models/LabelCreateLabelModel';
import { LabelLabelModel } from '../models/LabelLabelModel';
import { LabelUpdateLabelModel } from '../models/LabelUpdateLabelModel';

/**
 * no description
 */
export class LabelsApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     */
    public async labelsDefaultGet(_options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // Path Params
        const localVarPath = '/labels/default';

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
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param search Search text to filter labels by name
     * @param limit Maximum number of items to return (default: 10)
     * @param offset Number of items to skip for pagination (default: 0)
     */
    public async labelsGet(search?: string, limit?: number, offset?: number, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;




        // Path Params
        const localVarPath = '/labels';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (search !== undefined) {
            requestContext.setQueryParam("search", ObjectSerializer.serialize(search, "string", ""));
        }

        // Query Params
        if (limit !== undefined) {
            requestContext.setQueryParam("limit", ObjectSerializer.serialize(limit, "number", ""));
        }

        // Query Params
        if (offset !== undefined) {
            requestContext.setQueryParam("offset", ObjectSerializer.serialize(offset, "number", ""));
        }


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param id Label ID (UUID format)
     */
    public async labelsIdDelete(id: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'id' is not null or undefined
        if (id === null || id === undefined) {
            throw new RequiredError("LabelsApi", "labelsIdDelete", "id");
        }


        // Path Params
        const localVarPath = '/labels/{id}'
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
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param id Label ID (UUID format)
     */
    public async labelsIdGet(id: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'id' is not null or undefined
        if (id === null || id === undefined) {
            throw new RequiredError("LabelsApi", "labelsIdGet", "id");
        }


        // Path Params
        const localVarPath = '/labels/{id}'
            .replace('{' + 'id' + '}', encodeURIComponent(String(id)));

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
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param id Label ID (UUID format)
     * @param labelUpdateLabelModel Updated label data
     */
    public async labelsIdPut(id: string, labelUpdateLabelModel: LabelUpdateLabelModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'id' is not null or undefined
        if (id === null || id === undefined) {
            throw new RequiredError("LabelsApi", "labelsIdPut", "id");
        }


        // verify required parameter 'labelUpdateLabelModel' is not null or undefined
        if (labelUpdateLabelModel === null || labelUpdateLabelModel === undefined) {
            throw new RequiredError("LabelsApi", "labelsIdPut", "labelUpdateLabelModel");
        }


        // Path Params
        const localVarPath = '/labels/{id}'
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
            ObjectSerializer.serialize(labelUpdateLabelModel, "LabelUpdateLabelModel", ""),
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
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param labelCreateLabelModel Label creation data
     */
    public async labelsPost(labelCreateLabelModel: LabelCreateLabelModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'labelCreateLabelModel' is not null or undefined
        if (labelCreateLabelModel === null || labelCreateLabelModel === undefined) {
            throw new RequiredError("LabelsApi", "labelsPost", "labelCreateLabelModel");
        }


        // Path Params
        const localVarPath = '/labels';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(labelCreateLabelModel, "LabelCreateLabelModel", ""),
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

export class LabelsApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to labelsDefaultGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async labelsDefaultGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<Array<LabelLabelModel> >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: Array<LabelLabelModel> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<LabelLabelModel>", ""
            ) as Array<LabelLabelModel>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
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
            const body: Array<LabelLabelModel> = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "Array<LabelLabelModel>", ""
            ) as Array<LabelLabelModel>;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to labelsGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async labelsGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<DtoPaginatedResponseModelLabelLabelModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: DtoPaginatedResponseModelLabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "DtoPaginatedResponseModelLabelLabelModel", ""
            ) as DtoPaginatedResponseModelLabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid query parameters", body, response.headers);
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
            const body: DtoPaginatedResponseModelLabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "DtoPaginatedResponseModelLabelLabelModel", ""
            ) as DtoPaginatedResponseModelLabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to labelsIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async labelsIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
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
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Label not found", body, response.headers);
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
     * @params response Response returned by the server for a request to labelsIdGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async labelsIdGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<LabelLabelModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: LabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "LabelLabelModel", ""
            ) as LabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
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
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Label not found", body, response.headers);
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
            const body: LabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "LabelLabelModel", ""
            ) as LabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to labelsIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async labelsIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<LabelLabelModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: LabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "LabelLabelModel", ""
            ) as LabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid ID format or input data", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Label not found", body, response.headers);
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
            const body: LabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "LabelLabelModel", ""
            ) as LabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to labelsPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async labelsPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<LabelLabelModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: LabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "LabelLabelModel", ""
            ) as LabelLabelModel;
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
            const body: LabelLabelModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "LabelLabelModel", ""
            ) as LabelLabelModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
