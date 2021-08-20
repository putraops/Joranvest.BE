package commons

type DataTableResponse struct {
	Draw            int         `json:"draw"`
	RecordsTotal    int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
	DataRow         interface{} `json:"data"`
}
