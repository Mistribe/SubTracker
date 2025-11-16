package export

// LabelExportModel represents a label for export purposes
type LabelExportModel struct {
	Id            string  `json:"id" csv:"id" yaml:"id"`
	Name          string  `json:"name" csv:"name" yaml:"name"`
	Color         string  `json:"color" csv:"color" yaml:"color"`
	OwnerType     string  `json:"ownerType" csv:"ownerType" yaml:"ownerType"`
	OwnerFamilyId *string `json:"ownerFamilyId,omitempty" csv:"ownerFamilyId" yaml:"ownerFamilyId,omitempty"`
}

// ProviderExportModel represents a provider for export purposes
type ProviderExportModel struct {
	Id             string   `json:"id" csv:"id" yaml:"id"`
	Name           string   `json:"name" csv:"name" yaml:"name"`
	Description    *string  `json:"description,omitempty" csv:"description" yaml:"description,omitempty"`
	Url            *string  `json:"url,omitempty" csv:"url" yaml:"url,omitempty"`
	IconUrl        *string  `json:"iconUrl,omitempty" csv:"iconUrl" yaml:"iconUrl,omitempty"`
	PricingPageUrl *string  `json:"pricingPageUrl,omitempty" csv:"pricingPageUrl" yaml:"pricingPageUrl,omitempty"`
	Labels         []string `json:"labels" csv:"labels" yaml:"labels"`
}

// SubscriptionExportModel represents a subscription for export purposes
type SubscriptionExportModel struct {
	Id                 string   `json:"id" csv:"id" yaml:"id"`
	ProviderId         string   `json:"providerId" csv:"providerId" yaml:"providerId"`
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
	Labels             []string `json:"labels" csv:"labels" yaml:"labels"`
}
