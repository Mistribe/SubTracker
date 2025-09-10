// TODO: better import syntax?
import {BaseAPIRequestFactory, RequiredError, COLLECTION_FORMATS} from './baseapi';
import {Configuration} from '../configuration';
import {RequestContext, HttpMethod, ResponseContext, HttpFile, HttpInfo} from '../http/http';
import {ObjectSerializer} from '../models/ObjectSerializer';
import {ApiException} from './exception';
import {canConsumeForm, isCodeInRange} from '../util';
import {SecurityAuthentication} from '../auth/auth';


import { DtoPaginatedResponseModelSubscriptionSubscriptionModel } from '../models/DtoPaginatedResponseModelSubscriptionSubscriptionModel';
import { GinxHttpErrorResponse } from '../models/GinxHttpErrorResponse';
import { SubscriptionCreateSubscriptionModel } from '../models/SubscriptionCreateSubscriptionModel';
import { SubscriptionPatchSubscriptionModel } from '../models/SubscriptionPatchSubscriptionModel';
import { SubscriptionSubscriptionModel } from '../models/SubscriptionSubscriptionModel';
import { SubscriptionSubscriptionSummaryResponse } from '../models/SubscriptionSubscriptionSummaryResponse';
import { SubscriptionUpdateSubscriptionModel } from '../models/SubscriptionUpdateSubscriptionModel';

/**
 * no description
 */
export class SubscriptionsApiRequestFactory extends BaseAPIRequestFactory {

