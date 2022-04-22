package server

import (
	"net/http"

	"github.com/whiterale/go-prac/internal/repo"
)

type Metric struct {
	Type string
	Val  float64
}

func (m *Metric) GetValue() float64 {
	return m.Val
}

func (m *Metric) GetType() string {
	return m.Type
}

func getMetricFromURL(url string) (*Metric, error) {
	return nil, nil
}

type Server struct {
	Repo repo.Storer
}

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path[1:] // remove heading slash
	metric, err := getMetricFromURL(url)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}
	go s.Repo.Store(metric)
	w.Write([]byte(req.URL.Path))
}
