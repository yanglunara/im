package logic

import "math"

type weighted struct {
	fixedWeight   int64
	currentWeight int64
	currentConns  int64
	updated       int64
	addres        string
	hostName      string
	region        string
}

func (w *weighted) chosen() {
	w.currentConns++
}

func (w *weighted) calculateWeight(totalWeight, totalConns int64, gainWeight float64) {
	fixedWeight := float64(w.fixedWeight) * gainWeight
	totalWeight += int64(fixedWeight) - w.fixedWeight
	if totalConns > 0 {
		weightRatio := fixedWeight / float64(totalWeight)
		var connRatio float64
		if totalConns != 0 {
			connRatio = float64(w.currentConns) / float64(totalConns) * 0.5
		}
		diff := weightRatio - connRatio
		multipe := diff * float64(totalConns)
		floor := math.Floor(multipe)
		if floor-multipe >= -0.5 {
			w.currentWeight = int64(fixedWeight + floor)
		} else {
			w.currentWeight = int64(fixedWeight + math.Ceil(multipe))
		}

		if diff < 0 {
			if 1 > w.currentWeight {
				w.currentWeight = 1
			}

		} else {
			if 1<<20 < w.currentWeight {
				w.currentWeight = 1 << 20
			}
		}
	} else {
		w.currentWeight = 0
	}
}
