package auth

import (
	"time"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	//argon2id.CreateHash(string, params) always returns a hash string and an error
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	//argoin2id.ComparePasswordAndHash(string, string) always returns a bool and an error
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
	//jwt.NewWithClaims(signing method, claims) jwt.SigningMethodHS256 seems to be a default
		Issuer:		string(TokenTypeAccess),
		IssuedAt:	jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt:	jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		//jwt.NewNumericDate(time.Time) to set an issue date, and .Add(time.Duration) for the second to set an expiration date
		Subject:	userID.String(),
	})
	return token.SignedString(signingKey)
	//token.SignedString(key) requires the ket to be a []byte because the signing method is HS256 (HMAC signing method)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	//calling jwt.RegisteredClaims to make an empty struct
	token, err := jwt.ParseWithClaims(tokenString, &claimsStruct, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	//jwt.ParseWithClaims(string, claims, Keyfunc) takes a string, a call to a struct, and a key function either defined or named
	if err != nil {
		return uuid.Nil, err
	}
	userIDString, err := token.Claims.GetSubject()
	//token.Claims.GetSubject(string, error) generates a string and an error when called to a made token
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := token.Claims.GetIssuer()
	//token.Claims.GetIssuer(string, error) generates a string and an error when called to a made token
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	} //if the issuer string does not match the admin key
	id, err := uuid.Parse(userIDString)
	//uuid.Parse(string) decodes string into a UUID or returns an error if it cannot be parsed. 
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	//header.Get(string) created with an http.Header gets the first value associated with the given key
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authroization header")
	}
	return splitAuth[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
    apiKey := headers.Get("Authorization")
    if apiKey == "" {
        return "", ErrNoAuthHeaderIncluded
    }
    splitKey := strings.Split(apiKey, " ")
    if len(splitKey) < 2 || splitKey[0] != "ApiKey" {
        return "", errors.New("invalid authprization string")
    }
   return splitKey[1], nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	//creates a []byte of 32 entries, equal to 256 bit
	rand.Read(token)
	//rand.Read([]byte) fills byte with cryptographically secure random bytes
	return hex.EncodeToString(token)
	//hex.EncodeToString([]byte) returns the hexadecimal encoding of the byte
}