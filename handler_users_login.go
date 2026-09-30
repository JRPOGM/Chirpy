package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/JRPOGM/Chirpy/internal/auth"
	"github.com/JRPOGM/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password			string	`json:"password"`
		Email				string	`json:"email"`
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
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		//status code 500
		return
	}
	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:		user.ID,
		Token:		refreshToken,
		ExpiresAt:	time.Now().UTC().Add(time.Hour * 24 * 60),
	}) //database pulled from refresh_tokens.sql.go
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		//status code 500
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email:		user.Email,
			IsChirpyRed: user.IsChirpyRed,
		}, //status code 200
		Token:			accessToken,
		RefreshToken:	refreshToken,
	}) 
}