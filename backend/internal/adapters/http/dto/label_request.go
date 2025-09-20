package dto

import (
	"time"
)

type CreateLabelRequest struct {
	Id        *string            `json:"id,omitempty"`
	Name      string             `json:"name" binding:"required"`
	Color     string             `json:"color" binding:"required"`
	Owner     EditableOwnerModel `json:"owner" binding:"required"`
	CreatedAt *time.Time         `json:"created_at,omitempty" format:"date-time"`
}

type UpdateLabelRequest struct {
	Name      string     `json:"name" binding:"required"`
	Color     string     `json:"color" binding:"required"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}
