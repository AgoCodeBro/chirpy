package main

import (
	"net/http"

	"github.com/AgoCodeBro/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 500, "failed to get refresh token", err)
		return
	}

	err = cfg.db.RevokeToken(req.Context(), token)
	if err != nil {
		respondWithError(w, 500, "failed to revoke token", err)
		return
	}

	respondWithJson(w, 204, nil)
}
