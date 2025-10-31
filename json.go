package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondJsonError(w http.ResponseWriter, code int, errorMsg string, err error) {
	type responseError struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Println(err)
	}
	respBody := responseError{
		Error: errorMsg,
	}
	respondJson(w, code, respBody)
}

func respondJson(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Marshalling json failed: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
