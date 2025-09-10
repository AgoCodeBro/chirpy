package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) fileServerHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	bodyString := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	w.Write([]byte(bodyString))
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return
	}

	w.WriteHeader(200)
	cfg.fileServerHits.Store(0)
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to reset the users database", err)
		return
	}

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(handler)

}
