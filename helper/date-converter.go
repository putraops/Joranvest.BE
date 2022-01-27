package helper

import (
	"strings"
)

func ConvertMonthNameENGtoID(text string) string {
	switch strings.ToLower(text) {
	default:
		return "Unknown"
	case "january":
		return "Januari"
	case "february":
		return "Februari"
	case "march":
		return "Maret"
	case "april":
		return "April"
	case "may":
		return "Mei"
	case "june":
		return "Juni"
	case "july":
		return "Juli"
	case "august":
		return "Agustus"
	case "september":
		return "September"
	case "october":
		return "Oktober"
	case "november":
		return "Nopember"
	case "december":
		return "Desember"
	}
}
