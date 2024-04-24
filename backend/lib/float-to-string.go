package lib

import (
	"log"
	"strconv"
)

func FloatToString(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}

func StringToFloat(number string) float64 {
	n, err := strconv.ParseFloat(number, 64)

	if err != nil {
		log.Fatal(err)
	}

	return n
}
