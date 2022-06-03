package collectors

import "runtime"

type Runtime struct{}

func (r Runtime) Collect() []Metric {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	result := []Metric{
		{"Alloc", "gauge", float64(ms.Alloc)},
		{"BuckHashSys", "gauge", float64(ms.BuckHashSys)},
		{"Frees", "gauge", float64(ms.Frees)},
		{"GCCPUFraction", "gauge", float64(ms.GCCPUFraction)},
		{"GCSys", "gauge", float64(ms.GCSys)},
		{"HeapAlloc", "gauge", float64(ms.HeapAlloc)},
		{"HeapIdle", "gauge", float64(ms.HeapIdle)},
		{"HeapInuse", "gauge", float64(ms.HeapInuse)},
		{"HeapObjects", "gauge", float64(ms.HeapObjects)},
		{"HeapReleased", "gauge", float64(ms.HeapReleased)},
		{"HeapSys", "gauge", float64(ms.HeapSys)},
		{"LastGC", "gauge", float64(ms.LastGC)},
		{"Lookups", "gauge", float64(ms.Lookups)},
		{"MCacheInuse", "gauge", float64(ms.MCacheInuse)},
		{"MCacheSys", "gauge", float64(ms.MCacheSys)},
		{"MSpanInuse", "gauge", float64(ms.MSpanInuse)},
		{"MSpanSys", "gauge", float64(ms.MSpanSys)},
		{"Mallocs", "gauge", float64(ms.Mallocs)},
		{"NextGC", "gauge", float64(ms.NextGC)},
		{"NumForcedGC", "gauge", float64(ms.NumForcedGC)},
		{"NumGC", "gauge", float64(ms.NumGC)},
		{"OtherSys", "gauge", float64(ms.OtherSys)},
		{"PauseTotalNs", "gauge", float64(ms.PauseTotalNs)},
		{"StackInuse", "gauge", float64(ms.StackInuse)},
		{"StackSys", "gauge", float64(ms.StackSys)},
		{"Sys", "gauge", float64(ms.Sys)},
		{"TotalAlloc", "gauge", float64(ms.TotalAlloc)},
	}
	return result
}
