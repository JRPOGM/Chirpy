package main

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpIdentification(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	//http.Request.PathValue(string) returns the value for the named path wildcard in the ServeMux pattern that matched the request.
	chirpID, err := uuid.Parse(chirpIDString)
	//uuid.Parse(string) decodes string into a UUID or returns an error if it cannot be parsed. 
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		//status code 400
		return
	}
	allChirp, err := cfg.db.GetChirpID(r.Context(), chirpID)
	//gets GetChirpID from internal database chirps.sql.go
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find chirp ID", err)
		//status code 404
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:			allChirp.ID,
		CreatedAt:	allChirp.CreatedAt,
		UpdatedAt:	allChirp.UpdatedAt,
		Body:		allChirp.Body,
		UserID:		allChirp.UserID,
	}) //status code 200
}