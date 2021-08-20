package dto

//TechnicalAnalysisDto is a model that client use when updating a book
type TechnicalAnalysisDto struct {
	Id                 string  `json:"id" form:"id"`
	EmitenId           string  `json:"emiten_id" form:"emiten_id" binding:"required"`
	Signal             string  `json:"signal" form:"signal"`
	BandarmologyStatus string  `json:"bandarmology_status" form:"bandarmology_status"`
	StartBuy           float64 `json:"start_buy" form:"start_buy"`
	EndBuy             float64 `json:"end_buy" form:"end_buy"`
	StartSell          float64 `json:"start_sell" form:"start_sell"`
	EndSell            float64 `json:"end_sell" form:"end_sell"`
	StartCut           float64 `json:"start_cut" form:"start_cut"`
	EndCut             float64 `json:"end_cut" form:"end_cut"`
	StartRatio         float64 `json:"start_ratio" form:"start_ratio"`
	EndRatio           float64 `json:"end_ratio" form:"end_ratio"`
	Timeframe          string  `json:"timeframe" form:"timeframe"`
	ReasonToBuy        string  `json:"reason_to_buy" form:"reason_to_buy"`
	Status             string  `json:"status" form:"status"`
	Description        string  `json:"description" form:"description"`
	EntityId           string  `json:"-"`
	UpdatedBy          string
}
