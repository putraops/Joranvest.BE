package dto

//ArticleDto is a model that client use when updating
type ArticleDto struct {
	Id                string `json:"id" form:"id"`
	Title             string `json:"title" form:"title" binding:"required"`
	Body              string `json:"body" form:"body"`
	ArticleType       string `json:"article_type" form:"article_type"`
	ArticleCategoryId string `json:"article_category_id" form:"article_category_id"`
	Source            string `json:"source" form:"source"`
	Tag               string `json:"tag" form:"tag"`
	EntityId          string `json:"-"`
	UpdatedBy         string
}
