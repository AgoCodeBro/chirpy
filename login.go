package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AgoCodeBro/chirpy/internal/auth"
	"github.com/AgoCodeBro/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	type loginParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	loginJson := loginParams{}
	err := decoder.Decode(&loginJson)
	if err != nil {
		respondWithError(w, 500, "Failed to decode request", err)
		return
	}

	user, err := cfg.db.GetUser(req.Context(), loginJson.Email)
	if err != nil {
		respondWithError(w, 401, "email or password is incorrect", err)
		return
	}

	err = auth.CheckPasswordHash(loginJson.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "email or password is incorrect", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "failed to generate JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "failed to generate refresh token", err)
		return
	}

	err = cfg.registerRefreshToken(req.Context(), refreshToken, user.ID)
	if err != nil {
		respondWithError(w, 500, "failed to register refresh token", err)
		return
	}

	type jsonUserWithTokens struct {
		jsonUser
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	jsonableUser := jsonUserWithTokens{
		jsonUser: jsonUser{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJson(w, 200, jsonableUser)
}

func (cfg *apiConfig) registerRefreshToken(ctx context.Context, refreshToken string, userID uuid.UUID) error {
	duration := 60 * 24 * time.Hour
	expiresAt := time.Now().Add(duration)

	addRefreshTokenArgs := database.AddRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	_, err := cfg.db.AddRefreshToken(ctx, addRefreshTokenArgs)

	return err
}
