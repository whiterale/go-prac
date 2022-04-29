package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/whiterale/go-prac/internal/repo"
)

type Metric struct {
	kind       string
	name       string
	valueFloat float64
	valueInt   int64
}

func (m *Metric) ValueInt() int64 {
	return m.valueInt
}

func (m *Metric) ValueFloat() float64 {
	return m.valueFloat
}

func (m *Metric) Name() string {
	return m.name
}

func (m *Metric) Kind() string {
	return m.kind
}

type Server struct {
	Repo repo.Storer
}

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {

	name := chi.URLParam(req, "name")
	kind := chi.URLParam(req, "kind")
	rawValue := chi.URLParam(req, "value")

	log.Printf("name:%s, kind:%s, val:%s", name, kind, rawValue)

	if kind != "gauge" && kind != "counter" {
		http.Error(w, "Wrong metric kind", http.StatusNotImplemented)
		return
	}
	metric := &Metric{kind, name, 0, 0}
	if kind == "gauge" {
		value, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		metric.valueFloat = value
	}
	if kind == "counter" {
		value, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		metric.valueInt = value
	}

	log.Printf("%v+", metric)
	s.Repo.Store(metric)
}

func (s *Server) Value(w http.ResponseWriter, req *http.Request) {
	log.Print("value handler")
	name := chi.URLParam(req, "name")
	kind := chi.URLParam(req, "kind")

	log.Printf("%s, %s", name, kind)
	res, ok := s.Repo.Get(kind, name)
	log.Printf("%s", res)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
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
	r.Post("/update/{kind}/{name}/{value}", srv.Update)
	r.Get("/value/{kind}/{name}", srv.Value)
	r.Get("/", srv.Head)
	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
