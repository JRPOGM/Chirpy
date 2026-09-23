package main


import (
    "encoding/json"
    "net/http"
    "strings"
)


func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"body"`
    }
    type returnVals struct {
        CleanedBody string `json:"cleaned_body"`
    }
    decoder := json.NewDecoder(r.Body)
    //json.NewDecoder(r io.Reader) creates a decoder from a pointed Request source
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }
    const maxChirpLength = 140
    if len(params.Body) > maxChirpLength {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
        return
    }
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
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