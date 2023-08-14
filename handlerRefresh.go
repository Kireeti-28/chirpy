package main

import (
	"net/http"

	"github.com/kireeti-28/chirpy/internal/auth"
)

var refreshTokenSlice = []string{}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := auth.ValidateRefreshToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	exists := false
	for _, refToken := range refreshTokenSlice {
		if token == refToken {
			exists = true
		}
	}

	if !exists {
		respondWithError(w, http.StatusUnauthorized, "refresh token was revoked")
		return
	}

	newToken, err := auth.CreateAccessToken(cfg.jwtSecret, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	type Response struct {
		Token string `json:"token"`
	}

	response := Response{
		Token: newToken,
	}

	respondWithJSON(w, http.StatusOK, response)
}