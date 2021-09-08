package commons

type PaginationRequest struct {
	Page   int         `json:"page" form:"page"`
	Size   int         `json:"size" form:"size"`
	Search interface{} `json:"search"`
	Filter interface{} `json:"filter"`
}

type PaginationResponse struct {
	Total int         `json:"total" form:"total"`
	Data  interface{} `json:"data" form:"data"`
}
