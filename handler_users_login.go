package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/JRPOGM/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password			string	`json:"password"`
		Email				string	`json:"email"`
		ExpiresInSeconds	int		`json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token			string 	`json:"token"`
		RefreshToken	string	`json:"refresh_token"`
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
	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email:		user.Email,
		}, //status code 200
		Token: accessToken,
	}) 
}