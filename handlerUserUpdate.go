package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kireeti-28/chirpy/internal/auth"
	"github.com/kireeti-28/chirpy/internal/database"
)

func (cfg *apiConfig) userUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userIdStr, err := auth.ValidateAccessToken(cfg.jwtSecret, token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	userReq := database.UserReq{}
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	userIdInt, err := strconv.Atoi(userIdStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	user, err := cfg.DB.UpdateUser(userIdInt, userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	type UpdateResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	updateResp := UpdateResp{
		ID:    user.ID,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, updateResp)
}
