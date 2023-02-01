package config

import "github.com/golang-jwt/jwt/v4"

type JwtCustomClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}
