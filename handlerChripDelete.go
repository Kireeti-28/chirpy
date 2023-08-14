package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kireeti-28/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "chripId"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := auth.GetToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := auth.ValidateAccessToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	chrip, err := cfg.DB.GetChripById(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	userID, _ := strconv.Atoi(userId)

	if chrip.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "invalid user")
		return
	}

	err = cfg.DB.DeleteChirp(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}