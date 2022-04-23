package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/whiterale/go-prac/internal/repo"
)

type Metric struct {
	kind  string
	name  string
	value float64
}

func (m *Metric) ValueInt() int64 {
	return int64(m.value)
}

func (m *Metric) ValueFloat() float64 {
	return m.value
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
	value, err := strconv.ParseFloat(vars["value"], 64)
	if err != nil {
		log.Printf("Failed to parse float value for gauge mteric: %s", vars["value"])
		http.Error(w, "Bad request", 400)
		return
	}

	metric := &Metric{kind, name, value}
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
