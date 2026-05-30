package auth

import "github.com/golang-jwt/jwt/v4"

type ModelClaim struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type ModelRead struct {
	Token string `json:"token"`
}

type ModelCreate struct {
	Email string `json:"email"`
}
