import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration, ConfigurationOptions, mergeConfiguration } from '../configuration'
import type { Middleware } from '../middleware';
import { Observable, of, from } from '../rxjsStub';
import {mergeMap, map} from  '../rxjsStub';
import { CurrencyCurrencyRateModel } from '../models/CurrencyCurrencyRateModel';
import { CurrencyCurrencyRatesModel } from '../models/CurrencyCurrencyRatesModel';
import { DtoAmountModel } from '../models/DtoAmountModel';
import { DtoEditableOwnerModel } from '../models/DtoEditableOwnerModel';
import { DtoOwnerModel } from '../models/DtoOwnerModel';
import { DtoPaginatedResponseModelLabelLabelModel } from '../models/DtoPaginatedResponseModelLabelLabelModel';
import { DtoPaginatedResponseModelProviderProviderModel } from '../models/DtoPaginatedResponseModelProviderProviderModel';
import { DtoPaginatedResponseModelSubscriptionSubscriptionModel } from '../models/DtoPaginatedResponseModelSubscriptionSubscriptionModel';
import { FamilyCreateFamilyMemberModel } from '../models/FamilyCreateFamilyMemberModel';
import { FamilyCreateFamilyModel } from '../models/FamilyCreateFamilyModel';
import { FamilyFamilyAcceptInvitationRequest } from '../models/FamilyFamilyAcceptInvitationRequest';
import { FamilyFamilyDeclineInvitationRequest } from '../models/FamilyFamilyDeclineInvitationRequest';
import { FamilyFamilyInviteRequest } from '../models/FamilyFamilyInviteRequest';
import { FamilyFamilyInviteResponse } from '../models/FamilyFamilyInviteResponse';
import { FamilyFamilyMemberModel } from '../models/FamilyFamilyMemberModel';
import { FamilyFamilyModel } from '../models/FamilyFamilyModel';
import { FamilyFamilySeeInvitationResponse } from '../models/FamilyFamilySeeInvitationResponse';
import { FamilyPatchFamilyMemberModel } from '../models/FamilyPatchFamilyMemberModel';
import { FamilyPatchFamilyModel } from '../models/FamilyPatchFamilyModel';
import { FamilyUpdateFamilyMemberModel } from '../models/FamilyUpdateFamilyMemberModel';
import { FamilyUpdateFamilyModel } from '../models/FamilyUpdateFamilyModel';
import { GinxHttpErrorResponse } from '../models/GinxHttpErrorResponse';
import { LabelCreateLabelModel } from '../models/LabelCreateLabelModel';
import { LabelLabelModel } from '../models/LabelLabelModel';
import { LabelUpdateLabelModel } from '../models/LabelUpdateLabelModel';
import { ProviderCreatePlanModel } from '../models/ProviderCreatePlanModel';
import { ProviderCreatePriceModel } from '../models/ProviderCreatePriceModel';
import { ProviderCreateProviderModel } from '../models/ProviderCreateProviderModel';
import { ProviderPlanModel } from '../models/ProviderPlanModel';
import { ProviderPriceModel } from '../models/ProviderPriceModel';
import { ProviderProviderModel } from '../models/ProviderProviderModel';
import { ProviderUpdatePlanModel } from '../models/ProviderUpdatePlanModel';
import { ProviderUpdatePriceModel } from '../models/ProviderUpdatePriceModel';
import { ProviderUpdateProviderModel } from '../models/ProviderUpdateProviderModel';
import { SubscriptionCreateSubscriptionModel } from '../models/SubscriptionCreateSubscriptionModel';
import { SubscriptionEditableSubscriptionPayerModel } from '../models/SubscriptionEditableSubscriptionPayerModel';
import { SubscriptionLabelRefModel } from '../models/SubscriptionLabelRefModel';
import { SubscriptionPatchSubscriptionModel } from '../models/SubscriptionPatchSubscriptionModel';
import { SubscriptionSubscriptionFreeTrialModel } from '../models/SubscriptionSubscriptionFreeTrialModel';
import { SubscriptionSubscriptionModel } from '../models/SubscriptionSubscriptionModel';
import { SubscriptionSubscriptionPayerModel } from '../models/SubscriptionSubscriptionPayerModel';
import { SubscriptionSubscriptionPriceModel } from '../models/SubscriptionSubscriptionPriceModel';
import { SubscriptionSubscriptionSummaryResponse } from '../models/SubscriptionSubscriptionSummaryResponse';
import { SubscriptionSubscriptionSummaryTopLabelResponse } from '../models/SubscriptionSubscriptionSummaryTopLabelResponse';
import { SubscriptionSubscriptionSummaryTopProviderResponse } from '../models/SubscriptionSubscriptionSummaryTopProviderResponse';
import { SubscriptionSubscriptionSummaryUpcomingRenewalResponse } from '../models/SubscriptionSubscriptionSummaryUpcomingRenewalResponse';
import { SubscriptionUpdateSubscriptionModel } from '../models/SubscriptionUpdateSubscriptionModel';
import { UserUpdatePreferredCurrencyModel } from '../models/UserUpdatePreferredCurrencyModel';
import { UserUserPreferredCurrencyModel } from '../models/UserUserPreferredCurrencyModel';

