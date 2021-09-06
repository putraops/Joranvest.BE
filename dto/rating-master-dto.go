package dto

//RatingMasterDto is a model that client use when updating a book
type RatingMasterDto struct {
	Id        string `json:"id" form:"id"`
	RecordId  string `json:"record_id" form:"record_id" binding:"required"`
	UserId    string `json:"user_id" form:"user_id" binding:"required"`
	Rating    int    `json:"rating" form:"rating" binding:"required"`
	Comment   string `json:"comment" form:"comment"`
	EntityId  string `json:"-"`
	UpdatedBy string
}