    /**
     * Retrieve a paginated list of all subscriptions for the authenticated user
     * Get all subscriptions
     * @param search Search text
     * @param recurrencies Filter by recurrency types
     * @param fromDate Filter by start date (RFC3339)
     * @param toDate Filter by end date (RFC3339)
     * @param users Filter by user IDs
     * @param withInactive Include inactive subscriptions
     * @param providers Filter by provider IDs
     * @param limit Number of items per page (default: 10)
     * @param offset Page number (default: 0)
     */
    public async subscriptionsGet(search?: string, recurrencies?: Array<string>, fromDate?: string, toDate?: string, users?: Array<string>, withInactive?: boolean, providers?: Array<string>, limit?: number, offset?: number, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;










        // Path Params
        const localVarPath = '/subscriptions';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (search !== undefined) {
            requestContext.setQueryParam("search", ObjectSerializer.serialize(search, "string", ""));
        }

        // Query Params
        if (recurrencies !== undefined) {
            const serializedParams = ObjectSerializer.serialize(recurrencies, "Array<string>", "");
            for (const serializedParam of serializedParams) {
                requestContext.appendQueryParam("recurrencies", serializedParam);
            }
        }

        // Query Params
        if (fromDate !== undefined) {
            requestContext.setQueryParam("from_date", ObjectSerializer.serialize(fromDate, "string", ""));
        }

        // Query Params
        if (toDate !== undefined) {
            requestContext.setQueryParam("to_date", ObjectSerializer.serialize(toDate, "string", ""));
        }

        // Query Params
        if (users !== undefined) {
            const serializedParams = ObjectSerializer.serialize(users, "Array<string>", "");
            for (const serializedParam of serializedParams) {
                requestContext.appendQueryParam("users", serializedParam);
            }
        }

        // Query Params
        if (withInactive !== undefined) {
            requestContext.setQueryParam("with_inactive", ObjectSerializer.serialize(withInactive, "boolean", ""));
        }

        // Query Params
        if (providers !== undefined) {
            const serializedParams = ObjectSerializer.serialize(providers, "Array<string>", "");
            for (const serializedParam of serializedParams) {
                requestContext.appendQueryParam("providers", serializedParam);
            }
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
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param subscriptionPatchSubscriptionModel Complete subscription data
     */
    public async subscriptionsPatch(subscriptionPatchSubscriptionModel: SubscriptionPatchSubscriptionModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'subscriptionPatchSubscriptionModel' is not null or undefined
        if (subscriptionPatchSubscriptionModel === null || subscriptionPatchSubscriptionModel === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsPatch", "subscriptionPatchSubscriptionModel");
        }


        // Path Params
        const localVarPath = '/subscriptions';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PATCH);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(subscriptionPatchSubscriptionModel, "SubscriptionPatchSubscriptionModel", ""),
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
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param subscriptionCreateSubscriptionModel Subscription creation data
     */
    public async subscriptionsPost(subscriptionCreateSubscriptionModel: SubscriptionCreateSubscriptionModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'subscriptionCreateSubscriptionModel' is not null or undefined
        if (subscriptionCreateSubscriptionModel === null || subscriptionCreateSubscriptionModel === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsPost", "subscriptionCreateSubscriptionModel");
        }


        // Path Params
        const localVarPath = '/subscriptions';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.POST);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(subscriptionCreateSubscriptionModel, "SubscriptionCreateSubscriptionModel", ""),
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
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public async subscriptionsSubscriptionIdDelete(subscriptionId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'subscriptionId' is not null or undefined
        if (subscriptionId === null || subscriptionId === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSubscriptionIdDelete", "subscriptionId");
        }


        // Path Params
        const localVarPath = '/subscriptions/{subscriptionId}'
            .replace('{' + 'subscriptionId' + '}', encodeURIComponent(String(subscriptionId)));

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
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public async subscriptionsSubscriptionIdGet(subscriptionId: string, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'subscriptionId' is not null or undefined
        if (subscriptionId === null || subscriptionId === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSubscriptionIdGet", "subscriptionId");
        }


        // Path Params
        const localVarPath = '/subscriptions/{subscriptionId}'
            .replace('{' + 'subscriptionId' + '}', encodeURIComponent(String(subscriptionId)));

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
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     * @param subscriptionUpdateSubscriptionModel Updated subscription data
     */
    public async subscriptionsSubscriptionIdPut(subscriptionId: string, subscriptionUpdateSubscriptionModel: SubscriptionUpdateSubscriptionModel, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'subscriptionId' is not null or undefined
        if (subscriptionId === null || subscriptionId === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSubscriptionIdPut", "subscriptionId");
        }


        // verify required parameter 'subscriptionUpdateSubscriptionModel' is not null or undefined
        if (subscriptionUpdateSubscriptionModel === null || subscriptionUpdateSubscriptionModel === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSubscriptionIdPut", "subscriptionUpdateSubscriptionModel");
        }


        // Path Params
        const localVarPath = '/subscriptions/{subscriptionId}'
            .replace('{' + 'subscriptionId' + '}', encodeURIComponent(String(subscriptionId)));

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.PUT);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")


        // Body Params
        const contentType = ObjectSerializer.getPreferredMediaType([
            "application/json"
        ]);
        requestContext.setHeaderParam("Content-Type", contentType);
        const serializedBody = ObjectSerializer.stringify(
            ObjectSerializer.serialize(subscriptionUpdateSubscriptionModel, "SubscriptionUpdateSubscriptionModel", ""),
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
     * Returns summary information about subscriptions including total costs and upcoming renewals
     * Get subscription summary
     * @param topProviders Number of top providers to return
     * @param topLabels Number of top labels to return
     * @param upcomingRenewals Number of upcoming renewals to return
     * @param totalMonthly Include monthly total costs
     * @param totalYearly Include yearly total costs
     */
    public async subscriptionsSummaryGet(topProviders: number, topLabels: number, upcomingRenewals: number, totalMonthly: boolean, totalYearly: boolean, _options?: Configuration): Promise<RequestContext> {
        let _config = _options || this.configuration;

        // verify required parameter 'topProviders' is not null or undefined
        if (topProviders === null || topProviders === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSummaryGet", "topProviders");
        }


        // verify required parameter 'topLabels' is not null or undefined
        if (topLabels === null || topLabels === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSummaryGet", "topLabels");
        }


        // verify required parameter 'upcomingRenewals' is not null or undefined
        if (upcomingRenewals === null || upcomingRenewals === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSummaryGet", "upcomingRenewals");
        }


        // verify required parameter 'totalMonthly' is not null or undefined
        if (totalMonthly === null || totalMonthly === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSummaryGet", "totalMonthly");
        }


        // verify required parameter 'totalYearly' is not null or undefined
        if (totalYearly === null || totalYearly === undefined) {
            throw new RequiredError("SubscriptionsApi", "subscriptionsSummaryGet", "totalYearly");
        }


        // Path Params
        const localVarPath = '/subscriptions/summary';

        // Make Request Context
        const requestContext = _config.baseServer.makeRequestContext(localVarPath, HttpMethod.GET);
        requestContext.setHeaderParam("Accept", "application/json, */*;q=0.8")

        // Query Params
        if (topProviders !== undefined) {
            requestContext.setQueryParam("top_providers", ObjectSerializer.serialize(topProviders, "number", ""));
        }

        // Query Params
        if (topLabels !== undefined) {
            requestContext.setQueryParam("top_labels", ObjectSerializer.serialize(topLabels, "number", ""));
        }

        // Query Params
        if (upcomingRenewals !== undefined) {
            requestContext.setQueryParam("upcoming_renewals", ObjectSerializer.serialize(upcomingRenewals, "number", ""));
        }

        // Query Params
        if (totalMonthly !== undefined) {
            requestContext.setQueryParam("total_monthly", ObjectSerializer.serialize(totalMonthly, "boolean", ""));
        }

        // Query Params
        if (totalYearly !== undefined) {
            requestContext.setQueryParam("total_yearly", ObjectSerializer.serialize(totalYearly, "boolean", ""));
        }


        
        const defaultAuth: SecurityAuthentication | undefined = _config?.authMethods?.default
        if (defaultAuth?.applySecurityAuthentication) {
            await defaultAuth?.applySecurityAuthentication(requestContext);
        }

        return requestContext;
    }

}

export class SubscriptionsApiResponseProcessor {

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to subscriptionsGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<DtoPaginatedResponseModelSubscriptionSubscriptionModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: DtoPaginatedResponseModelSubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "DtoPaginatedResponseModelSubscriptionSubscriptionModel", ""
            ) as DtoPaginatedResponseModelSubscriptionSubscriptionModel;
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
            const body: DtoPaginatedResponseModelSubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "DtoPaginatedResponseModelSubscriptionSubscriptionModel", ""
            ) as DtoPaginatedResponseModelSubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to subscriptionsPatch
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsPatchWithHttpInfo(response: ResponseContext): Promise<HttpInfo<SubscriptionSubscriptionModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
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
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Subscription not found", body, response.headers);
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
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to subscriptionsPost
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsPostWithHttpInfo(response: ResponseContext): Promise<HttpInfo<SubscriptionSubscriptionModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("201", response.httpStatusCode)) {
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
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
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to subscriptionsSubscriptionIdDelete
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsSubscriptionIdDeleteWithHttpInfo(response: ResponseContext): Promise<HttpInfo<void >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("204", response.httpStatusCode)) {
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, undefined);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid subscription ID format", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Subscription not found", body, response.headers);
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
     * @params response Response returned by the server for a request to subscriptionsSubscriptionIdGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsSubscriptionIdGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<SubscriptionSubscriptionModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid subscription ID format", body, response.headers);
        }
        if (isCodeInRange("404", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Subscription not found", body, response.headers);
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
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to subscriptionsSubscriptionIdPut
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsSubscriptionIdPutWithHttpInfo(response: ResponseContext): Promise<HttpInfo<SubscriptionSubscriptionModel >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }
        if (isCodeInRange("400", response.httpStatusCode)) {
            const body: GinxHttpErrorResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "GinxHttpErrorResponse", ""
            ) as GinxHttpErrorResponse;
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Bad Request - Invalid input data or subscription ID", body, response.headers);
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
            throw new ApiException<GinxHttpErrorResponse>(response.httpStatusCode, "Subscription not found", body, response.headers);
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
            const body: SubscriptionSubscriptionModel = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionModel", ""
            ) as SubscriptionSubscriptionModel;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

    /**
     * Unwraps the actual response sent by the server from the response context and deserializes the response content
     * to the expected objects
     *
     * @params response Response returned by the server for a request to subscriptionsSummaryGet
     * @throws ApiException if the response code was not in [200, 299]
     */
     public async subscriptionsSummaryGetWithHttpInfo(response: ResponseContext): Promise<HttpInfo<SubscriptionSubscriptionSummaryResponse >> {
        const contentType = ObjectSerializer.normalizeMediaType(response.headers["content-type"]);
        if (isCodeInRange("200", response.httpStatusCode)) {
            const body: SubscriptionSubscriptionSummaryResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionSummaryResponse", ""
            ) as SubscriptionSubscriptionSummaryResponse;
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
            const body: SubscriptionSubscriptionSummaryResponse = ObjectSerializer.deserialize(
                ObjectSerializer.parse(await response.body.text(), contentType),
                "SubscriptionSubscriptionSummaryResponse", ""
            ) as SubscriptionSubscriptionSummaryResponse;
            return new HttpInfo(response.httpStatusCode, response.headers, response.body, body);
        }

        throw new ApiException<string | Blob | undefined>(response.httpStatusCode, "Unknown API Status Code!", await response.getBodyAsAny(), response.headers);
    }

}
