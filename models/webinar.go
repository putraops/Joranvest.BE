package models

import (
	"database/sql"
)

type Webinar struct {
	Id          string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive    bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   string       `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy string       `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  string       `gorm:"type:varchar(50)" json:"approved_by"`
	OwnerId     string       `gorm:"type:varchar(50)" json:"owner_id"`
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`

	WebinarCategoryId string       `gorm:"type:varchar(50);not null" json:"webinar_category_id"`
	Title             string       `gorm:"type:text;not null" json:"title"`
	Description       string       `gorm:"type:text" json:"description"`
	WebinarStartDate  sql.NullTime `gorm:"type:timestamp;default:null" json:"webinar_start_date"`
	WebinarEndDate    sql.NullTime `gorm:"type:timestamp;default:null" json:"webinar_end_date"`
	MinAge            int          `gorm:"type:int" json:"min_age"`
	WebinarLevel      string       `gorm:"type:varchar(50)" json:"webinar_level"`
	Price             float64      `gorm:"type:decimal(18,2)" json:"price"`
	Discount          float64      `gorm:"type:decimal(18,2)" json:"discount"`
	IsCertificate     bool         `gorm:"type:bool;default:0" json:"is_certificate"`
	Reward            int          `gorm:"type:int" json:"reward"`
	Status            int          `gorm:"type:int" json:"status"`
	SpeakerType       int          `gorm:"type:int;" json:"speaker_type"`
	Filepath          string       `gorm:"type:varchar(200)" json:"filepath"`
	FilepathThumbnail string       `gorm:"type:varchar(200)" json:"filepath_thumbnail"`
	Filename          string       `gorm:"type:varchar(200)" json:"filename"`
	Extension         string       `gorm:"type:varchar(10)" json:"extension"`
	Size              string       `gorm:"type:varchar(100)" json:"size"`

	WebinarSpeaker []WebinarSpeaker `gorm:"-" json:"webinar_speaker"`
}

func (Webinar) TableName() string {
	return "webinar"
}
