package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kireeti-28/chirpy/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode request")
		return
	}

	user, err := cfg.DB.GetUserByEmail(userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get user by email")
		return
	}

	err = auth.ComparePassword(user.Password, userReq.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := auth.CreateAccessToken(cfg.jwtSecret, fmt.Sprint(user.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := auth.CreateRefreshToken(cfg.jwtSecret, fmt.Sprint(user.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshTokenSlice = append(refreshTokenSlice, refreshToken) // in-memory

	type LoginResp struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	loginResp := LoginResp{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed: user.IsChirpyRed,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, loginResp)
}
