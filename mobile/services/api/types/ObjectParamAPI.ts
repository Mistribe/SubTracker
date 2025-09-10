import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration, ConfigurationOptions } from '../configuration'
import type { Middleware } from '../middleware';

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

import { ObservableCurrenciesApi } from "./ObservableAPI";
import { CurrenciesApiRequestFactory, CurrenciesApiResponseProcessor} from "../apis/CurrenciesApi";

export interface CurrenciesApiCurrenciesRatesGetRequest {
    /**
     * Conversion date in RFC3339 format (default: current time)
     * Defaults to: undefined
     * @type string
     * @memberof CurrenciesApicurrenciesRatesGet
     */
    date?: string
}

export interface CurrenciesApiCurrenciesSupportedGetRequest {
}

export class ObjectCurrenciesApi {
    private api: ObservableCurrenciesApi

    public constructor(configuration: Configuration, requestFactory?: CurrenciesApiRequestFactory, responseProcessor?: CurrenciesApiResponseProcessor) {
        this.api = new ObservableCurrenciesApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Get exchange rates for all currencies at a specific date
     * Get Currency Rates
     * @param param the request object
     */
    public currenciesRatesGetWithHttpInfo(param: CurrenciesApiCurrenciesRatesGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<CurrencyCurrencyRatesModel>> {
        return this.api.currenciesRatesGetWithHttpInfo(param.date,  options).toPromise();
    }

    /**
     * Get exchange rates for all currencies at a specific date
     * Get Currency Rates
     * @param param the request object
     */
    public currenciesRatesGet(param: CurrenciesApiCurrenciesRatesGetRequest = {}, options?: ConfigurationOptions): Promise<CurrencyCurrencyRatesModel> {
        return this.api.currenciesRatesGet(param.date,  options).toPromise();
    }

    /**
     * get details of all supported currencies
     * Get Supported Currencies
     * @param param the request object
     */
    public currenciesSupportedGetWithHttpInfo(param: CurrenciesApiCurrenciesSupportedGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<Array<string>>> {
        return this.api.currenciesSupportedGetWithHttpInfo( options).toPromise();
    }

    /**
     * get details of all supported currencies
     * Get Supported Currencies
     * @param param the request object
     */
    public currenciesSupportedGet(param: CurrenciesApiCurrenciesSupportedGetRequest = {}, options?: ConfigurationOptions): Promise<Array<string>> {
        return this.api.currenciesSupportedGet( options).toPromise();
    }

}

import { ObservableFamilyApi } from "./ObservableAPI";
import { FamilyApiRequestFactory, FamilyApiResponseProcessor} from "../apis/FamilyApi";

export interface FamilyApiFamiliesFamilyIdAcceptPostRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdAcceptPost
     */
    familyId: string
    /**
     * Invitation acceptance details
     * @type FamilyFamilyAcceptInvitationRequest
     * @memberof FamilyApifamiliesFamilyIdAcceptPost
     */
    familyFamilyAcceptInvitationRequest: FamilyFamilyAcceptInvitationRequest
}

export interface FamilyApiFamiliesFamilyIdDeclinePostRequest {
    /**
     * Family ID
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdDeclinePost
     */
    familyId: string
    /**
     * Decline invitation request
     * @type FamilyFamilyDeclineInvitationRequest
     * @memberof FamilyApifamiliesFamilyIdDeclinePost
     */
    familyFamilyDeclineInvitationRequest: FamilyFamilyDeclineInvitationRequest
}

export interface FamilyApiFamiliesFamilyIdDeleteRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdDelete
     */
    familyId: string
}

export interface FamilyApiFamiliesFamilyIdInvitationGetRequest {
    /**
     * Family ID
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdInvitationGet
     */
    familyId: string
    /**
     * Invitation code
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdInvitationGet
     */
    code: string
    /**
     * Family member ID
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdInvitationGet
     */
    familyMemberId: string
}

export interface FamilyApiFamiliesFamilyIdInvitePostRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdInvitePost
     */
    familyId: string
    /**
     * Invitation details including email, name, member ID and type (adult/kid)
     * @type FamilyFamilyInviteRequest
     * @memberof FamilyApifamiliesFamilyIdInvitePost
     */
    familyFamilyInviteRequest: FamilyFamilyInviteRequest
}

export interface FamilyApiFamiliesFamilyIdMembersFamilyMemberIdRevokePostRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersFamilyMemberIdRevokePost
     */
    familyId: string
    /**
     * Family Member ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersFamilyMemberIdRevokePost
     */
    familyMemberId: string
    /**
     * 
     * @type any
     * @memberof FamilyApifamiliesFamilyIdMembersFamilyMemberIdRevokePost
     */
    body?: any
}

export interface FamilyApiFamiliesFamilyIdMembersIdDeleteRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersIdDelete
     */
    familyId: string
    /**
     * Family member ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersIdDelete
     */
    id: string
}

export interface FamilyApiFamiliesFamilyIdMembersIdPutRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersIdPut
     */
    familyId: string
    /**
     * Family member ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersIdPut
     */
    id: string
    /**
     * Updated family member data
     * @type FamilyUpdateFamilyMemberModel
     * @memberof FamilyApifamiliesFamilyIdMembersIdPut
     */
    familyUpdateFamilyMemberModel: FamilyUpdateFamilyMemberModel
}

export interface FamilyApiFamiliesFamilyIdMembersPostRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdMembersPost
     */
    familyId: string
    /**
     * Family member creation data
     * @type FamilyCreateFamilyMemberModel
     * @memberof FamilyApifamiliesFamilyIdMembersPost
     */
    familyCreateFamilyMemberModel: FamilyCreateFamilyMemberModel
}

