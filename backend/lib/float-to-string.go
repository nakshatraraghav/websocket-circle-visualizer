package lib

import "strconv"

func FloatToString(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}
