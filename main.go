package main

import (
	"os"
	"log"
	"database/sql"
	"net/http"
	"sync/atomic"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/AgoCodeBro/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
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
		db       : dbQueries,
		platform : os.Getenv("PLATFORM"),
	}
	
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /api/healthz", http.HandlerFunc(readyHandler))
	mux.Handle("GET /admin/metrics", http.HandlerFunc(cfg.fileServerHitsHandler))
	mux.Handle("POST /admin/reset", http.HandlerFunc(cfg.metricsResetHandler))
	mux.Handle("POST /api/chirps", http.HandlerFunc(cfg.chirpHandler))
	mux.Handle("POST /api/users", http.HandlerFunc(cfg.createUserHandler))
	mux.Handle("GET /api/chirps", http.HandlerFunc(cfg.getAllChirps))
	mux.Handle("GET /api/chirps/{chirpID}", http.HandlerFunc(cfg.getChirp))

	srv := &http.Server{
		Addr    : port,
		Handler : mux,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

