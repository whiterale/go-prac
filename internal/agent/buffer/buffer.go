package buffer

import (
	"errors"
	"sync"

	"github.com/whiterale/go-prac/internal"
)

type Buffer struct {
	sync.Mutex
	metrics map[string]map[string]*internal.Metric
}

func Init() *Buffer {
	buffer := Buffer{metrics: make(map[string]map[string]*internal.Metric)}
	buffer.metrics["gauge"] = make(map[string]*internal.Metric)
	buffer.metrics["counter"] = make(map[string]*internal.Metric)
	return &buffer
}

func InitWithData(data map[string]map[string]*internal.Metric) *Buffer {
	buffer := Buffer{metrics: data}
	return &buffer
}

func (b *Buffer) GetRawMetrics() map[string]map[string]*internal.Metric {
	return b.metrics
}

func (b *Buffer) updateCounter(id string, delta int64) {
	b.Lock()
	defer b.Unlock()
	if m, ok := b.metrics["counter"][id]; ok {
		*m.Delta += delta
		return
	}
	b.metrics["counter"][id] = &internal.Metric{
		MType: "counter",
		Delta: &delta,
		ID:    id,
		Value: nil,
	}
}

func (b *Buffer) updateGauge(id string, value float64) {
	b.Lock()
	defer b.Unlock()
	if m, ok := b.metrics["gauge"][id]; ok {
		*m.Value = value
		return
	}
	b.metrics["gauge"][id] = &internal.Metric{
		MType: "gauge",
		Value: &value,
		ID:    id,
		Delta: nil,
	}
}

func (b *Buffer) Update(mtype string, id string, val interface{}) error {
	switch mtype {
	case "gauge":
		if value, ok := val.(float64); ok {
			b.updateGauge(id, value)
			return nil
		}
		return errors.New("update type mismatch")
	case "counter":
		if delta, ok := val.(int64); ok {
			b.updateCounter(id, int64(delta))
			return nil
		}
		return errors.New("update type mismatch")
	}
	return errors.New("unsupported metric type")
}

func (b *Buffer) Get(mtype string, id string) (*internal.Metric, bool) {
	switch mtype {
	case "gauge":
		metric, ok := b.metrics["gauge"][id]
		return metric, ok
	case "counter":
		metric, ok := b.metrics["counter"][id]
		return metric, ok
	}
	return nil, false
}

func (b *Buffer) Flush() {
	b.Lock()
	defer b.Unlock()
	b.metrics["gauge"] = make(map[string]*internal.Metric)
	b.metrics["counter"] = make(map[string]*internal.Metric)
}

func (b *Buffer) Dump() []*internal.Metric {
	res := make([]*internal.Metric, 0, (len(b.metrics["counter"]))+len(b.metrics["gauge"]))

	for _, v := range b.metrics["counter"] {
		res = append(res, v)
	}

	for _, v := range b.metrics["gauge"] {
		res = append(res, v)
	}
	return res
}
