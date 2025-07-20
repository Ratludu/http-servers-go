package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ratludu/http-servers-go/internal/auth"
)

type RefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	tokenData, err := cfg.db.GetRefreshToken(context.Background(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	newToken, err := auth.MakeJWT(tokenData.UserID, cfg.jwt, time.Duration(time.Hour))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	respondWithJSON(w, http.StatusOK, RefreshToken{
		Token: newToken,
	})

}
