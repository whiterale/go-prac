package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"

	"github.com/whiterale/go-prac/internal/server"
)

type gauge float64
type counter int64

type Metrics struct {
	gauges   map[string]gauge
	counters map[string]counter
}

func Init() *Metrics {
	metrics := Metrics{}
	metrics.counters = make(map[string]counter)
	metrics.gauges = make(map[string]gauge)
	return &metrics
}

func (m *Metrics) populateMemStats() {
	memstats := &runtime.MemStats{}
	runtime.ReadMemStats(memstats)

	m.gauges["Alloc"] = gauge(memstats.Alloc)
	m.gauges["BuckHashSys"] = gauge(memstats.BuckHashSys)
	m.gauges["Frees"] = gauge(memstats.Frees)
	m.gauges["GCCPUFraction"] = gauge(memstats.GCCPUFraction)
	m.gauges["GCSys"] = gauge(memstats.GCSys)

	m.gauges["HeapAlloc"] = gauge(memstats.HeapAlloc)
	m.gauges["HeapIdle"] = gauge(memstats.HeapIdle)
	m.gauges["HeapInuse"] = gauge(memstats.HeapInuse)
	m.gauges["HeapObjects"] = gauge(memstats.HeapObjects)
	m.gauges["HeapReleased"] = gauge(memstats.HeapReleased)
	m.gauges["HeapSys"] = gauge(memstats.HeapSys)

	m.gauges["LastGC"] = gauge(memstats.LastGC)
	m.gauges["Lookups"] = gauge(memstats.Lookups)

	m.gauges["MCacheInuse"] = gauge(memstats.MCacheInuse)
	m.gauges["MCacheSys"] = gauge(memstats.MCacheSys)
	m.gauges["MSpanInuse"] = gauge(memstats.MSpanInuse)
	m.gauges["MSpanSys"] = gauge(memstats.MSpanSys)

	m.gauges["Mallocs"] = gauge(memstats.Mallocs)

	m.gauges["NextGC"] = gauge(memstats.NextGC)

	m.gauges["NumForcedGC"] = gauge(memstats.NumForcedGC)
	m.gauges["NumGC"] = gauge(memstats.NumGC)

	m.gauges["OtherSys"] = gauge(memstats.OtherSys)
	m.gauges["PauseTotalNs"] = gauge(memstats.PauseTotalNs)
	m.gauges["StackInuse"] = gauge(memstats.StackInuse)
	m.gauges["StackSys"] = gauge(memstats.StackSys)
	m.gauges["Sys"] = gauge(memstats.Sys)
	m.gauges["TotalAlloc"] = gauge(memstats.TotalAlloc)
}

func (m *Metrics) populatePollCounter() {
	m.counters["PollCount"] += 1
}

func (m *Metrics) populateRandomValue() {
	m.gauges["RandomValue"] = gauge(rand.Float64())
}

func (m *Metrics) dumpCounters(format string) []string {
	var res []string
	for k, v := range m.counters {
		str := fmt.Sprintf(format, k, v)
		res = append(res, str)
	}
	return res
}

func (m *Metrics) dumpGauges(format string) []string {
	var res []string
	for k, v := range m.gauges {
		str := fmt.Sprintf(format, k, v)
		res = append(res, str)
	}
	return res
}

func (m *Metrics) dumpGaugesJSON() []*server.Metric {
	res := make([]*server.Metric, 0, len(m.gauges))
	for k, v := range m.gauges {
		delta := float64(v)
		metric := &server.Metric{
			MType: "gauge",
			ID:    k,
			Value: &delta,
		}
		res = append(res, metric)
	}
	return res
}

func (m *Metrics) dumpCountersJSON() []*server.Metric {
	res := make([]*server.Metric, 0, len(m.gauges))
	for k, v := range m.counters {
		delta := int64(v)
		metric := &server.Metric{
			MType: "counter",
			ID:    k,
			Delta: &delta,
		}
		res = append(res, metric)
	}
	return res
}

func (m *Metrics) Poll() {
	m.populateMemStats()
	m.populatePollCounter()
	m.populateRandomValue()
}

func (m *Metrics) Report(format string) error {

	gauges := m.dumpGaugesJSON()
	counters := m.dumpCountersJSON()

	var wg sync.WaitGroup
	for _, metric := range append(gauges, counters...) {
		wg.Add(1)

		go func(m *server.Metric) {
			defer wg.Done()
			payload, _ := json.Marshal(m)
			resp, err := http.Post("http://localhost:8080/update", "application/json", bytes.NewReader([]byte(payload)))
			if err != nil {
				log.Printf("failed to send metrics: %e", err)
				return
			}
			if resp.StatusCode != http.StatusCreated {
				log.Printf("Got %d, expected %d", resp.StatusCode, http.StatusCreated)
			}
			resp.Body.Close()
		}(metric)
	}
	wg.Wait()
	return nil
}
