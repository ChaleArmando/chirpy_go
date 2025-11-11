package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ChaleArmando/chirpy_go/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	refToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, "Refresh Token not found", err)
		return
	}

	if time.Now().Compare(refToken.ExpiresAt) != -1 || refToken.RevokedAt.Valid {
		respondJsonError(w, http.StatusUnauthorized, "Refresh Token no longer valid", errors.New("refresh token no longer valid, token expired or was revoked"))
		return
	}

	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(refToken.UserID, cfg.secret, expirationTime)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, "Failed creating Access Token", err)
		return
	}

	jsonResp := response{
		Token: accessToken,
	}
	respondJson(w, http.StatusOK, jsonResp)
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	_, err = cfg.dbQueries.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "Failed to update Refresh Token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
