package dto

// LabelExportModel represents a label for export purposes
type LabelExportModel struct {
	Id            string  `json:"id" csv:"id" yaml:"id"`
	Name          string  `json:"name" csv:"name" yaml:"name"`
	Color         string  `json:"color" csv:"color" yaml:"color"`
	OwnerType     string  `json:"ownerType" csv:"ownerType" yaml:"ownerType"`
	OwnerFamilyId *string `json:"ownerFamilyId,omitempty" csv:"ownerFamilyId" yaml:"ownerFamilyId,omitempty"`
}
