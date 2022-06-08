package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/whiterale/go-prac/internal"
	"github.com/whiterale/go-prac/internal/agent/buffer"
)

type Updater interface {
	Update(string, string, interface{}) error
	Dump() []*internal.Metric
	Get(string, string) (*internal.Metric, bool)
	Flush()
}

type Server struct {
	Storage Updater
}

func makeFloat(raw string) (float64, error) {
	return strconv.ParseFloat(raw, 64)
}

func makeInt(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	mtype := chi.URLParam(req, "mtype")
	rawValue := chi.URLParam(req, "value")

	switch mtype {
	case "counter":
		value, err := makeInt(rawValue)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Error parsing delta for counter: %v", err)
			return
		}
		err = s.Storage.Update(mtype, id, value)
		if err != nil {
			http.Error(w, "Bad request", http.StatusInternalServerError)
			log.Printf("Error updating counter: %v", err)
			return
		}
		return
	case "gauge":
		value, err := makeFloat(rawValue)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err = s.Storage.Update(mtype, id, value)
		if err != nil {
			http.Error(w, "Bad request", http.StatusInternalServerError)
			log.Printf("Error updating gauge: %v", err)
			return
		}
		return
	default:
		http.Error(w, "Unsuppoerted metric type", http.StatusNotImplemented)
		log.Printf("Unsupported metric type: %s", mtype)
		return
	}
}

func (s *Server) UpdateJSON(w http.ResponseWriter, req *http.Request) {
	var m internal.Metric
	var err error
	err = json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.Storage.Update(m.MType, m.ID, m.GetValue())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) ValueJSON(w http.ResponseWriter, req *http.Request) {

	var m internal.Metric
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric, ok := s.Storage.Get(m.MType, m.ID)
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(metric)
}

func (s *Server) Value(w http.ResponseWriter, req *http.Request) {
	mtype := chi.URLParam(req, "mtype")
	id := chi.URLParam(req, "id")

	res, ok := s.Storage.Get(mtype, id)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res.PlainText()))
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func (s *Server) Head(w http.ResponseWriter, req *http.Request) {
	metrics := s.Storage.Dump()
	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].ID < metrics[j].ID
	})
	for _, m := range metrics {
		fmt.Fprintf(w, "%s %s\n", m.ID, m.PlainText())
	}
}

func Listen() {
	srv := Server{Storage: buffer.Init()}
	r := chi.NewRouter()

	r.Post("/update/{mtype}/{id}/{value}", srv.Update)
	r.Get("/value/{mtype}/{id}", srv.Value)
	r.Get("/", srv.Head)

	r.Post("/update/", srv.UpdateJSON)
	r.Post("/value/", srv.ValueJSON)

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
