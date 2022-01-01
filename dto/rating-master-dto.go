package dto

//RatingMasterDto is a model that client use when updating a book
type RatingMasterDto struct {
	Id            string `json:"id" form:"id"`
	UserId        string `json:"user_id" form:"user_id" binding:"required"`
	ObjectRatedId string `json:"object_rated_id" form:"object_rated_id" binding:"required"`
	ReferenceId   string `json:"reference_id" form:"reference_id"`
	Rating        int    `json:"rating" form:"rating" binding:"required"`
	Comment       string `json:"comment" form:"comment"`
	EntityId      string `json:"-"`
	UpdatedBy     string
}
