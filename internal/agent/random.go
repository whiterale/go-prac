package agent

import "math/rand"

type Random struct{}

func (r *Random) Collect() []AgentMetric {
	return []AgentMetric{{
		MType: "gauge",
		ID:    "RandomValue",
		Value: rand.Float64(),
	}}
}
