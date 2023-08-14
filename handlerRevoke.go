package main

import (
	"net/http"

	"github.com/kireeti-28/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	_, err = auth.ValidateRefreshToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	temp := []string{}
	for _, refToken := range refreshTokenSlice {
		if token != refToken {
			temp = append(temp, refToken)
		}
	}
	refreshTokenSlice = temp

	respondWithJSON(w, http.StatusOK, nil)
}