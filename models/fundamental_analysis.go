package models

import (
	"database/sql"
)

type FundamentalAnalysis struct {
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

	EmitenId       string       `gorm:"type:varchar(50);not null" json:"emiten_id"`
	CurrentPrice   float64      `gorm:"type:decimal(18,2)" json:"current_price"`
	NormalPrice    float64      `gorm:"type:decimal(18,2)" json:"normal_price"`
	MarginOfSafety float64      `gorm:"type:decimal(18,2)" json:"margin_of_safety"`
	ResearchDate   sql.NullTime `gorm:"type:timestamp" json:"research_date"`
	ResearchData   string       `gorm:"type:text" json:"research_data"`
	Status         string       `gorm:"type:varchar(50)" json:"status"`
	Description    string       `gorm:"type:text" json:"description"`

	// OrderDetail []OrderDetail `json:"order_detail"`
	FundamentalAnalysisTag []FundamentalAnalysisTag `gorm:"-" json:"fundamental_analysis_tag"`
	Emiten                 Emiten                   `gorm:"foreignkey:EmitenId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"emiten"`
}

func (FundamentalAnalysis) TableName() string {
	return "fundamental_analysis"
}
