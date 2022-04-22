package main

import (
	"log"
	"net/http"

	"github.com/whiterale/go-prac/internal/repo"
	"github.com/whiterale/go-prac/internal/server"
)

func main() {
	srv := server.Server{Repo: repo.DevNull{}}
	http.HandleFunc("/update/", srv.Update)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
