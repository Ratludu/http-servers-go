package main

import (
	"context"
	"net/http"

	"github.com/ratludu/http-servers-go/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	_, err = cfg.db.RevokeToken(context.Background(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
