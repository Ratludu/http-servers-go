package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ratludu/http-servers-go/internal/auth"
	"github.com/ratludu/http-servers-go/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

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
	token, err := auth.MakeJWT(dbUser.ID, cfg.jwt, time.Duration(expiry))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate token", err)
		return

	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate refresh token", err)
		return

	}

	expiryRefresh := 60 * 24 * time.Hour
	_, err = cfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Duration(expiryRefresh)),
	})

	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
		ChirpyRed:    dbUser.IsChirpyRed.Bool,
	}

	respondWithJSON(w, http.StatusOK, user)

}
