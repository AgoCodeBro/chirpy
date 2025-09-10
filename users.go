package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AgoCodeBro/chirpy/internal/auth"
	"github.com/AgoCodeBro/chirpy/internal/database"
	"github.com/google/uuid"
)

type jsonUser struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Falied to decode request", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Failed to hash password", err)
		return
	}

	createUserArgs := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	}

	result, err := c.db.CreateUser(req.Context(), createUserArgs)
	if err != nil {
		respondWithError(w, 500, "Falied to decode request", err)
		return
	}

	jsonableResult := jsonUser{
		ID:          result.ID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		Email:       result.Email,
		IsChirpyRed: result.IsChirpyRed,
	}

	respondWithJson(w, 201, jsonableResult)

}
