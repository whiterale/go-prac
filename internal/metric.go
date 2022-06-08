package internal

import (
	"fmt"
)

type Metric struct {
	MType string   `json:"type"`
	ID    string   `json:"id"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}

func (m *Metric) GetValue() interface{} {
	switch m.MType {
	case "gauge":
		return *m.Value
	case "counter":
		return *m.Delta
	default:
		return nil
	}
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

func (m *Metric) PlainText() string {
	if m.MType == "gauge" {
		return fmt.Sprintf("%g", *m.Value)
	}
	if m.MType == "counter" {
		return fmt.Sprintf("%d", *m.Delta)
	}
	return ""
}
