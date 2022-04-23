package repo

import (
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

func (im *InMemory) Store(m NameValuer) error {
	im.Lock()
	defer im.Unlock()
	switch m.Kind() {
	case "gauge":
		im.gauges[m.Name()] = m.ValueFloat()
	case "counter":
		im.counters[m.Name()] += 1
	}
	return nil
}
