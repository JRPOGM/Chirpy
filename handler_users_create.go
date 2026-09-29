package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/JRPOGM/Chirpy/internal/auth"
	"github.com/JRPOGM/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
	Password	string		`json:"password"`
}
//create a struct copy of the users table
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
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
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}//status code 500
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:			params.Email,
		HashedPassword:	hashedPassword,
	})
	//cfg.db.CreateUser calls the create query function in users.sql.go file
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		//status code 500
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:			user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email:		user.Email,
		},
	})
	//StatusCreated is status code 201
}