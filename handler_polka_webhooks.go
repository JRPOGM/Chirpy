package main

import (
	"net/http"
	"encoding/json"
	"errors"
	"database/sql"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		//status code 500
		return
	}
	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Parameter event wrong", err)
		//status code 204
		return
	}
	_, err = cfg.db.UpgradeRedChirps(r.Context(), params.Data.UserID)
	//cfg.db.UpgradeRedChirps taken from users.sql.go
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//errors.Is
			respondWithError(w, http.StatusNotFound, "Could not find user", err)
			//status code 404
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not upgrade chirp", err)
		//status code 500
		return
	}
	w.WriteHeader(http.StatusNoContent)
	//status code 204
}