package reporters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/whiterale/go-prac/internal"
)

type JSON struct {
	Host string
}

func (j *JSON) Report(metrics []*internal.Metric) error {
	var wg sync.WaitGroup
	for _, m := range metrics {
		wg.Add(1)
		u := fmt.Sprintf("http://%s/update/", j.Host)
		payload, err := json.Marshal(m)
		if err != nil {
			log.Printf("Failed to marshal metric %v: %v", m, err)
			continue
		}
		go func(payload []byte) {
			defer wg.Done()

			resp, err := http.Post(u, "application/json", bytes.NewReader(payload))
			if err != nil {
				log.Printf("failed to send metrics: %e", err)
				return
			}
			resp.Body.Close()
		}(payload)
	}
	wg.Wait()
	return nil
}
