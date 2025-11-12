package main

import (
	"encoding/json"
	"net/http"

	"github.com/ChaleArmando/chirpy_go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, "couldn't find api key", err)
		return
	}
	if cfg.polkaKey != apiKey {
		respondJsonError(w, http.StatusUnauthorized, "api key is incorrect", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed decoding parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.dbQueries.GetUser(r.Context(), params.Data.UserID)
	if err != nil {
		respondJsonError(w, http.StatusNotFound, "failed finding user", err)
		return
	}

	err = cfg.dbQueries.UpgradeChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed updating user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
