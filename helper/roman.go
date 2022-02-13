package helper

func NumberMonthToRoman(number int) string {
	if number == 1 {
		return "I"
	} else if number == 2 {
		return "II"
	} else if number == 3 {
		return "III"
	} else if number == 4 {
		return "IV"
	} else if number == 5 {
		return "V"
	} else if number == 6 {
		return "VI"
	} else if number == 7 {
		return "VII"
	} else if number == 8 {
		return "VIII"
	} else if number == 9 {
		return "IX"
	} else if number == 10 {
		return "X"
	} else if number == 11 {
		return "XI"
	}
	return "XII"
}
