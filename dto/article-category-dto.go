package dto

//ArticleCategoryDto is a model that client use when updating a book
type ArticleCategoryDto struct {
	Id                string `json:"id" form:"id"`
	Title             string `json:"title" form:"title" binding:"required"`
	SubTitle          string `json:"sub_title" form:"sub_title"`
	Body              string `json:"body" form:"body" binding:"required"`
	Source            string `json:"source" form:"source"`
	ArticleType       string `json:"article_type" form:"article_type"`
	ArticleCategoryId string `json:"article_category_id" form:"article_category_id"`
	Description       string `json:"description" form:"description"`
	Tag               string `json:"tag" form:"tag"`
	EntityId          string `json:"-"`
	UpdatedBy         string
}
