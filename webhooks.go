package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/AgoCodeBro/chirpy/internal/auth"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	api_key := os.Getenv("POLKA_KEY")

	reqKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, 401, "failed to get api key from request", err)
		return
	}

	if api_key != reqKey {
		respondWithError(w, 401, "invalid api key", nil)
		return
	}

	type dataJson struct {
		UserId string `json:"user_id"`
	}
	type reqJson struct {
		Event string   `json:"event"`
		Data  dataJson `json:"data"`
	}

	reqBody := reqJson{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 500, "could not decode request body", err)
		return
	}

	if reqBody.Event == "user.upgraded" {
		userId, err := uuid.Parse(reqBody.Data.UserId)
		if err != nil {
			respondWithError(w, 500, "failed to parse user id", err)
			return
		}

		err = cfg.db.UpgradeUser(req.Context(), userId)
		if err != nil {
			respondWithError(w, 404, "user not found", err)
			return
		}
	}

	respondWithJson(w, 204, nil)
}
