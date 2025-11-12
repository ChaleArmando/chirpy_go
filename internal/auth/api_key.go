package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("auth header not found")
	}
	splitAuthHeader := strings.Fields(apiKey)
	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "ApiKey" {
		return "", errors.New("authorization header malformed")
	}

	return strings.TrimSpace(splitAuthHeader[1]), nil
}
