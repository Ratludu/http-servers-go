package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ratludu/http-servers-go/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expiry   *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode the json", err)
		return
	}

	dbUser, err := cfg.db.FindUserFromEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPasswords)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	expiry := time.Hour
	if params.Expiry != nil {
		expiry = time.Duration(*params.Expiry) * time.Second
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwt, time.Duration(expiry))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token", err)
		return

	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, user)

}
