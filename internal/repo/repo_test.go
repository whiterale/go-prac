package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type GaugeTest struct{}

func (g *GaugeTest) ValueInt() int64 {
	panic("This should never happen")
}

func (g *GaugeTest) ValueFloat() float64 {
	return float64(3.14)
}

func (g *GaugeTest) Name() string {
	return "gaugeTest"
}

func (g *GaugeTest) Kind() string {
	return "gauge"
}

type CounterTest struct{}

func (c *CounterTest) ValueInt() int64 {
	return int64(42)
}

func (c *CounterTest) ValueFloat() float64 {
	panic("This should never happen")
}

func (c *CounterTest) Name() string {
	return "counterTest"
}

func (c *CounterTest) Kind() string {
	return "counter"
}

func TestInMemory_Store(t *testing.T) {
	inMem := InitInMemory()
	gaugeMetric := GaugeTest{}
	counterMetric := CounterTest{}

	inMem.Store(&gaugeMetric)
	gaugeVal, gaugeOk := inMem.gauges["gaugeTest"]
	assert.True(t, gaugeOk)
	assert.Equal(t, 3.14, gaugeVal)

	inMem.Store(&counterMetric)
	countVal, countOk := inMem.counters["counterTest"]
	assert.True(t, countOk)
	assert.Equal(t, int64(42), countVal)

}
