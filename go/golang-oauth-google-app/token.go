package main

import (
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

type UserClaims struct {
	jwt.StandardClaims
	SessionID string `json:"session_id"`
}

// It is recommended taht you use a custom Valid method for the UserClaims struct
func (uc *UserClaims) Valid() error {
	if len(uc.SessionID) == 0 {
		return fmt.Errorf("ERROR - JWT token does not have a session id")
	}
	if !uc.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return fmt.Errorf("ERROR - JWT token has expired")
	}
	return nil
}

// Get a signed jwt token based on a session id
func createToken(sessionID string) (string, error) {
	claim := &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(sessionExpiry).Unix(),
		},
		SessionID: sessionID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// The key to sign this token needs to be of 32 chars based on the algo chosen here
	return token.SignedString(jwtSigningSecretKey)
}

// Validate the jwt token and get the session id from the token
func parseToken(token string) (string, error) {
	if len(token) == 0 {
		return "", nil
	}
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("ERROR - The signing algo of the token and what we expected do not match, so cant use the shared key to verify signature")
		}
		return jwtSigningSecretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("ERROR - Encountered an error when parsing the token - %v", err)
	}
	if !t.Valid {
		return "", fmt.Errorf("ERROR - Token is not valid or has expired")
	}
	claims := t.Claims.(*UserClaims)
	return claims.SessionID, nil
}
