package metric

import (
	"fmt"
)

type Metric struct {
	MType string
	ID    string
	Value *float64
	Delta *int64
}

func (m *Metric) String() string {
	if m.MType == "gauge" {
		return fmt.Sprintf("%s/%s/%g", m.MType, m.ID, *m.Value)
	}
	if m.MType == "counter" {
		return fmt.Sprintf("%s/%s/%d", m.MType, m.ID, *m.Delta)
	}
	return ""
}
