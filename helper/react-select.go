package helper

type ReactSelectItem struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	ParentId string `json:"parent_id"`
}

type ReactSelectItemGroup struct {
	Value    string      `json:"value"`
	Label    string      `json:"label"`
	ParentId string      `json:"parent_id"`
	Children interface{} `json:"children"`
}

type ReactSelectRequest struct {
	Q     string   `json:"q" form:"q"`
	Page  int      `json:"page" form:"page"`
	Size  int      `json:"size" form:"size"`
	Field []string `json:"field" form:"field"`
}

type ReactSelectResponse struct {
	Results []ReactSelectItem `json:"results"`
	Count   int               `json:"count"`
}

type ReactSelectGroupResponse struct {
	Results []ReactSelectItemGroup `json:"results"`
	Count   int                    `json:"count"`
}
