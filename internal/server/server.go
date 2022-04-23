package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/whiterale/go-prac/internal/repo"
)

type Gauge struct {
	name  string
	value float64
}

func (g *Gauge) Value() interface{} {
	return g.value
}

func (g *Gauge) Name() string {
	return g.name
}

type Counter struct {
	name  string
	value int64
}

func (c *Counter) Value() interface{} {
	return c.value
}

func (c *Counter) Name() string {
	return c.name
}

type Server struct {
	Repo repo.Storer
}

func (s *Server) CounterUpdate(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metric := &Counter{}

	metric.name = vars["name"]

	value, err := strconv.ParseInt(vars["value"], 10, 64)
	if err != nil {
		log.Printf("Failed to parse int value for counter metric: %s", vars["value"])
	}
	metric.value = value
	s.Repo.Store("counter", metric)
	return
}

func (s *Server) GaugeUpdate(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metric := &Gauge{}

	metric.name = vars["name"]

	value, err := strconv.ParseFloat(vars["value"], 64)
	if err != nil {
		log.Printf("Failed to parse float value for gauge mteric: %s", vars["value"])
	}
	metric.value = value
	s.Repo.Store("gauge", metric)
	return
}

func Listen() {
	srv := Server{Repo: repo.DevNull{}}

	updateRouter := mux.NewRouter()
	updateRouter.HandleFunc("/update/counter/{name:[a-zA-Z]+}/{value:[0-9]+}", srv.CounterUpdate)
	updateRouter.HandleFunc("/update/gauge/{name:[a-zA-Z]+}/{value:[0-9]+\\.[0-9]*}", srv.GaugeUpdate)
	http.Handle("/", updateRouter)
	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
