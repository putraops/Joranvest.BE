package models

import (
	"database/sql"
)

type TechnicalAnalysis struct {
	Id         string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive   bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked   bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault  bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt  sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy  string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt  sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy  string       `gorm:"type:varchar(50)" json:"updated_by"`
	ApprovedAt sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	OwnerName  string       `gorm:"type:varchar(50)" json:"owner_name"`
	ApprovedBy string       `gorm:"type:varchar(50)" json:"approved_by"`
	EntityId   string       `gorm:"type:varchar(50);null" json:"entity_id"`

	EmitenId           string  `gorm:"type:varchar(50);not null" json:"emiten_id"`
	Signal             string  `gorm:"type:varchar(50)" json:"signal"`
	BandarmologyStatus string  `gorm:"type:varchar(50)" json:"bandarmology_status"`
	StartBuy           float64 `gorm:"type:decimal(18,2)" json:"start_buy"`
	EndBuy             float64 `gorm:"type:decimal(18,2)" json:"end_buy"`
	StartSell          float64 `gorm:"type:decimal(18,2)" json:"start_sell"`
	EndSell            float64 `gorm:"type:decimal(18,2)" json:"end_sell"`
	StartCut           float64 `gorm:"type:decimal(18,2)" json:"start_cut"`
	EndCut             float64 `gorm:"type:decimal(18,2)" json:"end_cut"`
	StartRatio         float64 `gorm:"type:decimal(18,2)" json:"start_ratio"`
	EndRatio           float64 `gorm:"type:decimal(18,2)" json:"end_ratio"`
	Timeframe          string  `gorm:"type:varchar(50)" json:"timeframe"`
	ReasonToBuy        string  `gorm:"type:text" json:"reason_to_buy"`
	Status             string  `gorm:"type:varchar(50)" json:"status"`
	Description        string  `gorm:"type:text" json:"description"`

	Emiten Emiten `gorm:"foreignkey:EmitenId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"emiten"`
}

func (TechnicalAnalysis) TableName() string {
	return "technical_analysis"
}
