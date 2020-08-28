package main

import (
	"Go-Tracker/api"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.Use(AddJsonContentTypeHeaderMiddleware)
	r.PathPrefix("/api/v1")
	r.HandleFunc("/storage", api.IdempotentIncrementStorageApi).Methods("POST")
	r.HandleFunc("/storage/{key}", api.IncrementStorageWithKeyApi).Methods("PUT")
	r.HandleFunc("/storage/{key}/increment", api.IncrementStorageWithKeyApi).Methods("PUT")
	r.HandleFunc("/storage/{key}/decrement", api.DecrementStorageApi).Methods("PUT")
	r.HandleFunc("/storage/{key}", api.CheckStorageApi).Methods("GET")
	server := &http.Server {
		Handler: r,
		Addr: "127.0.0.1:1234",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func AddJsonContentTypeHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "api") {
			log.Debugln("Detected API endpoint, adding global json content-type.")
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

