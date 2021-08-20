package helper

import (
	"strings"
)

func StringifyToArray(req string) []string {
	field := req
	field = strings.Replace(field, `"`, ``, -1)
	field = strings.Replace(field, `[`, ``, -1)
	field = strings.Replace(field, `]`, ``, -1)
	return strings.Split(field, ",")
}
