package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	sAuthorId := r.URL.Query().Get("author_id")
	sSortBy := r.URL.Query().Get("sort")
	if sAuthorId != "" {

		authorId, err := strconv.Atoi(sAuthorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		chirps, err := cfg.DB.GetChirpByAuthorId(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		sort.Slice(chirps, func(i, j int) bool {
			if sSortBy == "desc" {
				return chirps[i].ID > chirps[j].ID
			}
			return chirps[i].ID < chirps[j].ID
		})

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sSortBy == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChripById(w http.ResponseWriter, r *http.Request) {
	chripId, err := strconv.Atoi(chi.URLParam(r, "chirpId"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to read id from url params")
		return
	}

	chirp, err := cfg.DB.GetChripById(chripId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       chirp.ID,
		Body:     chirp.Body,
		AuthorID: chirp.AuthorID,
	})
}
