package repo

import "log"

type NameValuer interface {
	Value() interface{}
	Name() string
}

type Storer interface {
	Store(string, NameValuer) error
}

type DevNull struct{}

func (d DevNull) Store(metricType string, m NameValuer) error {

	switch metricType {
	case "gauge":
		value := m.Value().(float64)
		name := m.Name()
		log.Printf("Storing %s=%f into devnull", name, value)
	case "counter":
		value := m.Value().(int64)
		name := m.Name()
		log.Printf("Storing %s=%d into devnull", name, value)
	}
	return nil
}
