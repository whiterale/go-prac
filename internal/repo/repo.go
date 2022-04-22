package repo

import "log"

type Metric interface {
	GetType() string
	GetValue() float64
}

type Storer interface {
	Store(Metric) error
}

type DevNull struct{}

func (d DevNull) Store(m Metric) error {
	log.Printf("Storing %s=%f into devnull", m.GetType(), m.GetValue())
	return nil
}
