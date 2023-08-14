package auth

import (
	"errors"
	"net/http"
	"strings"
)


func GetApiKey(headers http.Header) (string, error) {
	authSlice := strings.Split(headers.Get("Authorization"), " ")

	if authSlice[0] != "ApiKey" {
		return "", errors.New("malformed apikey")
	}

	return authSlice[1], nil
}