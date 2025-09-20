package dto

import (
	"time"

	"github.com/mistribe/subtracker/internal/domain/label"
)

type LabelModel struct {
	// @Description Unique identifier for the label (UUID format)
	Id string `json:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	// @Description Display name of the label
	Name string  `json:"name" binding:"required" example:"Entertainment" minLength:"1" maxLength:"100"`
	Key  *string `json:"key,omitempty"`
	// @Description Hexadecimal color code for visual representation of the label
	Color string `json:"color" binding:"required" example:"#FF5733" pattern:"^#[0-9A-Fa-f]{6}$"`
	// @Description Ownership information specifying whether this label belongs to a user or family
	Owner OwnerModel `json:"owner" binding:"required"`
	// @Description ISO 8601 timestamp indicating when the label was originally created
	CreatedAt time.Time `json:"created_at" binding:"required" format:"date-time" example:"2023-01-15T10:30:00Z"`
	// @Description ISO 8601 timestamp indicating when the label was last modified
	UpdatedAt time.Time `json:"updated_at" binding:"required" format:"date-time" example:"2023-01-20T14:45:30Z"`
	// @Description Entity tag used for optimistic concurrency control to prevent conflicting updates
	Etag string `json:"etag" binding:"required" example:"W/\"123456789\""`
}

func NewLabelModel(source label.Label) LabelModel {
	return LabelModel{
		Id:        source.Id().String(),
		Name:      source.Name(),
		Key:       source.Key(),
		Color:     source.Color(),
		Owner:     NewOwnerModel(source.Owner()),
		CreatedAt: source.CreatedAt(),
		UpdatedAt: source.UpdatedAt(),
		Etag:      source.ETag(),
	}
}
