package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/whiterale/go-prac/internal"
	"github.com/whiterale/go-prac/internal/agent/buffer"
)

func TestServer_Update(t *testing.T) {
	var resp *http.Response
	var err error
	buf := buffer.Init()
	srv := Server{Storage: buf}
	r := chi.NewRouter()
	r.Post("/update/{mtype}/{id}/{value}", srv.Update)
	r.Get("/value/{mtype}/{id}", srv.Value)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err = http.Post(ts.URL+"/update/gauge/g1/100.1", "text/plain", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	assert.NoError(t, err)

	val, ok := buf.Get("gauge", "g1")
	assert.True(t, ok)
	assert.Equal(t, float64(100.1), *val.Value)

	resp, err = http.Post(ts.URL+"/update/counter/c1/10", "text/plain", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	assert.NoError(t, err)

	resp, err = http.Post(ts.URL+"/update/counter/c1/10", "text/plain", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	assert.NoError(t, err)

	val, ok = buf.Get("counter", "c1")
	assert.True(t, ok)
	assert.Equal(t, int64(20), *val.Delta)
}

const (
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
	updateURL = "/update"
	valueURL  = "/value"
)

func TestServer_JSON(t *testing.T) {
	var resp *http.Response
	var m *internal.Metric
	var err error
	var ok bool
	buf := buffer.Init()
	srv := Server{Storage: buf}
	r := chi.NewRouter()

	r.Post("/update/{mtype}/{id}/{value}", srv.Update)
	r.Get("/value/{mtype}/{id}", srv.Value)

	r.Post("/update", srv.UpdateJSON)
	r.Post("/value", srv.ValueJSON)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, payload := range []string{CounterJSON, CounterJSON, GaugeJSON} {
		_, err = http.Post(ts.URL+updateURL, "application/json", bytes.NewReader([]byte(payload)))
		assert.NoError(t, err)
	}
	gauge, ok := buf.Get("gauge", "some-gauge")
	assert.True(t, ok)

	counter, ok := buf.Get("counter", "some-counter")
	assert.True(t, ok)

	resp, _ = http.Post(ts.URL+valueURL, "application/json", bytes.NewReader([]byte(GetGaugeJSON)))
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&m))
	assert.Equal(t, *m.Value, *gauge.Value)
	assert.Equal(t, m.MType, gauge.MType)

	resp, _ = http.Post(ts.URL+valueURL, "application/json", bytes.NewReader([]byte(GetCounterJSON)))
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&m))
	assert.Equal(t, *m.Delta, *counter.Delta)
	assert.Equal(t, m.MType, counter.MType)
	defer ts.Close()
}
