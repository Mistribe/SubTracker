package dto

// SubscriptionExportModel represents a subscription for export purposes
type SubscriptionExportModel struct {
	Id                 string   `json:"id" csv:"id" yaml:"id"`
	ProviderKey        string   `json:"providerKey" csv:"providerKey" yaml:"providerKey"`
	FriendlyName       *string  `json:"friendlyName,omitempty" csv:"friendlyName" yaml:"friendlyName,omitempty"`
	StartDate          string   `json:"startDate" csv:"startDate" yaml:"startDate"`
	EndDate            *string  `json:"endDate,omitempty" csv:"endDate" yaml:"endDate,omitempty"`
	Recurrency         string   `json:"recurrency" csv:"recurrency" yaml:"recurrency"`
	CustomRecurrency   *int32   `json:"customRecurrency,omitempty" csv:"customRecurrency" yaml:"customRecurrency,omitempty"`
	Amount             float64  `json:"amount" csv:"amount" yaml:"amount"`
	Currency           string   `json:"currency" csv:"currency" yaml:"currency"`
	OwnerType          string   `json:"ownerType" csv:"ownerType" yaml:"ownerType"`
	FreeTrialStartDate *string  `json:"freeTrialStartDate,omitempty" csv:"freeTrialStartDate" yaml:"freeTrialStartDate,omitempty"`
	FreeTrialEndDate   *string  `json:"freeTrialEndDate,omitempty" csv:"freeTrialEndDate" yaml:"freeTrialEndDate,omitempty"`
	Payer              *string  `json:"payer,omitempty" csv:"payer" yaml:"payer,omitempty" enums:"family,family_member"`
	PayerMemberId      *string  `json:"payerMemberId,omitempty" csv:"payerMemberId" yaml:"payerMemberId,omitempty"`
	FamilyUsers        []string `json:"familyUsers" csv:"familyUsers" yaml:"familyUsers"`
	Labels             []string `json:"labels" csv:"labels" yaml:"labels"`
}