export interface FamilyApiFamiliesFamilyIdPutRequest {
    /**
     * Family ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof FamilyApifamiliesFamilyIdPut
     */
    familyId: string
    /**
     * Updated family data
     * @type FamilyUpdateFamilyModel
     * @memberof FamilyApifamiliesFamilyIdPut
     */
    familyUpdateFamilyModel: FamilyUpdateFamilyModel
}

export interface FamilyApiFamiliesMeGetRequest {
}

export interface FamilyApiFamiliesPatchRequest {
    /**
     * Family update data with members
     * @type FamilyPatchFamilyModel
     * @memberof FamilyApifamiliesPatch
     */
    familyPatchFamilyModel: FamilyPatchFamilyModel
}

export interface FamilyApiFamiliesPostRequest {
    /**
     * Family creation data
     * @type FamilyCreateFamilyModel
     * @memberof FamilyApifamiliesPost
     */
    familyCreateFamilyModel: FamilyCreateFamilyModel
}

export class ObjectFamilyApi {
    private api: ObservableFamilyApi

    public constructor(configuration: Configuration, requestFactory?: FamilyApiRequestFactory, responseProcessor?: FamilyApiResponseProcessor) {
        this.api = new ObservableFamilyApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param param the request object
     */
    public familiesFamilyIdAcceptPostWithHttpInfo(param: FamilyApiFamiliesFamilyIdAcceptPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.familiesFamilyIdAcceptPostWithHttpInfo(param.familyId, param.familyFamilyAcceptInvitationRequest,  options).toPromise();
    }

    /**
     * Accepts an invitation to join a family using the provided invitation code
     * Accept a family invitation
     * @param param the request object
     */
    public familiesFamilyIdAcceptPost(param: FamilyApiFamiliesFamilyIdAcceptPostRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.familiesFamilyIdAcceptPost(param.familyId, param.familyFamilyAcceptInvitationRequest,  options).toPromise();
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param param the request object
     */
    public familiesFamilyIdDeclinePostWithHttpInfo(param: FamilyApiFamiliesFamilyIdDeclinePostRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.familiesFamilyIdDeclinePostWithHttpInfo(param.familyId, param.familyFamilyDeclineInvitationRequest,  options).toPromise();
    }

    /**
     * Endpoint to decline an invitation to join a family
     * Decline family invitation
     * @param param the request object
     */
    public familiesFamilyIdDeclinePost(param: FamilyApiFamiliesFamilyIdDeclinePostRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.familiesFamilyIdDeclinePost(param.familyId, param.familyFamilyDeclineInvitationRequest,  options).toPromise();
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param param the request object
     */
    public familiesFamilyIdDeleteWithHttpInfo(param: FamilyApiFamiliesFamilyIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.familiesFamilyIdDeleteWithHttpInfo(param.familyId,  options).toPromise();
    }

    /**
     * Permanently delete a family and all its members
     * Delete family by ID
     * @param param the request object
     */
    public familiesFamilyIdDelete(param: FamilyApiFamiliesFamilyIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.familiesFamilyIdDelete(param.familyId,  options).toPromise();
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param param the request object
     */
    public familiesFamilyIdInvitationGetWithHttpInfo(param: FamilyApiFamiliesFamilyIdInvitationGetRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilySeeInvitationResponse>> {
        return this.api.familiesFamilyIdInvitationGetWithHttpInfo(param.familyId, param.code, param.familyMemberId,  options).toPromise();
    }

    /**
     * Get information about a family invitation using invitation code
     * View family invitation details
     * @param param the request object
     */
    public familiesFamilyIdInvitationGet(param: FamilyApiFamiliesFamilyIdInvitationGetRequest, options?: ConfigurationOptions): Promise<FamilyFamilySeeInvitationResponse> {
        return this.api.familiesFamilyIdInvitationGet(param.familyId, param.code, param.familyMemberId,  options).toPromise();
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param param the request object
     */
    public familiesFamilyIdInvitePostWithHttpInfo(param: FamilyApiFamiliesFamilyIdInvitePostRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyInviteResponse>> {
        return this.api.familiesFamilyIdInvitePostWithHttpInfo(param.familyId, param.familyFamilyInviteRequest,  options).toPromise();
    }

    /**
     * Creates an invitation for a new member to join the family
     * Invite a new member to the family
     * @param param the request object
     */
    public familiesFamilyIdInvitePost(param: FamilyApiFamiliesFamilyIdInvitePostRequest, options?: ConfigurationOptions): Promise<FamilyFamilyInviteResponse> {
        return this.api.familiesFamilyIdInvitePost(param.familyId, param.familyFamilyInviteRequest,  options).toPromise();
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param param the request object
     */
    public familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(param: FamilyApiFamiliesFamilyIdMembersFamilyMemberIdRevokePostRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.familiesFamilyIdMembersFamilyMemberIdRevokePostWithHttpInfo(param.familyId, param.familyMemberId, param.body,  options).toPromise();
    }

    /**
     * Revokes a member from the family
     * Revoke family member
     * @param param the request object
     */
    public familiesFamilyIdMembersFamilyMemberIdRevokePost(param: FamilyApiFamiliesFamilyIdMembersFamilyMemberIdRevokePostRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.familiesFamilyIdMembersFamilyMemberIdRevokePost(param.familyId, param.familyMemberId, param.body,  options).toPromise();
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param param the request object
     */
    public familiesFamilyIdMembersIdDeleteWithHttpInfo(param: FamilyApiFamiliesFamilyIdMembersIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.familiesFamilyIdMembersIdDeleteWithHttpInfo(param.familyId, param.id,  options).toPromise();
    }

    /**
     * Permanently delete a family member from a family
     * Delete family member by ID
     * @param param the request object
     */
    public familiesFamilyIdMembersIdDelete(param: FamilyApiFamiliesFamilyIdMembersIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.familiesFamilyIdMembersIdDelete(param.familyId, param.id,  options).toPromise();
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param param the request object
     */
    public familiesFamilyIdMembersIdPutWithHttpInfo(param: FamilyApiFamiliesFamilyIdMembersIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        return this.api.familiesFamilyIdMembersIdPutWithHttpInfo(param.familyId, param.id, param.familyUpdateFamilyMemberModel,  options).toPromise();
    }

    /**
     * Update an existing family member\'s information such as name and kid status
     * Update family member by ID
     * @param param the request object
     */
    public familiesFamilyIdMembersIdPut(param: FamilyApiFamiliesFamilyIdMembersIdPutRequest, options?: ConfigurationOptions): Promise<FamilyFamilyModel> {
        return this.api.familiesFamilyIdMembersIdPut(param.familyId, param.id, param.familyUpdateFamilyMemberModel,  options).toPromise();
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param param the request object
     */
    public familiesFamilyIdMembersPostWithHttpInfo(param: FamilyApiFamiliesFamilyIdMembersPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        return this.api.familiesFamilyIdMembersPostWithHttpInfo(param.familyId, param.familyCreateFamilyMemberModel,  options).toPromise();
    }

    /**
     * Add a new member to an existing family
     * Add a new family member
     * @param param the request object
     */
    public familiesFamilyIdMembersPost(param: FamilyApiFamiliesFamilyIdMembersPostRequest, options?: ConfigurationOptions): Promise<FamilyFamilyModel> {
        return this.api.familiesFamilyIdMembersPost(param.familyId, param.familyCreateFamilyMemberModel,  options).toPromise();
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param param the request object
     */
    public familiesFamilyIdPutWithHttpInfo(param: FamilyApiFamiliesFamilyIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        return this.api.familiesFamilyIdPutWithHttpInfo(param.familyId, param.familyUpdateFamilyModel,  options).toPromise();
    }

    /**
     * Update family information such as name and other details
     * Update a family
     * @param param the request object
     */
    public familiesFamilyIdPut(param: FamilyApiFamiliesFamilyIdPutRequest, options?: ConfigurationOptions): Promise<FamilyFamilyModel> {
        return this.api.familiesFamilyIdPut(param.familyId, param.familyUpdateFamilyModel,  options).toPromise();
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     * @param param the request object
     */
    public familiesMeGetWithHttpInfo(param: FamilyApiFamiliesMeGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        return this.api.familiesMeGetWithHttpInfo( options).toPromise();
    }

    /**
     * Retrieve the user\'s family
     * Get user\'s family
     * @param param the request object
     */
    public familiesMeGet(param: FamilyApiFamiliesMeGetRequest = {}, options?: ConfigurationOptions): Promise<FamilyFamilyModel> {
        return this.api.familiesMeGet( options).toPromise();
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param param the request object
     */
    public familiesPatchWithHttpInfo(param: FamilyApiFamiliesPatchRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        return this.api.familiesPatchWithHttpInfo(param.familyPatchFamilyModel,  options).toPromise();
    }

    /**
     * Update or create a family with specified members. If family doesn\'t exist, it will be created.
     * Patch family with members
     * @param param the request object
     */
    public familiesPatch(param: FamilyApiFamiliesPatchRequest, options?: ConfigurationOptions): Promise<FamilyFamilyModel> {
        return this.api.familiesPatch(param.familyPatchFamilyModel,  options).toPromise();
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param param the request object
     */
    public familiesPostWithHttpInfo(param: FamilyApiFamiliesPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<FamilyFamilyModel>> {
        return this.api.familiesPostWithHttpInfo(param.familyCreateFamilyModel,  options).toPromise();
    }

    /**
     * Create a new family with the authenticated user as the owner and initial member
     * Create a new family
     * @param param the request object
     */
    public familiesPost(param: FamilyApiFamiliesPostRequest, options?: ConfigurationOptions): Promise<FamilyFamilyModel> {
        return this.api.familiesPost(param.familyCreateFamilyModel,  options).toPromise();
    }

}

import { ObservableHealthApi } from "./ObservableAPI";
import { HealthApiRequestFactory, HealthApiResponseProcessor} from "../apis/HealthApi";

export interface HealthApiHealthzLiveGetRequest {
}

export class ObjectHealthApi {
    private api: ObservableHealthApi

    public constructor(configuration: Configuration, requestFactory?: HealthApiRequestFactory, responseProcessor?: HealthApiResponseProcessor) {
        this.api = new ObservableHealthApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Returns the health status of the application
     * Health check endpoint
     * @param param the request object
     */
    public healthzLiveGetWithHttpInfo(param: HealthApiHealthzLiveGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<{ [key: string]: string; }>> {
        return this.api.healthzLiveGetWithHttpInfo( options).toPromise();
    }

    /**
     * Returns the health status of the application
     * Health check endpoint
     * @param param the request object
     */
    public healthzLiveGet(param: HealthApiHealthzLiveGetRequest = {}, options?: ConfigurationOptions): Promise<{ [key: string]: string; }> {
        return this.api.healthzLiveGet( options).toPromise();
    }

}

import { ObservableLabelsApi } from "./ObservableAPI";
import { LabelsApiRequestFactory, LabelsApiResponseProcessor} from "../apis/LabelsApi";

export interface LabelsApiLabelsDefaultGetRequest {
}

export interface LabelsApiLabelsGetRequest {
    /**
     * Search text to filter labels by name
     * Defaults to: undefined
     * @type string
     * @memberof LabelsApilabelsGet
     */
    search?: string
    /**
     * Maximum number of items to return (default: 10)
     * Defaults to: undefined
     * @type number
     * @memberof LabelsApilabelsGet
     */
    limit?: number
    /**
     * Number of items to skip for pagination (default: 0)
     * Defaults to: undefined
     * @type number
     * @memberof LabelsApilabelsGet
     */
    offset?: number
}

export interface LabelsApiLabelsIdDeleteRequest {
    /**
     * Label ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof LabelsApilabelsIdDelete
     */
    id: string
}

export interface LabelsApiLabelsIdGetRequest {
    /**
     * Label ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof LabelsApilabelsIdGet
     */
    id: string
}

export interface LabelsApiLabelsIdPutRequest {
    /**
     * Label ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof LabelsApilabelsIdPut
     */
    id: string
    /**
     * Updated label data
     * @type LabelUpdateLabelModel
     * @memberof LabelsApilabelsIdPut
     */
    labelUpdateLabelModel: LabelUpdateLabelModel
}

export interface LabelsApiLabelsPostRequest {
    /**
     * Label creation data
     * @type LabelCreateLabelModel
     * @memberof LabelsApilabelsPost
     */
    labelCreateLabelModel: LabelCreateLabelModel
}

export class ObjectLabelsApi {
    private api: ObservableLabelsApi

    public constructor(configuration: Configuration, requestFactory?: LabelsApiRequestFactory, responseProcessor?: LabelsApiResponseProcessor) {
        this.api = new ObservableLabelsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     * @param param the request object
     */
    public labelsDefaultGetWithHttpInfo(param: LabelsApiLabelsDefaultGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<Array<LabelLabelModel>>> {
        return this.api.labelsDefaultGetWithHttpInfo( options).toPromise();
    }

    /**
     * Retrieves a list of default system labels available to all users
     * Get default labels
     * @param param the request object
     */
    public labelsDefaultGet(param: LabelsApiLabelsDefaultGetRequest = {}, options?: ConfigurationOptions): Promise<Array<LabelLabelModel>> {
        return this.api.labelsDefaultGet( options).toPromise();
    }

    /**
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param param the request object
     */
    public labelsGetWithHttpInfo(param: LabelsApiLabelsGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<DtoPaginatedResponseModelLabelLabelModel>> {
        return this.api.labelsGetWithHttpInfo(param.search, param.limit, param.offset,  options).toPromise();
    }

    /**
     * Retrieve a paginated list of labels with optional filtering by owner type and search text
     * Get all labels
     * @param param the request object
     */
    public labelsGet(param: LabelsApiLabelsGetRequest = {}, options?: ConfigurationOptions): Promise<DtoPaginatedResponseModelLabelLabelModel> {
        return this.api.labelsGet(param.search, param.limit, param.offset,  options).toPromise();
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param param the request object
     */
    public labelsIdDeleteWithHttpInfo(param: LabelsApiLabelsIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.labelsIdDeleteWithHttpInfo(param.id,  options).toPromise();
    }

    /**
     * Permanently delete a label by its unique identifier
     * Delete label by ID
     * @param param the request object
     */
    public labelsIdDelete(param: LabelsApiLabelsIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.labelsIdDelete(param.id,  options).toPromise();
    }

    /**
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param param the request object
     */
    public labelsIdGetWithHttpInfo(param: LabelsApiLabelsIdGetRequest, options?: ConfigurationOptions): Promise<HttpInfo<LabelLabelModel>> {
        return this.api.labelsIdGetWithHttpInfo(param.id,  options).toPromise();
    }

    /**
     * Retrieve a single label by its unique identifier
     * Get label by ID
     * @param param the request object
     */
    public labelsIdGet(param: LabelsApiLabelsIdGetRequest, options?: ConfigurationOptions): Promise<LabelLabelModel> {
        return this.api.labelsIdGet(param.id,  options).toPromise();
    }

    /**
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param param the request object
     */
    public labelsIdPutWithHttpInfo(param: LabelsApiLabelsIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<LabelLabelModel>> {
        return this.api.labelsIdPutWithHttpInfo(param.id, param.labelUpdateLabelModel,  options).toPromise();
    }

    /**
     * Update an existing label\'s name and color by its unique identifier
     * Update label by ID
     * @param param the request object
     */
    public labelsIdPut(param: LabelsApiLabelsIdPutRequest, options?: ConfigurationOptions): Promise<LabelLabelModel> {
        return this.api.labelsIdPut(param.id, param.labelUpdateLabelModel,  options).toPromise();
    }

    /**
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param param the request object
     */
    public labelsPostWithHttpInfo(param: LabelsApiLabelsPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<LabelLabelModel>> {
        return this.api.labelsPostWithHttpInfo(param.labelCreateLabelModel,  options).toPromise();
    }

    /**
     * Create a new label with specified name, color, and owner information
     * Create a new label
     * @param param the request object
     */
    public labelsPost(param: LabelsApiLabelsPostRequest, options?: ConfigurationOptions): Promise<LabelLabelModel> {
        return this.api.labelsPost(param.labelCreateLabelModel,  options).toPromise();
    }

}

import { ObservableProvidersApi } from "./ObservableAPI";
import { ProvidersApiRequestFactory, ProvidersApiResponseProcessor} from "../apis/ProvidersApi";

export interface ProvidersApiProvidersGetRequest {
    /**
     * Search term
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersGet
     */
    search?: string
    /**
     * Offset (default: 0)
     * Defaults to: undefined
     * @type number
     * @memberof ProvidersApiprovidersGet
     */
    offset?: number
    /**
     * Limit per request (default: 10)
     * Defaults to: undefined
     * @type number
     * @memberof ProvidersApiprovidersGet
     */
    limit?: number
}

export interface ProvidersApiProvidersPostRequest {
    /**
     * Provider creation data
     * @type ProviderCreateProviderModel
     * @memberof ProvidersApiprovidersPost
     */
    providerCreateProviderModel: ProviderCreateProviderModel
}

export interface ProvidersApiProvidersProviderIdDeleteRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdDelete
     */
    providerId: string
}

export interface ProvidersApiProvidersProviderIdGetRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdGet
     */
    providerId: string
}

export interface ProvidersApiProvidersProviderIdPlansPlanIdDeleteRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdDelete
     */
    providerId: string
    /**
     * Plan ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdDelete
     */
    planId: string
}

export interface ProvidersApiProvidersProviderIdPlansPlanIdPricesPostRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPost
     */
    providerId: string
    /**
     * Plan ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPost
     */
    planId: string
    /**
     * Price creation data
     * @type ProviderCreatePriceModel
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPost
     */
    providerCreatePriceModel: ProviderCreatePriceModel
}

export interface ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdDeleteRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdDelete
     */
    providerId: string
    /**
     * Plan ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdDelete
     */
    planId: string
    /**
     * Price ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdDelete
     */
    priceId: string
}

export interface ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdPutRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdPut
     */
    providerId: string
    /**
     * Plan ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdPut
     */
    planId: string
    /**
     * Price ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdPut
     */
    priceId: string
    /**
     * Updated price data
     * @type ProviderUpdatePriceModel
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPricesPriceIdPut
     */
    providerUpdatePriceModel: ProviderUpdatePriceModel
}

export interface ProvidersApiProvidersProviderIdPlansPlanIdPutRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPut
     */
    providerId: string
    /**
     * Plan ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPut
     */
    planId: string
    /**
     * Updated plan data
     * @type ProviderUpdatePlanModel
     * @memberof ProvidersApiprovidersProviderIdPlansPlanIdPut
     */
    providerUpdatePlanModel: ProviderUpdatePlanModel
}

export interface ProvidersApiProvidersProviderIdPlansPostRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPlansPost
     */
    providerId: string
    /**
     * Plan creation data
     * @type ProviderCreatePlanModel
     * @memberof ProvidersApiprovidersProviderIdPlansPost
     */
    providerCreatePlanModel: ProviderCreatePlanModel
}

export interface ProvidersApiProvidersProviderIdPutRequest {
    /**
     * Provider ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof ProvidersApiprovidersProviderIdPut
     */
    providerId: string
    /**
     * Updated provider data
     * @type ProviderUpdateProviderModel
     * @memberof ProvidersApiprovidersProviderIdPut
     */
    providerUpdateProviderModel: ProviderUpdateProviderModel
}

export class ObjectProvidersApi {
    private api: ObservableProvidersApi

    public constructor(configuration: Configuration, requestFactory?: ProvidersApiRequestFactory, responseProcessor?: ProvidersApiResponseProcessor) {
        this.api = new ObservableProvidersApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param param the request object
     */
    public providersGetWithHttpInfo(param: ProvidersApiProvidersGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<DtoPaginatedResponseModelProviderProviderModel>> {
        return this.api.providersGetWithHttpInfo(param.search, param.offset, param.limit,  options).toPromise();
    }

    /**
     * Retrieve a paginated list of all providers with their plans and prices
     * Get all providers
     * @param param the request object
     */
    public providersGet(param: ProvidersApiProvidersGetRequest = {}, options?: ConfigurationOptions): Promise<DtoPaginatedResponseModelProviderProviderModel> {
        return this.api.providersGet(param.search, param.offset, param.limit,  options).toPromise();
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param param the request object
     */
    public providersPostWithHttpInfo(param: ProvidersApiProvidersPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderProviderModel>> {
        return this.api.providersPostWithHttpInfo(param.providerCreateProviderModel,  options).toPromise();
    }

    /**
     * Create a new service provider with labels and owner information
     * Create a new provider
     * @param param the request object
     */
    public providersPost(param: ProvidersApiProvidersPostRequest, options?: ConfigurationOptions): Promise<ProviderProviderModel> {
        return this.api.providersPost(param.providerCreateProviderModel,  options).toPromise();
    }

    /**
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param param the request object
     */
    public providersProviderIdDeleteWithHttpInfo(param: ProvidersApiProvidersProviderIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.providersProviderIdDeleteWithHttpInfo(param.providerId,  options).toPromise();
    }

    /**
     * Permanently delete a provider and all its associated plans and prices
     * Delete provider by ID
     * @param param the request object
     */
    public providersProviderIdDelete(param: ProvidersApiProvidersProviderIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.providersProviderIdDelete(param.providerId,  options).toPromise();
    }

    /**
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param param the request object
     */
    public providersProviderIdGetWithHttpInfo(param: ProvidersApiProvidersProviderIdGetRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderProviderModel>> {
        return this.api.providersProviderIdGetWithHttpInfo(param.providerId,  options).toPromise();
    }

    /**
     * Retrieve a single provider with all its plans and prices by ID
     * Get provider by ID
     * @param param the request object
     */
    public providersProviderIdGet(param: ProvidersApiProvidersProviderIdGetRequest, options?: ConfigurationOptions): Promise<ProviderProviderModel> {
        return this.api.providersProviderIdGet(param.providerId,  options).toPromise();
    }

    /**
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdDeleteWithHttpInfo(param: ProvidersApiProvidersProviderIdPlansPlanIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.providersProviderIdPlansPlanIdDeleteWithHttpInfo(param.providerId, param.planId,  options).toPromise();
    }

    /**
     * Permanently delete a provider plan and all its associated prices
     * Delete provider plan by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdDelete(param: ProvidersApiProvidersProviderIdPlansPlanIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.providersProviderIdPlansPlanIdDelete(param.providerId, param.planId,  options).toPromise();
    }

    /**
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPricesPostWithHttpInfo(param: ProvidersApiProvidersProviderIdPlansPlanIdPricesPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderPriceModel>> {
        return this.api.providersProviderIdPlansPlanIdPricesPostWithHttpInfo(param.providerId, param.planId, param.providerCreatePriceModel,  options).toPromise();
    }

    /**
     * Create a new pricing option for a specific provider plan
     * Create a new provider price
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPricesPost(param: ProvidersApiProvidersProviderIdPlansPlanIdPricesPostRequest, options?: ConfigurationOptions): Promise<ProviderPriceModel> {
        return this.api.providersProviderIdPlansPlanIdPricesPost(param.providerId, param.planId, param.providerCreatePriceModel,  options).toPromise();
    }

    /**
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(param: ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.providersProviderIdPlansPlanIdPricesPriceIdDeleteWithHttpInfo(param.providerId, param.planId, param.priceId,  options).toPromise();
    }

    /**
     * Permanently delete a specific price from a provider plan
     * Delete provider price by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPricesPriceIdDelete(param: ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.providersProviderIdPlansPlanIdPricesPriceIdDelete(param.providerId, param.planId, param.priceId,  options).toPromise();
    }

    /**
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(param: ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderPriceModel>> {
        return this.api.providersProviderIdPlansPlanIdPricesPriceIdPutWithHttpInfo(param.providerId, param.planId, param.priceId, param.providerUpdatePriceModel,  options).toPromise();
    }

    /**
     * Update an existing price for a specific provider plan
     * Update provider price by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPricesPriceIdPut(param: ProvidersApiProvidersProviderIdPlansPlanIdPricesPriceIdPutRequest, options?: ConfigurationOptions): Promise<ProviderPriceModel> {
        return this.api.providersProviderIdPlansPlanIdPricesPriceIdPut(param.providerId, param.planId, param.priceId, param.providerUpdatePriceModel,  options).toPromise();
    }

    /**
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPutWithHttpInfo(param: ProvidersApiProvidersProviderIdPlansPlanIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderPlanModel>> {
        return this.api.providersProviderIdPlansPlanIdPutWithHttpInfo(param.providerId, param.planId, param.providerUpdatePlanModel,  options).toPromise();
    }

    /**
     * Update an existing provider plan\'s information
     * Update provider plan by ID
     * @param param the request object
     */
    public providersProviderIdPlansPlanIdPut(param: ProvidersApiProvidersProviderIdPlansPlanIdPutRequest, options?: ConfigurationOptions): Promise<ProviderPlanModel> {
        return this.api.providersProviderIdPlansPlanIdPut(param.providerId, param.planId, param.providerUpdatePlanModel,  options).toPromise();
    }

    /**
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param param the request object
     */
    public providersProviderIdPlansPostWithHttpInfo(param: ProvidersApiProvidersProviderIdPlansPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderPlanModel>> {
        return this.api.providersProviderIdPlansPostWithHttpInfo(param.providerId, param.providerCreatePlanModel,  options).toPromise();
    }

    /**
     * Create a new subscription plan for an existing provider
     * Create a new provider plan
     * @param param the request object
     */
    public providersProviderIdPlansPost(param: ProvidersApiProvidersProviderIdPlansPostRequest, options?: ConfigurationOptions): Promise<ProviderPlanModel> {
        return this.api.providersProviderIdPlansPost(param.providerId, param.providerCreatePlanModel,  options).toPromise();
    }

    /**
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param param the request object
     */
    public providersProviderIdPutWithHttpInfo(param: ProvidersApiProvidersProviderIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<ProviderProviderModel>> {
        return this.api.providersProviderIdPutWithHttpInfo(param.providerId, param.providerUpdateProviderModel,  options).toPromise();
    }

    /**
     * Update an existing provider\'s basic information
     * Update provider by ID
     * @param param the request object
     */
    public providersProviderIdPut(param: ProvidersApiProvidersProviderIdPutRequest, options?: ConfigurationOptions): Promise<ProviderProviderModel> {
        return this.api.providersProviderIdPut(param.providerId, param.providerUpdateProviderModel,  options).toPromise();
    }

}

import { ObservableSubscriptionsApi } from "./ObservableAPI";
import { SubscriptionsApiRequestFactory, SubscriptionsApiResponseProcessor} from "../apis/SubscriptionsApi";

export interface SubscriptionsApiSubscriptionsGetRequest {
    /**
     * Search text
     * Defaults to: undefined
     * @type string
     * @memberof SubscriptionsApisubscriptionsGet
     */
    search?: string
    /**
     * Filter by recurrency types
     * Defaults to: undefined
     * @type Array&lt;string&gt;
     * @memberof SubscriptionsApisubscriptionsGet
     */
    recurrencies?: Array<string>
    /**
     * Filter by start date (RFC3339)
     * Defaults to: undefined
     * @type string
     * @memberof SubscriptionsApisubscriptionsGet
     */
    fromDate?: string
    /**
     * Filter by end date (RFC3339)
     * Defaults to: undefined
     * @type string
     * @memberof SubscriptionsApisubscriptionsGet
     */
    toDate?: string
    /**
     * Filter by user IDs
     * Defaults to: undefined
     * @type Array&lt;string&gt;
     * @memberof SubscriptionsApisubscriptionsGet
     */
    users?: Array<string>
    /**
     * Include inactive subscriptions
     * Defaults to: undefined
     * @type boolean
     * @memberof SubscriptionsApisubscriptionsGet
     */
    withInactive?: boolean
    /**
     * Filter by provider IDs
     * Defaults to: undefined
     * @type Array&lt;string&gt;
     * @memberof SubscriptionsApisubscriptionsGet
     */
    providers?: Array<string>
    /**
     * Number of items per page (default: 10)
     * Defaults to: undefined
     * @type number
     * @memberof SubscriptionsApisubscriptionsGet
     */
    limit?: number
    /**
     * Page number (default: 0)
     * Defaults to: undefined
     * @type number
     * @memberof SubscriptionsApisubscriptionsGet
     */
    offset?: number
}

export interface SubscriptionsApiSubscriptionsPatchRequest {
    /**
     * Complete subscription data
     * @type SubscriptionPatchSubscriptionModel
     * @memberof SubscriptionsApisubscriptionsPatch
     */
    subscriptionPatchSubscriptionModel: SubscriptionPatchSubscriptionModel
}

export interface SubscriptionsApiSubscriptionsPostRequest {
    /**
     * Subscription creation data
     * @type SubscriptionCreateSubscriptionModel
     * @memberof SubscriptionsApisubscriptionsPost
     */
    subscriptionCreateSubscriptionModel: SubscriptionCreateSubscriptionModel
}

export interface SubscriptionsApiSubscriptionsSubscriptionIdDeleteRequest {
    /**
     * Subscription ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof SubscriptionsApisubscriptionsSubscriptionIdDelete
     */
    subscriptionId: string
}

export interface SubscriptionsApiSubscriptionsSubscriptionIdGetRequest {
    /**
     * Subscription ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof SubscriptionsApisubscriptionsSubscriptionIdGet
     */
    subscriptionId: string
}

export interface SubscriptionsApiSubscriptionsSubscriptionIdPutRequest {
    /**
     * Subscription ID (UUID format)
     * Defaults to: undefined
     * @type string
     * @memberof SubscriptionsApisubscriptionsSubscriptionIdPut
     */
    subscriptionId: string
    /**
     * Updated subscription data
     * @type SubscriptionUpdateSubscriptionModel
     * @memberof SubscriptionsApisubscriptionsSubscriptionIdPut
     */
    subscriptionUpdateSubscriptionModel: SubscriptionUpdateSubscriptionModel
}

export interface SubscriptionsApiSubscriptionsSummaryGetRequest {
    /**
     * Number of top providers to return
     * Defaults to: undefined
     * @type number
     * @memberof SubscriptionsApisubscriptionsSummaryGet
     */
    topProviders: number
    /**
     * Number of top labels to return
     * Defaults to: undefined
     * @type number
     * @memberof SubscriptionsApisubscriptionsSummaryGet
     */
    topLabels: number
    /**
     * Number of upcoming renewals to return
     * Defaults to: undefined
     * @type number
     * @memberof SubscriptionsApisubscriptionsSummaryGet
     */
    upcomingRenewals: number
    /**
     * Include monthly total costs
     * Defaults to: undefined
     * @type boolean
     * @memberof SubscriptionsApisubscriptionsSummaryGet
     */
    totalMonthly: boolean
    /**
     * Include yearly total costs
     * Defaults to: undefined
     * @type boolean
     * @memberof SubscriptionsApisubscriptionsSummaryGet
     */
    totalYearly: boolean
}

export class ObjectSubscriptionsApi {
    private api: ObservableSubscriptionsApi

    public constructor(configuration: Configuration, requestFactory?: SubscriptionsApiRequestFactory, responseProcessor?: SubscriptionsApiResponseProcessor) {
        this.api = new ObservableSubscriptionsApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Retrieve a paginated list of all subscriptions for the authenticated user
     * Get all subscriptions
     * @param param the request object
     */
    public subscriptionsGetWithHttpInfo(param: SubscriptionsApiSubscriptionsGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<DtoPaginatedResponseModelSubscriptionSubscriptionModel>> {
        return this.api.subscriptionsGetWithHttpInfo(param.search, param.recurrencies, param.fromDate, param.toDate, param.users, param.withInactive, param.providers, param.limit, param.offset,  options).toPromise();
    }

    /**
     * Retrieve a paginated list of all subscriptions for the authenticated user
     * Get all subscriptions
     * @param param the request object
     */
    public subscriptionsGet(param: SubscriptionsApiSubscriptionsGetRequest = {}, options?: ConfigurationOptions): Promise<DtoPaginatedResponseModelSubscriptionSubscriptionModel> {
        return this.api.subscriptionsGet(param.search, param.recurrencies, param.fromDate, param.toDate, param.users, param.withInactive, param.providers, param.limit, param.offset,  options).toPromise();
    }

    /**
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param param the request object
     */
    public subscriptionsPatchWithHttpInfo(param: SubscriptionsApiSubscriptionsPatchRequest, options?: ConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        return this.api.subscriptionsPatchWithHttpInfo(param.subscriptionPatchSubscriptionModel,  options).toPromise();
    }

    /**
     * Update or create a subscription with complete details. If subscription doesn\'t exist, it will be created.
     * Patch subscription
     * @param param the request object
     */
    public subscriptionsPatch(param: SubscriptionsApiSubscriptionsPatchRequest, options?: ConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        return this.api.subscriptionsPatch(param.subscriptionPatchSubscriptionModel,  options).toPromise();
    }

    /**
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param param the request object
     */
    public subscriptionsPostWithHttpInfo(param: SubscriptionsApiSubscriptionsPostRequest, options?: ConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        return this.api.subscriptionsPostWithHttpInfo(param.subscriptionCreateSubscriptionModel,  options).toPromise();
    }

    /**
     * Create a new subscription with provider, plan, pricing, and payment information
     * Create a new subscription
     * @param param the request object
     */
    public subscriptionsPost(param: SubscriptionsApiSubscriptionsPostRequest, options?: ConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        return this.api.subscriptionsPost(param.subscriptionCreateSubscriptionModel,  options).toPromise();
    }

    /**
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param param the request object
     */
    public subscriptionsSubscriptionIdDeleteWithHttpInfo(param: SubscriptionsApiSubscriptionsSubscriptionIdDeleteRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.subscriptionsSubscriptionIdDeleteWithHttpInfo(param.subscriptionId,  options).toPromise();
    }

    /**
     * Permanently delete an existing subscription
     * Delete subscription by ID
     * @param param the request object
     */
    public subscriptionsSubscriptionIdDelete(param: SubscriptionsApiSubscriptionsSubscriptionIdDeleteRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.subscriptionsSubscriptionIdDelete(param.subscriptionId,  options).toPromise();
    }

    /**
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param param the request object
     */
    public subscriptionsSubscriptionIdGetWithHttpInfo(param: SubscriptionsApiSubscriptionsSubscriptionIdGetRequest, options?: ConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        return this.api.subscriptionsSubscriptionIdGetWithHttpInfo(param.subscriptionId,  options).toPromise();
    }

    /**
     * Retrieve a single subscription with all its details including provider, plan, and pricing information
     * Get subscription by ID
     * @param param the request object
     */
    public subscriptionsSubscriptionIdGet(param: SubscriptionsApiSubscriptionsSubscriptionIdGetRequest, options?: ConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        return this.api.subscriptionsSubscriptionIdGet(param.subscriptionId,  options).toPromise();
    }

    /**
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param param the request object
     */
    public subscriptionsSubscriptionIdPutWithHttpInfo(param: SubscriptionsApiSubscriptionsSubscriptionIdPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionModel>> {
        return this.api.subscriptionsSubscriptionIdPutWithHttpInfo(param.subscriptionId, param.subscriptionUpdateSubscriptionModel,  options).toPromise();
    }

    /**
     * Update an existing subscription\'s details including provider, plan, pricing, and payment information
     * Update subscription by ID
     * @param param the request object
     */
    public subscriptionsSubscriptionIdPut(param: SubscriptionsApiSubscriptionsSubscriptionIdPutRequest, options?: ConfigurationOptions): Promise<SubscriptionSubscriptionModel> {
        return this.api.subscriptionsSubscriptionIdPut(param.subscriptionId, param.subscriptionUpdateSubscriptionModel,  options).toPromise();
    }

    /**
     * Returns summary information about subscriptions including total costs and upcoming renewals
     * Get subscription summary
     * @param param the request object
     */
    public subscriptionsSummaryGetWithHttpInfo(param: SubscriptionsApiSubscriptionsSummaryGetRequest, options?: ConfigurationOptions): Promise<HttpInfo<SubscriptionSubscriptionSummaryResponse>> {
        return this.api.subscriptionsSummaryGetWithHttpInfo(param.topProviders, param.topLabels, param.upcomingRenewals, param.totalMonthly, param.totalYearly,  options).toPromise();
    }

    /**
     * Returns summary information about subscriptions including total costs and upcoming renewals
     * Get subscription summary
     * @param param the request object
     */
    public subscriptionsSummaryGet(param: SubscriptionsApiSubscriptionsSummaryGetRequest, options?: ConfigurationOptions): Promise<SubscriptionSubscriptionSummaryResponse> {
        return this.api.subscriptionsSummaryGet(param.topProviders, param.topLabels, param.upcomingRenewals, param.totalMonthly, param.totalYearly,  options).toPromise();
    }

}

import { ObservableUsersApi } from "./ObservableAPI";
import { UsersApiRequestFactory, UsersApiResponseProcessor} from "../apis/UsersApi";

export interface UsersApiUsersDeleteRequest {
}

export interface UsersApiUsersPreferredCurrencyGetRequest {
}

export interface UsersApiUsersPreferredCurrencyPutRequest {
    /**
     * Bearer token
     * Defaults to: undefined
     * @type string
     * @memberof UsersApiusersPreferredCurrencyPut
     */
    authorization: string
    /**
     * Profile update parameters
     * @type UserUpdatePreferredCurrencyModel
     * @memberof UsersApiusersPreferredCurrencyPut
     */
    userUpdatePreferredCurrencyModel: UserUpdatePreferredCurrencyModel
}

export class ObjectUsersApi {
    private api: ObservableUsersApi

    public constructor(configuration: Configuration, requestFactory?: UsersApiRequestFactory, responseProcessor?: UsersApiResponseProcessor) {
        this.api = new ObservableUsersApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * Deletes the authenticated user\'s account
     * Delete user
     * @param param the request object
     */
    public usersDeleteWithHttpInfo(param: UsersApiUsersDeleteRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.usersDeleteWithHttpInfo( options).toPromise();
    }

    /**
     * Deletes the authenticated user\'s account
     * Delete user
     * @param param the request object
     */
    public usersDelete(param: UsersApiUsersDeleteRequest = {}, options?: ConfigurationOptions): Promise<void> {
        return this.api.usersDelete( options).toPromise();
    }

    /**
     * Returns the preferred currency for the authenticated user
     * Get user preferred currency
     * @param param the request object
     */
    public usersPreferredCurrencyGetWithHttpInfo(param: UsersApiUsersPreferredCurrencyGetRequest = {}, options?: ConfigurationOptions): Promise<HttpInfo<UserUserPreferredCurrencyModel>> {
        return this.api.usersPreferredCurrencyGetWithHttpInfo( options).toPromise();
    }

    /**
     * Returns the preferred currency for the authenticated user
     * Get user preferred currency
     * @param param the request object
     */
    public usersPreferredCurrencyGet(param: UsersApiUsersPreferredCurrencyGetRequest = {}, options?: ConfigurationOptions): Promise<UserUserPreferredCurrencyModel> {
        return this.api.usersPreferredCurrencyGet( options).toPromise();
    }

    /**
     * Updates the preferred currency for the authenticated user
     * Update user preferred currency
     * @param param the request object
     */
    public usersPreferredCurrencyPutWithHttpInfo(param: UsersApiUsersPreferredCurrencyPutRequest, options?: ConfigurationOptions): Promise<HttpInfo<void>> {
        return this.api.usersPreferredCurrencyPutWithHttpInfo(param.authorization, param.userUpdatePreferredCurrencyModel,  options).toPromise();
    }

    /**
     * Updates the preferred currency for the authenticated user
     * Update user preferred currency
     * @param param the request object
     */
    public usersPreferredCurrencyPut(param: UsersApiUsersPreferredCurrencyPutRequest, options?: ConfigurationOptions): Promise<void> {
        return this.api.usersPreferredCurrencyPut(param.authorization, param.userUpdatePreferredCurrencyModel,  options).toPromise();
    }

}
