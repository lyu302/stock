package arithmetic

import (
	"math"
	"sort"
)

var (
	NormalMadK float64 = 1.8
)

func Mad(values []float64) (float64, float64) {
	var (
		mid float64
		madMid float64
	)
	sort.Float64s(values)

	count := len(values)
	if count % 2 == 0 {
		mid = (values[count/2 - 1] + values[count/2]) / 2
	} else {
		mid = values[count/2]
	}

	midValues := make([]float64, 0)
	for _, value := range values {
		midValue := math.Abs(value - mid)
		midValues = append(midValues, midValue)
	}
	sort.Float64s(midValues)
	if count % 2 == 0 {
		madMid = (midValues[count/2 - 1] + midValues[count/2]) / 2
	} else {
		madMid = midValues[count/2]
	}

	return mid, madMid
}

func MadWithDuration(values []float64, duration int) [][]int {
	indexes := make([][]int, 0)

	sortValues := make([]float64, 0)
	sortValues = append(sortValues, values...)

	mid, madMid := Mad(sortValues)
	//upValue := mid + NormalMadK * madMid
	lowValue := mid - NormalMadK * madMid

	curIndexes := make([]int, 0)

	for index, value := range values {
		//if value > upValue || value < lowValue {
		if value < lowValue {
			curIndexes = append(curIndexes, index)
			if index == len(values) - 1 && len(curIndexes) >= duration {
				indexes = append(indexes, curIndexes)
			}
		} else {
			if len(curIndexes) >= duration {
				indexes = append(indexes, curIndexes)
			}
			curIndexes = nil
		}
	}

	return indexes
}