package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ratludu/http-servers-go/internal/auth"
	"github.com/ratludu/http-servers-go/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode the json", err)
		return
	}

	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Email:           params.Email,
		HashedPasswords: hashedPw,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user in database", err)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)

}
