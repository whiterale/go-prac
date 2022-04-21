package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

type gauge float64
type counter int64

type Metrics struct {
	memstats *runtime.MemStats
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
	if _, ok := m.counters["PollCount"]; ok {
		m.counters["PollCount"]++
	} else {
		m.counters["PollCount"] = 1
	}
}

func (m *Metrics) populateRandomValue() {
	m.gauges["RandomValue"] = gauge(rand.Float64())
}

func (m *Metrics) Poll() {
	m.populateMemStats()
	m.populatePollCounter()
	m.populateRandomValue()

}

func (m *Metrics) Report(format string) error {
	var urls []string
	// TODO: refactor this to metrics.Dump(...) or smth
	for k, v := range m.gauges {
		// TODO: move host to const/config/envvar
		url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%.2f", k, v)
		urls = append(urls, url)
	}

	for k, v := range m.counters {
		url := fmt.Sprintf("https://localhost:8080/update/counter/%s/%d", k, v)
		urls = append(urls, url)
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		u := url
		go func() {
			defer wg.Done()
			resp, err := http.Post(u, "text/plain", nil)
			if err != nil {
				log.Printf("failed to send metrics: %e", err)
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()
	return nil
}
