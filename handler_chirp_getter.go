package main


import (
    "net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
    allChirps, err :=cfg.db.GetChirps(r.Context())
	//calls GetChirps from chirps.sql.go
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
        //status code 500
        return
    }
	chirps := []Chirp{}
	//creates an array of the Chirp struct
	for _, allChirp := range allChirps {
		chirps = append(chirps, Chirp{
			ID:         allChirp.ID,
        	CreatedAt:  allChirp.CreatedAt,
        	UpdatedAt:  allChirp.UpdatedAt,
        	Body:       allChirp.Body,
        	UserID:     allChirp.UserID,
		})
	} //appends an array to an array
    respondWithJSON(w, http.StatusOK, chirps)
	//status code 200
}