import { CurrenciesApiRequestFactory, CurrenciesApiResponseProcessor} from "../apis/CurrenciesApi";
export class ObservableCurrenciesApi {
    private requestFactory: CurrenciesApiRequestFactory;
    private responseProcessor: CurrenciesApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: CurrenciesApiRequestFactory,
        responseProcessor?: CurrenciesApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new CurrenciesApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new CurrenciesApiResponseProcessor();
    }

    /**
     * Get exchange rates for all currencies at a specific date
     * Get Currency Rates
     * @param [date] Conversion date in RFC3339 format (default: current time)
     */
    public currenciesRatesGetWithHttpInfo(date?: string, _options?: ConfigurationOptions): Observable<HttpInfo<CurrencyCurrencyRatesModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.currenciesRatesGet(date, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.currenciesRatesGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Get exchange rates for all currencies at a specific date
     * Get Currency Rates
     * @param [date] Conversion date in RFC3339 format (default: current time)
     */
    public currenciesRatesGet(date?: string, _options?: ConfigurationOptions): Observable<CurrencyCurrencyRatesModel> {
        return this.currenciesRatesGetWithHttpInfo(date, _options).pipe(map((apiResponse: HttpInfo<CurrencyCurrencyRatesModel>) => apiResponse.data));
    }

    /**
     * get details of all supported currencies
     * Get Supported Currencies
     */
    public currenciesSupportedGetWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<Array<string>>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.currenciesSupportedGet(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.currenciesSupportedGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * get details of all supported currencies
     * Get Supported Currencies
     */
    public currenciesSupportedGet(_options?: ConfigurationOptions): Observable<Array<string>> {
        return this.currenciesSupportedGetWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<Array<string>>) => apiResponse.data));
    }

}

