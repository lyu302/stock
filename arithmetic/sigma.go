package arithmetic

import (
	"math"
)

func Sigma3(values []float64) []int {
	rst := make([]int, 0)

	indexesWithGroup := Sigma3WithDuration(values, 1)
	for _, indexes := range indexesWithGroup {
		rst = append(rst, indexes...)
	}

	return rst
}

func Sigma3WithDuration(values []float64, duration int) [][]int {
	var (
		sigma float64 = 3
		indexes = make([][]int, 0)
	)

	mean := mean(values)
	std := std(values, mean)

	curIndexes := make([]int, 0)

	for index, value := range values {
		if value > (mean + sigma * std) || value <(mean - sigma * std) {
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

func mean(values []float64)float64 {
	var sum float64 = 0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

func std(values []float64, mean float64) float64 {
	var variance float64 = 0
	for _, value := range values {
		variance += math.Pow(value - mean, 2)
	}

	return math.Sqrt(variance / float64(len(values)))
}

