package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
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

	// resp, err = http.Get(ts.URL)
	// resp.Body.Close()
	// assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	// assert.NoError(t, err)

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

	// resp, err = http.Get(ts.URL + "/value/gauge/g1")

	// body, _ := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()
	// assert.NoError(t, err)
	// assert.Equal(t, "100.1", string(body))
}
