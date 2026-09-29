package main

import (
	"net/http"
	"encoding/json"
	"github.com/JRPOGM/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string	`json:"password"`
		Email		string	`json:"email"`
	}
	type response struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		//status code 500
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	//cfg.db.GetUserByEmail pulled from users.sql.go
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		//status code 401
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	//params.Password found within struct, user.HashedPassword pulled from users.sql.go where GetUserByEmail is from
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Password and hash do not match", err)
		//status code 401
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email:		user.Email,
		},
	}) //status code 200
}