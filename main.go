package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	fs := http.FileServer(http.Dir("."))
	fileServerHandler := http.StripPrefix("/app", fs)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpyHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())

}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirpyHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnValue struct {
		CleanedBody string `json:"cleaned_body"`
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

	respBody := returnValue{
		CleanedBody: replaceBadWords(params.Body),
	}
	respondJson(w, http.StatusOK, respBody)
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
