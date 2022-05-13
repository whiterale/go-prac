package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/whiterale/go-prac/internal/repo"
)

type Metric struct {
	MType string   `json:"type"`
	ID    string   `json:"id"`
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
}

func (m *Metric) ValueInt() int64 {
	return *m.Delta
}

func (m *Metric) ValueFloat() float64 {
	return *m.Value
}

func (m *Metric) Name() string {
	return m.ID
}

func (m *Metric) Kind() string {
	return m.MType
}

type Server struct {
	Repo repo.Storer
}

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	var metric Metric
	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Repo.Store(&metric)
}

func (s *Server) Value(w http.ResponseWriter, req *http.Request) {
	var metric Metric
	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Repo.Get -> func (repo.Storer).Get(string, string) (interface{}, bool)
	rawValue, ok := s.Repo.Get(metric.MType, metric.ID)

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if metric.MType == "gauge" {
		value, _ := strconv.ParseFloat(rawValue, 64)
		metric.Value = &value
	}

	if metric.MType == "counter" {
		value, _ := strconv.ParseInt(rawValue, 10, 64)
		metric.Delta = &value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(metric)
}

func (s *Server) Head(w http.ResponseWriter, req *http.Request) {
	metrics := s.Repo.GetAll()
	log.Printf("%+v", metrics)
	head.Execute(w, metrics)
}

var head *template.Template
var headSrc = `{{range $index, $element := .}}
{{$index}}={{$element}}
{{end}}
`

func Listen() {
	srv := Server{Repo: repo.InitInMemory()}
	head = template.Must(template.New("head").Parse(headSrc))
	r := chi.NewRouter()
	r.Post("/update", srv.Update)
	r.Get("/value", srv.Value)
	r.Get("/", srv.Head)
	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
