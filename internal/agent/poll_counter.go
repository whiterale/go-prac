package agent

type PollCounter struct{}

func (pc *PollCounter) Collect() []AgentMetric {
	return []AgentMetric{{
		MType: "counter",
		Value: int64(1),
		ID:    "PollCounter",
	}}
}
