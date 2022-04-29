package repo

import (
	"fmt"
	"log"
	"sync"
)

type NameValuer interface {
	ValueInt() int64
	ValueFloat() float64
	Name() string
	Kind() string
}

type Storer interface {
	Store(NameValuer) error
	Get(string, string) (string, bool)
	GetAll() map[string]string
}

type DevNull struct{}

func (d DevNull) Store(m NameValuer) error {
	switch m.Kind() {
	case "gauge":
		log.Printf("Storing %s:%s=%f into devnull", m.Kind(), m.Name(), m.ValueFloat())
	case "counter":
		log.Printf("Storing %s:%s=%d into devnull", m.Kind(), m.Name(), m.ValueInt())
	}
	return nil
}

type InMemory struct {
	sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
}

func InitInMemory() *InMemory {
	counters := make(map[string]int64)
	gauges := make(map[string]float64)
	return &InMemory{
		counters: counters,
		gauges:   gauges,
	}
}

func (im *InMemory) GetAll() map[string]string {
	res := make(map[string]string)
	for k, v := range im.gauges {
		res[k] = fmt.Sprintf("%g", v)
	}
	for k, v := range im.counters {
		res[k] = fmt.Sprintf("%d", v)
	}
	return res
}

func (im *InMemory) Store(m NameValuer) error {
	im.Lock()
	defer im.Unlock()
	metricName := m.Name()
	switch m.Kind() {
	case "gauge":
		im.gauges[metricName] = m.ValueFloat()
	case "counter":
		im.counters[metricName] += m.ValueInt()
	}
	return nil
}

func (im *InMemory) Get(kind string, name string) (string, bool) {
	if kind == "gauge" {
		value, ok := im.gauges[name]
		return fmt.Sprintf("%g", value), ok
	}
	if kind == "counter" {
		value, ok := im.counters[name]
		return fmt.Sprintf("%d", value), ok
	}
	return "", false
}
