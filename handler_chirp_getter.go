package main


import (
    "net/http"
	"sort"
	"github.com/google/uuid"
	"github.com/JRPOGM/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
    var allChirps []database.Chirp
	authorID, err := authorIDRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		//status code 400
		return
	}
	if authorID != uuid.Nil {
		allChirps, err = cfg.db.GetChirpAuthor(r.Context(), authorID)
		//calls GetChirpAuthor from chirps.sql.go
	} else {
		allChirps, err = cfg.db.GetChirps(r.Context())
		//calls GetChirps from chirps.sql.go
	}
	if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
        //status code 500
        return
    }
	sorting := "asc"
	sortingParams := r.URL.Query().Get("sort")
	if sortingParams == "desc" {
		sorting = "desc"
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
	sort.Slice(chirps, func(i, j int) bool {
		if sorting == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
    respondWithJSON(w, http.StatusOK, chirps)
	//status code 200
}

func authorIDRequest(r *http.Request) (uuid.UUID, error) {
	authorID := r.URL.Query().Get("author_id")
	//returns a string that contains the provided Get() value
	if authorID == "" {
		return uuid.Nil, nil
	}
	auID, err := uuid.Parse(authorID)
	if err != nil {
		return uuid.Nil, err
	}
	return auID, nil
}