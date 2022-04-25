package metrics

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_populateMemStats(t *testing.T) {

	m := Init()
	m.populateMemStats()

	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)

	for _, key := range []string{"Alloc", "HeapAlloc", "LastGC", "OtherSys"} {
		_, ok := m.gauges[key]
		assert.True(t, ok)
	}
}
