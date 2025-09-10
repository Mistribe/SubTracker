import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration, PromiseConfigurationOptions, wrapOptions } from '../configuration'
import { PromiseMiddleware, Middleware, PromiseMiddlewareWrapper } from '../middleware';

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
import { ObservableCurrenciesApi } from './ObservableAPI';

import { CurrenciesApiRequestFactory, CurrenciesApiResponseProcessor} from "../apis/CurrenciesApi";
export class PromiseCurrenciesApi {
    private api: ObservableCurrenciesApi

    public constructor(
        configuration: Configuration,
        requestFactory?: CurrenciesApiRequestFactory,
        responseProcessor?: CurrenciesApiResponseProcessor
    ) {
        this.api = new ObservableCurrenciesApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Get exchange rates for all currencies at a specific date
     * Get Currency Rates
     * @param [date] Conversion date in RFC3339 format (default: current time)
     */
    public currenciesRatesGetWithHttpInfo(date?: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<CurrencyCurrencyRatesModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.currenciesRatesGetWithHttpInfo(date, observableOptions);
        return result.toPromise();
    }

    /**
     * Get exchange rates for all currencies at a specific date
     * Get Currency Rates
     * @param [date] Conversion date in RFC3339 format (default: current time)
     */
    public currenciesRatesGet(date?: string, _options?: PromiseConfigurationOptions): Promise<CurrencyCurrencyRatesModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.currenciesRatesGet(date, observableOptions);
        return result.toPromise();
    }

    /**
     * get details of all supported currencies
     * Get Supported Currencies
     */
    public currenciesSupportedGetWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<Array<string>>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.currenciesSupportedGetWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * get details of all supported currencies
     * Get Supported Currencies
     */
    public currenciesSupportedGet(_options?: PromiseConfigurationOptions): Promise<Array<string>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.currenciesSupportedGet(observableOptions);
        return result.toPromise();
    }


}



import { ObservableFamilyApi } from './ObservableAPI';

import { FamilyApiRequestFactory, FamilyApiResponseProcessor} from "../apis/FamilyApi";
export class PromiseFamilyApi {
    private api: ObservableFamilyApi

