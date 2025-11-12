package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ChaleArmando/chirpy_go/internal/auth"
	"github.com/ChaleArmando/chirpy_go/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed decoding parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed hashing password", err)
		return
	}

	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), userParams)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed creating user", err)
		return
	}
	newUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondJson(w, http.StatusCreated, newUser)
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed decoding parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed hashing password", err)
		return
	}

	userParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	}
	user, err := cfg.dbQueries.UpdateUser(r.Context(), userParams)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed updating user", err)
		return
	}
	dbUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondJson(w, http.StatusOK, dbUser)
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed decoding parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	passwordMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !passwordMatch {
		respondJsonError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed creating Access Token", err)
		return
	}

	refToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed creating Refresh Token", err)
		return
	}
	refTokenParams := database.CreateRefreshTokenParams{
		Token:  refToken,
		UserID: user.ID,
	}
	dbRefreshToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), refTokenParams)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, "failed saving Refresh Token in DB", err)
		return
	}

	jsonResp := response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: dbRefreshToken.Token,
	}
	respondJson(w, http.StatusOK, jsonResp)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("execute only in development environment"))
		return
	}
	cfg.fileserverHits.Swap(0)
	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to reset users: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File Hits reset to 0 and Database reset to initial state"))
}