import { FamilyApiRequestFactory, FamilyApiResponseProcessor} from "../apis/FamilyApi";
export class ObservableFamilyApi {
    private requestFactory: FamilyApiRequestFactory;
    private responseProcessor: FamilyApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: FamilyApiRequestFactory,
        responseProcessor?: FamilyApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new FamilyApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new FamilyApiResponseProcessor();
    }

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param familyId Family ID (UUID format)
     * @param familyFamilyAcceptInvitationRequest Invitation acceptance details
     */
    public familiesFamilyIdAcceptPostWithHttpInfo(familyId: string, familyFamilyAcceptInvitationRequest: FamilyFamilyAcceptInvitationRequest, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdAcceptPost(familyId, familyFamilyAcceptInvitationRequest, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdAcceptPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param familyId Family ID (UUID format)
     * @param familyFamilyAcceptInvitationRequest Invitation acceptance details
     */
    public familiesFamilyIdAcceptPost(familyId: string, familyFamilyAcceptInvitationRequest: FamilyFamilyAcceptInvitationRequest, _options?: ConfigurationOptions): Observable<void> {
        return this.familiesFamilyIdAcceptPostWithHttpInfo(familyId, familyFamilyAcceptInvitationRequest, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param familyId Family ID
     * @param familyFamilyDeclineInvitationRequest Decline invitation request
     */
    public familiesFamilyIdDeclinePostWithHttpInfo(familyId: string, familyFamilyDeclineInvitationRequest: FamilyFamilyDeclineInvitationRequest, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdDeclinePost(familyId, familyFamilyDeclineInvitationRequest, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdDeclinePostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param familyId Family ID
     * @param familyFamilyDeclineInvitationRequest Decline invitation request
     */
    public familiesFamilyIdDeclinePost(familyId: string, familyFamilyDeclineInvitationRequest: FamilyFamilyDeclineInvitationRequest, _options?: ConfigurationOptions): Observable<void> {
        return this.familiesFamilyIdDeclinePostWithHttpInfo(familyId, familyFamilyDeclineInvitationRequest, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param familyId Family ID (UUID format)
     */
    public familiesFamilyIdDeleteWithHttpInfo(familyId: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdDelete(familyId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param familyId Family ID (UUID format)
     */
    public familiesFamilyIdDelete(familyId: string, _options?: ConfigurationOptions): Observable<void> {
        return this.familiesFamilyIdDeleteWithHttpInfo(familyId, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param familyId Family ID
     * @param code Invitation code
     * @param familyMemberId Family member ID
     */
    public familiesFamilyIdInvitationGetWithHttpInfo(familyId: string, code: string, familyMemberId: string, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilySeeInvitationResponse>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdInvitationGet(familyId, code, familyMemberId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdInvitationGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param familyId Family ID
     * @param code Invitation code
     * @param familyMemberId Family member ID
     */
    public familiesFamilyIdInvitationGet(familyId: string, code: string, familyMemberId: string, _options?: ConfigurationOptions): Observable<FamilyFamilySeeInvitationResponse> {
        return this.familiesFamilyIdInvitationGetWithHttpInfo(familyId, code, familyMemberId, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilySeeInvitationResponse>) => apiResponse.data));
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param familyId Family ID (UUID format)
     * @param familyFamilyInviteRequest Invitation details including email, name, member ID and type (adult/kid)
     */
    public familiesFamilyIdInvitePostWithHttpInfo(familyId: string, familyFamilyInviteRequest: FamilyFamilyInviteRequest, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyInviteResponse>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdInvitePost(familyId, familyFamilyInviteRequest, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdInvitePostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param familyId Family ID (UUID format)
     * @param familyFamilyInviteRequest Invitation details including email, name, member ID and type (adult/kid)
     */
    public familiesFamilyIdInvitePost(familyId: string, familyFamilyInviteRequest: FamilyFamilyInviteRequest, _options?: ConfigurationOptions): Observable<FamilyFamilyInviteResponse> {
        return this.familiesFamilyIdInvitePostWithHttpInfo(familyId, familyFamilyInviteRequest, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilyInviteResponse>) => apiResponse.data));
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param familyId Family ID (UUID format)
     * @param familyMemberId Family Member ID (UUID format)
     * @param [body]
     */
    public familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(familyId: string, familyMemberId: string, body?: any, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdMembersFamilyMemberIdRevokePost(familyId, familyMemberId, body, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param familyId Family ID (UUID format)
     * @param familyMemberId Family Member ID (UUID format)
     * @param [body]
     */
    public familiesFamilyIdMembersFamilyMemberIdRevokePost(familyId: string, familyMemberId: string, body?: any, _options?: ConfigurationOptions): Observable<void> {
        return this.familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(familyId, familyMemberId, body, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     */
    public familiesFamilyIdMembersIdDeleteWithHttpInfo(familyId: string, id: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdMembersIdDelete(familyId, id, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdMembersIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     */
    public familiesFamilyIdMembersIdDelete(familyId: string, id: string, _options?: ConfigurationOptions): Observable<void> {
        return this.familiesFamilyIdMembersIdDeleteWithHttpInfo(familyId, id, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     * @param familyUpdateFamilyMemberModel Updated family member data
     */
    public familiesFamilyIdMembersIdPutWithHttpInfo(familyId: string, id: string, familyUpdateFamilyMemberModel: FamilyUpdateFamilyMemberModel, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdMembersIdPut(familyId, id, familyUpdateFamilyMemberModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdMembersIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     * @param familyUpdateFamilyMemberModel Updated family member data
     */
    public familiesFamilyIdMembersIdPut(familyId: string, id: string, familyUpdateFamilyMemberModel: FamilyUpdateFamilyMemberModel, _options?: ConfigurationOptions): Observable<FamilyFamilyModel> {
        return this.familiesFamilyIdMembersIdPutWithHttpInfo(familyId, id, familyUpdateFamilyMemberModel, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilyModel>) => apiResponse.data));
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param familyId Family ID (UUID format)
     * @param familyCreateFamilyMemberModel Family member creation data
     */
    public familiesFamilyIdMembersPostWithHttpInfo(familyId: string, familyCreateFamilyMemberModel: FamilyCreateFamilyMemberModel, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdMembersPost(familyId, familyCreateFamilyMemberModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdMembersPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param familyId Family ID (UUID format)
     * @param familyCreateFamilyMemberModel Family member creation data
     */
    public familiesFamilyIdMembersPost(familyId: string, familyCreateFamilyMemberModel: FamilyCreateFamilyMemberModel, _options?: ConfigurationOptions): Observable<FamilyFamilyModel> {
        return this.familiesFamilyIdMembersPostWithHttpInfo(familyId, familyCreateFamilyMemberModel, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilyModel>) => apiResponse.data));
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param familyId Family ID (UUID format)
     * @param familyUpdateFamilyModel Updated family data
     */
    public familiesFamilyIdPutWithHttpInfo(familyId: string, familyUpdateFamilyModel: FamilyUpdateFamilyModel, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesFamilyIdPut(familyId, familyUpdateFamilyModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesFamilyIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param familyId Family ID (UUID format)
     * @param familyUpdateFamilyModel Updated family data
     */
    public familiesFamilyIdPut(familyId: string, familyUpdateFamilyModel: FamilyUpdateFamilyModel, _options?: ConfigurationOptions): Observable<FamilyFamilyModel> {
        return this.familiesFamilyIdPutWithHttpInfo(familyId, familyUpdateFamilyModel, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilyModel>) => apiResponse.data));
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     */
    public familiesMeGetWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesMeGet(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesMeGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     */
    public familiesMeGet(_options?: ConfigurationOptions): Observable<FamilyFamilyModel> {
        return this.familiesMeGetWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<FamilyFamilyModel>) => apiResponse.data));
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param familyPatchFamilyModel Family update data with members
     */
    public familiesPatchWithHttpInfo(familyPatchFamilyModel: FamilyPatchFamilyModel, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesPatch(familyPatchFamilyModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesPatchWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param familyPatchFamilyModel Family update data with members
     */
    public familiesPatch(familyPatchFamilyModel: FamilyPatchFamilyModel, _options?: ConfigurationOptions): Observable<FamilyFamilyModel> {
        return this.familiesPatchWithHttpInfo(familyPatchFamilyModel, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilyModel>) => apiResponse.data));
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param familyCreateFamilyModel Family creation data
     */
    public familiesPostWithHttpInfo(familyCreateFamilyModel: FamilyCreateFamilyModel, _options?: ConfigurationOptions): Observable<HttpInfo<FamilyFamilyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.familiesPost(familyCreateFamilyModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.familiesPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param familyCreateFamilyModel Family creation data
     */
    public familiesPost(familyCreateFamilyModel: FamilyCreateFamilyModel, _options?: ConfigurationOptions): Observable<FamilyFamilyModel> {
        return this.familiesPostWithHttpInfo(familyCreateFamilyModel, _options).pipe(map((apiResponse: HttpInfo<FamilyFamilyModel>) => apiResponse.data));
    }

}

import { HealthApiRequestFactory, HealthApiResponseProcessor} from "../apis/HealthApi";
export class ObservableHealthApi {
    private requestFactory: HealthApiRequestFactory;
    private responseProcessor: HealthApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: HealthApiRequestFactory,
        responseProcessor?: HealthApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new HealthApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new HealthApiResponseProcessor();
    }

    /**
     * Returns the health status of the application
     * Health check endpoint
     */
    public healthzLiveGetWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<{ [key: string]: string; }>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.healthzLiveGet(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.healthzLiveGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Returns the health status of the application
     * Health check endpoint
     */
    public healthzLiveGet(_options?: ConfigurationOptions): Observable<{ [key: string]: string; }> {
        return this.healthzLiveGetWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<{ [key: string]: string; }>) => apiResponse.data));
    }

}

import { LabelsApiRequestFactory, LabelsApiResponseProcessor} from "../apis/LabelsApi";
export class ObservableLabelsApi {
    private requestFactory: LabelsApiRequestFactory;
    private responseProcessor: LabelsApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: LabelsApiRequestFactory,
        responseProcessor?: LabelsApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new LabelsApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new LabelsApiResponseProcessor();
    }

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     */
    public labelsDefaultGetWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<Array<LabelLabelModel>>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.labelsDefaultGet(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.labelsDefaultGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     */
    public labelsDefaultGet(_options?: ConfigurationOptions): Observable<Array<LabelLabelModel>> {
        return this.labelsDefaultGetWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<Array<LabelLabelModel>>) => apiResponse.data));
    }

    /**
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param [search] Search text to filter labels by name
     * @param [limit] Maximum number of items to return (default: 10)
     * @param [offset] Number of items to skip for pagination (default: 0)
     */
    public labelsGetWithHttpInfo(search?: string, limit?: number, offset?: number, _options?: ConfigurationOptions): Observable<HttpInfo<DtoPaginatedResponseModelLabelLabelModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.labelsGet(search, limit, offset, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.labelsGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param [search] Search text to filter labels by name
     * @param [limit] Maximum number of items to return (default: 10)
     * @param [offset] Number of items to skip for pagination (default: 0)
     */
    public labelsGet(search?: string, limit?: number, offset?: number, _options?: ConfigurationOptions): Observable<DtoPaginatedResponseModelLabelLabelModel> {
        return this.labelsGetWithHttpInfo(search, limit, offset, _options).pipe(map((apiResponse: HttpInfo<DtoPaginatedResponseModelLabelLabelModel>) => apiResponse.data));
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdDeleteWithHttpInfo(id: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.labelsIdDelete(id, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.labelsIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdDelete(id: string, _options?: ConfigurationOptions): Observable<void> {
        return this.labelsIdDeleteWithHttpInfo(id, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdGetWithHttpInfo(id: string, _options?: ConfigurationOptions): Observable<HttpInfo<LabelLabelModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.labelsIdGet(id, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.labelsIdGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdGet(id: string, _options?: ConfigurationOptions): Observable<LabelLabelModel> {
        return this.labelsIdGetWithHttpInfo(id, _options).pipe(map((apiResponse: HttpInfo<LabelLabelModel>) => apiResponse.data));
    }

    /**
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param id Label ID (UUID format)
     * @param labelUpdateLabelModel Updated label data
     */
    public labelsIdPutWithHttpInfo(id: string, labelUpdateLabelModel: LabelUpdateLabelModel, _options?: ConfigurationOptions): Observable<HttpInfo<LabelLabelModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.labelsIdPut(id, labelUpdateLabelModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.labelsIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param id Label ID (UUID format)
     * @param labelUpdateLabelModel Updated label data
     */
    public labelsIdPut(id: string, labelUpdateLabelModel: LabelUpdateLabelModel, _options?: ConfigurationOptions): Observable<LabelLabelModel> {
        return this.labelsIdPutWithHttpInfo(id, labelUpdateLabelModel, _options).pipe(map((apiResponse: HttpInfo<LabelLabelModel>) => apiResponse.data));
    }

    /**
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param labelCreateLabelModel Label creation data
     */
    public labelsPostWithHttpInfo(labelCreateLabelModel: LabelCreateLabelModel, _options?: ConfigurationOptions): Observable<HttpInfo<LabelLabelModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.labelsPost(labelCreateLabelModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.labelsPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param labelCreateLabelModel Label creation data
     */
    public labelsPost(labelCreateLabelModel: LabelCreateLabelModel, _options?: ConfigurationOptions): Observable<LabelLabelModel> {
        return this.labelsPostWithHttpInfo(labelCreateLabelModel, _options).pipe(map((apiResponse: HttpInfo<LabelLabelModel>) => apiResponse.data));
    }

}

import { ProvidersApiRequestFactory, ProvidersApiResponseProcessor} from "../apis/ProvidersApi";
export class ObservableProvidersApi {
    private requestFactory: ProvidersApiRequestFactory;
    private responseProcessor: ProvidersApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: ProvidersApiRequestFactory,
        responseProcessor?: ProvidersApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new ProvidersApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new ProvidersApiResponseProcessor();
    }

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param [search] Search term
     * @param [offset] Offset (default: 0)
     * @param [limit] Limit per request (default: 10)
     */
    public providersGetWithHttpInfo(search?: string, offset?: number, limit?: number, _options?: ConfigurationOptions): Observable<HttpInfo<DtoPaginatedResponseModelProviderProviderModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersGet(search, offset, limit, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param [search] Search term
     * @param [offset] Offset (default: 0)
     * @param [limit] Limit per request (default: 10)
     */
    public providersGet(search?: string, offset?: number, limit?: number, _options?: ConfigurationOptions): Observable<DtoPaginatedResponseModelProviderProviderModel> {
        return this.providersGetWithHttpInfo(search, offset, limit, _options).pipe(map((apiResponse: HttpInfo<DtoPaginatedResponseModelProviderProviderModel>) => apiResponse.data));
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param providerCreateProviderModel Provider creation data
     */
    public providersPostWithHttpInfo(providerCreateProviderModel: ProviderCreateProviderModel, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderProviderModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersPost(providerCreateProviderModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param providerCreateProviderModel Provider creation data
     */
    public providersPost(providerCreateProviderModel: ProviderCreateProviderModel, _options?: ConfigurationOptions): Observable<ProviderProviderModel> {
        return this.providersPostWithHttpInfo(providerCreateProviderModel, _options).pipe(map((apiResponse: HttpInfo<ProviderProviderModel>) => apiResponse.data));
    }

    /**
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdDeleteWithHttpInfo(providerId: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdDelete(providerId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdDelete(providerId: string, _options?: ConfigurationOptions): Observable<void> {
        return this.providersProviderIdDeleteWithHttpInfo(providerId, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdGetWithHttpInfo(providerId: string, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderProviderModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdGet(providerId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdGet(providerId: string, _options?: ConfigurationOptions): Observable<ProviderProviderModel> {
        return this.providersProviderIdGetWithHttpInfo(providerId, _options).pipe(map((apiResponse: HttpInfo<ProviderProviderModel>) => apiResponse.data));
    }

    /**
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     */
    public providersProviderIdPlansPlanIdDeleteWithHttpInfo(providerId: string, planId: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPlansPlanIdDelete(providerId, planId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPlansPlanIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     */
    public providersProviderIdPlansPlanIdDelete(providerId: string, planId: string, _options?: ConfigurationOptions): Observable<void> {
        return this.providersProviderIdPlansPlanIdDeleteWithHttpInfo(providerId, planId, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerCreatePriceModel Price creation data
     */
    public providersProviderIdPlansPlanIdPricesPostWithHttpInfo(providerId: string, planId: string, providerCreatePriceModel: ProviderCreatePriceModel, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderPriceModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPlansPlanIdPricesPost(providerId, planId, providerCreatePriceModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPlansPlanIdPricesPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerCreatePriceModel Price creation data
     */
    public providersProviderIdPlansPlanIdPricesPost(providerId: string, planId: string, providerCreatePriceModel: ProviderCreatePriceModel, _options?: ConfigurationOptions): Observable<ProviderPriceModel> {
        return this.providersProviderIdPlansPlanIdPricesPostWithHttpInfo(providerId, planId, providerCreatePriceModel, _options).pipe(map((apiResponse: HttpInfo<ProviderPriceModel>) => apiResponse.data));
    }

    /**
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     */
    public providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(providerId: string, planId: string, priceId: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPlansPlanIdPricesPriceIdDelete(providerId, planId, priceId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     */
    public providersProviderIdPlansPlanIdPricesPriceIdDelete(providerId: string, planId: string, priceId: string, _options?: ConfigurationOptions): Observable<void> {
        return this.providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(providerId, planId, priceId, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     * @param providerUpdatePriceModel Updated price data
     */
    public providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(providerId: string, planId: string, priceId: string, providerUpdatePriceModel: ProviderUpdatePriceModel, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderPriceModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPlansPlanIdPricesPriceIdPut(providerId, planId, priceId, providerUpdatePriceModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     * @param providerUpdatePriceModel Updated price data
     */
    public providersProviderIdPlansPlanIdPricesPriceIdPut(providerId: string, planId: string, priceId: string, providerUpdatePriceModel: ProviderUpdatePriceModel, _options?: ConfigurationOptions): Observable<ProviderPriceModel> {
        return this.providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(providerId, planId, priceId, providerUpdatePriceModel, _options).pipe(map((apiResponse: HttpInfo<ProviderPriceModel>) => apiResponse.data));
    }

    /**
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerUpdatePlanModel Updated plan data
     */
    public providersProviderIdPlansPlanIdPutWithHttpInfo(providerId: string, planId: string, providerUpdatePlanModel: ProviderUpdatePlanModel, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderPlanModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPlansPlanIdPut(providerId, planId, providerUpdatePlanModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPlansPlanIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerUpdatePlanModel Updated plan data
     */
    public providersProviderIdPlansPlanIdPut(providerId: string, planId: string, providerUpdatePlanModel: ProviderUpdatePlanModel, _options?: ConfigurationOptions): Observable<ProviderPlanModel> {
        return this.providersProviderIdPlansPlanIdPutWithHttpInfo(providerId, planId, providerUpdatePlanModel, _options).pipe(map((apiResponse: HttpInfo<ProviderPlanModel>) => apiResponse.data));
    }

    /**
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param providerId Provider ID (UUID format)
     * @param providerCreatePlanModel Plan creation data
     */
    public providersProviderIdPlansPostWithHttpInfo(providerId: string, providerCreatePlanModel: ProviderCreatePlanModel, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderPlanModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPlansPost(providerId, providerCreatePlanModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPlansPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param providerId Provider ID (UUID format)
     * @param providerCreatePlanModel Plan creation data
     */
    public providersProviderIdPlansPost(providerId: string, providerCreatePlanModel: ProviderCreatePlanModel, _options?: ConfigurationOptions): Observable<ProviderPlanModel> {
        return this.providersProviderIdPlansPostWithHttpInfo(providerId, providerCreatePlanModel, _options).pipe(map((apiResponse: HttpInfo<ProviderPlanModel>) => apiResponse.data));
    }

    /**
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param providerId Provider ID (UUID format)
     * @param providerUpdateProviderModel Updated provider data
     */
    public providersProviderIdPutWithHttpInfo(providerId: string, providerUpdateProviderModel: ProviderUpdateProviderModel, _options?: ConfigurationOptions): Observable<HttpInfo<ProviderProviderModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.providersProviderIdPut(providerId, providerUpdateProviderModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.providersProviderIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param providerId Provider ID (UUID format)
     * @param providerUpdateProviderModel Updated provider data
     */
    public providersProviderIdPut(providerId: string, providerUpdateProviderModel: ProviderUpdateProviderModel, _options?: ConfigurationOptions): Observable<ProviderProviderModel> {
        return this.providersProviderIdPutWithHttpInfo(providerId, providerUpdateProviderModel, _options).pipe(map((apiResponse: HttpInfo<ProviderProviderModel>) => apiResponse.data));
    }

}

import { SubscriptionsApiRequestFactory, SubscriptionsApiResponseProcessor} from "../apis/SubscriptionsApi";
export class ObservableSubscriptionsApi {
    private requestFactory: SubscriptionsApiRequestFactory;
    private responseProcessor: SubscriptionsApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: SubscriptionsApiRequestFactory,
        responseProcessor?: SubscriptionsApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new SubscriptionsApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new SubscriptionsApiResponseProcessor();
    }

    /**
     * Retrieve a paginated list of all subscriptions for the authenticated user
     * Get all subscriptions
     * @param [search] Search text
     * @param [recurrencies] Filter by recurrency types
     * @param [fromDate] Filter by start date (RFC3339)
     * @param [toDate] Filter by end date (RFC3339)
     * @param [users] Filter by user IDs
     * @param [withInactive] Include inactive subscriptions
     * @param [providers] Filter by provider IDs
     * @param [limit] Number of items per page (default: 10)
     * @param [offset] Page number (default: 0)
     */
    public subscriptionsGetWithHttpInfo(search?: string, recurrencies?: Array<string>, fromDate?: string, toDate?: string, users?: Array<string>, withInactive?: boolean, providers?: Array<string>, limit?: number, offset?: number, _options?: ConfigurationOptions): Observable<HttpInfo<DtoPaginatedResponseModelSubscriptionSubscriptionModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsGet(search, recurrencies, fromDate, toDate, users, withInactive, providers, limit, offset, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve a paginated list of all subscriptions for the authenticated user
     * Get all subscriptions
     * @param [search] Search text
     * @param [recurrencies] Filter by recurrency types
     * @param [fromDate] Filter by start date (RFC3339)
     * @param [toDate] Filter by end date (RFC3339)
     * @param [users] Filter by user IDs
     * @param [withInactive] Include inactive subscriptions
     * @param [providers] Filter by provider IDs
     * @param [limit] Number of items per page (default: 10)
     * @param [offset] Page number (default: 0)
     */
    public subscriptionsGet(search?: string, recurrencies?: Array<string>, fromDate?: string, toDate?: string, users?: Array<string>, withInactive?: boolean, providers?: Array<string>, limit?: number, offset?: number, _options?: ConfigurationOptions): Observable<DtoPaginatedResponseModelSubscriptionSubscriptionModel> {
        return this.subscriptionsGetWithHttpInfo(search, recurrencies, fromDate, toDate, users, withInactive, providers, limit, offset, _options).pipe(map((apiResponse: HttpInfo<DtoPaginatedResponseModelSubscriptionSubscriptionModel>) => apiResponse.data));
    }

    /**
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param subscriptionPatchSubscriptionModel Complete subscription data
     */
    public subscriptionsPatchWithHttpInfo(subscriptionPatchSubscriptionModel: SubscriptionPatchSubscriptionModel, _options?: ConfigurationOptions): Observable<HttpInfo<SubscriptionSubscriptionModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsPatch(subscriptionPatchSubscriptionModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsPatchWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param subscriptionPatchSubscriptionModel Complete subscription data
     */
    public subscriptionsPatch(subscriptionPatchSubscriptionModel: SubscriptionPatchSubscriptionModel, _options?: ConfigurationOptions): Observable<SubscriptionSubscriptionModel> {
        return this.subscriptionsPatchWithHttpInfo(subscriptionPatchSubscriptionModel, _options).pipe(map((apiResponse: HttpInfo<SubscriptionSubscriptionModel>) => apiResponse.data));
    }

    /**
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param subscriptionCreateSubscriptionModel Subscription creation data
     */
    public subscriptionsPostWithHttpInfo(subscriptionCreateSubscriptionModel: SubscriptionCreateSubscriptionModel, _options?: ConfigurationOptions): Observable<HttpInfo<SubscriptionSubscriptionModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsPost(subscriptionCreateSubscriptionModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsPostWithHttpInfo(rsp)));
            }));
    }

    /**
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param subscriptionCreateSubscriptionModel Subscription creation data
     */
    public subscriptionsPost(subscriptionCreateSubscriptionModel: SubscriptionCreateSubscriptionModel, _options?: ConfigurationOptions): Observable<SubscriptionSubscriptionModel> {
        return this.subscriptionsPostWithHttpInfo(subscriptionCreateSubscriptionModel, _options).pipe(map((apiResponse: HttpInfo<SubscriptionSubscriptionModel>) => apiResponse.data));
    }

    /**
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdDeleteWithHttpInfo(subscriptionId: string, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsSubscriptionIdDelete(subscriptionId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsSubscriptionIdDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdDelete(subscriptionId: string, _options?: ConfigurationOptions): Observable<void> {
        return this.subscriptionsSubscriptionIdDeleteWithHttpInfo(subscriptionId, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdGetWithHttpInfo(subscriptionId: string, _options?: ConfigurationOptions): Observable<HttpInfo<SubscriptionSubscriptionModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsSubscriptionIdGet(subscriptionId, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsSubscriptionIdGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdGet(subscriptionId: string, _options?: ConfigurationOptions): Observable<SubscriptionSubscriptionModel> {
        return this.subscriptionsSubscriptionIdGetWithHttpInfo(subscriptionId, _options).pipe(map((apiResponse: HttpInfo<SubscriptionSubscriptionModel>) => apiResponse.data));
    }

    /**
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     * @param subscriptionUpdateSubscriptionModel Updated subscription data
     */
    public subscriptionsSubscriptionIdPutWithHttpInfo(subscriptionId: string, subscriptionUpdateSubscriptionModel: SubscriptionUpdateSubscriptionModel, _options?: ConfigurationOptions): Observable<HttpInfo<SubscriptionSubscriptionModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsSubscriptionIdPut(subscriptionId, subscriptionUpdateSubscriptionModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsSubscriptionIdPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     * @param subscriptionUpdateSubscriptionModel Updated subscription data
     */
    public subscriptionsSubscriptionIdPut(subscriptionId: string, subscriptionUpdateSubscriptionModel: SubscriptionUpdateSubscriptionModel, _options?: ConfigurationOptions): Observable<SubscriptionSubscriptionModel> {
        return this.subscriptionsSubscriptionIdPutWithHttpInfo(subscriptionId, subscriptionUpdateSubscriptionModel, _options).pipe(map((apiResponse: HttpInfo<SubscriptionSubscriptionModel>) => apiResponse.data));
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
    public subscriptionsSummaryGetWithHttpInfo(topProviders: number, topLabels: number, upcomingRenewals: number, totalMonthly: boolean, totalYearly: boolean, _options?: ConfigurationOptions): Observable<HttpInfo<SubscriptionSubscriptionSummaryResponse>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.subscriptionsSummaryGet(topProviders, topLabels, upcomingRenewals, totalMonthly, totalYearly, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.subscriptionsSummaryGetWithHttpInfo(rsp)));
            }));
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
    public subscriptionsSummaryGet(topProviders: number, topLabels: number, upcomingRenewals: number, totalMonthly: boolean, totalYearly: boolean, _options?: ConfigurationOptions): Observable<SubscriptionSubscriptionSummaryResponse> {
        return this.subscriptionsSummaryGetWithHttpInfo(topProviders, topLabels, upcomingRenewals, totalMonthly, totalYearly, _options).pipe(map((apiResponse: HttpInfo<SubscriptionSubscriptionSummaryResponse>) => apiResponse.data));
    }

}

import { UsersApiRequestFactory, UsersApiResponseProcessor} from "../apis/UsersApi";
export class ObservableUsersApi {
    private requestFactory: UsersApiRequestFactory;
    private responseProcessor: UsersApiResponseProcessor;
    private configuration: Configuration;

    public constructor(
        configuration: Configuration,
        requestFactory?: UsersApiRequestFactory,
        responseProcessor?: UsersApiResponseProcessor
    ) {
        this.configuration = configuration;
        this.requestFactory = requestFactory || new UsersApiRequestFactory(configuration);
        this.responseProcessor = responseProcessor || new UsersApiResponseProcessor();
    }

    /**
     * Deletes the authenticated user\'s account
     * Delete user
     */
    public usersDeleteWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.usersDelete(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.usersDeleteWithHttpInfo(rsp)));
            }));
    }

    /**
     * Deletes the authenticated user\'s account
     * Delete user
     */
    public usersDelete(_options?: ConfigurationOptions): Observable<void> {
        return this.usersDeleteWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

    /**
     * Returns the preferred currency for the authenticated user
     * Get user preferred currency
     */
    public usersPreferredCurrencyGetWithHttpInfo(_options?: ConfigurationOptions): Observable<HttpInfo<UserUserPreferredCurrencyModel>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.usersPreferredCurrencyGet(_config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.usersPreferredCurrencyGetWithHttpInfo(rsp)));
            }));
    }

    /**
     * Returns the preferred currency for the authenticated user
     * Get user preferred currency
     */
    public usersPreferredCurrencyGet(_options?: ConfigurationOptions): Observable<UserUserPreferredCurrencyModel> {
        return this.usersPreferredCurrencyGetWithHttpInfo(_options).pipe(map((apiResponse: HttpInfo<UserUserPreferredCurrencyModel>) => apiResponse.data));
    }

    /**
     * Updates the preferred currency for the authenticated user
     * Update user preferred currency
     * @param authorization Bearer token
     * @param userUpdatePreferredCurrencyModel Profile update parameters
     */
    public usersPreferredCurrencyPutWithHttpInfo(authorization: string, userUpdatePreferredCurrencyModel: UserUpdatePreferredCurrencyModel, _options?: ConfigurationOptions): Observable<HttpInfo<void>> {
        const _config = mergeConfiguration(this.configuration, _options);

        const requestContextPromise = this.requestFactory.usersPreferredCurrencyPut(authorization, userUpdatePreferredCurrencyModel, _config);
        // build promise chain
        let middlewarePreObservable = from<RequestContext>(requestContextPromise);
        for (const middleware of _config.middleware) {
            middlewarePreObservable = middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => middleware.pre(ctx)));
        }

        return middlewarePreObservable.pipe(mergeMap((ctx: RequestContext) => _config.httpApi.send(ctx))).
            pipe(mergeMap((response: ResponseContext) => {
                let middlewarePostObservable = of(response);
                for (const middleware of _config.middleware.reverse()) {
                    middlewarePostObservable = middlewarePostObservable.pipe(mergeMap((rsp: ResponseContext) => middleware.post(rsp)));
                }
                return middlewarePostObservable.pipe(map((rsp: ResponseContext) => this.responseProcessor.usersPreferredCurrencyPutWithHttpInfo(rsp)));
            }));
    }

    /**
     * Updates the preferred currency for the authenticated user
     * Update user preferred currency
     * @param authorization Bearer token
     * @param userUpdatePreferredCurrencyModel Profile update parameters
     */
    public usersPreferredCurrencyPut(authorization: string, userUpdatePreferredCurrencyModel: UserUpdatePreferredCurrencyModel, _options?: ConfigurationOptions): Observable<void> {
        return this.usersPreferredCurrencyPutWithHttpInfo(authorization, userUpdatePreferredCurrencyModel, _options).pipe(map((apiResponse: HttpInfo<void>) => apiResponse.data));
    }

}
