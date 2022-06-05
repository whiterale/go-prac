package collectors

import "math/rand"

type Random struct{}

func (r *Random) Collect() []Metric {
	return []Metric{{
		MType: "gauge",
		ID:    "RandomValue",
		Value: rand.Float64(),
	}}
}
