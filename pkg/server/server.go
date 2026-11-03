package server

import (
	"TODO_List/pkg/api"
	"log"
	"net/http"
	"time"
)

func Server(l *log.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	api.Init(mux)

	return &http.Server{
		Addr:         ":7540",
		Handler:      mux,
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
}
