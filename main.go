package main

import (
	"database/sql"
	"github.com/AgoCodeBro/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	dbQueries := database.New(db)

	port := ":8080"
	mux := http.NewServeMux()
	cfg := apiConfig{
		db:       dbQueries,
		platform: os.Getenv("PLATFORM"),
		secret:   os.Getenv("SECRET"),
	}

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /api/healthz", http.HandlerFunc(readyHandler))
	mux.Handle("GET /admin/metrics", http.HandlerFunc(cfg.fileServerHitsHandler))
	mux.Handle("POST /admin/reset", http.HandlerFunc(cfg.metricsResetHandler))
	mux.Handle("POST /api/chirps", http.HandlerFunc(cfg.chirpHandler))
	mux.Handle("POST /api/users", http.HandlerFunc(cfg.createUserHandler))
	mux.Handle("PUT /api/users", http.HandlerFunc(cfg.changeCredentialsHandler))
	mux.Handle("GET /api/chirps", http.HandlerFunc(cfg.getChirps))
	mux.Handle("GET /api/chirps/{chirpID}", http.HandlerFunc(cfg.getChirp))
	mux.Handle("DELETE /api/chirps/{chirpID}", http.HandlerFunc(cfg.deleteChirp))
	mux.Handle("POST /api/login", http.HandlerFunc(cfg.loginHandler))
	mux.Handle("POST /api/refresh", http.HandlerFunc(cfg.refreshHandler))
	mux.Handle("POST /api/revoke", http.HandlerFunc(cfg.revokeHandler))
	mux.Handle("POST /api/polka/webhooks", http.HandlerFunc(cfg.polkaWebhookHandler))

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
