package models

import (
	"database/sql"
)

type WebinarRecording struct {
	Id          string        `gorm:"type:varchar(50);not null;primary_key" json:"id"`
	IsActive    *bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    *bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   *bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   *sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   string        `gorm:"type:varchar(50);not null" json:"created_by"`
	UpdatedAt   *sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   string        `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt *sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy string        `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt  *sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  *string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId    *string       `gorm:"type:varchar(50);null" json:"entity_id"`
	WebinarId   string        `gorm:"type:varchar(50);not null" json:"webinar_id"`
	VideoUrl    *string       `gorm:"type:text" json:"video_url"`
	PathUrl     string        `gorm:"type:text;not null;uniqueIndex:idx_webinar_recording_path_url" json:"path_url"`
	Webinar     *Webinar      `gorm:"foreignkey:WebinarId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
}

func (WebinarRecording) TableName() string {
	return "webinar_recording"
}
