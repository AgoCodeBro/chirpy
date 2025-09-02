package main

import (
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
)

func (c *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	type reqParams struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Falied to decode request", err)
		return
	}


	result, err := c.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "Falied to decode request", err)
		return
	}

	type jsonUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	jsonableResult := jsonUser{
		ID        : result.ID,
		CreatedAt : result.CreatedAt,
		UpdatedAt : result.UpdatedAt,
		Email     : result.Email,
	}

	respondWithJson(w, 201, jsonableResult)
	
}


	

