package commons

type DataTableRequest struct {
	Draw                  int                   `json:"draw" form:"draw"`
	Start                 int                   `json:"start" form:"start"`
	Length                int                   `json:"length" form:"length"`
	DataTableColumn       []DataTableColumn     `json:"columns"`
	DataTableOrder        []DataTableOrder      `json:"order"`
	DataTableDefaultOrder DataTableDefaultOrder `json:"default_order"`
	Search                DataTableSearch       `json:"search"`
	Filter                []DataTableFilter     `json:"filter"`
}

type DataTableColumn struct {
	Data       string `json:"data" form:"data"`
	Name       string `json:"name" form:"name"`
	Searchable bool   `json:"searchable" form:"searchable"`
	Orderable  bool   `json:"orderable" form:"orderable"`
}

type DataTableSearch struct {
	Value string `json:"value" form:"value"`
	Regex bool   `json:"regex" form:"regex"`
}

type DataTableOrder struct {
	Column int    `json:"column" form:"column"`
	Dir    string `json:"dir" form:"dir"`
}

type DataTableDefaultOrder struct {
	Column string `json:"column" form:"column"`
	Dir    string `json:"dir" form:"dir"`
}

type DataTableFilter struct {
	Column string `json:"column" form:"column"`
	Value  string `json:"value" form:"value"`
}
