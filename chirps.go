package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AgoCodeBro/chirpy/internal/auth"
	"github.com/AgoCodeBro/chirpy/internal/database"
	"github.com/google/uuid"
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
		Body string `json:"body"`
	}

	authString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "no auth header found", err)
		return
	}

	authId, err := auth.ValidateJWT(authString, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "invalid JWT", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	chirpJson := chirpParams{}
	err = decoder.Decode(&chirpJson)
	if err != nil {
		respondWithError(w, 500, "Failed to decode request", err)
		return
	}

	if len(chirpJson.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	cleanChirp := cleanChirpBody(chirpJson.Body)

	createChirpArgs := database.CreateChirpParams{
		Body:   cleanChirp,
		UserID: authId,
	}

	result, err := cfg.db.CreateChirp(req.Context(), createChirpArgs)
	if err != nil {
		respondWithError(w, 500, "Failed to save post", err)
		return
	}

	jsonableResult := jsonableChirpStruct{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body:      result.Body,
		UserID:    result.UserID,
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
			ID:        resultChirp.ID,
			CreatedAt: resultChirp.CreatedAt,
			UpdatedAt: resultChirp.UpdatedAt,
			Body:      resultChirp.Body,
			UserID:    resultChirp.UserID,
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
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body:      result.Body,
		UserID:    result.UserID,
	}

	respondWithJson(w, 200, jsonableChirp)

}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "failed to get access token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "failed to authenticate user", err)
		return
	}

	chirpId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 500, "failed to parse chirp id", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, 403, "not posted by user", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirp.ID)
	respondWithJson(w, 204, nil)
}
