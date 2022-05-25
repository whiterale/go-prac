package metric

import (
	"errors"
	"sync"
)

type Buffer struct {
	sync.Mutex
	metrics map[string]map[string]*Metric
}

func Init() *Buffer {
	buffer := Buffer{metrics: make(map[string]map[string]*Metric)}
	buffer.metrics["gauge"] = make(map[string]*Metric)
	buffer.metrics["counter"] = make(map[string]*Metric)
	return &buffer
}

func (b *Buffer) updateCounter(id string, delta int64) {
	b.Lock()
	defer b.Unlock()
	if metric, ok := b.metrics["counter"][id]; ok {
		*metric.Delta += delta
		return
	}
	b.metrics["counter"][id] = &Metric{
		MType: "counter",
		Delta: &delta,
		ID:    id,
		Value: nil,
	}
}

func (b *Buffer) updateGauge(id string, value float64) {
	b.Lock()
	defer b.Unlock()
	if metric, ok := b.metrics["gauge"][id]; ok {
		*metric.Value = value
		return
	}
	b.metrics["gauge"][id] = &Metric{
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
		if delta, ok := val.(int); ok {
			b.updateCounter(id, int64(delta))
			return nil
		}
		return errors.New("update type mismatch")
	}
	return errors.New("unsupported metric type")
}

func (b *Buffer) Get(mtype string, id string) (*Metric, bool) {
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
	b.metrics["gauge"] = make(map[string]*Metric)
	b.metrics["counter"] = make(map[string]*Metric)
}
