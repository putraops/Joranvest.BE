package commons

type PaginationRequest struct {
	Page   int               `json:"page" form:"page"`
	Size   int               `json:"size" form:"size"`
	Search interface{}       `json:"search"`
	Filter map[string]string `json:"filter"`
	Order  map[string]string `json:"order"`
}

type Pagination2ndRequest struct {
	Page   int               `json:"page" form:"page"`
	Size   int               `json:"size" form:"size"`
	Search interface{}       `json:"search"`
	Filter []Filter          `json:"filter"`
	Order  map[string]string `json:"order"`
}

type PaginationResponse struct {
	Total int         `json:"total" form:"total"`
	Data  interface{} `json:"data" form:"data"`
}

type Filter struct {
	Field    string `json:"field" form:"field"`
	Operator string `json:"operator" form:"operator"`
	Value    string `json:"value" form:"value"`
}
