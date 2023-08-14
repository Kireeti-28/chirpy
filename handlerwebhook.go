package main

import (
	"encoding/json"
	"net/http"

	"github.com/kireeti-28/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserId int `json:"user_id"`
	}

	type Parameters struct {
		Event string `json:"event"`
		Data `json:"data"`
	}

	_, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	paramters := Parameters{}
	err = decoder.Decode(&paramters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if paramters.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}
	err = cfg.DB.UpdateMemberShipRed(paramters.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
