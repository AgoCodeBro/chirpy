package main

import (
	"log"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	port := ":8080"
	mux := http.NewServeMux()
	cfg := apiConfig{}
	
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /api/healthz", http.HandlerFunc(readyHandler))
	mux.Handle("GET /admin/metrics", http.HandlerFunc(cfg.fileServerHitsHandler))
	mux.Handle("POST /admin/reset", http.HandlerFunc(cfg.metricsResetHandler))

	srv := &http.Server{
		Addr    : port,
		Handler : mux,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}

func readyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) fileServerHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	bodyString := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",  cfg.fileServerHits.Load())
	w.Write([]byte(bodyString))
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	cfg.fileServerHits.Store(0)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	handler := func (w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(handler)

}
