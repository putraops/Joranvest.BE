package models

import (
	"database/sql"
)

type WebinarSpeaker struct {
	Id         string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive   bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy  string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt  sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy string       `gorm:"type:varchar(50)" json:"approved_by"`
	OwnerId    string       `gorm:"type:varchar(50)" json:"owner_id"`
	EntityId   string       `gorm:"type:varchar(50);null" json:"entity_id"`

	WebinarId string `gorm:"type:varchar(50);not null" json:"webinar_id"`
	SpeakerId string `gorm:"type:varchar(50);not null" json:"speaker_id"`

	Webinar Webinar         `gorm:"foreignkey:WebinarId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"webinar"`
	Speaker ApplicationUser `gorm:"foreignkey:SpeakerId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"speaker"`
}

func (WebinarSpeaker) TableName() string {
	return "webinar_speaker"
}
