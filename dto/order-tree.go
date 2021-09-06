package dto

type OrderTreeDto struct {
	RecordId   string `json:"record_id" form:"record_id"`
	ParentId   string `json:"parent_id" form:"parent_id"`
	OrderIndex int    `json:"order_index" form:"order_index"`
}
