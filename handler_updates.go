package main

import (
	"net/http"
	"encoding/json"
	"github.com/JRPOGM/Chirpy/internal/database"
	"github.com/JRPOGM/Chirpy/internal/auth"
)

func (cfg * apiConfig) handlerUpdates(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string	`json:"password"`
		Email		string	`json:"email"`
	}
	type response struct {
		User
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
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		//status code 500
		return
	}
	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		//status code 500
		return
	}
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:				userID,
		Email:          params.Email,
		HashedPassword: hashed,
	}) //cfg.db.UpdateUser and database.UpdateUserParams pulled from users.sql.go
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user information", err)
		//status code 500
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt:	user.UpdatedAt,
			Email: 		user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	}) //status code 200
}