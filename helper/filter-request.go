package helper

type ProductFilter struct {
	CategoryId  string
	ProductType string
}

type Operator struct {
	Operator string
	Field    []string
	Value    interface{}
}

type DataFilter struct {
	Request   []Operator
	TableName string
}
