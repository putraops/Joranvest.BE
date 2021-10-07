package dto

//ArticleCategoryDto is a model that client use when updating
type ArticleCategoryDto struct {
	Id                string `json:"id" form:"id"`
	Name              string `json:"name" form:"name" binding:"required"`
	Description       string `json:"description" form:"description"`
	ParentId          string `json:"parent_id" form:"parent_id"`
	Source            string `json:"source" form:"source"`
	ArticleType       string `json:"article_type" form:"article_type"`
	ArticleCategoryId string `json:"article_category_id" form:"article_category_id"`
	Tag               string `json:"tag" form:"tag"`
	EntityId          string `json:"-"`
	UpdatedBy         string
}
