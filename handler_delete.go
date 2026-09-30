package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/JRPOGM/Chirpy/internal/database"
	"github.com/JRPOGM/Chirpy/internal/auth"
)

func (cfg * apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"id"`
	}
	type response struct {
		User
	}

	user, err := cfg.db.DeleteChirp(r.Context(), params.ID)
	if err != nil {

	}
}