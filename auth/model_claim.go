package auth

import "github.com/golang-jwt/jwt/v4"

type AuthClaimModel struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
