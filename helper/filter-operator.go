package helper

import "strings"

type FilterOperator interface {
	GetOperator(pwd []byte) string
}

func GetOperator(text string) string {
	switch strings.ToLower(text) {
	case "eq":
		return "="
	case "neq":
		return "<>"
	case "lt":
		return "<"
	case "lte":
		return "<="
	case "gt":
		return ">"
	case "gte":
		return ">="
	default:
		return "="
	}
}
