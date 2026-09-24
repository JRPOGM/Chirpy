package main


import (
    "encoding/json"
    "net/http"
    "errors"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/JRPOGM/Chirpy/internal/database"
)

type Chirp struct {
    ID          uuid.UUID   `json:"id"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
    Body        string      `json:"body"`
    UserID      uuid.UUID   `json:"uder_id"`
}

func (cfg *apiConfig) handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body    string          `json:"body"`
        UserID  uuid.UUID          `json:"user_id"`
    }
    decoder := json.NewDecoder(r.Body)
    //json.NewDecoder(r io.Reader) creates a decoder from a pointed Request source
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        //status code 500
        return
    }
    cleaned, err := validateChirp(params.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error(), err)
        //status code 400
        return
    }
    chirp, err :=cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
        Body:   cleaned,
        UserID: params.UserID,
    }) //calls CreateChirp and the CreateChirpParams from the internal database models.go and chirps.sql.go
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
        //status code 500
        return
    }
    respondWithJSON(w, http.StatusCreated, Chirp{
        ID:         chirp.ID,
        CreatedAt:  chirp.CreatedAt,
        UpdatedAt   chirp.UpdatedAt,
        Body:       chirp.Body,
        UserID:     chirp.UserID,
    }) //status code 201
}

func validateChirp(body string) (string, error) { 
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        return "", errors.New("Chirp is too long")
    }
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
    //strings.Split(s, sep string) slices an 's' string separated by the 'sep' separator
	for i, word := range words {
		loweredWord := strings.ToLower(word)
        // strings.ToLower(string) maps all letters to their lowercase version
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
    //strings.Join(elems []string, sep string) concatenares elements of the first argument into a single string
    //sep separator is placed between elements in the resulting string
	return cleaned
}