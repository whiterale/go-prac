package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/whiterale/go-prac/internal"
	"github.com/whiterale/go-prac/internal/agent/buffer"
	"github.com/whiterale/go-prac/internal/agent/collectors"
)

type testRep struct {
	contentText []string
	contentRaw  map[string]*internal.Metric
}

func (r *testRep) Report(metrics []*internal.Metric) error {
	contentText := make([]string, len(metrics))
	contentRaw := make(map[string]*internal.Metric)
	for i, m := range metrics {
		contentText[i] = m.String()
		contentRaw[m.ID] = m
	}
	r.contentText = contentText
	r.contentRaw = contentRaw
	return nil
}

func TestAgent(t *testing.T) {

	buffer := buffer.Init()
	reporter := &testRep{
		contentText: make([]string, 0),
		contentRaw:  make(map[string]*internal.Metric),
	}
	agent := Init(
		reporter,
		buffer,
		[]Collector{&collectors.Random{}, &collectors.PollCounter{}, &collectors.Runtime{}},
	)
	agent.Poll()
	agent.Poll()
	agent.Report()

	var keys []string
	for k := range reporter.contentRaw {
		keys = append(keys, k)
	}
	expectedKeys := []string{
		"PollCounter",
		"MCacheSys",
		"Mallocs",
		"StackInuse",
		"HeapObjects",
		"NextGC",
		"Frees",
		"OtherSys",
		"MSpanSys",
		"HeapInuse",
		"HeapSys",
		"Lookups",
		"NumForcedGC",
		"RandomValue",
		"NumGC",
		"StackSys",
		"HeapAlloc",
		"BuckHashSys",
		"GCCPUFraction",
		"GCSys",
		"Alloc",
		"HeapReleased",
		"LastGC",
		"MCacheInuse",
		"MSpanInuse",
		"HeapIdle",
		"Sys",
		"TotalAlloc",
		"PauseTotalNs",
	}
	assert.ElementsMatch(t, expectedKeys, keys)
	assert.Equal(t, int64(2), int64(*reporter.contentRaw["PollCounter"].Delta))
}
