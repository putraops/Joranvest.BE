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
	Q     string
	Page  int
	Size  int
	Field []string
}

type ReactSelectResponse struct {
	Results []ReactSelectItem `json:"results"`
	Count   int               `json:"count"`
}

type ReactSelectGroupResponse struct {
	Results []ReactSelectItemGroup `json:"results"`
	Count   int                    `json:"count"`
}
