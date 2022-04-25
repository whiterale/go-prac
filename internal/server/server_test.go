package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/whiterale/go-prac/internal/repo"
)

func TestServer_Update(t *testing.T) {

	srv := Server{Repo: repo.InitInMemory()}
	head = template.Must(template.New("head").Parse(headSrc))
	r := chi.NewRouter()
	r.Post("/update/{kind}/{name}/{value}", srv.Update)
	r.Get("/value/{kind}/{name}", srv.Value)

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NoError(t, err)

	resp, err = http.Post(ts.URL+"/update/gauge/g1/100.1", "text/plain", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, err)

	resp, err = http.Get(ts.URL + "/value/gauge/g1")
	body, _ := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "100.1", string(body))
}
