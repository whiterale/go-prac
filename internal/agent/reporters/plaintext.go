package reporters

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/whiterale/go-prac/internal"
)

type PlainText struct {
	Host string
}

func (pt *PlainText) Report(metrics []*internal.Metric) error {
	var wg sync.WaitGroup
	for _, m := range metrics {
		wg.Add(1)
		u := fmt.Sprintf("%s/update/%s", pt.Host, m.String())
		go func() {
			defer wg.Done()
			resp, err := http.Post(u, "text/plain", nil)
			if err != nil {
				log.Printf("failed to send metrics, url=%s: %e", u, err)
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()
	return nil
}
