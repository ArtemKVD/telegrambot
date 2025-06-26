package calc

import (
	"log"
	"strconv"
)

func Kforset(gender, heights, weights string) int {
	height, err := strconv.Atoi(heights)
	if err != nil {
		log.Print("convert error")
	}

	weight, err := strconv.Atoi(weights)
	if err != nil {
		log.Print("convert error")
	}
	if gender == "м" {
		return int(90 + 13.4*float64(weight) + 4.8*float64(height))
	} else {
		return int(450 + 9.2*float64(weight) + 3.1*float64(height))
	}
}
func Kforlost(gender, heights, weights string) int {
	height, err := strconv.Atoi(heights)
	if err != nil {
		log.Print("convert error")
	}

	weight, err := strconv.Atoi(weights)
	if err != nil {
		log.Print("convert error")
	}
	if gender == "м" {
		return int(90 + 13.4*float64(weight) + 4.8*float64(height) - 2.5*float64(weight))
	} else {
		return int(450 + 9.2*float64(weight) + 3.1*float64(height) - 2*float64(weight))
	}
}
func Kforget(gender, heights, weights string) int {
	height, err := strconv.Atoi(heights)
	if err != nil {
		log.Print("convert error")
	}

	weight, err := strconv.Atoi(weights)
	if err != nil {
		log.Print("convert error")
	}
	if gender == "м" {
		return int(90 + 13.4*float64(weight) + 4.8*float64(height) + 2.5*float64(weight))
	} else {
		return int(450 + 9.2*float64(weight) + 3.1*float64(height) + 2*float64(weight))
	}
}
