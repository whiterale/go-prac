package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/whiterale/go-prac/internal/repo"
)

const (
	updateURL = "/update"
	valueURL  = "/value"

	CounterJSON = `{
		"id": "some-counter",
		"type": "counter",
		"delta": 42
	}`
	GaugeJSON = `{
		"id": "some-gauge",
		"type": "gauge",
		"value": 420.69
	}`

	GetCounterJSON = `{
		"id": "some-counter",
		"type": "counter"
	}`

	GetGaugeJSON = `{
		"id": "some-gauge",
		"type": "gauge"
	}`
)

// TODO: refactor to table tests
func TestServer_Update(t *testing.T) {

	var metric Metric

	repo := repo.InitInMemory()
	srv := Server{repo}
	r := chi.NewRouter()
	r.Post(updateURL, srv.Update)
	r.Post(valueURL, srv.Value)
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NoError(t, err)

	for _, payload := range []string{CounterJSON, CounterJSON, GaugeJSON} {
		_, err = http.Post(ts.URL+updateURL, "application/json", bytes.NewReader([]byte(payload)))
		assert.NoError(t, err)
	}

	cntVal, ok := srv.Repo.Get("counter", "some-counter")
	assert.True(t, ok)
	assert.Equal(t, cntVal, fmt.Sprintf("%d", 42*2))

	gaugeVal, ok := srv.Repo.Get("gauge", "some-gauge")
	assert.True(t, ok)
	assert.Equal(t, gaugeVal, fmt.Sprintf("%g", 420.69))

	res, _ := http.Post(ts.URL+valueURL, "application/json", bytes.NewReader([]byte(GetGaugeJSON)))
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&metric))

	assert.Equal(t, metric.ID, "some-gauge")
	assert.Equal(t, 420.69, *metric.Value)

	res, _ = http.Post(ts.URL+valueURL, "application/json", bytes.NewReader([]byte(GetCounterJSON)))
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&metric))

	assert.Equal(t, metric.ID, "some-counter")
	assert.Equal(t, int64(42*2), *metric.Delta)

}
