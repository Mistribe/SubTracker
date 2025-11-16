package dto

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
