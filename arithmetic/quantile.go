package arithmetic

import (
	"sort"
)

var (
	NormalK = 1.5
)

func Quantile(values []float64, q float64) float64 {
	sort.Float64s(values)

	count := len(values)
	index := q * float64(count - 1)
	lns := int(index)
	delta := index - float64(lns)

	if count == 0 {
		return 0
	}

	if lns == count - 1 {
		return values[lns]
	} else {
		return (1 - delta) * values[lns] + delta * values[lns + 1]
	}
}

func QuantileWithDuration(values []float64, duration int) [][]int {
	indexes := make([][]int, 0)

	sortValues := make([]float64, 0)
	sortValues = append(sortValues, values...)

	quantile1 := Quantile(sortValues, 0.25)
	quantile3 := Quantile(sortValues, 0.75)
	upValue := quantile3 + NormalK * (quantile3 - quantile1)
	lowValue := quantile1 - NormalK * (quantile3 - quantile1)

	curIndexes := make([]int, 0)

	for index, value := range values {
		if value >upValue || value < lowValue {
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