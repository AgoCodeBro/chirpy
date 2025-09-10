package main

import (
	"net/http"
	"encoding/json"

	"github.com/AgoCodeBro/chirpy/internal/auth"
	"github.com/AgoCodeBro/chirpy/internal/database"
)

func (cfg *apiConfig) changeCredentialsHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "invalid access token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "invalid access token", err)
		return
	}

	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	body := reqParams{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&body)
	if err != nil {
		respondWithError(w, 500, "failed to decode request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w, 500, "failed to hash password", err)
		return
	}

	changePasswordArgs := database.ChangePasswordParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	}

	returnedUser, err := cfg.db.ChangePassword(req.Context(), changePasswordArgs)
	if err != nil {
		respondWithError(w, 500, "failed to store new credentials", err)
		return
	}

	jsonableUser := jsonUser{
		ID:             returnedUser.ID,
		CreatedAt:      returnedUser.CreatedAt,
		UpdatedAt:      returnedUser.UpdatedAt,
		Email:          returnedUser.Email,
	}

	respondWithJson(w, 200, jsonableUser)
}