    public constructor(
        configuration: Configuration,
        requestFactory?: FamilyApiRequestFactory,
        responseProcessor?: FamilyApiResponseProcessor
    ) {
        this.api = new ObservableFamilyApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param familyId Family ID (UUID format)
     * @param familyFamilyAcceptInvitationRequest Invitation acceptance details
     */
    public familiesFamilyIdAcceptPostWithHttpInfo(familyId: string, familyFamilyAcceptInvitationRequest: FamilyFamilyAcceptInvitationRequest, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdAcceptPostWithHttpInfo(familyId, familyFamilyAcceptInvitationRequest, observableOptions);
        return result.toPromise();
    }

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param familyId Family ID (UUID format)
     * @param familyFamilyAcceptInvitationRequest Invitation acceptance details
     */
    public familiesFamilyIdAcceptPost(familyId: string, familyFamilyAcceptInvitationRequest: FamilyFamilyAcceptInvitationRequest, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdAcceptPost(familyId, familyFamilyAcceptInvitationRequest, observableOptions);
        return result.toPromise();
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param familyId Family ID
     * @param familyFamilyDeclineInvitationRequest Decline invitation request
     */
    public familiesFamilyIdDeclinePostWithHttpInfo(familyId: string, familyFamilyDeclineInvitationRequest: FamilyFamilyDeclineInvitationRequest, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdDeclinePostWithHttpInfo(familyId, familyFamilyDeclineInvitationRequest, observableOptions);
        return result.toPromise();
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param familyId Family ID
     * @param familyFamilyDeclineInvitationRequest Decline invitation request
     */
    public familiesFamilyIdDeclinePost(familyId: string, familyFamilyDeclineInvitationRequest: FamilyFamilyDeclineInvitationRequest, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdDeclinePost(familyId, familyFamilyDeclineInvitationRequest, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param familyId Family ID (UUID format)
     */
    public familiesFamilyIdDeleteWithHttpInfo(familyId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdDeleteWithHttpInfo(familyId, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param familyId Family ID (UUID format)
     */
    public familiesFamilyIdDelete(familyId: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdDelete(familyId, observableOptions);
        return result.toPromise();
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param familyId Family ID
     * @param code Invitation code
     * @param familyMemberId Family member ID
     */
    public familiesFamilyIdInvitationGetWithHttpInfo(familyId: string, code: string, familyMemberId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilySeeInvitationResponse>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdInvitationGetWithHttpInfo(familyId, code, familyMemberId, observableOptions);
        return result.toPromise();
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param familyId Family ID
     * @param code Invitation code
     * @param familyMemberId Family member ID
     */
    public familiesFamilyIdInvitationGet(familyId: string, code: string, familyMemberId: string, _options?: PromiseConfigurationOptions): Promise<FamilyFamilySeeInvitationResponse> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdInvitationGet(familyId, code, familyMemberId, observableOptions);
        return result.toPromise();
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param familyId Family ID (UUID format)
     * @param familyFamilyInviteRequest Invitation details including email, name, member ID and type (adult/kid)
     */
    public familiesFamilyIdInvitePostWithHttpInfo(familyId: string, familyFamilyInviteRequest: FamilyFamilyInviteRequest, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyInviteResponse>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdInvitePostWithHttpInfo(familyId, familyFamilyInviteRequest, observableOptions);
        return result.toPromise();
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param familyId Family ID (UUID format)
     * @param familyFamilyInviteRequest Invitation details including email, name, member ID and type (adult/kid)
     */
    public familiesFamilyIdInvitePost(familyId: string, familyFamilyInviteRequest: FamilyFamilyInviteRequest, _options?: PromiseConfigurationOptions): Promise<FamilyFamilyInviteResponse> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdInvitePost(familyId, familyFamilyInviteRequest, observableOptions);
        return result.toPromise();
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param familyId Family ID (UUID format)
     * @param familyMemberId Family Member ID (UUID format)
     * @param [body]
     */
    public familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(familyId: string, familyMemberId: string, body?: any, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(familyId, familyMemberId, body, observableOptions);
        return result.toPromise();
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param familyId Family ID (UUID format)
     * @param familyMemberId Family Member ID (UUID format)
     * @param [body]
     */
    public familiesFamilyIdMembersFamilyMemberIdRevokePost(familyId: string, familyMemberId: string, body?: any, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersFamilyMemberIdRevokePost(familyId, familyMemberId, body, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     */
    public familiesFamilyIdMembersIdDeleteWithHttpInfo(familyId: string, id: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersIdDeleteWithHttpInfo(familyId, id, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     */
    public familiesFamilyIdMembersIdDelete(familyId: string, id: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersIdDelete(familyId, id, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     * @param familyUpdateFamilyMemberModel Updated family member data
     */
    public familiesFamilyIdMembersIdPutWithHttpInfo(familyId: string, id: string, familyUpdateFamilyMemberModel: FamilyUpdateFamilyMemberModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersIdPutWithHttpInfo(familyId, id, familyUpdateFamilyMemberModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param familyId Family ID (UUID format)
     * @param id Family member ID (UUID format)
     * @param familyUpdateFamilyMemberModel Updated family member data
     */
    public familiesFamilyIdMembersIdPut(familyId: string, id: string, familyUpdateFamilyMemberModel: FamilyUpdateFamilyMemberModel, _options?: PromiseConfigurationOptions): Promise<FamilyFamilyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersIdPut(familyId, id, familyUpdateFamilyMemberModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param familyId Family ID (UUID format)
     * @param familyCreateFamilyMemberModel Family member creation data
     */
    public familiesFamilyIdMembersPostWithHttpInfo(familyId: string, familyCreateFamilyMemberModel: FamilyCreateFamilyMemberModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersPostWithHttpInfo(familyId, familyCreateFamilyMemberModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param familyId Family ID (UUID format)
     * @param familyCreateFamilyMemberModel Family member creation data
     */
    public familiesFamilyIdMembersPost(familyId: string, familyCreateFamilyMemberModel: FamilyCreateFamilyMemberModel, _options?: PromiseConfigurationOptions): Promise<FamilyFamilyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdMembersPost(familyId, familyCreateFamilyMemberModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param familyId Family ID (UUID format)
     * @param familyUpdateFamilyModel Updated family data
     */
    public familiesFamilyIdPutWithHttpInfo(familyId: string, familyUpdateFamilyModel: FamilyUpdateFamilyModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdPutWithHttpInfo(familyId, familyUpdateFamilyModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param familyId Family ID (UUID format)
     * @param familyUpdateFamilyModel Updated family data
     */
    public familiesFamilyIdPut(familyId: string, familyUpdateFamilyModel: FamilyUpdateFamilyModel, _options?: PromiseConfigurationOptions): Promise<FamilyFamilyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesFamilyIdPut(familyId, familyUpdateFamilyModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     */
    public familiesMeGetWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesMeGetWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     */
    public familiesMeGet(_options?: PromiseConfigurationOptions): Promise<FamilyFamilyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesMeGet(observableOptions);
        return result.toPromise();
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param familyPatchFamilyModel Family update data with members
     */
    public familiesPatchWithHttpInfo(familyPatchFamilyModel: FamilyPatchFamilyModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesPatchWithHttpInfo(familyPatchFamilyModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param familyPatchFamilyModel Family update data with members
     */
    public familiesPatch(familyPatchFamilyModel: FamilyPatchFamilyModel, _options?: PromiseConfigurationOptions): Promise<FamilyFamilyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesPatch(familyPatchFamilyModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param familyCreateFamilyModel Family creation data
     */
    public familiesPostWithHttpInfo(familyCreateFamilyModel: FamilyCreateFamilyModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesPostWithHttpInfo(familyCreateFamilyModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param familyCreateFamilyModel Family creation data
     */
    public familiesPost(familyCreateFamilyModel: FamilyCreateFamilyModel, _options?: PromiseConfigurationOptions): Promise<FamilyFamilyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.familiesPost(familyCreateFamilyModel, observableOptions);
        return result.toPromise();
    }


}



import { ObservableHealthApi } from './ObservableAPI';

import { HealthApiRequestFactory, HealthApiResponseProcessor} from "../apis/HealthApi";
export class PromiseHealthApi {
    private api: ObservableHealthApi

    public constructor(
        configuration: Configuration,
        requestFactory?: HealthApiRequestFactory,
        responseProcessor?: HealthApiResponseProcessor
    ) {
        this.api = new ObservableHealthApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns the health status of the application
     * Health check endpoint
     */
    public healthzLiveGetWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<{ [key: string]: string; }>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.healthzLiveGetWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Returns the health status of the application
     * Health check endpoint
     */
    public healthzLiveGet(_options?: PromiseConfigurationOptions): Promise<{ [key: string]: string; }> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.healthzLiveGet(observableOptions);
        return result.toPromise();
    }


}



import { ObservableLabelsApi } from './ObservableAPI';

import { LabelsApiRequestFactory, LabelsApiResponseProcessor} from "../apis/LabelsApi";
export class PromiseLabelsApi {
    private api: ObservableLabelsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: LabelsApiRequestFactory,
        responseProcessor?: LabelsApiResponseProcessor
    ) {
        this.api = new ObservableLabelsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     */
    public labelsDefaultGetWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<Array<LabelLabelModel>>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsDefaultGetWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     */
    public labelsDefaultGet(_options?: PromiseConfigurationOptions): Promise<Array<LabelLabelModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsDefaultGet(observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param [search] Search text to filter labels by name
     * @param [limit] Maximum number of items to return (default: 10)
     * @param [offset] Number of items to skip for pagination (default: 0)
     */
    public labelsGetWithHttpInfo(search?: string, limit?: number, offset?: number, _options?: PromiseConfigurationOptions): Promise<HttpInfo<DtoPaginatedResponseModelLabelLabelModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsGetWithHttpInfo(search, limit, offset, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param [search] Search text to filter labels by name
     * @param [limit] Maximum number of items to return (default: 10)
     * @param [offset] Number of items to skip for pagination (default: 0)
     */
    public labelsGet(search?: string, limit?: number, offset?: number, _options?: PromiseConfigurationOptions): Promise<DtoPaginatedResponseModelLabelLabelModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsGet(search, limit, offset, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdDeleteWithHttpInfo(id: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsIdDeleteWithHttpInfo(id, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdDelete(id: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsIdDelete(id, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdGetWithHttpInfo(id: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<LabelLabelModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsIdGetWithHttpInfo(id, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param id Label ID (UUID format)
     */
    public labelsIdGet(id: string, _options?: PromiseConfigurationOptions): Promise<LabelLabelModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsIdGet(id, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param id Label ID (UUID format)
     * @param labelUpdateLabelModel Updated label data
     */
    public labelsIdPutWithHttpInfo(id: string, labelUpdateLabelModel: LabelUpdateLabelModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<LabelLabelModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsIdPutWithHttpInfo(id, labelUpdateLabelModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param id Label ID (UUID format)
     * @param labelUpdateLabelModel Updated label data
     */
    public labelsIdPut(id: string, labelUpdateLabelModel: LabelUpdateLabelModel, _options?: PromiseConfigurationOptions): Promise<LabelLabelModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsIdPut(id, labelUpdateLabelModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param labelCreateLabelModel Label creation data
     */
    public labelsPostWithHttpInfo(labelCreateLabelModel: LabelCreateLabelModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<LabelLabelModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsPostWithHttpInfo(labelCreateLabelModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param labelCreateLabelModel Label creation data
     */
    public labelsPost(labelCreateLabelModel: LabelCreateLabelModel, _options?: PromiseConfigurationOptions): Promise<LabelLabelModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.labelsPost(labelCreateLabelModel, observableOptions);
        return result.toPromise();
    }


}



import { ObservableProvidersApi } from './ObservableAPI';

import { ProvidersApiRequestFactory, ProvidersApiResponseProcessor} from "../apis/ProvidersApi";
export class PromiseProvidersApi {
    private api: ObservableProvidersApi

    public constructor(
        configuration: Configuration,
        requestFactory?: ProvidersApiRequestFactory,
        responseProcessor?: ProvidersApiResponseProcessor
    ) {
        this.api = new ObservableProvidersApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param [search] Search term
     * @param [offset] Offset (default: 0)
     * @param [limit] Limit per request (default: 10)
     */
    public providersGetWithHttpInfo(search?: string, offset?: number, limit?: number, _options?: PromiseConfigurationOptions): Promise<HttpInfo<DtoPaginatedResponseModelProviderProviderModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersGetWithHttpInfo(search, offset, limit, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param [search] Search term
     * @param [offset] Offset (default: 0)
     * @param [limit] Limit per request (default: 10)
     */
    public providersGet(search?: string, offset?: number, limit?: number, _options?: PromiseConfigurationOptions): Promise<DtoPaginatedResponseModelProviderProviderModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersGet(search, offset, limit, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param providerCreateProviderModel Provider creation data
     */
    public providersPostWithHttpInfo(providerCreateProviderModel: ProviderCreateProviderModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderProviderModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersPostWithHttpInfo(providerCreateProviderModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param providerCreateProviderModel Provider creation data
     */
    public providersPost(providerCreateProviderModel: ProviderCreateProviderModel, _options?: PromiseConfigurationOptions): Promise<ProviderProviderModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersPost(providerCreateProviderModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdDeleteWithHttpInfo(providerId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdDeleteWithHttpInfo(providerId, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdDelete(providerId: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdDelete(providerId, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdGetWithHttpInfo(providerId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderProviderModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdGetWithHttpInfo(providerId, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param providerId Provider ID (UUID format)
     */
    public providersProviderIdGet(providerId: string, _options?: PromiseConfigurationOptions): Promise<ProviderProviderModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdGet(providerId, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     */
    public providersProviderIdPlansPlanIdDeleteWithHttpInfo(providerId: string, planId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdDeleteWithHttpInfo(providerId, planId, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     */
    public providersProviderIdPlansPlanIdDelete(providerId: string, planId: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdDelete(providerId, planId, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerCreatePriceModel Price creation data
     */
    public providersProviderIdPlansPlanIdPricesPostWithHttpInfo(providerId: string, planId: string, providerCreatePriceModel: ProviderCreatePriceModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderPriceModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPricesPostWithHttpInfo(providerId, planId, providerCreatePriceModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerCreatePriceModel Price creation data
     */
    public providersProviderIdPlansPlanIdPricesPost(providerId: string, planId: string, providerCreatePriceModel: ProviderCreatePriceModel, _options?: PromiseConfigurationOptions): Promise<ProviderPriceModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPricesPost(providerId, planId, providerCreatePriceModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     */
    public providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(providerId: string, planId: string, priceId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(providerId, planId, priceId, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     */
    public providersProviderIdPlansPlanIdPricesPriceIdDelete(providerId: string, planId: string, priceId: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPricesPriceIdDelete(providerId, planId, priceId, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     * @param providerUpdatePriceModel Updated price data
     */
    public providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(providerId: string, planId: string, priceId: string, providerUpdatePriceModel: ProviderUpdatePriceModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderPriceModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(providerId, planId, priceId, providerUpdatePriceModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param priceId Price ID (UUID format)
     * @param providerUpdatePriceModel Updated price data
     */
    public providersProviderIdPlansPlanIdPricesPriceIdPut(providerId: string, planId: string, priceId: string, providerUpdatePriceModel: ProviderUpdatePriceModel, _options?: PromiseConfigurationOptions): Promise<ProviderPriceModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPricesPriceIdPut(providerId, planId, priceId, providerUpdatePriceModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerUpdatePlanModel Updated plan data
     */
    public providersProviderIdPlansPlanIdPutWithHttpInfo(providerId: string, planId: string, providerUpdatePlanModel: ProviderUpdatePlanModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderPlanModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPutWithHttpInfo(providerId, planId, providerUpdatePlanModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param providerId Provider ID (UUID format)
     * @param planId Plan ID (UUID format)
     * @param providerUpdatePlanModel Updated plan data
     */
    public providersProviderIdPlansPlanIdPut(providerId: string, planId: string, providerUpdatePlanModel: ProviderUpdatePlanModel, _options?: PromiseConfigurationOptions): Promise<ProviderPlanModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPlanIdPut(providerId, planId, providerUpdatePlanModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param providerId Provider ID (UUID format)
     * @param providerCreatePlanModel Plan creation data
     */
    public providersProviderIdPlansPostWithHttpInfo(providerId: string, providerCreatePlanModel: ProviderCreatePlanModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderPlanModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPostWithHttpInfo(providerId, providerCreatePlanModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param providerId Provider ID (UUID format)
     * @param providerCreatePlanModel Plan creation data
     */
    public providersProviderIdPlansPost(providerId: string, providerCreatePlanModel: ProviderCreatePlanModel, _options?: PromiseConfigurationOptions): Promise<ProviderPlanModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPlansPost(providerId, providerCreatePlanModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param providerId Provider ID (UUID format)
     * @param providerUpdateProviderModel Updated provider data
     */
    public providersProviderIdPutWithHttpInfo(providerId: string, providerUpdateProviderModel: ProviderUpdateProviderModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ProviderProviderModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPutWithHttpInfo(providerId, providerUpdateProviderModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param providerId Provider ID (UUID format)
     * @param providerUpdateProviderModel Updated provider data
     */
    public providersProviderIdPut(providerId: string, providerUpdateProviderModel: ProviderUpdateProviderModel, _options?: PromiseConfigurationOptions): Promise<ProviderProviderModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.providersProviderIdPut(providerId, providerUpdateProviderModel, observableOptions);
        return result.toPromise();
    }


}



import { ObservableSubscriptionsApi } from './ObservableAPI';

import { SubscriptionsApiRequestFactory, SubscriptionsApiResponseProcessor} from "../apis/SubscriptionsApi";
export class PromiseSubscriptionsApi {
    private api: ObservableSubscriptionsApi

    public constructor(
        configuration: Configuration,
        requestFactory?: SubscriptionsApiRequestFactory,
        responseProcessor?: SubscriptionsApiResponseProcessor
    ) {
        this.api = new ObservableSubscriptionsApi(configuration, requestFactory, responseProcessor);
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
    public subscriptionsGetWithHttpInfo(search?: string, recurrencies?: Array<string>, fromDate?: string, toDate?: string, users?: Array<string>, withInactive?: boolean, providers?: Array<string>, limit?: number, offset?: number, _options?: PromiseConfigurationOptions): Promise<HttpInfo<DtoPaginatedResponseModelSubscriptionSubscriptionModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsGetWithHttpInfo(search, recurrencies, fromDate, toDate, users, withInactive, providers, limit, offset, observableOptions);
        return result.toPromise();
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
    public subscriptionsGet(search?: string, recurrencies?: Array<string>, fromDate?: string, toDate?: string, users?: Array<string>, withInactive?: boolean, providers?: Array<string>, limit?: number, offset?: number, _options?: PromiseConfigurationOptions): Promise<DtoPaginatedResponseModelSubscriptionSubscriptionModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsGet(search, recurrencies, fromDate, toDate, users, withInactive, providers, limit, offset, observableOptions);
        return result.toPromise();
    }

    /**
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param subscriptionPatchSubscriptionModel Complete subscription data
     */
    public subscriptionsPatchWithHttpInfo(subscriptionPatchSubscriptionModel: SubscriptionPatchSubscriptionModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsPatchWithHttpInfo(subscriptionPatchSubscriptionModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param subscriptionPatchSubscriptionModel Complete subscription data
     */
    public subscriptionsPatch(subscriptionPatchSubscriptionModel: SubscriptionPatchSubscriptionModel, _options?: PromiseConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsPatch(subscriptionPatchSubscriptionModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param subscriptionCreateSubscriptionModel Subscription creation data
     */
    public subscriptionsPostWithHttpInfo(subscriptionCreateSubscriptionModel: SubscriptionCreateSubscriptionModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsPostWithHttpInfo(subscriptionCreateSubscriptionModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param subscriptionCreateSubscriptionModel Subscription creation data
     */
    public subscriptionsPost(subscriptionCreateSubscriptionModel: SubscriptionCreateSubscriptionModel, _options?: PromiseConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsPost(subscriptionCreateSubscriptionModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdDeleteWithHttpInfo(subscriptionId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSubscriptionIdDeleteWithHttpInfo(subscriptionId, observableOptions);
        return result.toPromise();
    }

    /**
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdDelete(subscriptionId: string, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSubscriptionIdDelete(subscriptionId, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdGetWithHttpInfo(subscriptionId: string, _options?: PromiseConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSubscriptionIdGetWithHttpInfo(subscriptionId, observableOptions);
        return result.toPromise();
    }

    /**
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     */
    public subscriptionsSubscriptionIdGet(subscriptionId: string, _options?: PromiseConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSubscriptionIdGet(subscriptionId, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     * @param subscriptionUpdateSubscriptionModel Updated subscription data
     */
    public subscriptionsSubscriptionIdPutWithHttpInfo(subscriptionId: string, subscriptionUpdateSubscriptionModel: SubscriptionUpdateSubscriptionModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSubscriptionIdPutWithHttpInfo(subscriptionId, subscriptionUpdateSubscriptionModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param subscriptionId Subscription ID (UUID format)
     * @param subscriptionUpdateSubscriptionModel Updated subscription data
     */
    public subscriptionsSubscriptionIdPut(subscriptionId: string, subscriptionUpdateSubscriptionModel: SubscriptionUpdateSubscriptionModel, _options?: PromiseConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSubscriptionIdPut(subscriptionId, subscriptionUpdateSubscriptionModel, observableOptions);
        return result.toPromise();
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
    public subscriptionsSummaryGetWithHttpInfo(topProviders: number, topLabels: number, upcomingRenewals: number, totalMonthly: boolean, totalYearly: boolean, _options?: PromiseConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionSummaryResponse>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSummaryGetWithHttpInfo(topProviders, topLabels, upcomingRenewals, totalMonthly, totalYearly, observableOptions);
        return result.toPromise();
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
    public subscriptionsSummaryGet(topProviders: number, topLabels: number, upcomingRenewals: number, totalMonthly: boolean, totalYearly: boolean, _options?: PromiseConfigurationOptions): Promise<SubscriptionSubscriptionSummaryResponse> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.subscriptionsSummaryGet(topProviders, topLabels, upcomingRenewals, totalMonthly, totalYearly, observableOptions);
        return result.toPromise();
    }


}



import { ObservableUsersApi } from './ObservableAPI';

import { UsersApiRequestFactory, UsersApiResponseProcessor} from "../apis/UsersApi";
export class PromiseUsersApi {
    private api: ObservableUsersApi

    public constructor(
        configuration: Configuration,
        requestFactory?: UsersApiRequestFactory,
        responseProcessor?: UsersApiResponseProcessor
    ) {
        this.api = new ObservableUsersApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Deletes the authenticated user\'s account
     * Delete user
     */
    public usersDeleteWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.usersDeleteWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Deletes the authenticated user\'s account
     * Delete user
     */
    public usersDelete(_options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.usersDelete(observableOptions);
        return result.toPromise();
    }

    /**
     * Returns the preferred currency for the authenticated user
     * Get user preferred currency
     */
    public usersPreferredCurrencyGetWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<UserUserPreferredCurrencyModel>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.usersPreferredCurrencyGetWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     * Returns the preferred currency for the authenticated user
     * Get user preferred currency
     */
    public usersPreferredCurrencyGet(_options?: PromiseConfigurationOptions): Promise<UserUserPreferredCurrencyModel> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.usersPreferredCurrencyGet(observableOptions);
        return result.toPromise();
    }

    /**
     * Updates the preferred currency for the authenticated user
     * Update user preferred currency
     * @param authorization Bearer token
     * @param userUpdatePreferredCurrencyModel Profile update parameters
     */
    public usersPreferredCurrencyPutWithHttpInfo(authorization: string, userUpdatePreferredCurrencyModel: UserUpdatePreferredCurrencyModel, _options?: PromiseConfigurationOptions): Promise<HttpInfo<void>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.usersPreferredCurrencyPutWithHttpInfo(authorization, userUpdatePreferredCurrencyModel, observableOptions);
        return result.toPromise();
    }

    /**
     * Updates the preferred currency for the authenticated user
     * Update user preferred currency
     * @param authorization Bearer token
     * @param userUpdatePreferredCurrencyModel Profile update parameters
     */
    public usersPreferredCurrencyPut(authorization: string, userUpdatePreferredCurrencyModel: UserUpdatePreferredCurrencyModel, _options?: PromiseConfigurationOptions): Promise<void> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.usersPreferredCurrencyPut(authorization, userUpdatePreferredCurrencyModel, observableOptions);
        return result.toPromise();
    }


}



