package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(req)
	name, kind := vars["name"], vars["kind"]

	if kind != "gauge" && kind != "counter" {
		http.Error(w, "Wrong metric kind", 501)
		return
	}
	metric := &Metric{kind, name, 0, 0}
	if kind == "gauge" {
		value, err := strconv.ParseFloat(vars["value"], 64)
		if err != nil {
			log.Printf("Failed to parse float value for gauge metric: %s", vars["value"])
			http.Error(w, "Bad request", 400)
			return
		}
		metric.valueFloat = value
	}
	if kind == "counter" {
		value, err := strconv.ParseInt(vars["value"], 10, 64)
		if err != nil {
			log.Printf("Failed to parse float value for gauge metric: %s", vars["value"])
			http.Error(w, "Bad request", 400)
			return
		}
		metric.valueInt = value
	}

	log.Printf("%v+", metric)
	s.Repo.Store(metric)
}

func Listen() {
	srv := Server{Repo: repo.InitInMemory()}

	updateRouter := mux.NewRouter()
	updateRouter.HandleFunc("/update/{kind}/{name}/{value}", srv.Update)
	http.Handle("/", updateRouter)
	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
