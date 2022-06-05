package reporters

import (
	"log"

	"github.com/whiterale/go-prac/internal"
)

type Log struct{}

func (lr *Log) Report(metrics []*internal.Metric) error {
	for _, m := range metrics {
		log.Printf("Report: %s\n", m)
	}
	return nil
}
