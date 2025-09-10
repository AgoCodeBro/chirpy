package main

import (
	"net/http"
	"time"

	"github.com/AgoCodeBro/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {
	type jsonableToken struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "failed to get refresh token", err)
		return
	}

	refreshToken, err := cfg.db.GetToken(req.Context(), token)
	if err != nil {
		respondWithError(w, 401, "invalid refresh token", err)
		return
	} else if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "expired token", err)
		return
	} else if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "revoked token", nil)
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "failed to create jwt", err)
		return
	}

	result := jsonableToken{
		Token: accessToken,
	}

	respondWithJson(w, 200, result)
}
