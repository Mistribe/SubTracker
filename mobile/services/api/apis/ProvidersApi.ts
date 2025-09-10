// TODO: better import syntax?
import {BaseAPIRequestFactory, RequiredError, COLLECTION_FORMATS} from './baseapi';
import {Configuration} from '../configuration';
import {RequestContext, HttpMethod, ResponseContext, HttpFile, HttpInfo} from '../http/http';
import {ObjectSerializer} from '../models/ObjectSerializer';
import {ApiException} from './exception';
import {canConsumeForm, isCodeInRange} from '../util';
import {SecurityAuthentication} from '../auth/auth';


import { DtoPaginatedResponseModelProviderProviderModel } from '../models/DtoPaginatedResponseModelProviderProviderModel';
import { GinxHttpErrorResponse } from '../models/GinxHttpErrorResponse';
import { ProviderCreatePlanModel } from '../models/ProviderCreatePlanModel';
import { ProviderCreatePriceModel } from '../models/ProviderCreatePriceModel';
import { ProviderCreateProviderModel } from '../models/ProviderCreateProviderModel';
import { ProviderPlanModel } from '../models/ProviderPlanModel';
import { ProviderPriceModel } from '../models/ProviderPriceModel';
import { ProviderProviderModel } from '../models/ProviderProviderModel';
import { ProviderUpdatePlanModel } from '../models/ProviderUpdatePlanModel';
import { ProviderUpdatePriceModel } from '../models/ProviderUpdatePriceModel';
import { ProviderUpdateProviderModel } from '../models/ProviderUpdateProviderModel';

/**
 * no description
 */
