package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ChaleArmando/chirpy_go/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed decoding parameters", err)
		return
	}
	if len(params.Body) > 140 {
		respondJsonError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   replaceBadWords(params.Body),
		UserID: params.UserId,
	}
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed creating chirp", err)
		return
	}

	newChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondJson(w, http.StatusCreated, newChirp)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed to transform id into uuid, invalid id", err)
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondJsonError(w, http.StatusNotFound, "failed to get chirp", err)
		return
	}

	jsonChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}

	respondJson(w, http.StatusOK, jsonChirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed to get chirps", err)
		return
	}
	jsonChirps := []Chirp{}
	for _, dbChirp := range chirps {
		jsonChirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		}
		jsonChirps = append(jsonChirps, jsonChirp)
	}

	respondJson(w, http.StatusOK, jsonChirps)
}

func replaceBadWords(respBody string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(respBody)
	lowerWords := strings.Fields(strings.ToLower(respBody))
	for id, word := range lowerWords {
		if slices.Contains(badWords, word) {
			words[id] = "****"
		}
	}
	return strings.Join(words, " ")
}
