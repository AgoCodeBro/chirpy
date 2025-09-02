package main

import (
	"strings"
	"time"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/AgoCodeBro/chirpy/internal/database"
)

type jsonableChirpStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, req *http.Request) {
	type chirpParams struct {
		Body   string `json:"body"`
		UserId string `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	chirpJson := chirpParams{}
	err := decoder.Decode(&chirpJson)
	if err != nil {
		respondWithError(w, 500, "Failed to decode request", err)
		return
	}
	
	if len(chirpJson.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)		
		return
	}

	cleanChirp := cleanChirpBody(chirpJson.Body)
	parsedUUID, err := uuid.Parse(chirpJson.UserId)
	if err != nil {
		respondWithError(w, 500, "Failed to parse uuid", err)
		return
	}

	createChirpArgs := database.CreateChirpParams {
		Body   : cleanChirp,
		UserID : parsedUUID,
	}
	
	result, err := cfg.db.CreateChirp(req.Context(), createChirpArgs)
	if err != nil {
		respondWithError(w, 500, "Failed to save post", err)
		return
	}
	
	jsonableResult := jsonableChirpStruct {
		ID        : result.ID,
		CreatedAt : result.CreatedAt,
		UpdatedAt : result.UpdatedAt,
		Body      : result.Body,
		UserID    : result.UserID,
	}

	respondWithJson(w, 201, jsonableResult)
}
		
func cleanChirpBody(msg string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")
	for i, word := range words {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				words[i] = "****"
			}
		}
	}

	msg = strings.Join(words, " ")

	return msg
}


func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, req *http.Request) {
	result, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "Failed to get chirps", err)
		return
	}

	jsonableChirps := make([]jsonableChirpStruct, len(result))
	for i, resultChirp := range result {
		jsonableChirps[i] = jsonableChirpStruct{
			ID        : resultChirp.ID,
			CreatedAt : resultChirp.CreatedAt,
			UpdatedAt : resultChirp.UpdatedAt,
			Body      : resultChirp.Body,
			UserID    : resultChirp.UserID,
		}

	}

	respondWithJson(w, 200, jsonableChirps)
}
			
func (cfg *apiConfig) getChirp(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 500, "Failed to parse ID", err)
		return
	}

	result, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	jsonableChirp := jsonableChirpStruct{
		ID        : result.ID,
		CreatedAt : result.CreatedAt,
		UpdatedAt : result.UpdatedAt,
		Body      : result.Body,
		UserID    : result.UserID,
	}

	respondWithJson(w, 200, jsonableChirp)

}

		









