package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/JRPOGM/Chirpy/internal/auth"
)

func (cfg * apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	//http.Request.PathValue(string) returns the value for the named path wildcard in the ServeMux pattern that matched the request.
	chirpID, err := uuid.Parse(chirpIDString)
	//uuid.Parse(string) decodes string into a UUID or returns an error if it cannot be parsed. 
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		//status code 400
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	//pulls GetBearerToken from internal/auth.go
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get token", err)
		//status code 401
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate secret", err)
		//status code 401
		return
	}
	dataChirp, err := cfg.db.GetChirpID(r.Context(), chirpID)
	//cfg.db.GetChirp pulled from chirps.sql.go
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get chirp", err)
		//status code 404
		return
	}
	if dataChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Cannot delete chirp", err)
		//status code 403
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	//cfg.db.DeleteChirp pulled from chirps.sql.go
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp", err)
		//status code 500
		return
	}
	w.WriteHeader(http.StatusNoContent) //status code 204
}