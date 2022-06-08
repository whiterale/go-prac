package collectors

type PollCounter struct{}

func (pc *PollCounter) Collect() []Metric {
	return []Metric{{
		MType: "counter",
		Value: int64(1),
		ID:    "PollCount",
	}}
}
