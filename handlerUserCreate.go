package main

import (
	"encoding/json"
	"net/http"

	"github.com/kireeti-28/chirpy/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	data := database.UserReq{}
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode request")
		return
	}

	user, err := cfg.DB.CreateUser(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type UserResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	userResp := UserResp{
		ID:    user.ID,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusCreated, userResp)
}
