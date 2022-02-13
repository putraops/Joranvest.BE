package view_models

import (
	"time"
)

type BaseViewModel struct {
	Id                    string     `json:"id"`
	IsActive              bool       `json:"is_active"`
	IsLocked              bool       `json:"is_locked"`
	IsDefault             bool       `json:"is_default"`
	CreatedAt             *time.Time `json:"created_at"`
	CreatedBy             string     `json:"created_by"`
	CreatedUserFullname   string     `json:"created_user_fullname"`
	UpdatedAt             *time.Time `json:"updated_at"`
	UpdatedBy             *string    `json:"updated_by"`
	UpdatedUserFullname   *string    `json:"updated_user_fullname"`
	SubmittedAt           *time.Time `json:"submitted_at"`
	SubmittedBy           *string    `json:"submitted_by"`
	SubmittedUserFullname *string    `json:"submitted_user_fullname"`
	ApprovedAt            *time.Time `json:"approved_at"`
	ApprovedBy            *string    `json:"approved_by"`
	ApprovedUserFullname  *string    `json:"approved_user_fullname"`
}