export class ProvidersApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param search Search term
     * @param offset Offset (default: 0)
     * @param limit Limit per request (default: 10)
     */
    public async providersGet(search?: string, offset?: number, limit?: number, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;




        // Path Params
        const localVarPath = '/providers';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (search !== undefined) {
            requestContext.setQueryParam("search", ObjectSerializer.serialize(search, "string", ""));
        }

        // Query Params
        if (offset !== undefined) {
            requestContext.setQueryParam("offset", ObjectSerializer.serialize(offset, "number", ""));
        }

        // Query Params
        if (limit !== undefined) {
            requestContext.setQueryParam("limit", ObjectSerializer.serialize(limit, "number", ""));
        }


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param providerCreateProviderModel Provider creation data
     */
    public async providersPost(providerCreateProviderModel: ProviderCreateProviderModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerCreateProviderModel' is not null or undefined
        if (providerCreateProviderModel === null || providerCreateProviderModel === undefined) {
            throw new RequiredError("ProvidersApi", "providersPost", "providerCreateProviderModel");
        }


        // Path Params
        const localVarPath = '/providers';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(providerCreateProviderModel, "ProviderCreateProviderModel", ""),
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
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public async providersProviderIdDelete(providerId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdDelete", "providerId");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)));

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
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public async providersProviderIdGet(providerId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdGet", "providerId");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)));

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
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     */
    public async providersProviderIdPlansPlanIdDelete(providerId: string, planId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdDelete", "providerId");
        }


        // verify required parameter 'planId' is not null or undefined
        if (planId === null || planId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdDelete", "planId");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}/plans/{planId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)))
            .replace('{' + 'planId' + '}', encodeURIComponent(String(planId)));

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
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerCreatePriceModel Price creation data
     */
    public async providersProviderIdPlansPlanIdPricesPost(providerId: string, planId: string, providerCreatePriceModel: ProviderCreatePriceModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPost", "providerId");
        }


        // verify required parameter 'planId' is not null or undefined
        if (planId === null || planId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPost", "planId");
        }


        // verify required parameter 'providerCreatePriceModel' is not null or undefined
        if (providerCreatePriceModel === null || providerCreatePriceModel === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPost", "providerCreatePriceModel");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}/plans/{planId}/prices'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)))
            .replace('{' + 'planId' + '}', encodeURIComponent(String(planId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(providerCreatePriceModel, "ProviderCreatePriceModel", ""),
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
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     */
    public async providersProviderIdPlansPlanIdPricesPriceIdDelete(providerId: string, planId: string, priceId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdDelete", "providerId");
        }


        // verify required parameter 'planId' is not null or undefined
        if (planId === null || planId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdDelete", "planId");
        }


        // verify required parameter 'priceId' is not null or undefined
        if (priceId === null || priceId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdDelete", "priceId");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}/plans/{planId}/prices/{priceId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)))
            .replace('{' + 'planId' + '}', encodeURIComponent(String(planId)))
            .replace('{' + 'priceId' + '}', encodeURIComponent(String(priceId)));

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
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     * @param providerUpdatePriceModel Updated price data
     */
    public async providersProviderIdPlansPlanIdPricesPriceIdPut(providerId: string, planId: string, priceId: string, providerUpdatePriceModel: ProviderUpdatePriceModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdPut", "providerId");
        }


        // verify required parameter 'planId' is not null or undefined
        if (planId === null || planId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdPut", "planId");
        }


        // verify required parameter 'priceId' is not null or undefined
        if (priceId === null || priceId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdPut", "priceId");
        }


        // verify required parameter 'providerUpdatePriceModel' is not null or undefined
        if (providerUpdatePriceModel === null || providerUpdatePriceModel === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPricesPriceIdPut", "providerUpdatePriceModel");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}/plans/{planId}/prices/{priceId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)))
            .replace('{' + 'planId' + '}', encodeURIComponent(String(planId)))
            .replace('{' + 'priceId' + '}', encodeURIComponent(String(priceId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(providerUpdatePriceModel, "ProviderUpdatePriceModel", ""),
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
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerUpdatePlanModel Updated plan data
     */
    public async providersProviderIdPlansPlanIdPut(providerId: string, planId: string, providerUpdatePlanModel: ProviderUpdatePlanModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPut", "providerId");
        }


        // verify required parameter 'planId' is not null or undefined
        if (planId === null || planId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPut", "planId");
        }


        // verify required parameter 'providerUpdatePlanModel' is not null or undefined
        if (providerUpdatePlanModel === null || providerUpdatePlanModel === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPlanIdPut", "providerUpdatePlanModel");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}/plans/{planId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)))
            .replace('{' + 'planId' + '}', encodeURIComponent(String(planId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(providerUpdatePlanModel, "ProviderUpdatePlanModel", ""),
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
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param providerId Provider ID (UUID format)
     * @param providerCreatePlanModel Plan creation data
     */
    public async providersProviderIdPlansPost(providerId: string, providerCreatePlanModel: ProviderCreatePlanModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPost", "providerId");
        }


        // verify required parameter 'providerCreatePlanModel' is not null or undefined
        if (providerCreatePlanModel === null || providerCreatePlanModel === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPlansPost", "providerCreatePlanModel");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}/plans'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(providerCreatePlanModel, "ProviderCreatePlanModel", ""),
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
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param providerId Provider ID (UUID format)
     * @param providerUpdateProviderModel Updated provider data
     */
    public async providersProviderIdPut(providerId: string, providerUpdateProviderModel: ProviderUpdateProviderModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'providerId' is not null or undefined
        if (providerId === null || providerId === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPut", "providerId");
        }


        // verify required parameter 'providerUpdateProviderModel' is not null or undefined
        if (providerUpdateProviderModel === null || providerUpdateProviderModel === undefined) {
            throw new RequiredError("ProvidersApi", "providersProviderIdPut", "providerUpdateProviderModel");
        }


        // Path Params
        const localVarPath = '/providers/{providerId}'
            .replace('{' + 'providerId' + '}', encodeURIComponent(String(providerId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(providerUpdateProviderModel, "ProviderUpdateProviderModel", ""),
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

export class ProvidersApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<DtoPaginatedResponseModelProviderProviderModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: DtoPaginatedResponseModelProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "DtoPaginatedResponseModelProviderProviderModel", ""
            ) as DtoPaginatedResponseModelProviderProviderModel;
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
            const body: DtoPaginatedResponseModelProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "DtoPaginatedResponseModelProviderProviderModel", ""
            ) as DtoPaginatedResponseModelProviderProviderModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderProviderModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: ProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderProviderModel", ""
            ) as ProviderProviderModel;
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
            const body: ProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderProviderModel", ""
            ) as ProviderProviderModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersProviderIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid provider ID format", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider not found", body, response.headers);
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
     * @params response Response returned by the server for a request to providersProviderIdGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderProviderModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderProviderModel", ""
            ) as ProviderProviderModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid provider ID format", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider not found", body, response.headers);
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
            const body: ProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderProviderModel", ""
            ) as ProviderProviderModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersProviderIdPlansPlanIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPlansPlanIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
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
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider or plan not found", body, response.headers);
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
     * @params response Response returned by the server for a request to providersProviderIdPlansPlanIdPricesPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPlansPlanIdPricesPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderPriceModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: ProviderPriceModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPriceModel", ""
            ) as ProviderPriceModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or IDs", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider or plan not found", body, response.headers);
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
            const body: ProviderPriceModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPriceModel", ""
            ) as ProviderPriceModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersProviderIdPlansPlanIdPricesPriceIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
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
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider, plan, or price not found", body, response.headers);
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
     * @params response Response returned by the server for a request to providersProviderIdPlansPlanIdPricesPriceIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderPriceModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ProviderPriceModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPriceModel", ""
            ) as ProviderPriceModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or IDs", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider, plan, or price not found", body, response.headers);
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
            const body: ProviderPriceModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPriceModel", ""
            ) as ProviderPriceModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersProviderIdPlansPlanIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPlansPlanIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderPlanModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ProviderPlanModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPlanModel", ""
            ) as ProviderPlanModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or IDs", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider or plan not found", body, response.headers);
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
            const body: ProviderPlanModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPlanModel", ""
            ) as ProviderPlanModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersProviderIdPlansPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPlansPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderPlanModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: ProviderPlanModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPlanModel", ""
            ) as ProviderPlanModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or provider ID", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider not found", body, response.headers);
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
            const body: ProviderPlanModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderPlanModel", ""
            ) as ProviderPlanModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to providersProviderIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async providersProviderIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<ProviderProviderModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: ProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderProviderModel", ""
            ) as ProviderProviderModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or provider ID", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Provider not found", body, response.headers);
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
            const body: ProviderProviderModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "ProviderProviderModel", ""
            ) as ProviderProviderModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
