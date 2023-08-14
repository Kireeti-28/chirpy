package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var InvalidAuthHeader = errors.New("Invalid Auth Header")

func CreateAccessToken(jwtSecret, userId string) (string, error) {
	signingKey := []byte(jwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		Subject: userId,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(time.Hour)).UTC()),
	})

	return token.SignedString(signingKey)
}

func CreateRefreshToken(jwtSecret, userId string) (string, error) {
	signingKey := []byte(jwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-refresh",
		Subject: userId,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(60 * 24 * time.Hour)).UTC()),
	})

	return token.SignedString(signingKey)
}

func GetToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	authSlice := strings.Split(authHeader, " ")

	if len(authSlice) < 2 || authSlice[0] != "Bearer" {
		return "", InvalidAuthHeader
	}

	return authSlice[1], nil
}

func ValidateAccessToken(tokenSecret, tokenString string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer == "chirpy-refresh" {
		return "", errors.New("wrong token")
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	expriesAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", err
	}

	if expriesAt.Before(time.Now().UTC()) {
		return "", errors.New("Token was expired")
	}

	return userId, nil
}
func ValidateRefreshToken(tokenSecret, tokenString string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer == "chirpy-access" {
		return "", errors.New("wrong token")
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	expriesAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", err
	}

	if expriesAt.Before(time.Now().UTC()) {
		return "", errors.New("Token was expired")
	}

	return userId, nil
}