package agent

import (
	"github.com/whiterale/go-prac/internal"
	"github.com/whiterale/go-prac/internal/agent/collectors"
)

type Agent struct {
	report     Reporter
	buffer     Updater
	collectors []Collector
}

type Collector interface {
	Collect() []collectors.Metric
}

type Updater interface {
	Update(string, string, interface{}) error
	Dump() []*internal.Metric
	Flush()
}

type Reporter interface {
	Report([]*internal.Metric) error
}

func Init(reporter Reporter, buffer Updater, collectors []Collector) *Agent {
	return &Agent{reporter, buffer, collectors}
}

func (agent *Agent) Poll() {
	for _, c := range agent.collectors {
		metrics := c.Collect()
		for _, m := range metrics {
			agent.buffer.Update(m.MType, m.ID, m.Value)
		}
	}
}

func (agent *Agent) Report() {
	metrics := agent.buffer.Dump()
	agent.report.Report(metrics)
	agent.buffer.Flush()
}
