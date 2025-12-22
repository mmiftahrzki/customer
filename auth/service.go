package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

type contextKey int
type service struct {
	signingKey []byte
	log        *logrus.Entry
}

const JWTContextKey contextKey = iota
const RequestHeaderAuthKey string = "Authorization"

var errEmptyAuth = errors.New("auth: authorization header not found")
var errInvalidAuth = errors.New("auth: authorization header invalid")

func newService(signingKey []byte) service {
	return service{
		signingKey: signingKey[:],
		log:        logger.GetLogger().WithField("component", "auth/service"),
	}
}

func extractAuthTokenStr(auth_value string) (string, error) {
	var token_str string

	if len(auth_value) == 0 {
		return token_str, errEmptyAuth
	}

	auth_value_fields := strings.Fields(auth_value)
	if len(auth_value_fields) != 2 || auth_value_fields[0] != "Bearer" {
		return token_str, errInvalidAuth
	}

	token_str = auth_value_fields[1]

	return token_str, nil
}

func (s *service) generateJWT(payload ModelCreate) (string, error) {
	registerdClaims := jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute))}
	claim := ModelClaim{
		Email:            payload.Email,
		RegisteredClaims: registerdClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedJWTString, err := token.SignedString(s.signingKey)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidKeyType) {
			return "", jwt.ErrInvalidKeyType
		}

		s.log.Error(err)

		return "", fmt.Errorf("error occured when try to stringify signed jwt: %w", err)
	}

	return signedJWTString, nil
}

func (s *service) getToken(tokenString string) (*jwt.Token, error) {
	keyFunc := func(t *jwt.Token) (any, error) {
		method, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		if method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		return s.signingKey[:], nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &ModelClaim{}, keyFunc)
	if err != nil {
		return token, err
	}

	return token, nil
}
